package logger

import "context"

// 日志模块
type Logger interface {
	Debug(ctx context.Context, message string, args ...any)
	Info(ctx context.Context, message string, args ...any)
	Warn(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, message string, args ...any)
	Fatal(ctx context.Context, message string, args ...any)
	Panic(ctx context.Context, message string, args ...any)

	Debugf(ctx context.Context, message string, args ...any)
	Infof(ctx context.Context, message string, args ...any)
	Warnf(ctx context.Context, message string, args ...any)
	Errorf(ctx context.Context, message string, args ...any)
	Fatalf(ctx context.Context, message string, args ...any)
	Panicf(ctx context.Context, message string, args ...any)
}

var Log Logger
