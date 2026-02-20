package logger

import "go.uber.org/zap"

type Logger interface {
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)
}

type zapAdapter struct {
	sugar *zap.SugaredLogger
}

func New() Logger {
	z, _ := zap.NewProduction()

	return &zapAdapter{
		sugar: z.Sugar(),
	}
}

func (l *zapAdapter) Info(msg string, keysAndValues ...any) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *zapAdapter) Error(msg string, keysAndValues ...any) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *zapAdapter) Fatal(msg string, keysAndValues ...any) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

func (l *zapAdapter) Warn(msg string, keysAndValues ...any) {
	l.sugar.Warnw(msg, keysAndValues...)
}
