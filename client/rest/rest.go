package rest

import (
	rConfig "github.com/slyerr/econfig/client/rest/config"
	"github.com/slyerr/econfig/core/rest/server"
	"github.com/slyerr/verifier"
)

func Start(port int) {
	verifier.I().GreaterNP(port, 0, "port")

	server.Start("econfig client's ", port, []server.Action{
		rConfig.Put,
	})
}
