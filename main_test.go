package logger

import (
	"testing"
)

func TestLoggerDebug(t *testing.T) {
	log := New("robokiller-ivr", "1.0")
	log.Set("foo", "bar")
	log.Debug("debug message")
}

func TestLoggerInfo(t *testing.T) {
	log := New("robokiller-ivr", "1.0")
	log.Set("foo", "bar")
	log.Info("info message")
}

func TestLoggerError(t *testing.T) {
	log := New("robokiller-ivr", "1.0")
	log.Set("foo", "bar")
	log.Error("error message")
}
