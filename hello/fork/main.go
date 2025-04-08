//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {
	pid, _, err := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)

	if err != 0 {
		fmt.Println("fork failed:", err)
		os.Exit(1)
	}

	if pid == 0 {
		// 子进程
		fmt.Println("Child process exiting...")
		os.Exit(0)
	} else {
		// 父进程
		fmt.Println("Parent process, child PID:", pid)
		fmt.Println("Sleeping for 60 seconds... Check `ps aux | grep Z` to see the zombie.")
		time.Sleep(60 * time.Second)

		// 你可以试试把下面这行取消注释，看看僵尸会不会消失：
		// var status syscall.WaitStatus
		// syscall.Wait4(int(pid), &status, 0, nil)
	}
}
