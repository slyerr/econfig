package server

import (
	"github.com/slyerr/econfig/core/zaps"
	"github.com/slyerr/econfig/server/mq"
	"github.com/slyerr/econfig/server/rest"
	"github.com/slyerr/econfig/server/store"
	"go.uber.org/zap"
)

func Start(logger zaps.LoggerConfig, port int) {
	zaps.Start(logger)

	store.Open()
	zap.S().Warn("econfig server's store db has been opened")

	mq.ConfigPush().Start()
	zap.S().Warn("econfig server's config push mq has been started")

	rest.Start(port)
}

func Stop() {
	store.Close()
	zaps.Stop()
}
