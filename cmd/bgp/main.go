package main

import (
	"github.com/contiv/bgp-vpp/plugins/bgp"
	"github.com/ligato/cn-infra/agent"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/datasync/kvdbsync/local"
	"github.com/ligato/cn-infra/datasync/resync"
	"github.com/ligato/cn-infra/db/keyval/etcd"
	"github.com/ligato/cn-infra/health/statuscheck"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/logging/logmanager"
	"github.com/ligato/vpp-agent/plugins/orchestrator"
)

type BgpAgent struct {
	LogManager *logmanager.Plugin
	BgpAgent   *bgp.BgpPlugin

	Orchestrator *orchestrator.Plugin
	ETCDDataSync *kvdbsync.Plugin
}

// New creates new SSwanVPPAgent instance.
func New() *BgpAgent {
	etcdDataSync := kvdbsync.NewPlugin(kvdbsync.UseKV(&etcd.DefaultPlugin))

	writers := datasync.KVProtoWriters{
		etcdDataSync,
	}
	statuscheck.DefaultPlugin.Transport = writers

	// Set watcher for KVScheduler.
	watchers := datasync.KVProtoWatchers{
		local.DefaultRegistry,
		etcdDataSync,
	}
	orchestrator.DefaultPlugin.Watcher = watchers

	return &BgpAgent{
		LogManager:   &logmanager.DefaultPlugin,
		Orchestrator: &orchestrator.DefaultPlugin,
		ETCDDataSync: etcdDataSync,
		BgpAgent:     &bgp.DefaultPlugin,
	}

}

func main() {
	// Create an instance of our plugin using its constructor.
	bgpAgent := New()

	// Create new agent with our plugin instance.
	a := agent.NewAgent(agent.AllPlugins(bgpAgent))

	// Run starts the agent with plugins, wait until shutdown
	// and then stops the agent and its plugins.
	if err := a.Run(); err != nil {
		logging.DefaultLogger.Error(err)
	}
}

// Init initializes main plugin.
func (ss *BgpAgent) Init() error {
	return nil
}

// AfterInit executes resync.
func (ss *BgpAgent) AfterInit() error {
	resync.DefaultPlugin.DoResync()
	return nil
}

// Close could close used resources.
func (ss *BgpAgent) Close() error {
	return nil
}

// String returns name of the plugin.
func (ss *BgpAgent) String() string {
	return "bgp-agent"
}
