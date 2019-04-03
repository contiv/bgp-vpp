// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	nodemodel "github.com/contiv/vpp/plugins/ksr/model/node"
	podmodel "github.com/contiv/vpp/plugins/ksr/model/pod"
	vppnodemodel "github.com/contiv/vpp/plugins/nodesync/vppnode"
	"github.com/gogo/protobuf/proto"
	"github.com/ligato/cn-infra/datasync"
	"reflect"

	"fmt"
	"strings"
)

type dataChangeProcessor interface {
	GetNames(key string) ([]string, error)
	GetValueProto() proto.Message
	AddRecord(ctc *ContivTelemetryCache, names []string, record proto.Message) error
	UpdateRecord(ctc *ContivTelemetryCache, names []string, oldRecord proto.Message, newRecord proto.Message) error
	DeleteRecord(ctc *ContivTelemetryCache, names []string) error
}

// dataChangeProcessor implementation for K8s pod data
type podChange struct{}

func (pc *podChange) GetNames(key string) ([]string, error) {
	pod, namespace, err := podmodel.ParsePodFromKey(key)
	return []string{pod, namespace}, err
}

func (pc *podChange) GetValueProto() proto.Message {
	return &podmodel.Pod{}
}

func (pc *podChange) AddRecord(ctc *ContivTelemetryCache, names []string, record proto.Message) error {
	ctc.Log.Infof("Adding K8s pod %s in namespace %s, podValue %+v", names[0], names[1], record)
	if pod, ok := record.(*podmodel.Pod); ok {
		return ctc.K8sCache.CreatePod(pod.Name, pod.Namespace, pod.Label, pod.IpAddress,
			pod.HostIpAddress, pod.Container)
	}
	return fmt.Errorf("bad record type %s", reflect.TypeOf(record))
}

func (pc *podChange) UpdateRecord(ctc *ContivTelemetryCache, names []string, _, record proto.Message) error {
	ctc.Log.Infof("Updating K8s pod %s in namespace %s, prevPodValue %+v", names[0], names[1], record)

	if pod, ok := record.(*podmodel.Pod); ok {
		return ctc.K8sCache.UpdatePod(pod.Name, pod.Namespace, pod.Label, pod.IpAddress,
			pod.HostIpAddress, pod.Container)
	}
	return fmt.Errorf("bad record type %s", reflect.TypeOf(record))
}

func (pc *podChange) DeleteRecord(ctc *ContivTelemetryCache, names []string) error {
	ctc.Log.Infof("Deleting K8s pod %s in namespace %s", names[0], names[1])
	return ctc.K8sCache.DeletePod(names[0])
}

// dataChangeProcessor implementation for K8s node data
type nodeChange struct{}

func (nc *nodeChange) GetNames(key string) ([]string, error) {
	node, err := nodemodel.ParseNodeFromKey(key)
	return []string{node}, err
}

func (nc *nodeChange) GetValueProto() proto.Message {
	return &nodemodel.Node{}
}

func (nc *nodeChange) AddRecord(ctc *ContivTelemetryCache, names []string, record proto.Message) error {
	ctc.Log.Infof("Adding K8s node %s, nodeValue %+v", names[0], record)

	if node, ok := record.(*nodemodel.Node); ok {
		return ctc.K8sCache.CreateK8sNode(node.Name, node.Pod_CIDR, node.Provider_ID, node.Addresses, node.NodeInfo)
	}
	return fmt.Errorf("bad record type %s", reflect.TypeOf(record))
}

func (nc *nodeChange) UpdateRecord(ctc *ContivTelemetryCache, names []string, _, newRecord proto.Message) error {
	ctc.Log.Infof("Updating K8s node %s, nodeValue %+v", names[0], newRecord)

	if node, ok := newRecord.(*nodemodel.Node); ok {
		return ctc.K8sCache.UpdateK8sNode(node.Name, node.Pod_CIDR, node.Provider_ID, node.Addresses, node.NodeInfo)
	}
	return fmt.Errorf("bad record type %s", reflect.TypeOf(newRecord))
}

