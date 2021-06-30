package plugins

import (
	"github.com/hashicorp/go-plugin"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PROXY_SERVER_PLUGIN",
	MagicCookieValue: "e524367f551015d0f6e1df2d3158dbe9e30945fb",
}

var VersionedPlugins = map[int]plugin.PluginSet{
	1: {
		"request_transport": &internal.TransportPlugin{},
		"response_modify":   &internal.ResponsePlugin{},
	},
}
