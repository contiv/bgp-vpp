package bgp

import (
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/datasync/kvdbsync/local"
	"github.com/ligato/cn-infra/db/keyval/etcd"
	"github.com/ligato/cn-infra/health/statuscheck"
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/ligato/vpp-agent/plugins/kvscheduler"
	"github.com/ligato/vpp-agent/plugins/orchestrator"
)

// DefaultPlugin is a default instance of IfPlugin.
var DefaultPlugin = *NewPlugin()

// NewPlugin creates a new Plugin with the provides Options
func NewPlugin(opts ...Option) *BgpPlugin {
	p := &BgpPlugin{}

	p.PluginName = "bgp-plugin"
	p.Rest = &rest.DefaultPlugin

	p.Orchestrator = &orchestrator.DefaultPlugin
	p.Scheduler = &kvscheduler.DefaultPlugin
	p.ETCDDataSync = kvdbsync.NewPlugin(kvdbsync.UseKV(&etcd.DefaultPlugin))
	//initializing interfaces that allow us to write to ETCD and watch for changes
	writers := datasync.KVProtoWriters{
		p.ETCDDataSync,
	}
	statuscheck.DefaultPlugin.Transport = writers // Set watcher for KVScheduler.
	watchers := datasync.KVProtoWatchers{
		local.DefaultRegistry,
		p.ETCDDataSync,
	}
	orchestrator.DefaultPlugin.Watcher = watchers

	//p.Watcher = watchers
	//p.Publisher = writers

	p.Setup()

	for _, o := range opts {
		o(p)
	}

	return p
}

// Option is a function that can be used in NewPlugin to customize Plugin.
type Option func(*BgpPlugin)

// UseDeps returns Option that can inject custom dependencies.
func UseDeps(f func(*Deps)) Option {
	return func(p *BgpPlugin) {
		f(&p.Deps)
	}
}
