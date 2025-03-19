package http

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttps(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		name := query.Get("name")
		if name == "" {
			writer.Write([]byte("url param 'name' missing"))
			return
		}
		writer.Write([]byte(fmt.Sprintf("hello, %s!", name)))
	})
	if err := http.ListenAndServeTLS("0.0.0.0:2230", "", "", mux); err != nil {
		t.Fatal(err)
	}
}
