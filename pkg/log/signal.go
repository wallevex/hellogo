// +build !windows

package log

import (
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
)

/**
 * HandleSignalChangeLogLevel 实现了根据信号量修改日志等级的方法
 * 使用方法 go HandleSignalChangeLogLevel()
 * 监听2个信号，kill -10 pid 升高日志等级
 * kill -12 pid 降低日志等级
 */
// Deprecated: use /debug/log/level
func HandleSignalChangeLogLevel() {
	down := make(chan os.Signal)
	up := make(chan os.Signal)
	signal.Notify(down, syscall.SIGUSR1)
	signal.Notify(up, syscall.SIGUSR2)
	for {
		select {
		case <-down:
			if lv := _globalLv.Level(); lv > zapcore.DebugLevel {
				_globalLv.SetLevel(lv - 1)
			}
		case <-up:
			if lv := _globalLv.Level(); lv < zapcore.FatalLevel {
				_globalLv.SetLevel(lv + 1)
			}
		}
	}
}
