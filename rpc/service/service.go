package service

import (
	"github.com/sirupsen/logrus"

	"github.com/werunclub/baymax/v2/rpc/client"
)

type Service struct {
	Cli      *client.Client
	Log      *logrus.Entry
	HandleID string
}

func NewService(client *client.Client, handleID string) Service {
	logger := logrus.WithField("HandleID", handleID)
	return Service{Cli: client, HandleID: handleID, Log: logger}
}
