package server

import (
	"net/http"

	"github.com/slyerr/econfig/core/rest"
	"github.com/slyerr/econfig/core/utils"
	"github.com/slyerr/verifier"
	"go.uber.org/zap"
	"goji.io/pat"
)

type Action func() (rest.Method, string, func(*http.Request, Params) (interface{}, error))

func (a Action) toHandleFunc() (*pat.Pattern, func(http.ResponseWriter, *http.Request)) {
	method, pattern, run := a()

	verifier.S().NotBlankNP(pattern, "actions's pattern")
	verifier.IF().NotNilNP(run, "actions's run")

	var pp *pat.Pattern
	switch method {
	case rest.MethodGet:
		pp = pat.Get(pattern)
	case rest.MethodPost:
		pp = pat.Post(pattern)
	case rest.MethodPut:
		pp = pat.Put(pattern)
	case rest.MethodDelete:
		pp = pat.Delete(pattern)
	default:
		panic("action's method must be one of the Get, POST, PUT, DELETE")
	}

	return pp, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", rest.ContentType)

		log := rest.NewLogger("server", string(method), r.URL.String(), zap.S().Debug)
		b, s, _ := utils.ReadCloserToString(r.Body)
		r.Body = b
		log.Req(s)

		var result *rest.Result
		if body, err := run(r, NewParams(r)); err != nil {
			result = rest.NewResultX(err)
		} else {
			result = rest.NewResultX(body)
		}

		resBody, err := result.ToBytes()
		if err != nil {
			zap.S().Error(err)

			result = rest.NewResult(http.StatusInternalServerError, err.Error(), nil)
			resBody, _ = result.ToBytes()
		}

		w.WriteHeader(result.Code)
		w.Write(resBody)
		log.Res(result.Code, string(resBody))
	}
}
