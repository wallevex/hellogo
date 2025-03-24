package hello

import (
	"fmt"
	"testing"
)

// A: 进入default分支，打印"DEFAULT, "
// B: 进入shutdown分支，打印"CLOSED, "
// C: 进入data分支，打印"HAS WRITTEN, "
// D: 程序会panic
// E: 程序可能panic，也可能打印"CLOSED, "
//
// 答案是E. 机制如下：
// 1. 从已被关闭的管道可以接收到零值
// 2. 向已被关闭的管道发送数据会panic
// 3. select会随机选取不会阻塞的case分支执行，如果所有case分支都阻塞才会执行default分支
func TestClose(t *testing.T) {
	data := make(chan int)
	shutdown := make(chan int)
	close(shutdown)
	close(data)

	select {
	case <-shutdown:
		fmt.Println("CLOSED, ")
	case data <- 1:
		fmt.Println("HAS WRITTEN, ")
	default:
		fmt.Println("DEFAULT, ")
	}
}
