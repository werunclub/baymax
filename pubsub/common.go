package pubsub

import (
	"os"
)

func getBrokerName() string {
	return os.Getenv("PUBSUB_BRORKER")
}
