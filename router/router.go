package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n-creativesystem/fwncs-contrib/prometheus"
	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/constant"
	"github.com/n-creativesystem/go-fwncs/sessions"
	"github.com/n-creativesystem/go-fwncs/sessions/redis"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	saml2 "github.com/russellhaering/gosaml2"
)

const (
	OidcContextKey = "oidc_key"
	AuthContextKey = "auth_key"
)

type MultiHostHandler map[string]http.Handler

func (mh MultiHostHandler) getCaPublic(w http.ResponseWriter, req *http.Request) {
	caName := config.NewCaName(config.File.Certificate)
	buf, err := caName.ReadCASSHPublicKey()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	header := w.Header()
	header.Set(constant.HeaderContentType, "application/octet-stream")
	header.Set(constant.HeaderContentDisposition, "attachment;filename=\"ca.key.pub\"")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func (mh MultiHostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := mh[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func (mh MultiHostHandler) Run(port int) error {
	l, err := utils.GetListen(port)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Handler: mh,
	}
	go func() {
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	return mh.run(srv)
}

// RunTLS is https
func (mh MultiHostHandler) RunTLS(port int, certFile, keyFile string) error {
	if certFile == "" {
		return errors.New("certFile is empty")
	}
	if keyFile == "" {
		return errors.New("keyFile is empty")
	}
	if !utils.IsExists(certFile) {
		return fmt.Errorf("%s is No such file or directory", certFile)
	}
	if !utils.IsExists(keyFile) {
		return fmt.Errorf("%s is No such file or directory", keyFile)
	}
	l, err := utils.GetListen(port)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Handler: mh,
	}
	go func() {
		if err := srv.ServeTLS(l, certFile, keyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	return mh.run(srv)
}

func (mh MultiHostHandler) run(s *http.Server) error {
	signals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGSTOP,
	}
	osNotify := make(chan os.Signal, 1)
	signal.Notify(osNotify, signals...)
	sig := <-osNotify
	log.Println(fmt.Sprintf("signal: %v", sig))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return s.Shutdown(ctx)
}

func New(conf config.Config, authfilename string, authConfig config.Auth, loadBalancer config.LoadBalancer) *fwncs.Router {
	var oidcConfig config.OidcConfig
	var samlConfig config.SamlConfig
	var sp *saml2.SAMLServiceProvider
	switch conf.AuthType {
	case "oidc":
		oidcConfig = config.NewOidcConfig(authfilename)
	case "saml":
		samlConfig = config.NewSamlConfig(authfilename)
		var err error
		sp, err = samlConfig.Load()
		if err != nil {
			panic(err)
		}
	}
	// multiHostHandler := make(MultiHostHandler, len(loadBalancers))
	redisEndpoints := conf.RedisEndpoints
	// for _, loadBalancer := range loadBalancers {
	router := fwncs.New()
	router.Use(GzipMiddleware())
	router.Use(prometheus.InstrumentHandlerInFlight, prometheus.InstrumentHandlerDuration, prometheus.InstrumentHandlerCounter, prometheus.InstrumentHandlerResponseSize)
	router.GET("= /metrics", fwncs.WrapHandler(promhttp.Handler()))
	router.Use(func(c fwncs.Context) {
		c.Set(AuthContextKey, authConfig)
		c.Next()
	}, sessions.Sessions("auth_server", redis.NewStore(&redis.RedisOptions{
		Username:  conf.RedisUsername,
		Password:  conf.RedisPassword,
		Endpoints: redisEndpoints,
		KeyPairs:  []byte("secure"),
	})))
	// if loadBalancer.SSH {
	// 	router.GET("= /ssh", func(c fwncs.Context) {
	// 		t := template.Must(template.ParseFiles("public/templates/ssh.html"))
	// 		t.Execute(c.Writer(), map[string]string{
	// 			"Host": c.QueryParam("host"),
	// 			"Port": c.QueryParam("port"),
	// 		})
	// 	})
	// 	router.GET("= /term", handlerSSH)
	// }
	switch conf.AuthType {
	case "oidc":
		router.Use(OpenIdConnectSetting(oidcConfig))
		auth := router.Group(authConfig.Path)
		{
			auth.GET("= "+authConfig.Login, oidcLogin(loadBalancer.DefaultURL))
			auth.GET("= "+authConfig.Callback, oidcCallback)
			auth.GET("= "+authConfig.Logout, oidcLogout(oidcConfig.Logout))
			auth.GET("= "+authConfig.UserInfo, oidcUserInfo)
		}
		setOidcLoadbalancer(router, loadBalancer)
		if loadBalancer.Portal {
			router.GET("= /portal", oidcPortal)
			router.ServeFiles("/js", http.Dir("public/js"))
			router.ServeFiles("/images", http.Dir("public/images"))
		}
	case "saml":
		auth := router.Group(authConfig.Path)
		{
			auth.Use(setServiceProvider(sp))
			auth.GET("= "+authConfig.Login, samlLogin)
			auth.POST("= "+authConfig.Callback, samlCallback)
			// auth.GET("= "+authConfig.Logout, oidcLogout(oidcConfig.Logout))
			// auth.GET("= "+authConfig.UserInfo, oidcUserInfo)
		}
		router.GET("= /portal", samlPortal)
	}
	// multiHostHandler[fmt.Sprintf("%s:%d", loadBalancer.Domain, conf.Port)] = router
	// }
	return router
}
