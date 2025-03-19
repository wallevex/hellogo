package hello

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"hellogo/pkg/log"
)

// nil channel
// 1. 向nil channel发送消息会阻塞
// 2. 从nil channel接收消息会阻塞
// 3. 关闭nil channel会panic
func TestNilChannel(t *testing.T) {
	var ch chan struct{}
	go func() {
		fmt.Println("before")
		ch <- struct{}{}
		fmt.Println("after")
	}()

	for {
	}
}

// The close built-in function closes a channel, which must be either
// bidirectional or send-only. It should be executed only by the sender,
// never the receiver, and has the effect of shutting down the channel after
// the last sent value is received. After the last value has been received
// from a closed channel c, any receive from c will succeed without
// blocking, returning the zero value for the channel element. The form
//
//	x, ok := <-c
//
// will also set ok to false for a closed channel.
//
// close channel
// 1. 已关闭的channel可以不断的读取消息，读到的是消息的零值，通过 _, ok := <- ch来判断channel是否已被关闭
// 2. 关闭已被关闭的channel会panic
func TestCloseChannel(t *testing.T) {
	ch := make(chan struct{})
	close(ch)
	for {
		_, ok := <-ch
		fmt.Println(ok)
		time.Sleep(1 * time.Second)
	}
}

func TestMergeChannels(t *testing.T) {
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					log.Info("a is done")
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					log.Info("b is done")
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}
