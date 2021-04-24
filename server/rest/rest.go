package rest

import (
	"github.com/slyerr/econfig/core/rest/server"
	"github.com/slyerr/econfig/server/rest/config"
	"github.com/slyerr/econfig/server/rest/host"
	"github.com/slyerr/verifier"
)

func Start(port int) {
	verifier.I().GreaterNP(port, 0, "port")

	server.Start("econfig server's ", port, []server.Action{
		config.Get,
		config.Put,
		config.Delete,
		host.Put,
		host.Delete,
	})
}
