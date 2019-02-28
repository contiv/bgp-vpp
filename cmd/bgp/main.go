package main

import (
	"github.com/contiv/bgp-vpp/plugins/bgp"
	"github.com/ligato/cn-infra/agent"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/logging/logmanager"
	"github.com/ligato/cn-infra/rpc/rest"
)


func main() {
	// Create an instance of our plugin using its constructor.
	p := bgp.NewPlugin()

	// Create new agent with our plugin instance.
	a := agent.NewAgent(agent.Plugins(p))

	// Run starts the agent with plugins, wait until shutdown
	// and then stops the agent and its plugins.
	if err := a.Run(); err != nil {
		logging.DefaultLogger.Error(err)
	}
}



// BgpAgent manages vswitch in contiv/vpp solution
type BgpAgent struct {
	LogManager   *logmanager.Plugin
	//Orchestrator *orchestrator.Plugin
	//KVScheduler  *kvscheduler.Scheduler
	Rest         *rest.Plugin
}




// New creates new OsseusAgent instance.
func New() *BgpAgent {

	return &BgpAgent{

	}

}