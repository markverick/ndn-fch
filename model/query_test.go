package model_test

import (
	"testing"

	"github.com/11th-ndn-hackathon/ndn-fch/model"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	assert := assert.New(t)

	list := model.ParseQueries("k=3&cap=udp&k=2&cap=http3&ipv4=1&ipv6=1&lon=121.4737&lat=31.2304")
	assert.Len(list, 2)
	for i, q := range list {
		if i == 0 {
			assert.Equal(3, q.Count)
			assert.Equal(model.TransportUDP, q.Transport)
		} else {
			assert.Equal(2, q.Count)
			assert.Equal(model.TransportH3, q.Transport)
		}
		assert.True(q.IPv4)
		assert.True(q.IPv6)
		assert.InDelta(121.4737, q.Position[0], 0.0001)
		assert.InDelta(31.2304, q.Position[1], 0.0001)
	}
}

type queryTestRouter struct {
	id     string
	prefix string
}

func (r queryTestRouter) ID() string                                 { return r.id }
func (queryTestRouter) Position() model.LonLat                       { return model.LonLat{} }
func (r queryTestRouter) Prefix() string                             { return r.prefix }
func (queryTestRouter) ConnectString(model.TransportIPFamily) string { return "router.example:6363" }
func (queryTestRouter) Neighbors() map[string]int                    { return nil }

func TestQueryNetwork(t *testing.T) {
	available := []model.RouterAvail{
		{
			Router: queryTestRouter{id: "nfd", prefix: "/ndn/test"},
			Available: map[model.TransportIPFamily]bool{
				{Transport: model.TransportUDP, Family: model.IPv4}: true,
			},
		},
		{
			Router: queryTestRouter{id: "ndnd", prefix: "/ndnd/test"},
			Available: map[model.TransportIPFamily]bool{
				{Transport: model.TransportUDP, Family: model.IPv4}: true,
			},
		},
	}

	legacy := model.ParseQueries("cap=udp&k=2")[0].Execute(available)
	assert.Equal(t, []string{"nfd"}, routerIDs(legacy))

	nfd := model.ParseQueries("cap=udp&k=2&network=ndn")[0].Execute(available)
	assert.Equal(t, []string{"nfd"}, routerIDs(nfd))

	ndnd := model.ParseQueries("cap=udp&k=2&network=ndnd")[0].Execute(available)
	assert.Equal(t, []string{"ndnd"}, routerIDs(ndnd))
}

func routerIDs(routers []model.RouterAvail) (ids []string) {
	for _, router := range routers {
		ids = append(ids, router.ID())
	}
	return ids
}
