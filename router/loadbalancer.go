package router

import (
	"github.com/n-creativesystem/go-fwncs"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal/plugins"
	"github.com/oidc-proxy-ecosystem/proxy-server/transport"
)

func setOidcLoadbalancer(router *fwncs.Router, loadBalancer config.LoadBalancer) {
	// for _, loadBalancer := range loadBalancers {
	// upstream配列を作成
	proxies := make(map[string]*fwncs.WeightProxyTarget, len(loadBalancer.Upstream))
	for _, balancer := range loadBalancer.Upstream {
		proxies[balancer.Name] = &fwncs.WeightProxyTarget{
			ProxyTarget: &fwncs.ProxyTarget{
				Name: balancer.Name,
				URL:  balancer.GetURL(),
				Meta: balancer.Meta,
			},
			Weight: balancer.Weight,
		}
	}
	var responsePlugin []*plugins.PluginClient
	var transportPlugin []*plugins.PluginClient
	for _, location := range loadBalancer.Locations {

		for _, name := range location.Plugins.Transport {
			if client, ok := plugins.Get(plugins.RequestTransport).Get(name); ok {
				transportPlugin = append(transportPlugin, client)
			}
		}

		for _, name := range location.Plugins.Response {
			if client, ok := plugins.Get(plugins.ResponseModify).Get(name); ok {
				responsePlugin = append(responsePlugin, client)
			}
		}
	}

	router.Use(SessionCheck)
	for _, location := range loadBalancer.Locations {
		balancer := fwncs.NewStaticWeightedRoundRobinBalancer(nil)
		for _, name := range location.Upstream {
			if proxy, ok := proxies[name]; ok {
				balancer.Add(proxy)
			}
		}
		var tokenHandler fwncs.HandlerFunc
		switch location.TokenType {
		case "id_token":
			tokenHandler = SetIdToken
		case "access_token":
			tokenHandler = SetAccessToken
		default:
			tokenHandler = SetIdToken
		}
		log := fwncs.DefaultLogger
		tr := transport.NewPluginTransport(transport.NewTLSSkipVerifyTransport(nil), transportPlugin, log)
		proxy := fwncs.ProxyWithConfig(fwncs.ProxyConfig{
			LoadBalancer:   balancer,
			Transport:      tr,
			Rewrite:        location.Rewrite,
			ModifyResponse: modifyResponsePlugins(responsePlugin),
		})
		router.Any(location.Path, tokenHandler, proxy)
	}
	// }
}
