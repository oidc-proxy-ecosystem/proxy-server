package database

import (
	"context"
	"crypto/tls"
)

type Config struct {
	UserName, Password string
	Endpoints          []string
	TLSConfig          *tls.Config
}

type Factory interface {
	Initialize(Config) error
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Close() error
}
