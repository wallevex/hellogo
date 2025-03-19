package driver

import (
	"context"
	"database/sql/driver"
)

type Pinger = driver.Pinger

type Driver interface {
	Open(dsn string) (Session, error)
}

type Session interface {
	WithNs(ns string) Session

	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte) error
	Del(ctx context.Context, key string) error

	Keys(ctx context.Context, re string) ([]string, error)

	Close() error
}

type MultiSession interface {
	MultiSet(ctx context.Context, keys []string, data [][]byte) error
	MultiGet(ctx context.Context, keys ...string) ([][]byte, error)
	MultiDel(ctx context.Context, keys ...string) error
}
