package utils

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/slyerr/verifier"
)

func CheckConfigKey(key string) (string, error) {
	err := verifier.S().NotBlankN(key, "config key")
	if err != nil {
		return "", err
	}

	return strings.ToUpper(strings.TrimSpace(key)), nil
}

func CleanConfigValue(value string) string {
	value = strings.TrimSpace(value)

	if len(value) == 0 {
		return "{}"
	}

	var dst bytes.Buffer
	err := json.Compact(&dst, []byte(value))
	if err != nil {
		return value
	}

	return dst.String()
}

func CleanHost(host string) string {
	host = strings.ToLower(strings.Trim(strings.TrimSpace(host), "/"))

	if strings.HasPrefix(host, "localhost") {
		host = strings.ReplaceAll(host, "localhost", "127.0.0.1")
	}

	return host
}

func CompleUrl(host string, url string) string {
	host = CleanHost(host)
	url = strings.TrimSpace(url)

	if len(host) == 0 || strings.Index(url, "://") > 0 {
		return strings.TrimSuffix(url, "/")
	}

	if strings.Index(host, "://") == -1 {
		host = "http://" + host
	}

	if len(url) == 0 {
		return host
	}

	return strings.Trim(host+"/"+strings.Trim(url, "/"), "/")
}
