package bgp

//go:generate protoc --proto_path=model --proto_path=$GOPATH/src --gogo_out=model model/bgp.proto

import (
	"log"

	"github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/infra"
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/ligato/vpp-agent/plugins/kvscheduler"
	"github.com/ligato/vpp-agent/plugins/orchestrator"
)

type BgpPlugin struct {
	Deps
	//figure out how to plug the gobgp library into this
}

//Deps is only for external dependencies
type Deps struct {
	//plugins - initialized in options.go NewPlugin()
	infra.PluginDeps
	Rest         *rest.Plugin
	Orchestrator *orchestrator.Plugin
	Scheduler    *kvscheduler.Scheduler
	ETCDDataSync *kvdbsync.Plugin

	//interface needed to write to ETCD - initialized in Init()
	//Watcher   datasync.KeyValProtoWatcher
	//Publisher datasync.KeyProtoValWriter
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
