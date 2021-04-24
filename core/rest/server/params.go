package server

import (
	"encoding/json"
	"net/http"

	"github.com/slyerr/econfig/core/utils"
	"goji.io/pat"
)

type Params struct {
	req *http.Request
}

func NewParams(req *http.Request) Params {
	return Params{req}
}

func (p Params) GetPathParam(name string) string {
	return pat.Param(p.req, name)
}

func (p Params) GetQueryParam(name string) string {
	return p.req.URL.Query().Get(name)
}

func (p Params) GetBody(body interface{}) error {
	bytes, err := p.GetBodyBytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, body)
}

func (p Params) GetBodyString() (string, error) {
	r, s, err := utils.ReadCloserToString(p.req.Body)
	if err != nil {
		return "", err
	}

	p.req.Body = r
	return s, nil
}

func (p Params) GetBodyBytes() ([]byte, error) {
	body, err := p.GetBodyString()
	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}
