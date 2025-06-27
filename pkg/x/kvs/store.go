package kvs

import (
	"context"
	"fmt"
	"hellogo/pkg/x/kvs/driver"
	"net/url"
	"strings"
	"time"
)

type Store struct {
	dsn     string
	ns      string
	timeout time.Duration

	driverName string
	driver     driver.Driver

	se  driver.Session
	mse driver.MultiSession
	p   driver.Pinger

	ctx   context.Context
	codec Codec
}

func newStore(driverName, dsn string, ns ...string) (*Store, error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", driverName)
	}

	sess, err := driveri.Open(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s.Open %s %w", driverName, dsn, err)
	}

	s := &Store{
		dsn:        dsn,
		driverName: driverName,
		timeout:    5 * time.Second,
		driver:     driveri,
		se:         sess,
		codec:      Raw,
	}
	s.mse, _ = s.se.(driver.MultiSession)
	s.p, _ = s.se.(driver.Pinger)
	if len(ns) > 0 {
		s = s.WithNs(ns[0]).(*Store)
	}

	i := strings.LastIndexByte(dsn, '?')
	if i < 1 { // no additional
		return s, nil
	}

	q, err := url.ParseQuery(dsn[i+1:])
	if err != nil {
		return nil, fmt.Errorf("ParseQuery %s %w", dsn[i+1:], err)
	}
	if timeout, err := time.ParseDuration(q.Get("timeout")); err != nil && timeout > time.Millisecond {
		s.timeout = timeout
	}
	if namespace := q.Get("namespace"); s.ns == "" {
		s.ns = namespace
	}

	return s, nil
}

func (s *Store) clone() *Store {
	cpy := *s
	return &cpy
}

func (s *Store) Ns() string {
	return s.ns
}

func (s *Store) WithNs(ns string) Kvs {
	if s.ns == ns {
		return s
	}

	cpy := s.clone()
	cpy.ns = ns
	cpy.se = s.se.WithNs(ns)
	cpy.mse, _ = cpy.se.(driver.MultiSession)

	return cpy
}

func (s *Store) WithCodec(codec Codec) Kvs {
	if s.codec == codec {
		return s
	}
	cpy := s.clone()
	cpy.codec = codec
	return cpy
}

func (s *Store) WithCtx(ctx context.Context) Kvs {
	if ctx == s.ctx {
		return s
	}
	cpy := s.clone()
	cpy.ctx = ctx
	return cpy
}

func (s *Store) ctxOrTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	if s.ctx != nil {
		return s.ctx, func() {}
	}
	return context.WithTimeout(context.Background(), timeout)
}

func (s *Store) GetContext(ctx context.Context, key string) (interface{}, error) {
	b, err := s.se.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	return s.codec.Unmarshal(b)
}

func (s *Store) SetContext(ctx context.Context, key string, v interface{}) error {
	b, err := s.codec.Marshal(v)
	if err != nil {
		return err
	}

	return s.se.Set(ctx, key, b)
}

func (s *Store) DelContext(ctx context.Context, key string) error {
	return s.se.Del(ctx, key)
}

func (s *Store) KeysContext(ctx context.Context, re string) ([]string, error) {
	return s.se.Keys(ctx, re)
}

func (s *Store) MultiGetContext(ctx context.Context, keys ...string) ([]interface{}, error) {
	if s.mse == nil {
		return s.multiGetWrap(ctx, keys...)
	}

	vs := make([]interface{}, len(keys), len(keys))
	bs, err := s.mse.MultiGet(ctx, keys...)
	if err != nil {
		return nil, err
	}

	for i := range bs {
		if bs[i] == nil {
			continue
		}

		v, err := s.codec.Unmarshal(bs[i])
		if err != nil && err != ErrNotFound {
			return nil, fmt.Errorf("%s unmarshal %x %w", s.codec.String(), bs[i], err)
		} else if err == ErrNotFound {
			continue
		}

		vs[i] = v
	}

	return vs, nil
}

func (s *Store) multiGetWrap(ctx context.Context, keys ...string) ([]interface{}, error) {
	vs := make([]interface{}, len(keys), len(keys))

	for i := range keys {
		b, err := s.se.Get(ctx, keys[i])
		if err != nil && err != ErrNotFound {
			return nil, err
		}

		if b == nil {
			continue
		}
		v, err := s.codec.Unmarshal(b)
		if err != nil && err != ErrNotFound {
			return nil, fmt.Errorf("%s unmarshal %x %w", s.codec.String(), b, err)
		} else if err == ErrNotFound {
			continue
		}
		vs[i] = v
	}

	return vs, nil
}

func (s *Store) MultiSetContext(ctx context.Context, pairs ...interface{}) error {
	if s.mse == nil {
		return s.multiSetWrap(ctx, pairs...)
	}

	keys := make([]string, 0, len(pairs)/2)
	bs := make([][]byte, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		keys = append(keys, pairs[i].(string))

		b, err := s.codec.Marshal(pairs[i+1])
		if err != nil {
			return err
		}
		bs = append(bs, b)
	}

	return s.mse.MultiSet(ctx, keys, bs)
}

func (s *Store) multiSetWrap(ctx context.Context, pairs ...interface{}) error {
	for i := 0; i < len(pairs); i += 2 {
		b, err := s.codec.Marshal(pairs[i+1])
		if err != nil {
			return err
		}

		err = s.se.Set(ctx, pairs[i].(string), b)
		if err != nil {
			return fmt.Errorf("set %s %v %w", pairs[i].(string), b, err)
		}
	}
	return nil
}

func (s *Store) MultiDelContext(ctx context.Context, keys ...string) error {
	if s.mse == nil {
		return s.multiDelWrap(ctx, keys...)
	}

	return s.mse.MultiDel(ctx, keys...)
}

func (s *Store) multiDelWrap(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		err := s.se.Del(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Ping(ctx context.Context) error {
	if s.p == nil {
		return nil
	}
	return s.p.Ping(ctx)
}

func (s *Store) Get(key string) (interface{}, error) {
	ctx, cancel := s.ctxOrTimeout(s.timeout)
	defer cancel()

	return s.GetContext(ctx, key)
}

func (s *Store) Set(key string, v interface{}) error {
	ctx, cancel := s.ctxOrTimeout(s.timeout)
	defer cancel()

	return s.SetContext(ctx, key, v)
}

func (s *Store) Del(key string) error {
	ctx, cancel := s.ctxOrTimeout(s.timeout)
	defer cancel()

	return s.DelContext(ctx, key)
}

func (s *Store) Keys(re string) ([]string, error) {
	ctx, cancel := s.ctxOrTimeout(s.timeout)
	defer cancel()

	return s.KeysContext(ctx, re)
}

func (s *Store) MultiGet(keys ...string) ([]interface{}, error) {
	ext := time.Duration(len(keys)) * 100 * time.Microsecond
	ctx, cancel := s.ctxOrTimeout(s.timeout + ext)
	defer cancel()

	return s.MultiGetContext(ctx, keys...)
}

func (s *Store) MultiSet(pairs ...interface{}) error {
	ext := time.Duration(len(pairs)/2) * 100 * time.Microsecond
	ctx, cancel := s.ctxOrTimeout(s.timeout + ext)
	defer cancel()

	return s.MultiSetContext(ctx, pairs...)
}

func (s *Store) MultiDel(keys ...string) error {
	ext := time.Duration(len(keys)) * 100 * time.Microsecond
	ctx, cancel := s.ctxOrTimeout(s.timeout + ext)
	defer cancel()

	return s.MultiDelContext(ctx, keys...)
}

func (s *Store) Close() error {
	return s.se.Close()
}
