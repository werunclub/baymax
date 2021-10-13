package log

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestSetLogurs(t *testing.T) {
	SetLogrus("info", "json", "stdout", false, "", 0, "")

	log.Infof("info log")
	log.Errorf("error log")
}
