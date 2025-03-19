package http

import (
	"fmt"
	"net/http"
	"testing"
)

func search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	text, ok := query["text"]
	if !ok || len(text) < 1 {
		w.Write([]byte("url param 'text' missing"))
		return
	}
	w.Write([]byte("search " + text[0] + " result: xxx.."))
}

func login(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id, ok := query["id"]
	if !ok || len(id) < 1 {
		w.Write([]byte("url param 'id' missing"))
		return
	}
	pwd, ok := query["pwd"]
	if !ok || len(pwd) < 1 {
		w.Write([]byte("url param 'pwd' missing"))
		return
	}
	w.Write([]byte(fmt.Sprintf("id=%s pwd=%s login successful!!", id[0], pwd[0])))
}

func TestMockGoogleSearchService(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", search)

	if err := http.ListenAndServe(":18888", mux); err != nil {
		t.Fatal(err)
	}
}

func TestMockGoogleUserService(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/user/login", login)

	if err := http.ListenAndServe(":28888", mux); err != nil {
		t.Fatal(err)
	}
}
