package log

import (
	"os"

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
		})

		// set static tag
		hook.SetTag("fluentd." + fluentdTag)

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
			TimestampFormat: "2006/01/02 - 15:04:05",
		}
	} else {
		formatter = &logrus.TextFormatter{
			TimestampFormat: "2006/01/02 - 15:04:05",
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
