package config

type ConfigStorer interface {
	Get(key string) (string, error)

	Put(key string, value string) (string, string, error)

	Delete(key string) (string, error)
}

var store ConfigStorer = &DbConfigStore{}

func Store() ConfigStorer {
	return store
}
