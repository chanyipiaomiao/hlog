package hlog

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	logger, err := New(&Option{
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

	logger.Debug(D{"hello": "world"}, "hello")
	logger.Info(D{"hello": "world"}, "hello")
	logger.Warn(D{"username": "warn"}, "呵呵")
	logger.Error(D{"username": "Error"}, "呵呵")
	//logger.Panic(D{"username": "Panic"}, "呵呵")
	//logger.Fatal(D{"username": "Fatal"}, "呵呵")

}

func TestNewSeparate(t *testing.T) {
	logger, err := NewSeparate(&Option{
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

	logger.Debug(D{"hello": "world"}, "hello")
	logger.Info(D{"hello": "world"}, "hello")
	logger.Warn(D{"username": "warn"}, "呵呵")
	logger.Error(D{"username": "Error"}, "呵呵")
	//logger.Panic(D{"username": "Panic"}, "呵呵")
	//logger.Fatal(D{"username": "Fatal"}, "呵呵")

}
