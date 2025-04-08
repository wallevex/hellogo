package hello

import (
	"encoding/json"
	"testing"
)

type S struct {
	A int
	B *int
	C float64
	D func() string
	E chan struct{}
}

func TestMarshal(t *testing.T) {
	s := S{
		A: 1,
		B: nil,
		C: 12.15,
		//D: func() string {
		//	return "NowCoder"
		//},
		E: make(chan struct{}),
	}

	_, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
}
