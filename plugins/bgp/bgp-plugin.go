package bgp

//go:generate protoc --proto_path=model --proto_path=$GOPATH/src --gogo_out=model model/bgp.proto
//go:generate descriptor-adapter --descriptor-name GlobalConf --value-type *model.GlobalConf --import "model" --output-dir "descriptor"
//go:generate descriptor-adapter --descriptor-name PeerConf --value-type *model.PeerConf --import "model" --output-dir "descriptor"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/contiv/bgp-vpp/plugins/bgp/descriptor"
	"github.com/contiv/vpp/plugins/ipnet/restapi"
	"github.com/contiv/vpp/plugins/ksr/model/node"
	"github.com/contiv/vpp/plugins/netctl/remote"
	"github.com/contiv/vpp/plugins/nodesync/vppnode"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/db/keyval"
	"github.com/ligato/cn-infra/infra"
	"github.com/ligato/cn-infra/rpc/rest"
	kvs "github.com/ligato/vpp-agent/plugins/kvscheduler/api"
	bgp_api "github.com/osrg/gobgp/api"
	gobgp "github.com/osrg/gobgp/pkg/server"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type BgpPlugin struct {
	Deps
	watchCloser chan string
	nlriMap     map[uint32]*any.Any
	nextHopMap map[uint32]string
	hasConfigChan chan bool
	hasConfig bool
	etcdInUse bool


}

//Deps is only for external dependencies
type Deps struct {
	//plugins - initialized in options.go NewPlugin()
	infra.PluginDeps
	Rest        *rest.Plugin
	BGPServer   *gobgp.BgpServer
	KVScheduler kvs.KVScheduler
	KVStore     keyval.KvProtoPlugin
}


const nodePrefix = "/vnf-agent/contiv-ksr/allocatedIDs/"
const getIpamDataCmd = "contiv/v1/ipam"
const nodePrefixTwo = "/vnf-agent/contiv-ksr/k8s/node/"
func (p *BgpPlugin) String() string {
	return "Starting BgpPlugin Application"
}
func (p *BgpPlugin) Init() error {
	if p.BGPServer == nil {
		p.BGPServer = gobgp.NewBgpServer()
		go p.BGPServer.Serve()
	}

	p.hasConfig = false
	p.etcdInUse = false
	p.nlriMap = make(map[uint32]*any.Any)
	p.nextHopMap = make(map[uint32]string)
	p.hasConfigChan = make(chan bool)

	gd := descriptor.NewGlobalConfDescriptor(p.Log, p.BGPServer, p.hasConfigChan)
	p.KVScheduler.RegisterKVDescriptor(gd)
	// register descriptor for bgp peer config

	pd := descriptor.NewPeerConfDescriptor(p.Log, p.BGPServer)
	p.KVScheduler.RegisterKVDescriptor(pd)
	p.watchCloser = make(chan string)

	watcher := p.Deps.KVStore.NewWatcher(nodePrefix)
	err := watcher.Watch(p.onChange, p.watchCloser, "")
	if err != nil {
		p.Log.Errorf("Failed to start the node watcher, error %s", err)
		return err
	}

	go p.addExisting()

	p.Log.Info("BGP Plugin initialized")
	return nil
}
func (p *BgpPlugin) Close() error {
	p.Log.Info("Closing Bgp Plugin")
	return nil
}

