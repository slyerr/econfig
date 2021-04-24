package config

import (
	"github.com/slyerr/econfig/core/utils"
	sto "github.com/slyerr/econfig/server/store"
	bolt "go.etcd.io/bbolt"
)

var configBucket = []byte("config")

type DbConfigStore struct {
}

func (s *DbConfigStore) Get(key string) (string, error) {
	key, err := utils.CheckConfigKey(key)
	if err != nil {
		return "", err
	}

	var value []byte
	err = sto.DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(configBucket)
		if bucket == nil {
			return nil
		}

		value = bucket.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", err
	}

	return utils.CleanConfigValue(string(value)), nil
}

func (s *DbConfigStore) Put(key string, value string) (string, string, error) {
	key, err := utils.CheckConfigKey(key)
	if err != nil {
		return "", "", err
	}

	value = utils.CleanConfigValue(value)

	err = sto.DB().Update(func(tx *bolt.Tx) error {
		// 如果 bucket 不存在则，创建一个 bucket
		bucket, err := tx.CreateBucketIfNotExists(configBucket)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), []byte(value))
	})
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

func (s *DbConfigStore) Delete(key string) (string, error) {
	key, err := utils.CheckConfigKey(key)
	if err != nil {
		return "", err
	}

	err = sto.DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(configBucket)
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(key))
	})
	if err != nil {
		return "", err
	}

	return key, nil
}
