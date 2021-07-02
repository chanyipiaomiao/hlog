package hlog

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	logger, err := New(&Option{
		LogPath:                "./logs/hlog.log",
		LogType:                JSON,
		LogLevel:               DebugLevel,
		MaxAge:                 7 * 24 * time.Hour,
		RotationTime:           24 * time.Hour,
		JSONPrettyPrint:        true,
		IsEnableRecordFileInfo: true,
		FileInfoField:          "caller",
	})

	if err != nil {
		t.Error(err)
		return
	}

	logger.Debug(nil, "hello: %s", "world")
	logger.Info(D{"hello": "world"}, "hello: %s", "world")
	logger.Warn(D{"username": "warn"}, "hello: %s", "world")
	logger.Error(D{"username": "Error"}, "hello: %s", "world")
	//logger.Panic(D{"username": "Panic"}, "呵呵")
	//logger.Fatal(D{"username": "Fatal"}, "呵呵")

}

func TestNewSeparate(t *testing.T) {
	logger, err := NewSeparate(&Option{
		LogPath:                "./logs/hlog.log",
		LogType:                Text,
		LogLevel:               DebugLevel,
		MaxAge:                 7 * 24 * time.Hour,
		RotationTime:           24 * time.Hour,
		IsEnableRecordFileInfo: true,
		FileInfoField:          "called",
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

func TestPrintf(t *testing.T) {
	Print("hello %s %s", "world", "hehe")
}
