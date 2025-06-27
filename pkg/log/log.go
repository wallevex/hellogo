package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
)

func NewZapLogger(filename string, maxSize, maxBackups, maxAge int, compress bool) *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	writer := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	atomic := zap.NewAtomicLevelAt(zap.InfoLevel)
	http.HandleFunc("/log/level", atomic.ServeHTTP)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(writer),
		atomic,
	)
	return zap.New(core, zap.AddCaller())
}
