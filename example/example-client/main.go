package main

import (
	"time"

	"github.com/slyerr/econfig/client"
	"github.com/slyerr/econfig/core/zaps"
	"go.uber.org/zap"
)

func main() {
	defer client.Stop()

	client.Start(client.Config{
		Logger: zaps.LoggerConfig{
			Level: zap.DebugLevel,
		},
		Port:           8077,
		ConfigKey:      "test",
		ServerHost:     "127.0.0.1:8088",
		ConfigSyncTick: int64(30 * time.Minute),
	})
}