func (nc *nodeChange) DeleteRecord(ctc *ContivTelemetryCache, names []string) error {
	ctc.Log.Infof("Deleting K8s node %s", names[0])
	return ctc.K8sCache.DeleteK8sNode(names[0])
}

// dataChangeProcessor implementation for vppNode data
type vppNodeChange struct{}

func (nic *vppNodeChange) GetNames(key string) ([]string, error) {
	nodeParts := strings.Split(key, "/")
	if len(nodeParts) != 2 {
		return nil, fmt.Errorf("invalid nodeLiveness key %s", key)
	}
	return []string{nodeParts[1]}, nil
}

func (nic *vppNodeChange) GetValueProto() proto.Message {
	return &vppnodemodel.VppNode{}
}

func (nic *vppNodeChange) AddRecord(ctc *ContivTelemetryCache, names []string, record proto.Message) error {
	ctc.Log.Infof("Adding VPP node, names %+v, nodeValue %+v", names, record)

	if ni, ok := record.(*vppnodemodel.VppNode); ok {
		return ctc.VppCache.CreateNode(ni.Id, ni.Name, ni.IpAddress)
	}
	return fmt.Errorf("bad record type %s", reflect.TypeOf(record))
}

func (nic *vppNodeChange) UpdateRecord(ctc *ContivTelemetryCache, names []string, _, newRecord proto.Message) error {
	ctc.Log.Infof("Updating VPP node, names %+v, nodeInfoValue %+v", names, newRecord)

	if ni, ok := newRecord.(*vppnodemodel.VppNode); ok {
		return ctc.VppCache.UpdateNode(ni.Id, ni.Name, ni.IpAddress)
	}
	return fmt.Errorf("bad record type %s", reflect.TypeOf(newRecord))
}

func (nic *vppNodeChange) DeleteRecord(ctc *ContivTelemetryCache, names []string) error {
	ctc.Log.Infof("Deleting VPP node, names %+v", names)
	return ctc.VppCache.DeleteNode(names[0])
}

// Update sends the update event passed as an argument to the ctc telemetryCache
// thread, where it is processed in the function below (update). )
func (ctc *ContivTelemetryCache) Update(dataChngEv datasync.ChangeEvent) error {
	ctc.dsUpdateChannel <- dataChngEv
	return nil
}

// Update processes a data sync change event associated with K8s State data.
// The change is applied into the cache and all subscribed watchers are
// notified.
// The function will forward any error returned by a watcher.
func (ctc *ContivTelemetryCache) update(dataChngEv datasync.ProtoWatchResp) error {
	err := error(nil)
	key := dataChngEv.GetKey()
	var dcp dataChangeProcessor

	// Determine which data is changing
	switch {
	case strings.HasPrefix(key, vppnodemodel.KeyPrefix):
		dcp = &vppNodeChange{}

	case strings.HasPrefix(key, nodemodel.KeyPrefix()):
		dcp = &nodeChange{}

	case strings.HasPrefix(key, podmodel.KeyPrefix()):
		dcp = &podChange{}

	default:
		return fmt.Errorf("unknown DATA CHANGE key %s", key)
	}

	names, err := dcp.GetNames(key)
	if err != nil {
		return err
	}

	// Determine the type of & perform the data change operation
	switch dataChngEv.GetChangeType() {
	case datasync.Delete:
		err = dcp.DeleteRecord(ctc, names)

	case datasync.Put:
		newRecord := dcp.GetValueProto()
		if err := dataChngEv.GetValue(newRecord); err != nil {
			err = fmt.Errorf("could not get new proto data for key %s, error %s", key, err)
			break
		}

		oldRecord := dcp.GetValueProto()
		exists, err := dataChngEv.GetPrevValue(oldRecord)
		if err != nil {
			err = fmt.Errorf("could not get previous proto data for key %s, error %s", key, err)
			break
		}

		if exists {
			err = dcp.UpdateRecord(ctc, names, oldRecord, newRecord)
		} else {
			err = dcp.AddRecord(ctc, names, newRecord)
		}

	default:
		err = fmt.Errorf("unknown event change type %+v", dataChngEv.GetChangeType())
	}

	return err
}
