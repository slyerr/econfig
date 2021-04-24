package client

import (
	"encoding/json"

	"github.com/slyerr/econfig/core/utils"
	sto "github.com/slyerr/econfig/server/store"
	"github.com/slyerr/verifier"
	bolt "go.etcd.io/bbolt"
)

var clientBucket = []byte("client")

type DbClientStore struct {
}

func (s *DbClientStore) Get(key string) ([]Client, error) {
	key, err := utils.CheckConfigKey(key)
	if err != nil {
		return nil, err
	}

	var values [][]byte
	err = sto.DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(clientBucket)
		if bucket == nil {
			return nil
		}

		subBucket := bucket.Bucket([]byte(key))
		if subBucket == nil {
			return nil
		}

		subBucket.ForEach(func(k, v []byte) error {
			values = append(values, v)
			return nil
		})

		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}

	var cc []Client
	for _, v := range values {
		c := Client{}
		err = json.Unmarshal(v, &c)
		if err != nil {
			return nil, err
		}

		cc = append(cc, c)
	}
	return cc, nil
}

func (s *DbClientStore) PutHost(key string, c Client) error {
	key, err := utils.CheckConfigKey(key)
	if err != nil {
		return err
	}

	if err := verifier.S().NotBlankN(c.Host, "client's host"); err != nil {
		return err
	}

	cbs, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return sto.DB().Update(func(tx *bolt.Tx) error {
		// 如果 bucket 不存在则，创建一个 bucket
		bucket, err := tx.CreateBucketIfNotExists(clientBucket)
		if err != nil {
			return err
		}

		subBucket, err := bucket.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return err
		}

		err = subBucket.Put([]byte(utils.CleanHost(c.Host)), cbs)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *DbClientStore) Delete(key string) error {
	key, err := utils.CheckConfigKey(key)
	if err != nil {
		return err
	}

	return sto.DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(clientBucket)
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(key))
	})
}

func (s *DbClientStore) DeleteHost(key string, host string) error {
	key, err := utils.CheckConfigKey(key)
	if err != nil {
		return err
	}

	if err := verifier.S().NotBlankN(host, "host"); err != nil {
		return err
	}

	return sto.DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(clientBucket)
		if bucket == nil {
			return nil
		}

		subBucket := bucket.Bucket([]byte(key))
		if subBucket == nil {
			return nil
		}

		return subBucket.Delete([]byte(utils.CleanHost(host)))
	})
}
