package config

import (
	"net/url"

	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
)

func mustUrl(rawUrl string) *url.URL {
	u, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}
	return u
}

type Upstream struct {
	Name   string                 `yaml:"name"`
	URL    string                 `yaml:"url"`
	Weight float64                `yaml:"weight"`
	Meta   map[string]interface{} `yaml:"meta"`
}

func (b Upstream) GetURL() *url.URL {
	return mustUrl(b.URL)
}

type Location struct {
	Path      string            `yaml:"path"`
	TokenType string            `yaml:"token_type"`
	Upstream  []string          `yaml:"upstream"`
	Rewrite   map[string]string `yaml:"rewrite"`
	Plugins   Plugins           `yaml:"plugins"`
}

type LoadBalancer struct {
	// Domain     string     `yaml:"domain"`
	Locations  []Location `yaml:"locations"`
	Upstream   []Upstream `yaml:"upstream"`
	Portal     bool       `yaml:"portal"`
	DefaultURL string     `yaml:"default_url"`
}

// type LoadBalancers []LoadBalancer

func NewLoadBalancer(filename string) LoadBalancer {
	var loadBalancer LoadBalancer
	utils.MustReadYaml(filename, &loadBalancer)
	if loadBalancer.DefaultURL == "" {
		if loadBalancer.Portal {
			loadBalancer.DefaultURL = "/portal"
		} else {
			loadBalancer.DefaultURL = "/"
		}
	}
	return loadBalancer
}
