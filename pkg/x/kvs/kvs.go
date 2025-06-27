package kvs

import (
	"context"
	"hellogo/pkg/x/kvs/driver"
	"strings"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound   = status.Error(codes.NotFound, "kvs: not found")
	ErrNotSupport = status.Error(codes.Unimplemented, "kvs: not support")
)

// Kvs is same as kvs.IKvs, but avoid import the impl
type Kvs interface {
	Ns() string
	WithNs(ns string) Kvs
	WithCodec(codec Codec) Kvs
	WithCtx(ctx context.Context) Kvs

	Get(string) (interface{}, error)
	Set(string, interface{}) error
	Del(string) error

	Keys(string) ([]string, error)

	MultiGet(keys ...string) ([]interface{}, error)
	MultiSet(pairs ...interface{}) error
	MultiDel(keys ...string) error

	Close() error
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("kvs: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("kvs: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// NewKvs instance Kvs by DSN
func NewKvs(DSN string, ns ...string) (Kvs, error) {
	var (
		driverName = strings.Split(strings.Split(DSN, "://")[0], ",")[0]
		dsn        = strings.TrimPrefix(DSN, driverName+",")
	)

	if driverName == "mem" {
		return MapStore, nil
	}
	return newStore(driverName, dsn, ns...)
}
