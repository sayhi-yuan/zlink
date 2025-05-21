package logger

import "context"

// 日志模块
type Logger interface {
	Debug(ctx context.Context, message string)
	Info(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
	Error(ctx context.Context, message string)
	Fatal(ctx context.Context, message string)
	Panic(ctx context.Context, message string)

	Debugf(ctx context.Context, message string, args ...any)
	Infof(ctx context.Context, message string, args ...any)
	Warnf(ctx context.Context, message string, args ...any)
	Errorf(ctx context.Context, message string, args ...any)
	Fatalf(ctx context.Context, message string, args ...any)
	Panicf(ctx context.Context, message string, args ...any)
}

var Log Logger

func Debug(ctx context.Context, message string) {
	if Log == nil {
		return
	}

	Log.Debug(ctx, message)
}

func Info(ctx context.Context, message string) {
	if Log == nil {
		return
	}

	Log.Info(ctx, message)
}

func Warn(ctx context.Context, message string) {
	if Log == nil {
		return
	}

	Log.Warn(ctx, message)
}

func Error(ctx context.Context, message string) {
	if Log == nil {
		return
	}

	Log.Error(ctx, message)
}

func Fatal(ctx context.Context, message string) {
	if Log == nil {
		return
	}

	Log.Fatal(ctx, message)
}

func Panic(ctx context.Context, message string) {
	if Log == nil {
		return
	}

	Log.Panic(ctx, message)
}

func Debugf(ctx context.Context, message string, args ...any) {
	if Log == nil {
		return
	}

	Log.Debugf(ctx, message, args...)
}

func Infof(ctx context.Context, message string, args ...any) {
	if Log == nil {
		return
	}

	Log.Infof(ctx, message, args...)
}

func Warnf(ctx context.Context, message string, args ...any) {
	if Log == nil {
		return
	}

	Log.Warnf(ctx, message, args...)
}

func Errorf(ctx context.Context, message string, args ...any) {
	if Log == nil {
		return
	}

	Log.Errorf(ctx, message, args...)
}

func Fatalf(ctx context.Context, message string, args ...any) {
	if Log == nil {
		return
	}

	Log.Fatalf(ctx, message, args...)
}

func Panicf(ctx context.Context, message string, args ...any) {
	if Log == nil {
		return
	}

	Log.Panicf(ctx, message, args...)
}
