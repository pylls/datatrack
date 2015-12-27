package database

import (
	"bytes"
	"datatrack/model"
	"encoding/gob"
	"errors"

	"github.com/boltdb/bolt"
)

// AddCategory adds a category to the database.
func AddCategory(c model.Category) (err error) {
	return DB.Batch(func(tx *bolt.Tx) error {
		db, err := tx.CreateBucketIfNotExists([]byte("category"))
		if err != nil {
			return err
		}
		encoded := new(bytes.Buffer)
		enc := gob.NewEncoder(encoded)
		err = enc.Encode(c)
		if err != nil {
			return err
		}

		return appendValueInList(encoded.String(), c.ID, db)
	})
}

// GetCategory returns the category with the provided identifier.
func GetCategory(id string) (c []model.Category, err error) {
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("category"))
		if b == nil {
			return errors.New("no category bucket")
		}
		list, err := getList(id, b)
		if err != nil {
			return err
		}

		c = make([]model.Category, 0, len(list))
		for i := 0; i < len(list); i++ {
			var tmp model.Category
			encoded := bytes.NewBuffer([]byte(list[i]))
			dec := gob.NewDecoder(encoded)
			err = dec.Decode(&tmp)
			if err != nil {
				return err
			}
			c = append(c, tmp)
		}
		return nil
	})
	return
}
