package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/slyerr/econfig/core/rest"
	"github.com/slyerr/econfig/core/utils"
	"go.uber.org/zap"
)

type RestClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

func NewRestClient() *RestClient {
	client := &http.Client{}
	return &RestClient{client: client}
}

func NewRestClient1(baseURL string) *RestClient {
	client := &http.Client{}
	return &RestClient{client: client, baseURL: baseURL}
}

func NewRestClient2(baseURL string, headers map[string]string) *RestClient {
	client := &http.Client{}
	return &RestClient{client: client, baseURL: baseURL, headers: headers}
}

func NewRestClientX(client *http.Client) *RestClient {
	return &RestClient{client: client}
}

func NewRestClientX2(client *http.Client, baseURL string) *RestClient {
	return &RestClient{client: client, baseURL: baseURL}
}

func NewRestClientX3(client *http.Client, baseURL string, headers map[string]string) *RestClient {
	return &RestClient{client: client, baseURL: baseURL, headers: headers}
}

func (c *RestClient) Get(url string, reqBody interface{}) (string, error) {
	return c.do(rest.MethodGet, url, reqBody, nil)
}

func (c *RestClient) Get2(url string, reqBody interface{}, resBody interface{}) error {
	_, err := c.do(rest.MethodGet, url, reqBody, resBody)
	return err
}

func (c *RestClient) Post(url string, reqBody interface{}) (string, error) {
	return c.do(rest.MethodPost, url, reqBody, nil)
}

func (c *RestClient) Post2(url string, reqBody interface{}, resBody interface{}) error {
	_, err := c.do(rest.MethodPost, url, reqBody, resBody)
	return err
}

func (c *RestClient) Put(url string, reqBody interface{}) (string, error) {
	return c.do(rest.MethodPut, url, reqBody, nil)
}

func (c *RestClient) Put2(url string, reqBody interface{}, resBody interface{}) error {
	_, err := c.do(rest.MethodPut, url, reqBody, resBody)
	return err
}

func (c *RestClient) Delete(url string, reqBody interface{}) (string, error) {
	return c.do(rest.MethodDelete, url, reqBody, nil)
}

func (c *RestClient) Delete2(url string, reqBody interface{}, resBody interface{}) error {
	_, err := c.do(rest.MethodDelete, url, reqBody, resBody)
	return err
}

func (c *RestClient) do(method rest.Method, url string, reqBody interface{}, resBody interface{}) (string, error) {
	// url
	baseURL := strings.TrimSpace(c.baseURL)
	if len(baseURL) > 0 {
		url = strings.Trim(baseURL, "/") + "/" + strings.Trim(url, "/")
	}

	log := rest.NewLogger("client", string(method), url, zap.S().Debug)

	// request
	var reqBodyReader io.Reader = nil
	switch req := reqBody.(type) {
	case string:
		reqBodyReader = strings.NewReader(req)
	case *string:
		reqBodyReader = strings.NewReader(*req)
	case io.Reader:
		reqBodyReader = req
	case *io.Reader:
		reqBodyReader = *req
	case []byte:
		reqBodyReader = bytes.NewReader(req)
	case *[]byte:
		reqBodyReader = bytes.NewReader(*req)
	default:
		if reqBody != nil {
			if _body, err := json.Marshal(reqBody); err != nil {
				log.Err(err)
				return "", err
			} else {
				reqBodyReader = bytes.NewReader(_body)
			}
		}
	}

	if reqBodyReader != nil {
		var reqBody string
		reqBodyReader, reqBody, _ = utils.ReaderToSrring(reqBodyReader)
		log.Req(reqBody)
	} else {
		log.Req("")
	}

	req, err := http.NewRequest(string(method), url, reqBodyReader)
	if err != nil {
		log.Err(err)
		return "", err
	}
	req.Close = true
	req.Header.Add("Content-Type", rest.ContentType)
	if c.headers != nil && len(c.headers) > 0 {
		for k, v := range c.headers {
			req.Header.Add(k, v)
		}
	}

	// do
	res, err := c.client.Do(req)
	if err != nil {
		log.Err(err)
		return "", err
	}

	statusCode := res.StatusCode
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.ErrC(err, statusCode)
		return "", err
	}

	// response
	if res.StatusCode < 200 || 299 < res.StatusCode {
		err := rest.NewError(res.StatusCode, res.Status)
		log.ErrC(err, statusCode)
		return "", err
	}

	if resBody != nil {
		if err := json.Unmarshal([]byte(bodyBytes), resBody); err != nil {
			log.ErrC(err, statusCode)
			return "", err
		}
	}

	body := string(bodyBytes)
	log.Res(res.StatusCode, body)

	return body, nil
}
