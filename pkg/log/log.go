package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	_globalMu   sync.RWMutex
	_globalL, _ = zap.NewDevelopment()
	_globalS    = _globalL.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Sugar() // use for s() only

	_globalLv zap.AtomicLevel
)

func ReplaceGlobals(logger *zap.Logger) func() {
	_globalMu.Lock()
	prev := _globalL
	_globalL = logger
	_globalS = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	_globalMu.Unlock()
	return func() { ReplaceGlobals(prev) }
}

func L() *zap.Logger {
	_globalMu.RLock()
	l := _globalL
	_globalMu.RUnlock()
	return l
}

func s() *zap.SugaredLogger {
	_globalMu.RLock()
	s := _globalS
	_globalMu.RUnlock()
	return s
}

const XLogAgentUrlEnv string = "XLOG_AGENT_URL"

var CoreLoki = zapcore.NewNopCore()

func init() {
	_globalLv = zap.NewAtomicLevel()
	http.HandleFunc("/debug/log/level", _globalLv.ServeHTTP)
}

func SetLogger(_, Dir, File string, Count int32, Size int64, Unit, Level string, compressT int64) error {
	roller := &lumberjack.Logger{
		Filename:   filepath.Join(Dir, File),
		MaxSize:    toMegaBytes(Size, Unit),
		MaxBackups: int(Count),
	}
	if compressT > 0 {
		roller.Compress = true
	}
	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewFileEncoderConfig()),
		zapcore.AddSync(roller),
		_globalLv,
	)

	_ = _globalLv.UnmarshalText([]byte(Level))

	ReplaceGlobals(L().WithOptions(zap.WrapCore(func(_ zapcore.Core) zapcore.Core {
		return zapcore.NewTee(CoreLoki, fileCore)
	})))
	return nil
}

func NewFileEncoderConfig() zapcore.EncoderConfig {
	enc := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if runtime.GOOS == "windows" {
		enc.EncodeLevel = zapcore.CapitalLevelEncoder
		enc.LineEnding = "\r\n"
	} else if runtime.GOOS == "darwin" {
		enc.LineEnding = "\r"
	}
	return enc
}

func toMegaBytes(Size int64, Unit string) int {
	switch Unit {
	case "KB":
		return 1
	case "MB":
		return int(Size)
	case "GB":
		return int(Size) * 1024
	case "TB":
		return int(Size) * 1024 * 1024
	}
	return 1
}

func Named(name string) Logger                  { return AddSpace(L()).Named(name) }
func With(pairs ...interface{}) Logger          { return AddSpace(L()).With(pairs...) }
func Debugf(format string, args ...interface{}) { s().Debugf(format, args...) }
func Infof(format string, args ...interface{})  { s().Infof(format, args...) }
func Warnf(format string, args ...interface{})  { s().Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { s().Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { s().Fatalf(format, args...) }
func Debug(args ...interface{})                 { s().Debug(join(args, " ")...) }
func Info(args ...interface{})                  { s().Info(join(args, " ")...) }
func Warn(args ...interface{})                  { s().Warn(join(args, " ")...) }
func Error(args ...interface{})                 { s().Error(join(args, " ")...) }
func Fatal(args ...interface{})                 { s().Fatal(join(args, " ")...) }
func Debugw(msg string, pairs ...interface{})   { s().Debugw(msg, pairs...) }
func Infow(msg string, pairs ...interface{})    { s().Infow(msg, pairs...) }
func Warnw(msg string, pairs ...interface{})    { s().Warnw(msg, pairs...) }
func Errorw(msg string, pairs ...interface{})   { s().Errorw(msg, pairs...) }
func Fatalw(msg string, pairs ...interface{})   { s().Fatalw(msg, pairs...) }

func join(a []interface{}, sep interface{}) []interface{} {
	switch len(a) {
	case 0:
		return []interface{}{}
	case 1:
		return a
	}

	r := make([]interface{}, 0, len(a)*2)
	for i := range a {
		r = append(r, a[i], " ")
	}

	return r[:len(r)-1]
}

type addSpace struct {
	*zap.SugaredLogger
}

func AddSpace(logger *zap.Logger) *addSpace {
	return &addSpace{SugaredLogger: logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()}
}

func (l *addSpace) Named(name string) Logger {
	return &addSpace{SugaredLogger: l.SugaredLogger.Named(name)}
}

func (l *addSpace) With(pairs ...interface{}) Logger {
	return &addSpace{SugaredLogger: l.SugaredLogger.With(pairs...)}
}

func (l *addSpace) Debugf(format string, args ...interface{}) {
	l.SugaredLogger.Debugf(format, args...)
}
func (l *addSpace) Infof(format string, args ...interface{}) { l.SugaredLogger.Infof(format, args...) }
func (l *addSpace) Warnf(format string, args ...interface{}) { l.SugaredLogger.Warnf(format, args...) }
func (l *addSpace) Errorf(format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args...)
}
func (l *addSpace) Fatalf(format string, args ...interface{}) {
	l.SugaredLogger.Fatalf(format, args...)
}

func (l *addSpace) Debug(args ...interface{}) { l.SugaredLogger.Debug(join(args, " ")...) }
func (l *addSpace) Info(args ...interface{})  { l.SugaredLogger.Info(join(args, " ")...) }
func (l *addSpace) Warn(args ...interface{})  { l.SugaredLogger.Warn(join(args, " ")...) }
func (l *addSpace) Error(args ...interface{}) { l.SugaredLogger.Error(join(args, " ")...) }
func (l *addSpace) Fatal(args ...interface{}) { l.SugaredLogger.Fatal(join(args, " ")...) }

func (l *addSpace) Debugw(msg string, pairs ...interface{}) { l.SugaredLogger.Debugw(msg, pairs...) }
func (l *addSpace) Infow(msg string, pairs ...interface{})  { l.SugaredLogger.Infow(msg, pairs...) }
func (l *addSpace) Warnw(msg string, pairs ...interface{})  { l.SugaredLogger.Warnw(msg, pairs...) }
func (l *addSpace) Errorw(msg string, pairs ...interface{}) { l.SugaredLogger.Errorw(msg, pairs...) }
func (l *addSpace) Fatalw(msg string, pairs ...interface{}) { l.SugaredLogger.Fatalw(msg, pairs...) }
