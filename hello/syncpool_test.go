package hello

import (
	"fmt"
	"sync"
	"testing"
)

type person struct {
	fullName string
}

var loop = 10000

// go test -v -bench=NoPool -benchmem
func BenchmarkNoPool(b *testing.B) {
	fn := func(firstName, lastName string) *person {
		return &person{fullName: fmt.Sprintf("%s %s", firstName, lastName)}
	}
	for i := 0; i < b.N; i++ {
		var g sync.WaitGroup
		g.Add(loop)
		for n := 0; n < loop; n++ {
			go func() {
				fn("foo", "bar")
				g.Done()
			}()
		}
		g.Wait()
	}
}

// go test -v -bench=SyncPoo -benchmem
func BenchmarkSyncPool(b *testing.B) {
	var personPool = sync.Pool{
		New: func() any {
			return &person{}
		},
	}
	fn := func(firstName, lastName string) *person {
		p := personPool.Get().(*person)
		p.fullName = fmt.Sprintf("%s %s", firstName, lastName)
		personPool.Put(p)
		return p
	}
	for i := 0; i < b.N; i++ {
		var g sync.WaitGroup
		g.Add(loop)
		for n := 0; n < loop; n++ {
			go func() {
				fn("foo", "bar")
				g.Done()
			}()
		}
		g.Wait()
	}
}
