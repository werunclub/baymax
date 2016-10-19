package log

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/evalphobia/logrus_fluent"
)

// 设置 logrus, 支持同步日志到 fluent
func SetLogrus(logLevel, logFormat, logOut string, fluentdEnable bool,
	fluentdHost string, fluentdPort int, fluentdTag string) error {

	if fluentdEnable {
		hook, err := logrus_fluent.New(fluentdHost, fluentdPort)
		if err != nil {
			return err
		}

		// set custom fire level
		hook.SetLevels([]logrus.Level{
			logrus.PanicLevel,
			logrus.ErrorLevel,
			logrus.DebugLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.FatalLevel,
		})

		// set static tag
		hook.SetTag(fluentdTag + ".fluentd")

		// ignore field
		hook.AddIgnore("context")

		// filter func
		hook.AddFilter("error", logrus_fluent.FilterError)

		logrus.AddHook(hook)
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
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
