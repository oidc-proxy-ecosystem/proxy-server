package config

import (
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
	"golang.org/x/oauth2"
)

type OidcConfig struct {
	Scopes       []string `yaml:"scopes"`
	Provider     string   `yaml:"provider"`
	ClientId     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	CallbackUrl  string   `yaml:"callback_url"`
	Logout       string   `yaml:"logout"`
	Audiences    []string `yaml:"audiences"`
}

// Audience
type Audience string

func (a Audience) String() string {
	return string(a)
}

type Audiences []string

func (a Audiences) SetValue() oauth2.AuthCodeOption {
	return oauth2.SetAuthURLParam("audience", a.String())
}

func (a Audiences) String() string {
	s := make([]string, len(a))
	for idx, audience := range a {
		s[idx] = audience
	}
	return strings.Join(s, " ")
}

type Authenticator struct {
	Provider   *oidc.Provider
	Config     oauth2.Config
	OidcConfig OidcConfig
}

func (auth *Authenticator) SetValues() []oauth2.AuthCodeOption {
	audiences := Audiences(auth.OidcConfig.Audiences)
	return []oauth2.AuthCodeOption{audiences.SetValue()}
}

func NewOidcConfig(filename string) OidcConfig {
	var oidcConfig OidcConfig
	utils.MustReadYamlExpand(filename, &oidcConfig)
	if len(oidcConfig.Scopes) == 0 {
		oidcConfig.Scopes = []string{"email", "openid", "offline_access", "profile"}
	}
	utils.Assert(oidcConfig.Provider == "", "provider is required")
	utils.Assert(oidcConfig.ClientId == "", "client_id is required")
	utils.Assert(oidcConfig.ClientSecret == "", "client_secret is required")
	utils.Assert(oidcConfig.CallbackUrl == "", "callback_url is required")
	return oidcConfig
}
