package router

import (
	"net/http"
	"net/url"

	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/sessions"
)

func oidcLogout(logoutUrl string) fwncs.HandlerFunc {
	u, err := url.Parse(logoutUrl)
	if err != nil {
		panic(err)
	}
	return func(c fwncs.Context) {
		session := sessions.Default(c)
		if err != nil {
			c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
			return
		}
		session.Options(sessions.Options{
			MaxAge: -1,
		})
		session.Clear()
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, u.String())
	}
}

// func samlLogout(c fwncs.Context) {
// 	sp := c.Get(spKey).(*saml2.SAMLServiceProvider)
// 	doc, err := sp.BuildLogoutResponseDocument("", "")
// 	sp.BuildLogoutURLRedirect("", doc)
// }
