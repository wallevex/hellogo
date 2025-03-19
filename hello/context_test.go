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

	sum, err := ComplexCal(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(sum)
}

// 使用Deadline context的标准形式
func ComplexCal(ctx context.Context, num int) (int, error) {
	ch := make(chan struct{}, 1)
	sum := 0
	go func() {
		for i := 1; i <= num; i++ {
			sum += i
		}
		time.Sleep(time.Second * 2)
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-ch:
		return sum, nil
	}
}
