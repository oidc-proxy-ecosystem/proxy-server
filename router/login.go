package router

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/sessions"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/oidc-proxy-ecosystem/proxy-server/transport"
	saml2 "github.com/russellhaering/gosaml2"
	"golang.org/x/oauth2"
)

func oidcLogin(defaultURL string) fwncs.HandlerFunc {
	return func(c fwncs.Context) {
		r := c.Request()
		ctx := context.Background()
		client := c.HttpClient(transport.NewDumpTransport(c.Logger(), nil))
		ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
		c.SetContext(ctx)
		q := r.URL.Query()

		// Generate random state
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
			return
		}
		state := base64.StdEncoding.EncodeToString(b)
		session := sessions.Default(c)
		session.Set("state", state)
		redirect := q.Get("redirect")
		if redirect == "" {
			redirect = defaultURL
		}
		session.Set("redirect", redirect)
		err = session.Save()
		if err != nil {
			c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
			return
		}
		authenticator := c.Get(IdpContextKey).(*config.Authenticator)
		u := authenticator.Config.AuthCodeURL(state, authenticator.SetValues()...)
		c.Redirect(http.StatusTemporaryRedirect, u)
	}
}

func samlLogin(c fwncs.Context) {
	sp := c.Get(spKey).(*saml2.SAMLServiceProvider)
	u, err := sp.BuildAuthURL("")
	if err != nil {
		c.Logger().Error(err)
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, u)
}
