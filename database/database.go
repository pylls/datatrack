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
		_, err = tx.CreateBucketIfNotExists([]byte("attribute map"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("category"))
		_, err = tx.CreateBucketIfNotExists([]byte("category"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("disclosed"))
		_, err = tx.CreateBucketIfNotExists([]byte("disclosed"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("disclosed map"))
		_, err = tx.CreateBucketIfNotExists([]byte("disclosed map"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("disclosure"))
		_, err = tx.CreateBucketIfNotExists([]byte("disclosure"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("disclosure map"))
		_, err = tx.CreateBucketIfNotExists([]byte("disclosure map"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("downstream origin"))
		_, err = tx.CreateBucketIfNotExists([]byte("downstream origin"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("downstream result"))
		_, err = tx.CreateBucketIfNotExists([]byte("downstream result"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("organization"))
		_, err = tx.CreateBucketIfNotExists([]byte("organization"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("user"))
		_, err = tx.CreateBucketIfNotExists([]byte("user"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("coordinate"))
		_, err = tx.CreateBucketIfNotExists([]byte("coordinate"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("coordinate time"))
		_, err = tx.CreateBucketIfNotExists([]byte("coordinate time"))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte("coordinate latitude"))
		_, err = tx.CreateBucketIfNotExists([]byte("coordinate latitude"))
		if err != nil {
			return err
		}

		return nil
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

func appendValueInMap(key, value, mapname string, bucket *bolt.Bucket) (err error) {
	themap, err := getMap(mapname, bucket)
	if err != nil {
		return err
	}

	data, exists := themap[key]
	if !exists {
		data = make([]string, 0, 1)
	}
	themap[key] = append(data, value)

	return writeMap(mapname, themap, bucket)
}

func writeMap(mapname string, themap map[string][]string, bucket *bolt.Bucket) (err error) {
	encoded := new(bytes.Buffer)
	enc := gob.NewEncoder(encoded)
	err = enc.Encode(themap)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(mapname), ephemeral.Encrypt(encoded.Bytes()))
}

func getMap(id string, bucket *bolt.Bucket) (themap map[string][]string, err error) {
	themap = make(map[string][]string)
	raw := bucket.Get([]byte(id))
	if raw == nil {
		return
	}
	encoded := bytes.NewBuffer(ephemeral.Decrypt(raw))
	dec := gob.NewDecoder(encoded)
	err = dec.Decode(&themap)
	return
}
