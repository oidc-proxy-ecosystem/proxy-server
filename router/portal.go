package router

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/sessions"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
)

func oidcPortal(c fwncs.Context) {
	mp, err, redirect := getUserInfo(c)
	if redirect {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		return
	}
	if err != nil {
		c.Logger().Error(err)
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
		return
	}
	picture, ok := mp["picture"].(string)
	if !ok {
		picture = "/images/user.png"
	}
	t := template.Must(template.ParseFiles("public/templates/layout.html", "public/templates/index.html"))
	if err = t.Execute(c.Writer(), struct {
		Title   string
		Url     string
		Picture string
		Menus   []*config.Menu
	}{
		Title:   "Portal",
		Picture: picture,
		Menus:   config.Menus,
	}); err != nil {
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
	}
}

func samlPortal(c fwncs.Context) {
	session := sessions.Default(c)
	// mp, err, redirect := getUserInfo(c)
	// if redirect {
	// 	c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
	// 	return
	// }
	// if err != nil {
	// 	c.Logger().Error(err)
	// 	c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
	// 	return
	// }
	var picture string = "/images/user.png"
	// picture, ok := mp["picture"].(string)
	// if !ok {
	// 	picture = "/images/user.png"
	// }
	value := session.Get("saml_user")
	c.Logger().Info(fmt.Sprintf("%#v", value))
	t := template.Must(template.ParseFiles("public/templates/layout.html", "public/templates/index.html"))
	if err := t.Execute(c.Writer(), struct {
		Title   string
		Url     string
		Picture string
		Menus   []*config.Menu
		Value   interface{}
	}{
		Title:   "Portal",
		Picture: picture,
		Menus:   config.Menus,
		Value:   value,
	}); err != nil {
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
	}
}
