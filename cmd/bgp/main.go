package main

import (
	"github.com/contiv/bgp-vpp/plugins/bgp"
	"github.com/ligato/cn-infra/agent"
	"github.com/ligato/cn-infra/logging"
)

func main() {
	// Create an instance of our plugin using its constructor.
	p := bgp.NewPlugin()

	// Create new agent with our plugin instance.
	a := agent.NewAgent(agent.AllPlugins(p))

	// Run starts the agent with plugins, wait until shutdown
	// and then stops the agent and its plugins.
	if err := a.Run(); err != nil {
		logging.DefaultLogger.Error(err)
	}
}
