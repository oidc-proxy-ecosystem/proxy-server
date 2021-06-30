package config

import (
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/oidc-proxy-ecosystem/proxy-server/plugins"
)

type SettingFile struct {
	Directory    string `envconfig:"DIRECTORY" default:""`
	Loadbalancer string `envconfig:"LOAD_BALANCER" default:"loadbalancer.yml"`
	Config       string `envconfig:"CONFIG" default:"config.yml"`
	Oidc         string `envconfig:"OIDC" default:"oidc.yml"`
	Saml         string `envconfig:"SAML" default:"saml.yml"`
	Auth         string `envconfig:"AUTH" default:"auth.yml"`
	Menu         string `envconfig:"MENU" default:"menu.yml"`
	Storage      string `envconfig:"Storage" default:"storage.yml"`
	Certificate  string `envconfig:"CA_FILE_NAME" default:"ca.yml"`
}

func (s SettingFile) ConvertPluginConfig() *plugins.Config {
	return &plugins.Config{
		Directory:    File.Directory,
		Loadbalancer: File.Loadbalancer,
		Config:       File.Config,
		Oidc:         File.Oidc,
		Saml:         File.Saml,
		Auth:         File.Auth,
		Menu:         File.Menu,
	}
}

var File SettingFile

func NewSettingFile() {
	envconfig.MustProcess("SETTING", &File)
	if File.Directory != "" {
		File.Loadbalancer = filepath.Join(File.Directory, File.Loadbalancer)
		File.Config = filepath.Join(File.Directory, File.Config)
		File.Oidc = filepath.Join(File.Directory, File.Oidc)
		File.Saml = filepath.Join(File.Directory, File.Saml)
		File.Auth = filepath.Join(File.Directory, File.Auth)
		File.Menu = filepath.Join(File.Directory, File.Menu)
		File.Storage = filepath.Join(File.Directory, File.Storage)
		File.Certificate = filepath.Join(File.Directory, File.Certificate)
	}
}
