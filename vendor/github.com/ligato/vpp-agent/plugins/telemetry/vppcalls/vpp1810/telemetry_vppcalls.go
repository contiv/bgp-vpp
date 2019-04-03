//  Copyright (c) 2019 Cisco and/or its affiliates.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at:
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package vpp1810

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	govppapi "git.fd.io/govpp.git/api"

	vpevppcalls "github.com/ligato/vpp-agent/plugins/govppmux/vppcalls"
	"github.com/ligato/vpp-agent/plugins/telemetry/vppcalls"
	"github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp1810/memclnt"
	"github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp1810/vpe"
)

func init() {
	var msgs []govppapi.Message
	msgs = append(msgs, memclnt.Messages...)
	msgs = append(msgs, vpe.Messages...)

	vppcalls.Versions["vpp1810"] = vppcalls.HandlerVersion{
		Msgs: msgs,
		New: func(ch govppapi.Channel) vppcalls.TelemetryVppAPI {
			vpeHandler := vpevppcalls.CompatibleVpeHandler(ch)
			return &TelemetryHandler{ch, vpeHandler}
		},
	}
}

type TelemetryHandler struct {
	ch govppapi.Channel
	vpevppcalls.VpeVppAPI
}

var (
	// Regular expression to parse output from `show memory`
	memoryRe = regexp.MustCompile(`Thread\s+(\d+)\s+(\w+).?\s+(\d+) objects, ([\dkmg\.]+) of ([\dkmg\.]+) used, ([\dkmg\.]+) free, ([\dkmg\.]+) reclaimed, ([\dkmg\.]+) overhead, ([\dkmg\.]+) capacity`)
)

// GetMemory retrieves `show memory` info.
func (h *TelemetryHandler) GetMemory() (*vppcalls.MemoryInfo, error) {
	data, err := h.RunCli("show memory")
	if err != nil {
		return nil, err
	}

	var threads []vppcalls.MemoryThread

	threadMatches := memoryRe.FindAllStringSubmatch(string(data), -1)
	for _, matches := range threadMatches {
		fields := matches[1:]
		if len(fields) != 9 {
			return nil, fmt.Errorf("invalid memory data for thread: %q", matches[0])
		}
		id, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return nil, err
		}
		thread := &vppcalls.MemoryThread{
			ID:        uint(id),
			Name:      fields[1],
			Objects:   strToUint64(fields[2]),
			Used:      strToUint64(fields[3]),
			Total:     strToUint64(fields[4]),
			Free:      strToUint64(fields[5]),
			Reclaimed: strToUint64(fields[6]),
			Overhead:  strToUint64(fields[7]),
			Size:      strToUint64(fields[8]),
		}
		threads = append(threads, *thread)
	}

	info := &vppcalls.MemoryInfo{
		Threads: threads,
	}

	return info, nil
}

var (
	// Regular expression to parse output from `show node counters`
	nodeCountersRe = regexp.MustCompile(`^\s+(\d+)\s+([\w-\/]+)\s+(.+)$`)
)

// GetNodeCounters retrieves node counters info.
func (h *TelemetryHandler) GetNodeCounters() (*vppcalls.NodeCounterInfo, error) {
	data, err := h.RunCli("show node counters")
	if err != nil {
		return nil, err
	}

	var counters []vppcalls.NodeCounter

	for i, line := range strings.Split(string(data), "\n") {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Check first line
		if i == 0 {
			fields := strings.Fields(line)
			// Verify header
			if len(fields) != 3 || fields[0] != "Count" {
				return nil, fmt.Errorf("invalid header for `show node counters` received: %q", line)
			}
			continue
		}

		// Parse lines using regexp
		matches := nodeCountersRe.FindStringSubmatch(line)
		if len(matches)-1 != 3 {
			return nil, fmt.Errorf("parsing failed for `show node counters` line: %q", line)
		}
		fields := matches[1:]

		counters = append(counters, vppcalls.NodeCounter{
			Count:  strToUint64(fields[0]),
			Node:   fields[1],
			Reason: fields[2],
		})
	}

	info := &vppcalls.NodeCounterInfo{
		Counters: counters,
	}

	return info, nil
}

var (
	// Regular expression to parse output from `show runtime`
	runtimeRe = regexp.MustCompile(`(?:-+\n)?(?:Thread (\d+) (\w+)(?: \(lcore \d+\))?\n)?` +
		`Time ([0-9\.e]+), average vectors/node ([0-9\.e]+), last (\d+) main loops ([0-9\.e]+) per node ([0-9\.e]+)\s+` +
		`vector rates in ([0-9\.e]+), out ([0-9\.e]+), drop ([0-9\.e]+), punt ([0-9\.e]+)\n` +
		`\s+Name\s+State\s+Calls\s+Vectors\s+Suspends\s+Clocks\s+Vectors/Call\s+` +
		`((?:[\w-:\.]+\s+\w+(?:[ -]\w+)*\s+\d+\s+\d+\s+\d+\s+[0-9\.e]+\s+[0-9\.e]+\s+)+)`)
	runtimeItemsRe = regexp.MustCompile(`([\w-:\.]+)\s+(\w+(?:[ -]\w+)*)\s+(\d+)\s+(\d+)\s+(\d+)\s+([0-9\.e]+)\s+([0-9\.e]+)\s+`)
)

