package main

import (
	"github.com/ligato/cn-infra/agent"
)


func main() {
	// Create an instance of our plugin using its constructor.
	p := NewBgpPlugin()

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
	Orchestrator *orchestrator.Plugin
	KVScheduler  *kvscheduler.Scheduler
	Rest         *rest.Plugin
}

// New creates new OsseusAgent instance.
func New() *BgpAgent {

	return &BgpAgent{

	}

}