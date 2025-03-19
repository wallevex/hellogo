package log

type Logger interface {
	Named(name string) Logger
	With(pairs ...interface{}) Logger

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugw(msg string, pairs ...interface{})
	Infow(msg string, pairs ...interface{})
	Warnw(msg string, pairs ...interface{})
	Errorw(msg string, pairs ...interface{})
	Fatalw(msg string, pairs ...interface{})
}
