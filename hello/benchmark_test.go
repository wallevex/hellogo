package hello

import (
	"testing"
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

// test单个程序：-run=<func name> 如go test -run=Fib40
// bench正则表达式：-bench=<regex expression> 如go test -bench=Fib*
// -benchtime=<basic benchmark time>: go test -bench=Fib40 -benchtime=5s
func BenchmarkFib1(b *testing.B)  { benchmarkFib(1, b) }
func BenchmarkFib2(b *testing.B)  { benchmarkFib(2, b) }
func BenchmarkFib3(b *testing.B)  { benchmarkFib(3, b) }
func BenchmarkFib10(b *testing.B) { benchmarkFib(10, b) }
func BenchmarkFib20(b *testing.B) { benchmarkFib(20, b) }
func BenchmarkFib40(b *testing.B) { benchmarkFib(40, b) }
