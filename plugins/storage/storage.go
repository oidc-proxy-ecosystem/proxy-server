package storage

import (
	"github.com/oidc-proxy-ecosystem/proxy-server/database"
	"github.com/oidc-proxy-ecosystem/proxy-server/database/redis"
)

var Client database.Factory

func New(typ string, conf database.Config) error {
	switch typ {
	case "redis":
		Client = redis.New()
	default:
		Client = redis.New()
	}
	return Client.Initialize(conf)
}
