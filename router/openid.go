package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/coreos/go-oidc"
	"github.com/form3tech-oss/jwt-go"
	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/sessions"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"golang.org/x/oauth2"
)

const (
	IdpContextKey = "idp_key"
)

func OpenIdConnectSetting(oidcConfig config.OidcConfig) fwncs.HandlerFunc {
	provider, err := oidc.NewProvider(context.Background(), oidcConfig.Provider)
	if err != nil {
		panic(err)
	}
	o2conf := oauth2.Config{
		ClientID:     oidcConfig.ClientId,
		ClientSecret: oidcConfig.ClientSecret,
		RedirectURL:  oidcConfig.CallbackUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       oidcConfig.Scopes,
	}
	auth := &config.Authenticator{
		Provider:   provider,
		Config:     o2conf,
		OidcConfig: oidcConfig,
	}
	return func(c fwncs.Context) {
		c.Set(OidcContextKey, oidcConfig)
		c.Set(IdpContextKey, auth)
		c.Next()
	}
}

func unAuthorized(c fwncs.Context) {
	c.Skip()
	auth := c.Get(AuthContextKey).(config.Auth)
	u := fmt.Sprintf("%s?redirect=%s", path.Join(auth.Path, auth.Login), c.RealPath())
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func SessionCheck(c fwncs.Context) {
	session := sessions.Default(c)
	if v, _ := session.Get("id_token").(string); v != "" {
		c.Next()
	} else {
		unAuthorized(c)
	}
}

func SetIdToken(c fwncs.Context) {
	req := c.Request()
	session := sessions.Default(c)
	if rawToken, ok := session.Get("id_token").(string); ok {
		token, _, err := new(jwt.Parser).ParseUnverified(rawToken, jwt.MapClaims{})
		if err != nil || token.Claims.Valid() != nil {
			err, _token := refresh(c)
			if err != nil {
				unAuthorized(c)
				return
			}
			rawToken = _token.Extra("id_token").(string)
		}
		req.Header.Set("Authorization", "Bearer "+rawToken)
	}
	c.SetRequest(req)
	c.Next()
}

func SetAccessToken(c fwncs.Context) {
	req := c.Request()
	session := sessions.Default(c)
	if rawToken, ok := session.Get("access_token").(string); ok {
		token, _, err := new(jwt.Parser).ParseUnverified(rawToken, jwt.MapClaims{})
		if err != nil || token.Claims.Valid() != nil {
			err, _token := refresh(c)
			if err != nil {
				unAuthorized(c)
				return
			}
			rawToken = _token.AccessToken
		}
		req.Header.Set("Authorization", "Bearer "+rawToken)
		c.SetRequest(req)
		c.Next()
	} else {
		unAuthorized(c)
	}
}

func refresh(c fwncs.Context) (error, *oauth2.Token) {
	ctx := context.Background()
	auth := c.Get(IdpContextKey).(*config.Authenticator)
	session := sessions.Default(c)
	refreshToken := session.Get("refresh_token").(string)
	src := auth.Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	token, err := src.Token()
	if err != nil {
		return err, nil
	}
	return SaveToken(session, token, auth), token
}

func SaveToken(session sessions.Session, token *oauth2.Token, auth *config.Authenticator) error {
	ctx := context.Background()
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return errors.New("No id_token field in oauth2 token.")
	}

	oidcConfig := &oidc.Config{
		ClientID: auth.Config.ClientID,
	}

	_, err := auth.Provider.Verifier(oidcConfig).Verify(ctx, rawIDToken)
	if err != nil {
		return fmt.Errorf("Failed to verify ID Token: %s", err.Error())
	}
	session.Set("token", token)
	session.Set("id_token", rawIDToken)
	session.Set("access_token", token.AccessToken)
	session.Set("refresh_token", token.RefreshToken)
	return session.Save()
}
