package kvs

import (
	"context"
	"github.com/pkg/errors"
	"strings"
	"sync"
)

var (
	MapStore = mapStore{
		nss: make(map[string]**sync.Map),
		mu:  new(sync.Mutex),
	}.WithNs("").(mapStore)
)

type mapStore struct {
	nss map[string]**sync.Map // ns -> map
	mu  *sync.Mutex

	ns string
	m  **sync.Map
}

func (s mapStore) Ns() string {
	return s.ns
}

func (s mapStore) WithNs(ns string) Kvs {
	s.mu.Lock()
	defer s.mu.Unlock()

	cpy := s
	m, ok := s.nss[ns]
	if !ok {
		p := &sync.Map{}
		m = &p
		s.nss[ns] = m
	}
	cpy.m = m
	cpy.ns = ns

	return cpy
}

func (s mapStore) WithCodec(Codec) Kvs {
	return s
}

func (s mapStore) WithCtx(context.Context) Kvs {
	return s
}

func (s mapStore) Get(key string) (interface{}, error) {
	data, ok := (*s.m).Load(key)
	if !ok {
		return nil, ErrNotFound
	}
	return data, nil
}

func (s mapStore) Set(key string, data interface{}) error {
	(*s.m).Store(key, data)
	return nil
}

func (s mapStore) Del(key string) error {
	(*s.m).Delete(key)
	return nil
}

func (s mapStore) Keys(re string) ([]string, error) {
	keys := make([]string, 0, 8)
	(*s.m).Range(func(key, _ interface{}) bool {
		if str := key.(string); strings.HasPrefix(str, re) {
			keys = append(keys, str)
		}
		return true
	})
	return keys, nil
}

func (s mapStore) MultiGet(keys ...string) ([]interface{}, error) {
	vs := make([]interface{}, 0, len(keys))

	for _, key := range keys {
		v, err := s.Get(key)
		if err != nil {
			return nil, err
		}

		vs = append(vs, v)
	}

	return vs, nil
}

func (s mapStore) MultiSet(pairs ...interface{}) (err error) {
	for i := 0; i < len(pairs); i += 2 {
		err = s.Set(pairs[i].(string), pairs[i+1])
		if err != nil {
			return errors.WithMessagef(err, "set %s %v", pairs[i].(string), pairs[i+1])
		}
	}
	return
}

func (s mapStore) MultiDel(keys ...string) error {
	for _, key := range keys {
		if err := s.Del(key); err != nil {
			return err
		}
	}

	return nil
}

func (s mapStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	*s.m = &sync.Map{}

	return nil
}

func (s mapStore) GetContext(_ context.Context, key string) (interface{}, error) {
	return s.Get(key)
}

func (s mapStore) SetContext(_ context.Context, key string, data interface{}) error {
	return s.Set(key, data)
}

func (s mapStore) DelContext(_ context.Context, key string) error {
	return s.Del(key)
}

func (s mapStore) KeysContext(_ context.Context, re string) ([]string, error) {
	return s.Keys(re)
}

func (s mapStore) MultiGetContext(_ context.Context, keys ...string) ([]interface{}, error) {
	return s.MultiGet(keys...)
}

func (s mapStore) MultiSetContext(_ context.Context, pairs ...interface{}) error {
	return s.MultiSet(pairs...)
}

func (s mapStore) MultiDelContext(_ context.Context, keys ...string) error {
	return s.MultiDel(keys...)
}

func (s mapStore) Ping(_ context.Context) error {
	return nil
}
