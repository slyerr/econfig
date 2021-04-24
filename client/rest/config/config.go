package config

import (
	"net/http"
	"strings"

	"github.com/slyerr/econfig/client/config"
	"github.com/slyerr/econfig/core/rest"
	"github.com/slyerr/econfig/core/rest/server"
)

func Put() (rest.Method, string, func(req *http.Request, params server.Params) (interface{}, error)) {
	return rest.MethodPut, strings.TrimSuffix(rest.ClientConfigUrlV1, "/"), func(req *http.Request, params server.Params) (interface{}, error) {
		b, err := params.GetBodyString()
		if err != nil {
			return nil, err
		}

		config.Set(b)
		return nil, nil
	}
}
