## hlog

基于logrus封装的日志库

## 特性

- 支持全局导入使用或是对象使用
- 支持json和普通的文本日志格式
- 支持各个级别的日志输出到同一个文件或者不同文件
- json日志支持美化输出
- json日志支持数据key
- 支持按照时间来轮转日志
- 支持打印方法名和行数

## 安装

`go get github.com/chanyipiaomiao/hlog`


## 例子

### 输出到单个文件

#### 使用对象的方式
```go
package main

import (
	"fmt"
	"github.com/chanyipiaomiao/hlog"
	"time"
)

func main() {
	logger, err := hlog.New(&hlog.Option{
		LogPath:            "/tmp/logs/hlog.log",
		LogType:            hlog.JSON,
		FileNameDateFormat: hlog.FileNameDateFormat,
		TimestampFormat:    hlog.TimestampFormat,
		LogLevel:           hlog.DebugLevel,
		MaxAge:             7 * 24 * time.Hour,
		RotationTime:       24 * time.Hour,
		JSONPrettyPrint: true,
		JSONDataKey: "data",
		ReportCaller: false,
	})

	if err != nil {
		hlog.StderrFatalf("error: %s", err)
	}

	logger.Debug(hlog.D{"hello": "world"}, "hello")
	logger.Info(hlog.D{"hello": "world"}, "hello")
	logger.Warn(hlog.D{"username": "warn"}, "呵呵")
	logger.Error(hlog.D{"username": "Error"}, "呵呵")

}

```

输出结果

```log
{
  "data": {
    "hello": "world"
  },
  "file": "/Users/xxx/workspace/hlog/hlog.go:257",
  "func": "github.com/chanyipiaomiao/hlog.(*Logger).Debug",
  "level": "debug",
  "msg": "hello",
  "time": "2020-05-21 17:53:47"
}
{
  "data": {
    "hello": "world"
  },
  "file": "/Users/xxx/workspace/hlog/hlog.go:261",
  "func": "github.com/chanyipiaomiao/hlog.(*Logger).Info",
  "level": "info",
  "msg": "hello",
  "time": "2020-05-21 17:53:47"
}
{
  "data": {
    "username": "warn"
  },
  "file": "/Users/xxx/workspace/hlog/hlog.go:265",
  "func": "github.com/chanyipiaomiao/hlog.(*Logger).Warn",
  "level": "warning",
  "msg": "呵呵",
  "time": "2020-05-21 17:53:47"
}
{
  "data": {
    "username": "Error"
  },
  "file": "/Users/xxx/workspace/hlog/hlog.go:269",
  "func": "github.com/chanyipiaomiao/hlog.(*Logger).Error",
  "level": "error",
  "msg": "呵呵",
  "time": "2020-05-21 17:53:47"
}
```

#### 使用全局导入的方式

```go
package main

import (
	"github.com/chanyipiaomiao/hlog"
	"time"
)

func main() {

	_, err := hlog.New(&hlog.Option{
		LogPath:            "/tmp/logs/hlog.log",
		LogType:            hlog.JSON,
		FileNameDateFormat: hlog.FileNameDateFormat,
		TimestampFormat:    hlog.TimestampFormat,
		LogLevel:           hlog.DebugLevel,
		MaxAge:             7 * 24 * time.Hour,
		RotationTime:       24 * time.Hour,
		JSONPrettyPrint: true,
		JSONDataKey: "data",
		ReportCaller: true,
	})

	if err != nil {
		hlog.StderrFatalf("error: %s", err)
	}

	hlog.Debug(hlog.D{"hello": "world"}, "hello")
	hlog.Info(hlog.D{"hello": "world"}, "hello")
	hlog.Warn(hlog.D{"username": "warn"}, "呵呵")
	hlog.Error(hlog.D{"username": "Error"}, "呵呵")
	//hlog.Panic(hlog.D{"username": "Panic"}, "呵呵")
	//hlog.Fatal(hlog.D{"username": "Fatal"}, "呵呵")

}

```

### 按照级别输出不同的文件

只需要把hlog.New 换成 hlog.NewSeparate 即可

```go
_, err := hlog.NewSeparate(&hlog.Option{
    LogPath:            "/tmp/logs/hlog.log",
    LogType:            hlog.JSON,
    FileNameDateFormat: hlog.FileNameDateFormat,
    TimestampFormat:    hlog.TimestampFormat,
    LogLevel:           hlog.DebugLevel,
    MaxAge:             7 * 24 * time.Hour,
    RotationTime:       24 * time.Hour,
    JSONPrettyPrint: true,
    JSONDataKey: "data",
    ReportCaller: true,
})
```
