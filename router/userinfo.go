package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/sessions"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"golang.org/x/oauth2"
)

func getUserInfo(c fwncs.Context) (map[string]interface{}, error, bool) {
	auth := c.Get(IdpContextKey).(*config.Authenticator)
	ctx := c.GetContext()
	session := sessions.Default(c)
	_token, ok := session.Get("token").(map[string]interface{})
	if !ok {
		return nil, errors.New("no token"), true
	}
	accessToken, _ := _token["access_token"].(string)
	tokenType, _ := _token["token_type"].(string)
	refreshToken, _ := _token["refresh_token"].(string)
	expiry, _ := _token["expiry"].(time.Time)
	token := oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}
	src := auth.Config.TokenSource(ctx, &token)
	u, err := auth.Provider.UserInfo(ctx, src)
	if err != nil {
		return nil, err, true
	}
	claims := map[string]interface{}{}
	if err := u.Claims(&claims); err != nil {
		return nil, err, false
	}

	mp := map[string]interface{}{
		"email":          u.Email,
		"email_verified": u.EmailVerified,
		"profile":        u.Profile,
		"sub":            u.Subject,
	}
	for key, value := range claims {
		mp[key] = value
	}
	return mp, nil, false
}

func oidcUserInfo(c fwncs.Context) {
	mp, err, redirect := getUserInfo(c)
	if redirect {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		return
	}
	if err != nil {
		c.Logger().Error(err)
		c.AbortWithStatusAndErrorMessage(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, &mp)
}
