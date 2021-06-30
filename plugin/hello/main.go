package main

import (
	"net/url"

	"github.com/hashicorp/go-plugin"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal/plugins"
	"github.com/oidc-proxy-ecosystem/proxy-server/shared"
)

type Hello struct {
}

var _ shared.Transport = (*Hello)(nil)

func (h *Hello) Transport(header map[string][]string, urlPath string, config map[string]string) shared.Result {
	header["X-Hello"] = []string{"world"}
	u, _ := url.Parse(urlPath)
	q := u.Query()
	q.Add("hello", "world")
	u.RawQuery = q.Encode()
	result := shared.Result{
		Header:       header,
		Status:       0,
		URL:          u.String(),
		ErrorMessage: "",
	}
	return result
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: {
				"request_transport": &internal.TransportPlugin{Impl: &Hello{}},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
