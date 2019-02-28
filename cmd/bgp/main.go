package main

import (
	"github.com/contiv/vpp/plugins/contivconf"
	"github.com/contiv/vpp/plugins/controller"
	"github.com/contiv/vpp/plugins/ipam"
	"github.com/contiv/vpp/plugins/ipv4net"
	"github.com/contiv/vpp/plugins/nodesync"
	"github.com/contiv/vpp/plugins/podmanager"
	"github.com/contiv/vpp/plugins/policy"
	"github.com/contiv/vpp/plugins/service"
	"github.com/contiv/vpp/plugins/statscollector"
	"os"

	"github.com/ligato/cn-infra/agent"

	//plugins that will need to be included in our app
	"github.com/ligato/cn-infra/logging/logmanager"
	"github.com/ligato/vpp-agent/plugins/orchestrator"
	"github.com/ligato/vpp-agent/plugins/kvscheduler"
	"github.com/ligato/vpp-agent/vendor/github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/rpc/rest"
)

// BgpAgent manages vswitch in contiv/vpp solution
type BgpAgent struct {
	LogManager  *logmanager.Plugin
	Orchestrator *orchestrator.Plugin
	KVScheduler *kvscheduler.Scheduler
	REST      *rest.Plugin
}