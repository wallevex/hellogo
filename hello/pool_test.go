package hello

import (
	"bytes"
	"sync"
	"testing"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var N = 100000

func BenchmarkBufferPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var g sync.WaitGroup
		g.Add(N)
		for j := 0; j < N; j++ {
			go func() {
				buf := bufferPool.Get().(*bytes.Buffer)
				buf.WriteString("hello world")
				bufferPool.Put(buf)
				g.Done()
			}()
		}
		g.Wait()
	}
}

func BenchmarkBufferNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var g sync.WaitGroup
		g.Add(N)
		for j := 0; j < N; j++ {
			go func() {
				buf := new(bytes.Buffer)
				buf.WriteString("hello world")
				g.Done()
			}()
		}
		g.Wait()
	}
}
