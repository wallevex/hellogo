package hello

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	a, err := Fibonacci(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(a)
}

// 长时间的批量数据处理需要支持context
func Fibonacci(ctx context.Context, n int) (int64, error) {
	var a int64 = 0
	var b int64 = 1
	for i := 2; i <= n; i++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			a, b = b, a+b
		}
	}
	return b, nil
}
