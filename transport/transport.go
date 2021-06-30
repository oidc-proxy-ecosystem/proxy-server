package transport

import (
	"fmt"
	"net/http"

	"github.com/n-creativesystem/go-fwncs"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal/plugins"
)

type pluginTransport struct {
	tr      http.RoundTripper
	clients []*plugins.PluginClient
	logger  fwncs.ILogger
}

func copyHeader(header http.Header) map[string][]string {
	resultHeader := make(map[string][]string, len(header))
	for key, values := range header {
		resultHeader[key] = make([]string, len(values))
		copy(resultHeader[key], values)
	}
	return resultHeader
}

func (t *pluginTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, client := range t.clients {
		rpcClient := client.GetRpcClient()
		raw, err := rpcClient.Dispense("request_transport")
		if err != nil {
			t.logger.Error(err)
			continue
		}
		protoVer := client.GetClient().NegotiatedVersion()
		switch protoVer {
		case 1:
			p := raw.(*internal.TransportGRPCClient)
			header, URL, err := p.Transport(req.Header, req.URL)
			if err != nil {
				t.logger.Warning(err)
			} else {
				req.Header = copyHeader(header)
				*req.URL = *URL
			}
		default:
			t.logger.Error(fmt.Sprintf("no support version: %d", protoVer))
		}
	}
	return t.tr.RoundTrip(req)
}

func NewPluginTransport(tr http.RoundTripper, clients []*plugins.PluginClient, log fwncs.ILogger) http.RoundTripper {
	if tr == nil {
		tr = http.DefaultTransport
	}
	return &pluginTransport{
		tr:      tr,
		clients: clients,
		logger:  log,
	}
}
