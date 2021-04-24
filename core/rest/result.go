package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResult(code int, msg string, data interface{}) *Result {
	if code <= 0 {
		code = http.StatusOK
	}

	msg = strings.TrimSpace(msg)

	return &Result{Code: code, Msg: msg, Data: data}
}

func NewResultX(x interface{}) *Result {
	switch t := x.(type) {
	case Result:
		return &t
	case *Result:
		return t
	case error:
		switch e := errors.Cause(t).(type) {
		case Error:
			return NewResult(e.Code, e.Msg, nil)
		case *Error:
			return NewResult(e.Code, e.Msg, nil)
		default:
			return NewResult(http.StatusInternalServerError, t.Error(), nil)
		}
	default:
		return NewResult(http.StatusOK, "", x)
	}
}

func (r *Result) ToBytes() ([]byte, error) {
	text, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return []byte(text), nil
}
