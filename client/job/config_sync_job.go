package job

import (
	"fmt"
	"time"

	"github.com/slyerr/econfig/client/config"
	"github.com/slyerr/econfig/core/rest/client"
	"github.com/slyerr/econfig/core/utils"
	"github.com/slyerr/verifier"
	"go.uber.org/zap"
)

type ConfigSyncJob struct {
	ticker     *time.Ticker
	rc         *client.RestClient
	configKey  string
	serverHost string
}

func NewConfigSyncJob() *ConfigSyncJob {
	return &ConfigSyncJob{
		rc: client.NewRestClient(),
	}
}

func (j *ConfigSyncJob) Start(tick time.Duration, configKey string, serverHost string) {
	if j.ticker != nil {
		return
	}

	verifier.I64().GreaterNP(int64(tick), 0, "tick")

	configKey, err := utils.CheckConfigKey(configKey)
	if err != nil {
		panic(err)
	}

	verifier.S().NotBlankNP(serverHost, "server host")

	j.ticker = time.NewTicker(tick)
	j.configKey = configKey
	j.serverHost = serverHost

	go func() {
		for range j.ticker.C {
			go j.do()
		}
	}()
}

func (j *ConfigSyncJob) do() {
	zap.S().Debugf("econfig client's config sync time: %+v", time.Now())

	err := config.Sync(j.configKey, j.serverHost)
	if err != nil {
		zap.S().Errorf("econfig client's config sync error: %+v", err)
	}

	fmt.Println("new config: " + config.Get())
}
