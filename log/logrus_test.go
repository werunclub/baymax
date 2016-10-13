package log

import (
	log "github.com/Sirupsen/logrus"
	"testing"
)

func TestSetLogurs(t *testing.T) {

	SetLogrus("info", "json", "stdout", false, "", 0, "")

	log.Infof("info log")
	log.Errorf("error log")
}
