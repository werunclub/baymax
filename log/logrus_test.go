package log

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetLogurs(t *testing.T) {
	SetLogrus("info", "json", "stdout", false, "", 0, "")

	logrus.Infof("info log")
	logrus.Errorf("error log")
}
