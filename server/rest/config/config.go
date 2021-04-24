package config

import (
	"net/http"

	"github.com/slyerr/econfig/core/rest"
	"github.com/slyerr/econfig/core/rest/server"
	"github.com/slyerr/econfig/server/mq"
	"github.com/slyerr/econfig/server/store/config"
)

func Get() (rest.Method, string, func(req *http.Request, params server.Params) (interface{}, error)) {
	return rest.MethodGet, rest.ServerConfigUrlV1 + ":key", func(req *http.Request, params server.Params) (interface{}, error) {
		if value, err := config.Store().Get(params.GetPathParam("key")); err != nil {
			return nil, err
		} else {
			return value, nil
		}
	}
}

func Put() (rest.Method, string, func(req *http.Request, params server.Params) (interface{}, error)) {
	return rest.MethodPut, rest.ServerConfigUrlV1 + ":key", func(req *http.Request, params server.Params) (interface{}, error) {
		value, err := params.GetBodyString()
		if err != nil {
			return nil, err
		}

		if k, v, err := config.Store().Put(params.GetPathParam("key"), value); err != nil {
			return nil, err
		} else {
			mq.ConfigPush().Produce(mq.ConfigPushMsg{Key: k, Value: v})
			return nil, nil
		}
	}
}

func Delete() (rest.Method, string, func(req *http.Request, params server.Params) (interface{}, error)) {
	return rest.MethodDelete, rest.ServerConfigUrlV1 + ":key", func(req *http.Request, params server.Params) (interface{}, error) {
		if k, err := config.Store().Delete(params.GetPathParam("key")); err != nil {
			return nil, err
		} else {
			mq.ConfigPush().Produce(mq.ConfigPushMsg{Key: k, Value: "{}"})
			return nil, nil
		}
	}
}
