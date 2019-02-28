package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	api "github.com/osrg/gobgp/api"
	gobgp "github.com/osrg/gobgp/pkg/server"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)
	s := gobgp.NewBgpServer()
	go s.Serve()

	// global configuration
	if err := s.StartBgp(context.Background(), &api.StartBgpRequest{
		Global: &api.Global{
			As:         64512,
			RouterId:   "10.20.0.2",
			ListenPort: -1, // gobgp won't listen on tcp:179
		},
	}); err != nil {
		log.Fatal(err)
	}

	// monitor the change of the peer state
	if err := s.MonitorPeer(context.Background(), &api.MonitorPeerRequest{}, func(p *api.Peer){log.Info(p)}); err != nil {
		log.Fatal(err)
	}

	// neighbor configuration
/*	wone := &api.Peer{
		Conf: &api.PeerConf{
			NeighborAddress: "10.20.0.10",
			PeerAs:          64512,
		},
	}

	wtwo := &api.Peer{
		Conf: &api.PeerConf{
			NeighborAddress: "10.20.0.11",
			PeerAs:          64512,
		},
	}
*/
	g := &api.Peer{
		Conf: &api.PeerConf{
			NeighborAddress: "10.20.0.100",
			PeerAs:          64512,
		},
	}

/*	if err := s.AddPeer(context.Background(), &api.AddPeerRequest{
		Peer: wone,
	}); err != nil {
		log.Fatal(err)
	}

	if err := s.AddPeer(context.Background(), &api.AddPeerRequest{
		Peer: wtwo,
	}); err != nil {
		log.Fatal(err)
	}*/

	if err := s.AddPeer(context.Background(), &api.AddPeerRequest{
		Peer: g,
	}); err != nil {
		log.Fatal(err)
	}

	// add routes
	nlri, _ := ptypes.MarshalAny(&api.IPAddressPrefix{
		Prefix:    "10.1.1.0",
		PrefixLen: 24,
	})

	a1, _ := ptypes.MarshalAny(&api.OriginAttribute{
		Origin: 0,
	})
	a2, _ := ptypes.MarshalAny(&api.NextHopAttribute{
		NextHop: "192.168.16.1",
	})
	attrs := []*any.Any{a1, a2}

	_, err := s.AddPath(context.Background(), &api.AddPathRequest{
		Path: &api.Path{
			Family:    &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attrs,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	nlri, _ = ptypes.MarshalAny(&api.IPAddressPrefix{
		Prefix:    "10.1.2.0",
		PrefixLen: 24,
	})

	a1, _ = ptypes.MarshalAny(&api.OriginAttribute{
		Origin: 0,
	})
	a2, _ = ptypes.MarshalAny(&api.NextHopAttribute{
		NextHop: "192.168.16.2",
	})
	attrs = []*any.Any{a1, a2}

	_, err = s.AddPath(context.Background(), &api.AddPathRequest{
		Path: &api.Path{
			Family: &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attrs,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	nlri, _ = ptypes.MarshalAny(&api.IPAddressPrefix{
		Prefix:    "10.1.3.0",
		PrefixLen: 24,
	})

	a1, _ = ptypes.MarshalAny(&api.OriginAttribute{
		Origin: 0,
	})
	a2, _ = ptypes.MarshalAny(&api.NextHopAttribute{
		NextHop: "192.168.16.3",
	})
	attrs = []*any.Any{a1, a2}

	_, err = s.AddPath(context.Background(), &api.AddPathRequest{
		Path: &api.Path{
			Family: &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attrs,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	//Infinite Loop so that program does not time out
	for {
		fmt.Println("Infinite Loop 1")
		time.Sleep(time.Second)
	}
	// do something useful here instead of exiting
	//time.Sleep(time.Minute * 3)
}