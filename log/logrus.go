package log

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/evalphobia/logrus_fluent"
)

// 设置 logrus, 支持同步日志到 fluent
func SetLogrus(logLevel, logFormat, logOut string, fluentdEnable bool,
	fluentdHost string, fluentdPort int, fluentdTag string) error {

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	if fluentdEnable {
		hook, err := logrus_fluent.NewWithConfig(logrus_fluent.Config{
			Host:                fluentdHost,
			Port:                fluentdPort,
			DefaultMessageField: logrus_fluent.MessageField,
			MarshalAsJSON:       true,
		})
		if err != nil {
			return err
		}

		// set custom fire level
		if level == logrus.DebugLevel {
			hook.SetLevels([]logrus.Level{
				logrus.PanicLevel,
				logrus.ErrorLevel,
				logrus.DebugLevel,
				logrus.WarnLevel,
				logrus.InfoLevel,
				logrus.FatalLevel,
			})
		} else {
			hook.SetLevels([]logrus.Level{
				logrus.InfoLevel,
				logrus.PanicLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
				logrus.FatalLevel,
			})
		}

		// set static tag
		hook.SetTag("go." + fluentdTag)

		// ignore field
		hook.AddIgnore("context")

		// filter func
		hook.AddFilter("error", logrus_fluent.FilterError)

		logrus.AddHook(hook)
	}

	logrus.SetLevel(level)

	if err := SetLogOut(logOut); err != nil {
		return err
	}

	// 设置格式
	SetLogFormatter(logFormat)

	return nil
}

func SetLogFormatter(formatString string) {
	var formatter logrus.Formatter

	if formatString == "json" {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		}
	} else {
		formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			ForceColors:     true,
			FullTimestamp:   true,
		}
	}
	logrus.SetFormatter(formatter)
}

// SetLogOut provide log stdout and stderr output
func SetLogOut(outString string) error {
	switch outString {
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "stderr":
		logrus.SetOutput(os.Stderr)
	default:
		f, err := os.OpenFile(outString, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			return err
		}

		logrus.SetOutput(f)
	}

	return nil
}

// 带 source 的 logrus
func SourcedLogrus() *logrus.Entry {
	log := logrus.StandardLogger()
	return SourceLogrus(logrus.NewEntry(log))
}

// 为 logrus 添加 source
func SourceLogrus(entry *logrus.Entry) *logrus.Entry {
	return entry.WithField("source", fileWithLineNum())
}

// 获取调用位置
func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!regexp.MustCompile(`baymax/log/.*.go`).MatchString(file) || regexp.MustCompile(`baymax/log/.*test.go`).MatchString(file)) {
			return fmt.Sprintf("%v:%v", file, line)
		}
	}
	return ""
}

// Logger is the interface for loggers used in the Prometheus components.
type Logger interface {
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})

	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})

	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})

	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})

	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	With(key string, value interface{}) Logger
}

type logger struct {
	entry *logrus.Entry
}

func (l logger) With(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

// Debug logs a message at level Debug on the standard logger.
func (l logger) Debug(args ...interface{}) {
	l.sourced().Debug(args...)
}

// Debug logs a message at level Debug on the standard logger.
func (l logger) Debugln(args ...interface{}) {
	l.sourced().Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (l logger) Debugf(format string, args ...interface{}) {
	l.sourced().Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Info(args ...interface{}) {
	l.sourced().Info(args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Infoln(args ...interface{}) {
	l.sourced().Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func (l logger) Infof(format string, args ...interface{}) {
	l.sourced().Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warn(args ...interface{}) {
	l.sourced().Warn(args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warnln(args ...interface{}) {
	l.sourced().Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func (l logger) Warnf(format string, args ...interface{}) {
	l.sourced().Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Error(args ...interface{}) {
	l.sourced().Error(args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Errorln(args ...interface{}) {
	l.sourced().Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (l logger) Errorf(format string, args ...interface{}) {
	l.sourced().Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatal(args ...interface{}) {
	l.sourced().Fatal(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatalln(args ...interface{}) {
	l.sourced().Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (l logger) Fatalf(format string, args ...interface{}) {
	l.sourced().Fatalf(format, args...)
}

// sourced adds a source field to the logger that contains
// the file name and line where the logging happened.
func (l logger) sourced() *logrus.Entry {
	return l.entry.WithField("source", fileWithLineNum())
}

var origLogger = logrus.StandardLogger()
var baseLogger = logger{entry: logrus.NewEntry(origLogger)}

// With adds a field to the logger.
func With(key string, value interface{}) Logger {
	return baseLogger.With(key, value)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	baseLogger.sourced().Debug(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	baseLogger.sourced().Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	baseLogger.sourced().Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	baseLogger.sourced().Info(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	baseLogger.sourced().Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	baseLogger.sourced().Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	baseLogger.sourced().Warn(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	baseLogger.sourced().Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	baseLogger.sourced().Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	baseLogger.sourced().Error(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	baseLogger.sourced().Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	baseLogger.sourced().Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	baseLogger.sourced().Fatal(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	baseLogger.sourced().Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	baseLogger.sourced().Fatalf(format, args...)
}

type errorLogWriter struct{}

func (errorLogWriter) Write(b []byte) (int, error) {
	baseLogger.sourced().Error(string(b))
	return len(b), nil
}
