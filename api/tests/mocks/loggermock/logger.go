package loggermock

import (
	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Info(msg string, keysAndValues ...any) {
	l.Called(msg, keysAndValues)
}

func (l *LoggerMock) Warn(msg string, keysAndValues ...any) {
	l.Called(msg, keysAndValues)
}

func (l *LoggerMock) Error(msg string, keysAndValues ...any) {
	l.Called(msg, keysAndValues)
}

func (l *LoggerMock) Fatal(msg string, keysAndValues ...any) {
	l.Called(msg, keysAndValues)
}
