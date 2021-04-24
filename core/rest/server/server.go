package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/slyerr/verifier"
	"go.uber.org/zap"
	"goji.io"
)

func Start(name string, port int, actions []Action) {
	verifier.I().GreaterNP(port, 0, "port")

	if len(actions) == 0 {
		return
	}

	mux := goji.NewMux()
	for _, action := range actions {
		p, f := action.toHandleFunc()
		mux.HandleFunc(p, f)

		zap.S().Warnf(name+"rest server registered action:%v %v",
			func(m map[string]struct{}) string {
				if len(m) == 0 {
					return "[]"
				}

				s := "["
				for k, _ := range m {
					s += k + ", "
				}
				s = strings.TrimSuffix(s, ", ") + "]"
				return s
			}(p.HTTPMethods()),
			p.String())
	}

	host := "localhost:" + strconv.Itoa(port)
	zap.S().Warnf(name+"rest server has been started: %v", host)
	http.ListenAndServe(host, mux)
}
