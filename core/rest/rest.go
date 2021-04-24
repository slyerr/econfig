package rest

type Method string

const (
	MethodGet    Method = "GET"
	MethodPost   Method = "POST"
	MethodPut    Method = "PUT"
	MethodDelete Method = "DELETE"
)

const ContentType = "application/json"

const (
	ServerHostUrlV1   = "/rest/v1/host/"
	ServerConfigUrlV1 = "/rest/v1/config/"
	ClientConfigUrlV1 = "/rest/v1/config/"
)
