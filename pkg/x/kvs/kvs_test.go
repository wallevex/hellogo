package kvs

import (
	"bytes"
	"context"
	"fmt"
	"hellogo/pkg/x/kvs/driver"
	"sync"
	"testing"
)

func init() {
	Register("test", testDriver{})
}

type testDriver struct{}

func (testDriver) Open(dsn string) (driver.Session, error) {
	fmt.Println("testDriver opened", dsn)
	return testSession{dsn: dsn, Map: new(sync.Map)}, nil
}

type testSession struct {
	dsn string
	ns  string
	*sync.Map
}

func (s testSession) WithNs(ns string) driver.Session {
	cpy := s
	cpy.ns = ns
	return cpy
}

func (s testSession) Get(ctx context.Context, key string) ([]byte, error) {
	data, ok := s.Map.Load(key)
	if !ok {
		return nil, ErrNotFound
	}
	return data.([]byte), nil
}

func (s testSession) Set(ctx context.Context, key string, data []byte) error {
	s.Map.Store(key, data)
	return nil
}

func (s testSession) Del(ctx context.Context, key string) error {
	s.Map.Delete(key)
	return nil
}

func (s testSession) Keys(context.Context, string) ([]string, error) {
	panic("implement me")
}

func (s testSession) Close() error {
	s.Map = nil
	return nil
}

func TestKvs(t *testing.T) {
	var (
		dsn  = "test://root:Pass@127.0.0.1:1234/test?charset=utf8"
		key  = "aaa"
		data = []byte("bbb")
	)

	store, err := NewKvs(dsn)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Get(key)
	if err != ErrNotFound {
		t.Fatal("Get", key, "should return ErrNotFound")
	}

	err = store.Set(key, data)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := store.Get(key)
	if !bytes.Equal(data, data2.([]byte)) {
		t.Fatal(data, "is not", data2)
	}

	err = store.Del(key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Get(key)
	if err != ErrNotFound {
		t.Fatal("Get", key, "should return ErrNotFound")
	}
}

type testJSON struct {
	A string  `json:"a"`
	B int     `json:"b"`
	C float64 `json:"c"`
}

func TestJsonCodec(t *testing.T) {
	var (
		dsn  = "test://root:Pass@127.0.0.1:1234/test?charset=utf8"
		key  = "aaa"
		data = &testJSON{A: "a", B: 1, C: 2.1}
	)

	store, err := NewKvs(dsn)
	if err != nil {
		t.Fatal(err)
	}
	store = store.WithCodec(JSON(data))

	_, err = store.Get(key)
	if err != ErrNotFound {
		t.Fatal("Get", key, "should return ErrNotFound")
	}

	err = store.Set(key, data)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := store.Get(key)
	if data2.(*testJSON).A != data.A ||
		data2.(*testJSON).B != data.B ||
		data2.(*testJSON).C != data.C {
		t.Fatal(data, "is not", data2)
	}

	err = store.Del(key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Get(key)
	if err != ErrNotFound {
		t.Fatal("Get", key, "should return ErrNotFound")
	}
}

func TestMemKvs(t *testing.T) {
	var (
		dsn  = "mem://root:Pass@127.0.0.1:1234/test?charset=utf8"
		key  = "aaa"
		data = []byte("bbb")
	)

	store, err := NewKvs(dsn)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Get(key)
	if err != ErrNotFound {
		t.Fatal("Get", key, "should return ErrNotFound but", err)
	}

	err = store.Set(key, data)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := store.Get(key)
	if !bytes.Equal(data, data2.([]byte)) {
		t.Fatal(data, "is not", data2)
	}

	err = store.Del(key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Get(key)
	if err != ErrNotFound {
		t.Fatal("Get", key, "should return ErrNotFound")
	}
}

func TestMemStoreJson(t *testing.T) {
	var (
		dsn  = "mem://root:Pass@127.0.0.1:1234/test?charset=utf8"
		key  = "aaa"
		data = testJSON{A: "a", B: 1, C: 2.1}
	)

	store, err := NewKvs(dsn)
	if err != nil {
		t.Fatal(err)
	}
	store = store.WithCodec(JSON(data))

	_, err = store.Get(key)
	if err != ErrNotFound {
		t.Fatal("Get", key, "should return ErrNotFound")
	}

	err = store.Set(key, data)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := store.Get(key)
	if data2.(testJSON).A != data.A ||
		data2.(testJSON).B != data.B ||
		data2.(testJSON).C != data.C {
		t.Fatal(data, "is not", data2)
	}

	err = store.Del(key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Get(key)
	if err != ErrNotFound {
		t.Fatal("Get", key, "should return ErrNotFound")
	}
}
