package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3)
	for i := 0; i < 10; i++ {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			fmt.Printf("hello world %d\n", i)
		}(i)
	}
	wg.Wait()
}
