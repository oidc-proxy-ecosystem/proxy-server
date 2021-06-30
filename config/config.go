package config

import (
	"os"

	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
)

type Plugins struct {
	Transport []string `yaml:"request_transport"`
	Response  []string `yaml:"response_modify"`
}

type Config struct {
	RedisEndpoints []string `yaml:"redis_endpoints"`
	RedisUsername  string   `yaml:"redis_username"`
	RedisPassword  string   `yaml:"redis_password"`
	CertFile       string   `yaml:"cert_file"`
	KeyFile        string   `yaml:"key_file"`
	Port           int      `yaml:"port"`
	LogLevel       string   `yaml:"log_level"`
	AuthType       string   `yaml:"-"`
	Plugins        Plugins  `yaml:"plugins"`
}

func NewConfig(filename string) Config {
	config := Config{}
	utils.MustReadYamlExpand(filename, &config)
	utils.Assert(len(config.RedisEndpoints) == 0, "redis endpoint is required")
	utils.Assert(config.Port == 0, "port is required")
	if config.AuthType == "" {
		config.AuthType = "oidc"
	}
	utils.Asserts("auth_type is oidc or saml only", config.AuthType == "oidc", config.AuthType == "saml")
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	os.Setenv("LOG_LEVEL", config.LogLevel)
	return config
}
