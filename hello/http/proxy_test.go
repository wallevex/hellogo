package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	targetUrl, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("proxy-response", "this is a message added by proxy")
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
	return proxy, nil
}

func TestReverseProxy(t *testing.T) {
	go TestMockGoogleSearchService(t)
	go TestMockGoogleUserService(t)

	searchProxy, err := NewProxy("http://127.0.0.1:18888")
	if err != nil {
		t.Fatal(err)
	}
	userProxy, err := NewProxy("http://127.0.0.1:28888")
	if err != nil {
		t.Fatal(err)
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchProxy.ServeHTTP(w, r)
	})
	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":18080", nil); err != nil {
		t.Fatal(err)
	}
}
