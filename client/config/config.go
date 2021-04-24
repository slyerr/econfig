package config

import (
	"encoding/json"

	"github.com/slyerr/econfig/core/rest"
	"github.com/slyerr/econfig/core/rest/client"
	"github.com/slyerr/econfig/core/utils"
	"github.com/slyerr/verifier"
)

var rc *client.RestClient = client.NewRestClient()

var value string

func Get() string {
	return utils.CleanConfigValue(value)
}

func GetX(i interface{}) error {
	return json.Unmarshal([]byte(Get()), i)
}

func Set(v string) {
	value = utils.CleanConfigValue(v)
}

func SetX(x interface{}) error {
	b, err := json.Marshal(x)
	if err != nil {
		return err
	}

	value = string(b)
	return nil
}

func Sync(configKey string, serverHost string) error {
	configKey, err := utils.CheckConfigKey(configKey)
	if err != nil {
		return err
	}

	if err = verifier.S().NotBlankN(serverHost, "server host"); err != nil {
		return err
	}

	result := &rest.Result{}
	if err := rc.Get2(utils.CompleUrl(serverHost, rest.ServerConfigUrlV1+configKey), nil, result); err != nil {
		return err
	}

	switch d := result.Data.(type) {
	case string:
		Set(d)
	case *string:
		Set(*d)
	default:
		SetX(result.Data)
	}

	return nil
}
