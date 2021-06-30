package config

import (
	"crypto/tls"
	"os"

	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
)

type StorageConfig struct {
	Type      string      `yaml:"type"`
	Endpoints []string    `yaml:"endpoints"`
	Username  string      `yaml:"username"`
	Password  string      `yaml:"password"`
	TLSConfig *tls.Config `yaml:"-"`
}

func NewStorageConfig(filename string) StorageConfig {
	var storageConfig StorageConfig
	if _, err := os.Stat(filename); err == nil {
		utils.MustReadYaml(filename, &storageConfig)
	} else {
		storageConfig = StorageConfig{
			Endpoints: []string{"localhost:6379"},
			Username:  "",
			Password:  "",
			TLSConfig: nil,
		}
	}
	return storageConfig
}
