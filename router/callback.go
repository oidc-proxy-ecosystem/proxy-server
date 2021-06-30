package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/sessions"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/oidc-proxy-ecosystem/proxy-server/transport"
	saml2 "github.com/russellhaering/gosaml2"
	"golang.org/x/oauth2"
)

func oidcCallback(c fwncs.Context) {
	log := c.Logger()
	r := c.Request()
	ctx := context.Background()
	client := c.HttpClient(transport.NewDumpTransport(c.Logger(), nil))
	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
	session := sessions.Default(c)
	authConfig := c.Get(AuthContextKey).(config.Auth)
	if r.URL.Query().Get("state") != session.Get("state") {
		log.Debug(fmt.Sprintf("request_state:%s", r.URL.Query().Get("state")))
		log.Debug(fmt.Sprintf("session_state:%s", session.Get("state")))
		c.Redirect(http.StatusTemporaryRedirect, authConfig.Login)
		return
	}
	auth := c.Get(IdpContextKey).(*config.Authenticator)
	token, err := auth.Config.Exchange(ctx, r.URL.Query().Get("code"), auth.SetValues()...)
	if err != nil {
		log.Critical(fmt.Sprintf("no token found: %v", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	session.Clear()
	err = SaveToken(session, token, auth)
	if err != nil {
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
		return
	}
	redirect := "/portal"
	c.Redirect(http.StatusTemporaryRedirect, redirect)
}

func samlCallback(c fwncs.Context) {
	sp := c.Get(spKey).(*saml2.SAMLServiceProvider)
	err := c.Request().ParseForm()
	if err != nil {
		c.Logger().Error(err)
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
		return
	}
	assertionInfo, err := sp.RetrieveAssertionInfo(c.FormValue("SAMLResponse"))
	if err != nil {
		c.Logger().Error(err)
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
		return
	}
	if assertionInfo.WarningInfo.InvalidTime {
		c.Logger().Error(errors.New("invalid time"))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if assertionInfo.WarningInfo.NotInAudience {
		c.Logger().Error(errors.New("not in audience"))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	session := sessions.Default(c)
	session.Clear()
	session.Set("saml_user", assertionInfo)
	if err := session.Save(); err != nil {
		c.Logger().Error(err)
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
		return
	}
	redirect, ok := session.Get("redirect").(string)
	if !ok {
		redirect = "/portal"
	}
	// 303で強制的にGETメソッドへ変更
	c.Redirect(http.StatusSeeOther, redirect)
}
