package rest

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Logger struct {
	name   string
	id     string
	method string
	url    string
	log    func(args ...interface{})
}

func NewLogger(name string, method string, url string, log func(args ...interface{})) Logger {
	return Logger{
		name:   cleanName(name),
		id:     uuid.NewString(),
		method: strings.ToUpper(strings.TrimSpace(method)),
		url:    strings.TrimSpace(url),
		log:    log,
	}
}

func (l Logger) Req(body string) {
	text := fmt.Sprintf("%v[%v] request\r\n%v %v\r\n%v\r\n",
		l.name, l.id,
		l.method, l.url,
		cleanBody(body),
	)

	l.log(text)
}

func (l Logger) Res(httpStatusCode int, body string) {
	text := fmt.Sprintf("%v[%v] response\r\n%v %v -> %v\r\n%v\r\n",
		l.name, l.id,
		l.method, l.url, httpStatusCode,
		cleanBody(body),
	)

	l.log(text)
}

func (l Logger) Err(err error) {
	text := fmt.Sprintf("%v[%v] error\r\n%v %v: %v\r\n",
		l.name, l.id,
		l.method, l.url,
		err,
	)

	l.log(text)
}

func (l Logger) ErrC(err error, httpStatusCode int) {
	text := fmt.Sprintf("%v[%v] error\r\n  %v %v -> %v: %v\r\n",
		l.name, l.id,
		l.method, l.url, httpStatusCode,
		err,
	)

	l.log(text)
}

func cleanName(name string) string {
	name = strings.TrimSpace(name)
	if name != "" {
		name = name + " "
	}

	return name
}

func cleanBody(body string) string {
	body = strings.TrimSpace(body)
	if body != "" {
		body = "\r\n" + body
	}

	return body
}
