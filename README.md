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
		JSONDataKey: "",
		IsEnableRecordFileInfo: true,
		FileInfoField: "call",
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
如果设置了 JSONDataKey: "data", 那么所有的字段都会是data字段的内嵌字段,除了默认的字段msg、time、level 
这是logrus的默认行为

```log
{
  "data": {
    "call": "main.go:29",
    "hello": "world"
  },
  "level": "debug",
  "msg": "hello",
  "time": "2020-05-22 18:16:17"
}
```

输出结果

json格式
```log
{
  "call": "main.go:29",
  "hello": "world",
  "level": "debug",
  "msg": "hello",
  "time": "2020-05-22 18:11:39"
}
{
  "call": "main.go:30",
  "hello": "world",
  "level": "info",
  "msg": "hello",
  "time": "2020-05-22 18:11:39"
}
{
  "call": "main.go:31",
  "level": "warning",
  "msg": "呵呵",
  "time": "2020-05-22 18:11:39",
  "username": "warn"
}
{
  "call": "main.go:32",
  "level": "error",
  "msg": "呵呵",
  "time": "2020-05-22 18:11:39",
  "username": "Error"
}
{
  "age": 18,
  "call": "hello.go:8",
  "level": "info",
  "msg": "修改年龄",
  "time": "2020-05-22 18:11:39"
}

```
文本格式
```log
time="2020-05-22 18:12:14" level=debug msg=hello call="main.go:29" hello=world
time="2020-05-22 18:12:14" level=info msg=hello call="main.go:30" hello=world
time="2020-05-22 18:12:14" level=warning msg="呵呵" call="main.go:31" username=warn
time="2020-05-22 18:12:14" level=error msg="呵呵" call="main.go:32" username=Error
time="2020-05-22 18:12:14" level=info msg="修改年龄" age=18 call="hello.go:8"
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

在其他包中导入即可使用
```go
import "github.com/chanyipiaomiao/hlog"
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
})
```
