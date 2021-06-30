package config

import "github.com/oidc-proxy-ecosystem/proxy-server/utils"

var Menus []*Menu

type Menu struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
	Current     bool   `yaml:"-"`
	Thumbnail   string `yaml:"thumbnail"`
}

func NewMenu(filename string) {
	utils.MustReadYaml(filename, &Menus)
	for _, menu := range Menus {
		if menu.Thumbnail == "" {
			menu.Thumbnail = "/images/no_image.png"
		}
	}
}
