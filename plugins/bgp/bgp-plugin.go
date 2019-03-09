package bgp

//go:generate protoc --proto_path=model --proto_path=$GOPATH/src --gogo_out=model model/bgp.proto

//go:generate descriptor-adapter --descriptor-name GlobalConf --value-type *model.GlobalConf --import "model" --output-dir "descriptor"
//go:generate descriptor-adapter --descriptor-name PeerConf --value-type *model.PeerConf --import "model" --output-dir "descriptor"

import (
	"log"

	"github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/infra"
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/ligato/vpp-agent/plugins/kvscheduler"
	"github.com/ligato/vpp-agent/plugins/orchestrator"
	gobgp "github.com/osrg/gobgp/pkg/server"
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
	BGPServer    *gobgp.BgpServer
	//interface needed to write to ETCD - initialized in Init()
	//Watcher   datasync.KeyValProtoWatcher
	//Publisher datasync.KeyProtoValWriter
}

func (p *BgpPlugin) String() string {
	return "HelloWorld"
}

func (p *BgpPlugin) Init() error {
	if p.Deps.BGPServer == nil {
		p.Deps.BGPServer = gobgp.NewBgpServer()
	}
	log.Println("Hello World!")
	return nil
}

func (p *BgpPlugin) Close() error {
	log.Println("Goodbye World!")
	return nil
}
