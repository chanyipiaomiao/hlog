package hlog

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"time"
)

const (
	FileNameDateFormat = "%Y%m%d"              // 日志文件名的默认日期格式
	TimestampFormat    = "2006-01-02 15:04:05" // 日志条目中的默认日期时间格式
	Text               = "text"                // 普通文本格式日志
	JSON               = "json"                // json格式日志
	DataKey            = "data"                // json日志条目中 数据字段都会作为该字段的嵌入字段
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

	// 文件名的日期格式
	FileNameDateFormat string

	// 日志中日期时间格式
	TimestampFormat string

	// 日志级别
	LogLevel Level

	// 日志最长保存多久
	MaxAge time.Duration

	// 日志默认多长时间轮转一次
	RotationTime time.Duration

	// 是否打印方法名和行号
	ReportCaller bool

	// json日志是否美化输出
	JSONPrettyPrint bool

	// json日志条目中 数据字段都会作为该字段的嵌入字段
	JSONDataKey string
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
	logrusLogger.SetOutput(ioutil.Discard)
	logrusLogger.SetReportCaller(option.ReportCaller)
	logrusLogger.Level = logrus.Level(option.LogLevel)

	switch option.LogType {
	case JSON:
		format := &logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
			PrettyPrint:     option.JSONPrettyPrint,
		}
		if option.JSONDataKey != "" {
			format.DataKey = option.JSONDataKey
		}
		logrusLogger.Formatter = format
	default:
		logrusLogger.Formatter = &logrus.TextFormatter{
			TimestampFormat: timestampFormat,
		}
	}

	return logrusLogger, nil
}

// New 返回Logger
// 日志类型是: 普通文本日志|JSON日志 全部级别都写入到同一个文件
func New(option *Option) (*Logger, error) {
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

func newRotatelog(option *Option, levelStr string) (*rotatelogs.RotateLogs, error) {
	var (
		err      error
		filename string
		writer   *rotatelogs.RotateLogs
	)

	filename = fmt.Sprintf("%s.%s", option.LogPath, levelStr)
	writer, err = rotatelogs.New(
		fmt.Sprintf("%s.%s", filename, fileNameDateFormat),
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(option.MaxAge),
		rotatelogs.WithRotationTime(option.RotationTime),
	)
	if err != nil {
		return nil, fmt.Errorf("rotatelogs.New error: %s", err)
	}

	return writer, nil
}

// NewSeparate 不同级别的日志输出到不同的文件
func NewSeparate(option *Option) (*Logger, error) {
	var (
		err          error
		logrusLogger *logrus.Logger
		debugWriter  *rotatelogs.RotateLogs
		infoWriter   *rotatelogs.RotateLogs
		warnWriter   *rotatelogs.RotateLogs
		errorWriter  *rotatelogs.RotateLogs
		fatalWriter  *rotatelogs.RotateLogs
		panicWriter  *rotatelogs.RotateLogs
		fileHook     *lfshook.LfsHook
	)

	if logrusLogger, err = newLogger(option); err != nil {
		return nil, err
	}

	if debugWriter, err = newRotatelog(option, "debug"); err != nil {
		return nil, err
	}

	if infoWriter, err = newRotatelog(option, "info"); err != nil {
		return nil, err
	}

	if warnWriter, err = newRotatelog(option, "warn"); err != nil {
		return nil, err
	}

	if errorWriter, err = newRotatelog(option, "error"); err != nil {
		return nil, err
	}

	if fatalWriter, err = newRotatelog(option, "fatal"); err != nil {
		return nil, err
	}

	if panicWriter, err = newRotatelog(option, "panic"); err != nil {
		return nil, err
	}

	fileHook = lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: debugWriter, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  infoWriter,
		logrus.WarnLevel:  warnWriter,
		logrus.ErrorLevel: errorWriter,
		logrus.FatalLevel: fatalWriter,
		logrus.PanicLevel: panicWriter,
	}, logrusLogger.Formatter)

	logrusLogger.Hooks.Add(fileHook)

	logger = &Logger{
		logrus: logrusLogger,
	}

	return logger, nil
}

func (l *Logger) Debug(dataFields D, message string) {
	l.logrus.WithFields(logrus.Fields(dataFields)).Debug(message)
}

func (l *Logger) Info(dataFields D, message string) {
	l.logrus.WithFields(logrus.Fields(dataFields)).Info(message)
}

func (l *Logger) Warn(dataFields D, message string) {
	l.logrus.WithFields(logrus.Fields(dataFields)).Warn(message)
}

func (l *Logger) Error(dataFields D, message string) {
	l.logrus.WithFields(logrus.Fields(dataFields)).Error(message)
}

func (l *Logger) Fatal(dataFields D, message string) {
	l.logrus.WithFields(logrus.Fields(dataFields)).Fatal(message)
}

func (l *Logger) Panic(dataFields D, message string) {
	l.logrus.WithFields(logrus.Fields(dataFields)).Panic(message)
}

func Debug(dataFields D, message string) {
	if logger.logrus == nil {
		logrus.Debug(dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Debug(message)
}

func Info(dataFields D, message string) {
	if logger.logrus == nil {
		logrus.Info(dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Info(message)
}

func Warn(dataFields D, message string) {
	if logger.logrus == nil {
		logrus.Warn(dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Warn(message)
}

func Error(dataFields D, message string) {
	if logger.logrus == nil {
		logrus.Error(dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Error(message)
}

func Fatal(dataFields D, message string) {
	if logger.logrus == nil {
		logrus.Fatal(dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Fatal(message)
}

func Panic(dataFields D, message string) {
	if logger.logrus == nil {
		logrus.Panic(dataFields, message)
		return
	}
	logger.logrus.WithFields(logrus.Fields(dataFields)).Panic(message)
}

func StderrFatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
