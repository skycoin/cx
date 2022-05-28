package bolt

import (
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func init() {
	db, err := bolt.Open("program_list.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("program"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func Add(key string, value []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("program"))
		if bucket == nil {
			return errors.New("bucket program was not found")
		}
		err := bucket.Put([]byte(key), value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func Get(key string) []byte {
	ret := []byte{}
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("program"))
		if bucket == nil {
			return errors.New("bucket program was not found")
		}
		ret = bucket.Get([]byte(key))
		if ret == nil {
			return errors.New("No value associated with key " + key)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return ret
}