//tutorial says keyval.protowatchresp but i couldnt find it
//in the  cn-infra/db/keyval/proto_watcher_api.go in our vendor folder there isnt a ProtoWatchResp
// but on the github for cn infra, there is one
func (p *BgpPlugin) onChange(resp datasync.ProtoWatchResp) {
	p.Log.Infof("onChange called, resp: %+v", resp)
	//key := resp.GetKey()
	if !p.hasConfig {
		return
	}
	//Getting ip
	value := &vppnode.VppNode{}
	changeType := resp.GetChangeType()
	if changeType == datasync.Delete {
		if prevValExist, err := resp.GetPrevValue(value); err != nil {
			p.Log.Errorf("could not get prev value in Delete op, error: %s", err)
			return
		} else if !prevValExist {
			p.Log.Errorf("prev value does not exist in Delete op")
			return
		}
	} else {
		if err := resp.GetValue(value); err != nil {
			p.Log.Errorf("get value error: %v", err)
			return
		}
	}
	if len(value.IpAddresses) == 0 {
		p.Log.Warnf("no IP address available for node %v", value.Id)
		return
	}

	//ip := value.IpAddresses[0]
	id := value.Id
	//ipParts := strings.Split(ip, "/")
	//ip = ipParts[0]


	switch changeType {
	case datasync.Put:
		for {
			if !p.etcdInUse {
				break
			}
		}
		p.etcdInUse = true
		p.add(id, value.Name)
		p.etcdInUse = false
	case datasync.Delete:
		for {
			if !p.etcdInUse {
				break
			}
		}
		p.etcdInUse= true
		p.delete(id)
		p.etcdInUse = false
	default:
		p.Log.Errorf("GetChangeType: %v", changeType)
	}
	ipam := restapi.NodeIPAMInfo{}
	err = json.Unmarshal(b, &ipam)
	if err != nil {
		p.Log.Errorf("failed to unmarshal IpamEntry, error: %s, buffer: %+v", err, b)
		return
	}
	//Setting Route info
	ip:= ipam.NodeIP
	podSubnetParts := strings.Split(ipam.PodSubnetThisNode, "/")
	prefixLen, err := strconv.ParseUint(podSubnetParts[1], 10, 32)
	if err != nil {
		p.Log.Errorf("failed to convert pod subnet mask %s on node %d to uint, error %s",
			ipam.PodSubnetThisNode, id, err)
		return
	}
	p.Log.Infof("PREFIX: %s, ", podSubnetParts[0])
	p.Log.Infof("PREFIXLEN: %d, ",uint32(prefixLen) )
	nlri, _ := ptypes.MarshalAny(&bgp_api.IPAddressPrefix{
		Prefix:    podSubnetParts[0],
		PrefixLen: uint32(prefixLen),
	})
	a1, _ := ptypes.MarshalAny(&bgp_api.OriginAttribute{
		Origin: 0,
	})
	a2, _ := ptypes.MarshalAny(&bgp_api.NextHopAttribute{
		NextHop: ip,
	})
	attrs := []*any.Any{a1, a2}
	p.Log.Infof("Put operation with NLRI: %v and Next Hop: %v", nlri, ip)
	_, err = p.Deps.BGPServer.AddPath(context.Background(), &bgp_api.AddPathRequest{
		Path: &bgp_api.Path{
			Family: &bgp_api.Family{Afi: bgp_api.Family_AFI_IP, Safi: bgp_api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attrs,
		},
	})
	if err != nil {
		p.Log.Errorf("AddPath: %v", err)
	}
	p.nlriMap[id] = nlri
	p.nextHopMap[id]=ip
}

func(p *BgpPlugin)delete(id uint32) {
	nlri := p.nlriMap[id]
	ip := p.nextHopMap[id]
	if nlri == nil || ip == "" {
		p.Log.Error("Node with id %d has not yet been added",id)
		return
	}
	a1, _ := ptypes.MarshalAny(&bgp_api.OriginAttribute{
		Origin: 0,
	})
	a2, _ := ptypes.MarshalAny(&bgp_api.NextHopAttribute{
		NextHop: ip,
	})
	attrs := []*any.Any{a1, a2}
	p.Log.Infof("Deleting Path with NLRI: %v", nlri)
	err := p.Deps.BGPServer.DeletePath(context.Background(), &bgp_api.DeletePathRequest{
		Path: &bgp_api.Path{
			Family: &bgp_api.Family{Afi: bgp_api.Family_AFI_IP, Safi: bgp_api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attrs,
		},
	})
	if err != nil {
		p.Log.Errorf("AddPath: %v", err)
	}
	delete(p.nlriMap, id)
	delete(p.nextHopMap, id)
}

func(p *BgpPlugin)add(id uint32, nodeName string) {
	broker := p.KVStore.NewBroker(nodePrefixTwo )
	ni := &node.Node{}
	broker.GetValue(nodeName,ni )
	var mgmtAddr string
	p.Log.Infof("%v\n",ni)
	for _, address := range ni.Addresses {
		if address.Type == node.NodeAddress_NodeInternalIP ||
			address.Type == node.NodeAddress_NodeExternalIP{
			mgmtAddr = address.Address
			break
		}
	}
	if mgmtAddr == "" {
		p.Log.Errorf("Could not find internal ip of %v in %v",nodeName, nodePrefix+nodeName )
		return
	}
	client, err := remote.CreateHTTPClient("../../cmd/bgp/http.client.conf")
	if err != nil {
		p.Log.Error("create http client error: %v", err)
		return
	}

	//Getting Ipam Info
	b, err := getNodeInfo(client, mgmtAddr, getIpamDataCmd)
	if err != nil {
		p.Log.Errorf("getnodeinfo error: %v", err)
		return
	}
	ipam := restapi.NodeIPAMInfo{}
	err = json.Unmarshal(b, &ipam)
	if err != nil {
		p.Log.Errorf("failed to unmarshal IpamEntry, error: %s, buffer: %+v", err, b)
		return
	}
	//Setting Route info
	ip:= ipam.NodeIP
	podSubnetParts := strings.Split(ipam.PodSubnetThisNode, "/")
	prefixLen, err := strconv.ParseUint(podSubnetParts[1], 10, 32)
	if err != nil {
		p.Log.Errorf("failed to convert pod subnet mask %s on node %d to uint, error %s",
			ipam.PodSubnetThisNode, id, err)
		return
	}
	p.Log.Infof("PREFIX: %s, ", podSubnetParts[0])
	p.Log.Infof("PREFIXLEN: %d, ",uint32(prefixLen) )
	nlri, _ := ptypes.MarshalAny(&bgp_api.IPAddressPrefix{
		Prefix:    podSubnetParts[0],
		PrefixLen: uint32(prefixLen),
	})
	a1, _ := ptypes.MarshalAny(&bgp_api.OriginAttribute{
		Origin: 0,
	})
	a2, _ := ptypes.MarshalAny(&bgp_api.NextHopAttribute{
		NextHop: ip,
	})
	attrs := []*any.Any{a1, a2}
	p.Log.Infof("Put operation with NLRI: %v and Next Hop: %v", nlri, ip)
	_, err = p.Deps.BGPServer.AddPath(context.Background(), &bgp_api.AddPathRequest{
		Path: &bgp_api.Path{
			Family: &bgp_api.Family{Afi: bgp_api.Family_AFI_IP, Safi: bgp_api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attrs,
		},
	})
	if err != nil {
		p.Log.Errorf("AddPath: %v", err)
	}
	p.nlriMap[id] = nlri
	p.nextHopMap[id]=ip
}

func(p *BgpPlugin)delete(id uint32) {
	nlri := p.nlriMap[id]
	ip := p.nextHopMap[id]
	if nlri == nil || ip == "" {
		p.Log.Error("Node with id %d has not yet been added",id)
		return
	}
	a1, _ := ptypes.MarshalAny(&bgp_api.OriginAttribute{
		Origin: 0,
	})
	a2, _ := ptypes.MarshalAny(&bgp_api.NextHopAttribute{
		NextHop: ip,
	})
	attrs := []*any.Any{a1, a2}
	p.Log.Infof("Deleting Path with NLRI: %v", nlri)
	err := p.Deps.BGPServer.DeletePath(context.Background(), &bgp_api.DeletePathRequest{
		Path: &bgp_api.Path{
			Family: &bgp_api.Family{Afi: bgp_api.Family_AFI_IP, Safi: bgp_api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attrs,
		},
	})
	if err != nil {
		p.Log.Errorf("AddPath: %v", err)
	}
	delete(p.nlriMap, id)
	delete(p.nextHopMap, id)
}

// getNodeInfo will make an http request for the given command and return an indented slice of bytes.
//copied from vpp/plugins/netctl/cmdimpl/nodes.go
func getNodeInfo(client *remote.HTTPClient, base string, cmd string) ([]byte, error) {
	res, err := client.Get(base, cmd)
	if err != nil {
		err := fmt.Errorf("getNodeInfo: url: %s Get Error: %s", cmd, err.Error())
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		err := fmt.Errorf("getNodeInfo: url: %s HTTP res.Status: %s", cmd, res.Status)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func (p *BgpPlugin) addExisting() {
	p.hasConfig = <-p.hasConfigChan
	for {
		if !p.etcdInUse {
			break
		}
	}
	p.etcdInUse = true
	kv:= p.KVStore.NewBroker("")

	itr, err := kv.ListValues(nodePrefix)
	if err != nil {
		fmt.Println("Failed to discover nodes in Contiv cluster")
		os.Exit(-1)
	}
	for {
		next, stop := itr.GetNext()
		if stop {
			break
		}

		vn := &vppnode.VppNode{}
		next.GetValue(vn)
		p.add(vn.Id, vn.Name)
	}
	p.etcdInUse = false
}