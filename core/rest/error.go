package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Error struct {
	Code int
	Msg  string
}

func NewError(code int, msg string) error {
	if code <= 0 {
		code = http.StatusInternalServerError
	}

	msg = strings.TrimSpace(msg)
	if len(msg) == 0 {
		msg = http.StatusText(code)

		if len(msg) == 0 {
			msg = "Unknown error"
		}
	}

	return errors.WithStack(Error{code, msg})
}

func NewArgsError(msg string) error {
	return NewError(http.StatusBadRequest, msg)
}

func (e Error) Error() string {
	return fmt.Sprintf("[%d] %v", e.Code, e.Msg)
}
