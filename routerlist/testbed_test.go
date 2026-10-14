package routerlist

import (
	"encoding/json"
	"testing"

	"github.com/11th-ndn-hackathon/ndn-fch/model"
)

func TestTestbedNodeWithoutNDNd(t *testing.T) {
	var node testbedNode
	err := json.Unmarshal([]byte(`{"shortname":"NODE","site":"https://node.example","ip_addresses":["192.0.2.1"],"position":[1,2],"prefix":"ndn:/ndn/test","neighbors":[]}`), &node)
	if err != nil {
		t.Fatalf("decode legacy testbed node: %v", err)
	}

	routers := node.Routers()
	if len(routers) != 1 {
		t.Fatalf("got %d routers, want one legacy router", len(routers))
	}
	if got := routers[0].ConnectString(model.TransportIPFamily{Transport: model.TransportUDP, Family: model.IPv4}); got != "node.example:6363" {
		t.Fatalf("legacy UDP endpoint = %q, want node.example:6363", got)
	}
}

func TestTestbedNodeWithNDNd(t *testing.T) {
	var node testbedNode
	err := json.Unmarshal([]byte(`{"shortname":"NODE","site":"https://node.example","ip_addresses":["192.0.2.1"],"position":[1,2],"prefix":"ndn:/ndn/test","neighbors":[],"ndnd":{"udp_port":6364,"prefix":"ndn:/ndnd/test"}}`), &node)
	if err != nil {
		t.Fatalf("decode NDNd testbed node: %v", err)
	}

	routers := node.Routers()
	if len(routers) != 2 {
		t.Fatalf("got %d routers, want NFD and NDNd", len(routers))
	}
	if got := routers[0].ConnectString(model.TransportIPFamily{Transport: model.TransportUDP, Family: model.IPv4}); got != "node.example:6363" {
		t.Fatalf("NFD UDP endpoint = %q, want node.example:6363", got)
	}
	if got := routers[1].ID(); got != "NODE-ndnd" {
		t.Fatalf("NDNd ID = %q, want NODE-ndnd", got)
	}
	if got := routers[1].ConnectString(model.TransportIPFamily{Transport: model.TransportUDP, Family: model.IPv4}); got != "node.example:6364" {
		t.Fatalf("NDNd UDP endpoint = %q, want node.example:6364", got)
	}
	if got := routers[1].Prefix(); got != "/ndnd/test" {
		t.Fatalf("NDNd prefix = %q, want /ndnd/test", got)
	}
}
