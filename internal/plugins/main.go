package plugins

import (
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type PluginClient struct {
	client    *plugin.Client
	rpcClient plugin.ClientProtocol
}

func (p *PluginClient) GetClient() *plugin.Client {
	return p.client
}

func (p *PluginClient) GetRpcClient() plugin.ClientProtocol {
	return p.rpcClient
}

type MapPlugins struct {
	mp map[string]*PluginClient
}

func (m *MapPlugins) Close() {
	for _, client := range m.mp {
		client.rpcClient.Close()
		client.client.Kill()
	}
}

func (m *MapPlugins) Add(name string, client *plugin.Client, rpcClient plugin.ClientProtocol) {
	if m.mp == nil {
		m.mp = make(map[string]*PluginClient)
	}
	m.mp[name] = &PluginClient{
		client:    client,
		rpcClient: rpcClient,
	}
}

func (m *MapPlugins) Get(name string) (*PluginClient, bool) {
	c, ok := m.mp[name]
	return c, ok
}

type kind string

func (k kind) String() string {
	return string(k)
}

const (
	RequestTransport kind = "request_transport"
	ResponseModify   kind = "response_modify"
)

var LoadPlugins map[kind]*MapPlugins = make(map[kind]*MapPlugins)

func Get(kind kind) *MapPlugins {
	return LoadPlugins[kind]
}

func Run(pluginNames []string, kind kind) error {
	if _, ok := LoadPlugins[kind]; !ok {
		LoadPlugins[kind] = new(MapPlugins)
	}
	mpPlugins := new(MapPlugins)
	for _, pluginName := range pluginNames {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig:  Handshake,
			VersionedPlugins: VersionedPlugins,
			Managed:          true,
			Logger:           hclog.Default(),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			Cmd:              exec.Command(fmt.Sprintf("ncs-%s", pluginName)),
		})
		rpcClient, err := client.Client()
		if err != nil {
			return err
		}
		mpPlugins.Add(pluginName, client, rpcClient)
	}
	LoadPlugins[kind] = mpPlugins
	return nil
}

func Close() {
	for _, mpPlugin := range LoadPlugins {
		mpPlugin.Close()
	}
}
