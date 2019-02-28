package bgp

import (
	"github.com/ligato/cn-infra/infra"
	"github.com/ligato/cn-infra/rpc/rest"
	"log"
)

type BgpPlugin struct {
	Deps
	//figure out how to plug the gobgp library into this
}

//Deps is only for external dependencies
type Deps struct {
	infra.PluginDeps
	Rest *rest.Plugin
	//KVSchedulser kvs.KVScheduler
}

func (p *BgpPlugin) String() string {
	return "HelloWorld"
}

func (p *BgpPlugin) Init() error {
	log.Println("Hello World!")
	return nil
}

func (p *BgpPlugin) Close() error {
	log.Println("Goodbye World!")
	return nil
}
