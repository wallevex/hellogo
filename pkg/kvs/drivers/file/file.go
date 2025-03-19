package file

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"hellogo/pkg/kvs"
	"hellogo/pkg/kvs/driver"
)

func init() {
	kvs.Register("file", &Driver{})
}

type Driver struct{}

func (d *Driver) Open(dsn string) (driver.Session, error) {
	return NewSession(dsn)
}

type Session struct {
	dir string
	ns  string

	full string
}

func NewSession(dsn string) (*Session, error) {
	dir := strings.TrimPrefix(dsn, "file://")
	ok, err := PathExists(dir)
	if err != nil {
		return nil, err
	} else if !ok {
		if err = os.MkdirAll(dir, 0644); err != nil {
			return nil, err
		}
	}
	return &Session{dir: dir, ns: "", full: dir}, nil
}

func (s Session) WithNs(ns string) driver.Session {
	if s.ns == ns {
		return s
	}

	cpy := s
	cpy.ns = ns
	cpy.full = filepath.Join(cpy.dir, cpy.ns)
	_ = os.Mkdir(cpy.full, 0644)
	return cpy
}

func (s Session) Get(ctx context.Context, key string) ([]byte, error) {
	// TODO:
	key = strings.Replace(key, ":", ".", -1)

	fp := filepath.Join(s.full, key)
	ok, err := PathExists(fp)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, kvs.ErrNotFound
	}
	return ioutil.ReadFile(fp)
}

func (s Session) Set(ctx context.Context, key string, data []byte) error {
	// TODO:
	key = strings.Replace(key, ":", ".", -1)

	fp := filepath.Join(s.full, key)
	ok, err := PathExists(fp)
	if err != nil {
		return err
	} else if ok {
		if err = os.Remove(fp); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(fp, data, 0644)
}

func (s Session) Del(ctx context.Context, key string) error {
	// TODO:
	key = strings.Replace(key, ":", ".", -1)

	return os.Remove(filepath.Join(s.full, key))
}

func (s Session) Keys(ctx context.Context, re string) ([]string, error) {
	fis, err := ioutil.ReadDir(s.full)
	if err != nil {
		return nil, err
	}

	matches := make([]string, 0, 8)
	for _, fi := range fis {
		if ok, err := filepath.Match(re, fi.Name()); err != nil {
			return nil, err
		} else if ok {
			// TODO:
			matches = append(matches, strings.Replace(fi.Name(), ".", ":", -1))
		}
	}

	return matches, nil
}

func (s Session) Close() error {
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
