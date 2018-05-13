package listener

import (
	"github.com/sczyh30/waffle-mesh/proxy/network"
	"github.com/sczyh30/waffle-mesh/proxy/network/config"
)

func NewPrivateListener(config config.ServerConfig) network.Listener {
	listener := network.NewListener(network.HTTP2, config)

	return listener
}
