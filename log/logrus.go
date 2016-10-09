package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/evalphobia/logrus_fluent"
)

func SetLogrusDefault(logLevel, logFormat string) {
	return SetLogrus(logLevel, logFormat, "stdout", "", 0, "")
}

func SetLogrus(logLevel, logFormat, logOut, fluentHost string, fluentPort int, fluentTag string) error {

	if logOut == "fluent" {
		hook, err := logrus_fluent.New(fluentHost, fluentPort)
		if err != nil {
			return err
		}

		// set custom fire level
		hook.SetLevels([]logrus.Level{
			logrus.PanicLevel,
			logrus.ErrorLevel,
		})

		// set static tag
		hook.SetTag(fluentTag + ".tag")

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

	if logFormat == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	return nil
}
