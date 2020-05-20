package hlog

import (
	"testing"
	"time"
)

func TestNewDefault(t *testing.T) {
	_, err := NewDefault(&Option{
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
	Info(D{"username": "张三"}, "添加成功")

}
