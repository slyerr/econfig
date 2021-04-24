package host

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/slyerr/econfig/core/rest"
	"github.com/slyerr/econfig/core/rest/server"
	"github.com/slyerr/econfig/server/store/client"
	"github.com/slyerr/verifier"
)

type PutBody struct {
	Port    int    `json:"port"`
	PushUrl string `json:"pushUrl"`
}

func Put() (rest.Method, string, func(req *http.Request, params server.Params) (interface{}, error)) {
	return rest.MethodPut, rest.ServerHostUrlV1 + ":key", func(req *http.Request, params server.Params) (interface{}, error) {
		b := &PutBody{}
		if err := params.GetBody(b); err != nil {
			return nil, err
		}

		if err := verifier.I().GreaterN(b.Port, 0, "port"); err != nil {
			return nil, err
		}

		ip, err := getIP(req)
		if err != nil {
			return nil, err
		}

		return nil, client.Store().PutHost(
			params.GetPathParam("key"),
			*&client.Client{Host: ip + ":" + strconv.Itoa(b.Port), PushUrl: b.PushUrl},
		)
	}
}

func Delete() (rest.Method, string, func(req *http.Request, params server.Params) (interface{}, error)) {
	return rest.MethodDelete, rest.ServerHostUrlV1 + ":key/:host", func(req *http.Request, params server.Params) (interface{}, error) {
		return nil, client.Store().DeleteHost(params.GetPathParam("key"), params.GetPathParam("host"))
	}
}

func getIP(req *http.Request) (string, error) {
	ips := req.Header.Get("X-Real-Ip")
	if ips == "" {
		ips = req.Header.Get("X-Forwarded-For")
	}
	if ips == "" {
		ips = req.RemoteAddr
	}

	ipss := strings.Split(ips, ",")
	if len(ipss) == 0 {
		return "", rest.NewError(http.StatusBadRequest, "unable to get IP addresses")
	}

	ip := strings.TrimSpace(ipss[0])
	ip = strings.Split(ip, ":")[0]
	if len(ip) == 0 {
		return "", rest.NewError(http.StatusBadRequest, "unable to get IP addresses")
	}

	return strings.TrimSpace(ip), nil
}
