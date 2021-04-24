package main

import (
	"github.com/slyerr/econfig/core/zaps"
	"github.com/slyerr/econfig/server"
	"go.uber.org/zap"
)

func main() {
	defer server.Stop()

	server.Start(zaps.LoggerConfig{Level: zap.DebugLevel}, 8088)
}
