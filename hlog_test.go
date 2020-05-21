package hlog

import (
	"testing"
	"time"
)

func TestNewDefault(t *testing.T) {
	_, err := New(&Option{
		LogPath:      "/tmp/logs/hlog.log",
		LogType:      JSON,
		LogLevel:     DebugLevel,
		MaxAge:       7 * 24 * time.Hour,
		RotationTime: 24 * time.Hour,
		ReportCaller: true,
	})

	if err != nil {
		t.Error(err)
		return
	}

	Debug(D{"hello": "world"}, "hello")
	Info(D{"hello": "world"}, "hello")
	Warn(D{"username": "warn"}, "呵呵")
	Error(D{"username": "Error"}, "呵呵")
	//Panic(D{"username": "Panic"}, "呵呵")
	//Fatal(D{"username": "Fatal"}, "呵呵")

}

func TestNewSeparate(t *testing.T) {
	_, err := NewSeparate(&Option{
		LogPath:      "/tmp/logs/hlog.log",
		LogType:      Text,
		LogLevel:     DebugLevel,
		MaxAge:       7 * 24 * time.Hour,
		RotationTime: 24 * time.Hour,
	})

	if err != nil {
		t.Error(err)
		return
	}

	Debug(D{"hello": "world"}, "hello")
	Info(D{"hello": "world"}, "hello")
	Warn(D{"username": "warn"}, "呵呵")
	Error(D{"username": "Error"}, "呵呵")
	//Panic(D{"username": "Panic"}, "呵呵")
	//Fatal(D{"username": "Fatal"}, "呵呵")

}
