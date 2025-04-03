package logger

import "testing"

func TestLogger(t *testing.T) {
	s := &Settings{
		Path:       "../log",
		Name:       "redis",
		Ext:        "log",
		TimeFormat: "2006/01/02 15:04:05",
	}

	Setup(s)
	Debug("hello debug")
	Info("hello info")
	Warn("hello warn")
	Error("hello error")
}
