package database

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/pylls/datatrack/ephemeral"
)

// DB is our database presumably created somewhere else.
var DB *bolt.DB

// Self is the constant that represents the user disclosing something him- or herself.
const Self = "USER"

// Start attempts to start the database from the provided file.
func Start(file string) (err error) {
	DB, err = bolt.Open(file, 0600, nil)
	if err != nil {
		return
	}
	return Setup()
}

// Close attempts to close the database.
func Close() (err error) {
	return DB.Close()
}

// Setup creates all buckets for bolt.
func Setup() (err error) {
	return DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("attribute"))
		_, err := tx.CreateBucketIfNotExists([]byte("attribute"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("attribute map"))
		if _, err = tx.CreateBucketIfNotExists([]byte("attribute map")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("category"))
		if _, err = tx.CreateBucketIfNotExists([]byte("category")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("disclosed"))
		if _, err = tx.CreateBucketIfNotExists([]byte("disclosed")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("disclosed map"))
		if _, err = tx.CreateBucketIfNotExists([]byte("disclosed map")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("disclosure"))
		if _, err = tx.CreateBucketIfNotExists([]byte("disclosure")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("disclosure map"))
		if _, err = tx.CreateBucketIfNotExists([]byte("disclosure map")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("downstream origin"))
		if _, err = tx.CreateBucketIfNotExists([]byte("downstream origin")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("downstream result"))
		if _, err = tx.CreateBucketIfNotExists([]byte("downstream result")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("organization"))
		if _, err = tx.CreateBucketIfNotExists([]byte("organization")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("user"))
		if _, err = tx.CreateBucketIfNotExists([]byte("user")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("coordinate"))
		if _, err = tx.CreateBucketIfNotExists([]byte("coordinate")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("coordinate time"))
		if _, err = tx.CreateBucketIfNotExists([]byte("coordinate time")); err != nil {
			return err
		}
		tx.DeleteBucket([]byte("coordinate latitude"))
		_, err = tx.CreateBucketIfNotExists([]byte("coordinate latitude"))

		return err
	})
}

func appendValueInList(value, name string, bucket *bolt.Bucket) (err error) {
	list, err := getList(name, bucket)
	if err != nil {
		return err
	}
	list = append(list, value)
	return writeList(name, list, bucket)
}

func writeList(name string, list []string, bucket *bolt.Bucket) (err error) {
	encoded := new(bytes.Buffer)
	enc := gob.NewEncoder(encoded)
	err = enc.Encode(list)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(name), ephemeral.Encrypt(encoded.Bytes()))
}

func getList(id string, bucket *bolt.Bucket) (list []string, err error) {
	list = make([]string, 0)
	raw := bucket.Get([]byte(id))
	if raw == nil {
		return
	}
	encoded := bytes.NewBuffer(ephemeral.Decrypt(raw))
	dec := gob.NewDecoder(encoded)
	err = dec.Decode(&list)
	return
}

func appendValueInMap(key, value, mapName string, bucket *bolt.Bucket) (err error) {
	m, err := getMap(mapName, bucket)
	if err != nil {
		return err
	}

	data, exists := m[key]
	if !exists {
		data = make([]string, 0, 1)
	}
	m[key] = append(data, value)

	return writeMap(mapName, m, bucket)
}

func writeMap(mapName string, m map[string][]string, bucket *bolt.Bucket) (err error) {
	encoded := new(bytes.Buffer)
	enc := gob.NewEncoder(encoded)
	err = enc.Encode(m)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(mapName), ephemeral.Encrypt(encoded.Bytes()))
}

func getMap(id string, bucket *bolt.Bucket) (m map[string][]string, err error) {
	m = make(map[string][]string)
	raw := bucket.Get([]byte(id))
	if raw == nil {
		return
	}
	encoded := bytes.NewBuffer(ephemeral.Decrypt(raw))
	dec := gob.NewDecoder(encoded)
	err = dec.Decode(&m)
	return
}
