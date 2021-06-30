package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-plugin"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/oidc-proxy-ecosystem/proxy-server/plugins"
	"github.com/oidc-proxy-ecosystem/proxy-server/shared"
	"google.golang.org/grpc"
)

type TransportPlugin struct {
	plugin.Plugin
	Impl shared.Transport
}

var _ plugin.GRPCPlugin = (*TransportPlugin)(nil)

func (p *TransportPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	plugins.RegisterTransportServer(s, &TransportGRPCServer{Impl: p.Impl})
	return nil
}

func (p *TransportPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &TransportGRPCClient{client: plugins.NewTransportClient(c)}, nil
}

type TransportGRPCServer struct {
	Impl shared.Transport
	plugins.UnimplementedTransportServer
}

var _ plugins.TransportServer = (*TransportGRPCServer)(nil)

func (m *TransportGRPCServer) Transport(ctx context.Context, inf *plugins.Interface) (*plugins.Reply, error) {
	header := make(map[string][]string, len(inf.Header))
	for key, value := range inf.Header {
		header[key] = make([]string, len(value.Value))
		copy(header[key], value.Value)
	}
	confMap := map[string]string{}
	confMap["directory"] = inf.Config.Directory
	confMap["config"] = inf.Config.Config
	confMap["auth"] = inf.Config.Auth
	confMap["loadbalancer"] = inf.Config.Loadbalancer
	confMap["menu"] = inf.Config.Menu
	confMap["oidc"] = inf.Config.Oidc
	confMap["saml"] = inf.Config.Saml

	r := m.Impl.Transport(header, inf.URL, confMap)
	replyHeader := make(map[string]*plugins.Values, len(r.Header))
	for key, values := range r.Header {
		v := &plugins.Values{
			Value: make([]string, len(values)),
		}
		copy(v.Value, values)
		replyHeader[key] = v
	}
	reply := &plugins.Reply{
		URL:          r.URL,
		Header:       replyHeader,
		Status:       r.Status,
		ErrorMessage: r.ErrorMessage,
	}
	return reply, nil
}

type TransportGRPCClient struct {
	client plugins.TransportClient
}

func (m *TransportGRPCClient) Transport(header http.Header, u *url.URL) (rHeader http.Header, rUrl *url.URL, rErr error) {
	defer func() {
		if err := recover(); err != nil {
			rErr = fmt.Errorf("panic: %v", err)
		}
		rUrl = new(url.URL)
		rHeader = header.Clone()
		*rUrl = *u
	}()
	rHeader = header.Clone()
	rUrl = new(url.URL)
	*rUrl = *u
	rErr = nil
	infHeader := make(map[string]*plugins.Values, len(header))
	for key, values := range header {
		v := &plugins.Values{
			Value: make([]string, len(values)),
		}
		copy(v.Value, values)
		infHeader[key] = v
	}
	inf := &plugins.Interface{
		URL:    u.String(),
		Header: infHeader,
		Config: config.File.ConvertPluginConfig(),
	}
	r, err := m.client.Transport(context.Background(), inf)
	if err != nil {
		rErr = err
		return
	}
	rHeader = make(http.Header, len(r.Header))
	for key, value := range r.Header {
		rHeader[key] = make([]string, len(value.Value))
		copy(rHeader[key], value.Value)
	}
	return
}
