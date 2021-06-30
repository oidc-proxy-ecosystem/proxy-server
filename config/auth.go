package config

import (
	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
)

type Auth struct {
	Path     string `yaml:"path"`
	Login    string `yaml:"login"`
	Callback string `yaml:"callback"`
	Logout   string `yaml:"logout"`
	UserInfo string `yaml:"userinfo"`
}

func NewAuthConfig(filename string) Auth {
	var auth Auth
	utils.MustReadYaml(filename, &auth)
	setDefault(&auth.Path, "/auth")
	setDefault(&auth.Login, "/login")
	setDefault(&auth.Callback, "/callback")
	setDefault(&auth.Logout, "/logout")
	setDefault(&auth.UserInfo, "/userinfo")

	checkPrefix(&auth.Path)
	checkPrefix(&auth.Login)
	checkPrefix(&auth.Callback)
	checkPrefix(&auth.Logout)
	checkPrefix(&auth.UserInfo)
	return auth
}

func checkPrefix(value *string) {
	v := *value
	if v[0] != '/' {
		*value = "/" + *value
	}
}

func setDefault(value *string, defaultValue string) {
	if *value == "" {
		*value = defaultValue
	}
}
