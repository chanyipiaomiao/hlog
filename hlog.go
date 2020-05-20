package hlog

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

const (
	FileNameDateFormat = "%Y%m%d"              // 日志文件名的默认日期格式
	TimestampFormat    = "2006-01-02 15:04:05" // 日志条目中的默认日期时间格式
	Text               = "text"                // 普通文本格式日志
	JSON               = "json"                // json格式日志
)

var (
	logger             *Logger
	fileNameDateFormat string // 日志文件名的日期格式
	timestampFormat    string // 日志条目中的日期时间格式
)

// Level type
type Level uint32

// 要写入日志的数据字段
type D map[string]interface{}

const (
	// 日志级别
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

type Option struct {
	// log 路径
	LogPath string

	// 日志类型 json|text
	LogType string

	// 文件名的日期格式 默认: %Y-%m-%d|%Y%m%d
	FileNameDateFormat string

	// 日志中日期时间格式 默认: 2006-01-02 15:04:05
	TimestampFormat string

	// 是否分离不同级别的日志 默认: true
	IsSeparateLevelLog bool

	// 日志级别 默认: log.InfoLevel
	LogLevel Level

	// 日志最长保存多久 默认: 15天
	MaxAge time.Duration

	// 日志默认多长时间轮转一次 默认: 24小时
	RotationTime time.Duration
}

type Logger struct {
	logrus *logrus.Logger
}

func GetLogger() *Logger {
	return logger
}

func newLogger(option *Option) (*logrus.Logger, error) {
	var (
		err          error
		logrusLogger *logrus.Logger
	)

	if err = makeDirAll(option.LogPath); err != nil {
		return nil, err
	}

	if !path.IsAbs(option.LogPath) {
		return nil, fmt.Errorf("LogPath please use absolute path: %s", option.LogPath)
	}

	if option.FileNameDateFormat == "" {
		fileNameDateFormat = FileNameDateFormat
	} else {
		fileNameDateFormat = option.FileNameDateFormat
	}

	if option.TimestampFormat == "" {
		timestampFormat = TimestampFormat
	} else {
		timestampFormat = option.TimestampFormat
	}

	logrusLogger = logrus.New()
	logrusLogger.Level = logrus.Level(option.LogLevel)

	if option.LogType == "json" {
		logrusLogger.Formatter = &logrus.JSONFormatter{TimestampFormat: timestampFormat}
	}

	switch option.LogType {
	case JSON:
		logrusLogger.Formatter = &logrus.JSONFormatter{TimestampFormat: timestampFormat}
	default:
		logrusLogger.Formatter = &logrus.TextFormatter{TimestampFormat: timestampFormat}
	}

	return logrusLogger, nil
}

// NewDefault 返回Logger
// 日志类型是: 普通文本日志|JSON日志 全部级别都写入到同一个文件
func NewDefault(option *Option) (*Logger, error) {
	var (
		err          error
		logrusLogger *logrus.Logger
		writer       *rotatelogs.RotateLogs
		fileHook     *lfshook.LfsHook
	)

	if logrusLogger, err = newLogger(option); err != nil {
		return nil, err
	}

	writer, err = rotatelogs.New(
		fmt.Sprintf("%s-%s", option.LogPath, fileNameDateFormat),
		rotatelogs.WithMaxAge(option.MaxAge),
		rotatelogs.WithRotationTime(option.RotationTime),
		rotatelogs.WithLinkName(option.LogPath),
	)
	if err != nil {
		return nil, err
	}

	fileHook = lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, logrusLogger.Formatter)

	logrusLogger.Hooks.Add(fileHook)

	logger = &Logger{
		logrus: logrusLogger,
	}
	return logger, nil
}

func Debug(dataFields D, message string) {
	if logger.logrus == nil {
		fmt.Printf("%v %s\n", dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Debug(message)
}

func Info(dataFields D, message string) {
	if logger.logrus == nil {
		fmt.Printf("%v %s\n", dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Info(message)
}

func Warn(dataFields D, message string) {
	if logger.logrus == nil {
		fmt.Printf("%v %s\n", dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Warn(message)
}

func Error(dataFields D, message string) {
	if logger.logrus == nil {
		fmt.Printf("%v %s\n", dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Error(message)
}

func Fatal(dataFields D, message string) {
	if logger.logrus == nil {
		fmt.Printf("%v %s\n", dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Fatal(message)
}

func Panic(dataFields D, message string) {
	if logger.logrus == nil {
		fmt.Printf("%v %s\n", dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Panic(message)
}