// GetRuntimeInfo retrieves how runtime info.
func (h *TelemetryHandler) GetRuntimeInfo() (*vppcalls.RuntimeInfo, error) {
	data, err := h.RunCli("show runtime")
	if err != nil {
		return nil, err
	}

	var threads []vppcalls.RuntimeThread

	threadMatches := runtimeRe.FindAllStringSubmatch(string(data), -1)
	for _, matches := range threadMatches {
		fields := matches[1:]
		if len(fields) != 12 {
			return nil, fmt.Errorf("invalid runtime data for thread: %q", matches[0])
		}
		thread := vppcalls.RuntimeThread{
			ID:                  uint(strToUint64(fields[0])),
			Name:                fields[1],
			Time:                strToFloat64(fields[2]),
			AvgVectorsPerNode:   strToFloat64(fields[3]),
			LastMainLoops:       strToUint64(fields[4]),
			VectorsPerMainLoop:  strToFloat64(fields[5]),
			VectorLengthPerNode: strToFloat64(fields[6]),
			VectorRatesIn:       strToFloat64(fields[7]),
			VectorRatesOut:      strToFloat64(fields[8]),
			VectorRatesDrop:     strToFloat64(fields[9]),
			VectorRatesPunt:     strToFloat64(fields[10]),
		}

		itemMatches := runtimeItemsRe.FindAllStringSubmatch(fields[11], -1)
		for _, matches := range itemMatches {
			fields := matches[1:]
			if len(fields) != 7 {
				return nil, fmt.Errorf("invalid runtime data for thread item: %q", matches[0])
			}
			thread.Items = append(thread.Items, vppcalls.RuntimeItem{
				Name:           fields[0],
				State:          fields[1],
				Calls:          strToUint64(fields[2]),
				Vectors:        strToUint64(fields[3]),
				Suspends:       strToUint64(fields[4]),
				Clocks:         strToFloat64(fields[5]),
				VectorsPerCall: strToFloat64(fields[6]),
			})
		}

		threads = append(threads, thread)
	}

	info := &vppcalls.RuntimeInfo{
		Threads: threads,
	}

	return info, nil
}

var (
	// Regular expression to parse output from `show buffers`
	buffersRe = regexp.MustCompile(`^\s+(\d+)\s+(\w+(?:[ \-]\w+)*)\s+(\d+)\s+(\d+)\s+([\dkmg\.]+)\s+([\dkmg\.]+)\s+(\d+)\s+(\d+).*$`)
)

// GetBuffersInfo retrieves buffers info
func (h *TelemetryHandler) GetBuffersInfo() (*vppcalls.BuffersInfo, error) {
	data, err := h.RunCli("show buffers")
	if err != nil {
		return nil, err
	}

	var items []vppcalls.BuffersItem

	for i, line := range strings.Split(string(data), "\n") {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Check first line
		if i == 0 {
			fields := strings.Fields(line)
			// Verify header
			if len(fields) != 8 || fields[0] != "Thread" {
				return nil, fmt.Errorf("invalid header for `show buffers` received: %q", line)
			}
			continue
		}

		// Parse lines using regexp
		matches := buffersRe.FindStringSubmatch(line)
		if len(matches)-1 != 8 {
			return nil, fmt.Errorf("parsing failed for `show buffers` line: %q", line)
		}
		fields := matches[1:]

		items = append(items, vppcalls.BuffersItem{
			ThreadID: uint(strToUint64(fields[0])),
			Name:     fields[1],
			Index:    uint(strToUint64(fields[2])),
			Size:     strToUint64(fields[3]),
			Alloc:    strToUint64(fields[4]),
			Free:     strToUint64(fields[5]),
			NumAlloc: strToUint64(fields[6]),
			NumFree:  strToUint64(fields[7]),
		})
	}

	info := &vppcalls.BuffersInfo{
		Items: items,
	}

	return info, nil
}

func strToFloat64(s string) float64 {
	// Replace 'k' (thousands) with 'e3' to make it parsable with strconv
	s = strings.Replace(s, "k", "e3", 1)
	s = strings.Replace(s, "m", "e6", 1)
	s = strings.Replace(s, "g", "e9", 1)

	num, err := strconv.ParseFloat(s, 10)
	if err != nil {
		return 0
	}
	return num
}

func strToUint64(s string) uint64 {
	return uint64(strToFloat64(s))
}
