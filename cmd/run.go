package cmd

import (
	"github.com/hashicorp/go-plugin"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal/plugins"
	"github.com/oidc-proxy-ecosystem/proxy-server/router"
	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
)

func loadPlugins(p config.Plugins) error {
	if err := plugins.Run(p.Transport, plugins.RequestTransport); err != nil {
		return err
	}
	if err := plugins.Run(p.Response, plugins.ResponseModify); err != nil {
		return err
	}
	return nil
}

func runProxy() error {
	conf := config.NewConfig(config.File.Config)
	loadBalancer := config.NewLoadBalancer(config.File.Loadbalancer)
	for _, location := range loadBalancer.Locations {
		if err := loadPlugins(location.Plugins); err != nil {
			return err
		}
	}
	defer plugins.Close()
	defer plugin.CleanupClients()
	authConfig := config.NewAuthConfig(config.File.Auth)
	config.NewMenu(config.File.Menu)
	var filename string
	// switch conf.AuthType {
	// case "oidc":
	filename = config.File.Oidc
	// case "saml":
	// 	filename = config.File.Saml
	// }
	router := router.New(conf, filename, authConfig, loadBalancer)
	var err error
	if conf.CertFile != "" && conf.KeyFile != "" && utils.IsExists(conf.CertFile) && utils.IsExists(conf.KeyFile) {
		err = router.RunTLS(conf.Port, conf.CertFile, conf.KeyFile)
	} else {
		err = router.Run(conf.Port)
	}
	return err
}
