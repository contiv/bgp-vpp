package bgp

import (
	"github.com/ligato/cn-infra/db/keyval/etcd"
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/ligato/vpp-agent/plugins/kvscheduler"
)

// DefaultPlugin is a default instance of IfPlugin.
var DefaultPlugin = *NewPlugin()

// NewPlugin creates a new Plugin with the provides Options
func NewPlugin(opts ...Option) *BgpPlugin {
	p := &BgpPlugin{}

	p.PluginName = "bgp-plugin"
	p.Rest = &rest.DefaultPlugin
	p.KVScheduler = &kvscheduler.DefaultPlugin

	p.KVStore = &etcd.DefaultPlugin
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
