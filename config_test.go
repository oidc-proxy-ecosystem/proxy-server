package main_test

import (
	"testing"

	"github.com/n-creativesystem/go-fwncs/tests"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	conf := config.NewConfig("tests/config.sample.yml")
	assert.Equal(t, []string{"localhost:6379"}, conf.RedisEndpoints)
	assert.Equal(t, "username", conf.RedisUsername)
	assert.Equal(t, "password", conf.RedisPassword)
	assert.Equal(t, "ssl/server.crt", conf.CertFile)
	assert.Equal(t, "ssl/server.key", conf.KeyFile)
	assert.Equal(t, 8443, conf.Port)
	assert.Equal(t, "info", conf.LogLevel)
}

func TestLoadBalancerConfig(t *testing.T) {
	loadbalancer := config.NewLoadBalancer("tests/loadbalancer.sample.yml")
	location := loadbalancer.Locations[0]
	assert.Equal(t, "/a/*all", location.Path)
	assert.Equal(t, "id_token", location.TokenType)
	assert.Equal(t, "/v1/$1", location.Rewrite["^/a/*"])
	upstream := loadbalancer.Upstream[0]
	assert.Equal(t, "backend1", upstream.Name)
	assert.Equal(t, "https://locahost:8080", upstream.GetURL().String())
	assert.Equal(t, 0.5, upstream.Weight)
	assert.Equal(t, true, location.Upstream[0] == upstream.Name)
}

func TestOidcConfig(t *testing.T) {
	scopes := map[string]bool{
		"email":          true,
		"openid":         true,
		"offline_access": true,
		"profile":        true,
	}
	oidcConfig := config.NewOidcConfig("tests/oidc.sample.yml")
	for _, scope := range oidcConfig.Scopes {
		assert.Condition(t, func() (success bool) {
			return scopes[scope]
		})
	}
	assert.Equal(t, "https://idp.com", oidcConfig.Provider)
	assert.Equal(t, "client_id", oidcConfig.ClientId)
	assert.Equal(t, "client_secret", oidcConfig.ClientSecret)
	assert.Equal(t, "https://localhost/auth/callback", oidcConfig.CallbackUrl)
	assert.Equal(t, "https://idp.com/v/logout?https://localhost/auth/login", oidcConfig.Logout)
	assert.Equal(t, "audiences", oidcConfig.Audiences[0])
}

func TestAuthConfig(t *testing.T) {
	tt := tests.TestFrames{
		{
			Name: "Auth Config",
			Fn: func(t *testing.T) {
				authConfig := config.NewAuthConfig("tests/auth.sample.yml")
				assert.Equal(t, "/path", authConfig.Path)
				assert.Equal(t, "/login", authConfig.Login)
				assert.Equal(t, "/callback", authConfig.Callback)
				assert.Equal(t, "/logout", authConfig.Logout)
				assert.Equal(t, "/userinfo", authConfig.UserInfo)
			},
		},
		{
			Name: "Auth Config No slash prefix",
			Fn: func(t *testing.T) {
				authConfig := config.NewAuthConfig("tests/auth.sample2.yml")
				assert.Equal(t, "/path", authConfig.Path)
				assert.Equal(t, "/login", authConfig.Login)
				assert.Equal(t, "/callback", authConfig.Callback)
				assert.Equal(t, "/logout", authConfig.Logout)
				assert.Equal(t, "/userinfo", authConfig.UserInfo)
			},
		},
		{
			Name: "Auth Config default value",
			Fn: func(t *testing.T) {
				authConfig := config.NewAuthConfig("tests/auth.sample3.yml")
				assert.Equal(t, "/auth", authConfig.Path)
				assert.Equal(t, "/login", authConfig.Login)
				assert.Equal(t, "/callback", authConfig.Callback)
				assert.Equal(t, "/logout", authConfig.Logout)
				assert.Equal(t, "/userinfo", authConfig.UserInfo)
			},
		},
	}
	tt.Run(t)
}
