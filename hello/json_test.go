package hello

import (
	"encoding/json"
	"sync"
	"sync/atomic"
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

func newGoRoutine(wg sync.WaitGroup, i *int64) {
	defer wg.Done()
	atomic.AddInt64(i, 1)
	return
}

func TestR(t *testing.T) {
	var wg sync.WaitGroup

	ans := int64(0)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go newGoRoutine(wg, &ans)
	}
	wg.Wait()
	recover()
}
