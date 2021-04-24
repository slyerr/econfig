package client

type Client struct {
	Host    string `json:"host"`
	PushUrl string `json:"pushUrl"`
}

type ClientStorer interface {
	Get(key string) ([]Client, error)

	PutHost(key string, c Client) error

	Delete(key string) error

	DeleteHost(key string, host string) error
}

var store ClientStorer = &DbClientStore{}

func Store() ClientStorer {
	return store
}
