package utils

import (
	"bytes"
	"io"
	"io/ioutil"
)

func ReaderToSrring(r io.Reader) (io.Reader, string, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, "", err
	}

	return bytes.NewReader(bs), string(bs), nil
}

func ReadCloserToString(r io.ReadCloser) (io.ReadCloser, string, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, "", err
	}

	return ioutil.NopCloser(bytes.NewBuffer(bs)), string(bs), nil
}
