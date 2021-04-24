package client

import (
	"time"

	"github.com/pkg/errors"
	"github.com/slyerr/econfig/client/config"
	"github.com/slyerr/econfig/client/job"
	clientRest "github.com/slyerr/econfig/client/rest"
	coreRest "github.com/slyerr/econfig/core/rest"
	"github.com/slyerr/econfig/core/rest/client"
	"github.com/slyerr/econfig/core/utils"
	"github.com/slyerr/econfig/core/zaps"
	"github.com/slyerr/verifier"
	"go.uber.org/zap"
)

type Config struct {
	Logger         zaps.LoggerConfig `json:"logger"`
	Port           int               `json:"port"`
	ConfigKey      string            `json:"configKey"`
	ServerHost     string            `json:"serverHost"`
	ConfigSyncTick int64             `json:"configSyncTick"`
}

func Start(c Config) {
	zaps.Start(c.Logger)

	sync(c.ConfigKey, c.ServerHost)
	zap.S().Warn("econfig client's config has been synchronized")

	job.ConfigSync().Start(time.Duration(c.ConfigSyncTick), c.ConfigKey, c.ServerHost)
	zap.S().Warn("econfig client's config sync job has been started")

	go func() {
		pushClient(c.Port, c.ConfigKey, c.ServerHost)
		zap.S().Warn("econfig client's has been pushed")
	}()

	clientRest.Start(c.Port)
}

func Stop() {
	zaps.Stop()
}

func sync(configKey string, serverHost string) {
	configKey, err := utils.CheckConfigKey(configKey)
	if err != nil {
		panic(err)
	}

	verifier.S().NotBlankNP(serverHost, "server host")

	err = config.Sync(configKey, serverHost)
	if err != nil {
		panic(errors.WithMessage(err, "econfig client's config sync error"))
	}
}

func pushClient(port int, configKey string, serverHost string) {
	verifier.I().GreaterNP(port, 0, "port")

	configKey, err := utils.CheckConfigKey(configKey)
	if err != nil {
		panic(err)
	}

	verifier.S().NotBlankNP(serverHost, "server host")

	type HostBody struct {
		Port    int    `json:"port"`
		PushUrl string `json:"pushUrl"`
	}

	_, err = client.NewRestClient().Put(utils.CompleUrl(serverHost, coreRest.ServerHostUrlV1+configKey), HostBody{Port: port, PushUrl: coreRest.ClientConfigUrlV1})
	if err != nil {
		panic(errors.WithMessage(err, "econfig client push to server error"))
	}
}
