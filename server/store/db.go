package store

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		panic("db has not been opened")
	}

	return db
}

func Open() {
	var err error
	db, err = bolt.Open("econfig.db", 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
}

func Close() error {
	return db.Close()
}
