package hello

import (
	"testing"
	"time"
)

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func benchmarkFib(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(i)
	}
}

// -run=<func name>: go test -run=Fib40
func TestFib40(t *testing.T) {
	start := time.Now()
	ans := fib(40)
	t.Logf("ans: %d, spend: %vns", ans, time.Now().Sub(start).Nanoseconds())
}

// -bench=<regex expression>: go test -bench=Fib*
// -benchtime=<basic benchmark time>: go test -bench=Fib40 -benchtime=5s
func BenchmarkFib1(b *testing.B)  { benchmarkFib(1, b) }
func BenchmarkFib2(b *testing.B)  { benchmarkFib(2, b) }
func BenchmarkFib3(b *testing.B)  { benchmarkFib(3, b) }
func BenchmarkFib10(b *testing.B) { benchmarkFib(10, b) }
func BenchmarkFib20(b *testing.B) { benchmarkFib(20, b) }
func BenchmarkFib40(b *testing.B) { benchmarkFib(40, b) }
