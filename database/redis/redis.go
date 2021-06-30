package redis

import (
	"context"
	"crypto/tls"

	"github.com/go-redis/redis/v8"
	"github.com/oidc-proxy-ecosystem/proxy-server/database"
	"github.com/oidc-proxy-ecosystem/proxy-server/utils"
)

type redisInterface interface {
	redis.Cmdable
	Close() error
}

func New() database.Factory {
	return &redisStore{}
}

type redisStore struct {
	username  string
	password  string
	endpoints []string
	cluster   bool
	tlsConfig *tls.Config
}

func (s *redisStore) Close() error {
	return nil
}

func (s *redisStore) Initialize(conf database.Config) error {
	s.endpoints = conf.Endpoints
	s.username = conf.UserName
	s.password = conf.Password
	s.tlsConfig = s.tlsConfig
	s.cluster = len(conf.Endpoints) > 1
	r, err := s.getClient()
	if err != nil {
		return err
	}
	defer r.Close()
	return nil
}

func (s *redisStore) loadClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:      s.endpoints[0],
		Username:  s.username,
		Password:  s.password,
		TLSConfig: s.tlsConfig,
	})
	err := client.Ping(context.Background()).Err()
	return client, err
}

func (s *redisStore) loadCluster() (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:     s.endpoints,
		Username:  s.username,
		Password:  s.password,
		TLSConfig: s.tlsConfig,
	})
	err := client.Ping(context.Background()).Err()
	return client, err
}

func (s *redisStore) getClient() (redisInterface, error) {
	if s.cluster {
		return s.loadCluster()
	} else {
		return s.loadClient()
	}
}

func (s *redisStore) Set(ctx context.Context, key string, buf []byte) error {
	client, err := s.getClient()
	if err != nil {
		return err
	}
	defer client.Close()
	value := utils.Base64Encode(buf)
	return client.Set(ctx, key, value, 0).Err()
}

func (s *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == string(redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return utils.Base64Decode(value)
}
