package main

import (
	"net/http"
	_ "net/http/pprof"
)

func Debug(listen string) {
	if err := http.ListenAndServe(listen, nil); err != nil {
		panic(err)
	}
}
