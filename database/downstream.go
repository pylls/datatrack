package database

import (
	"sync"

	"github.com/boltdb/bolt"
	"github.com/pylls/datatrack/model"
)

// AddDownstream adds a downstream relationship.
func AddDownstream(d model.Downstream) (err error) {
	return DB.Batch(func(tx *bolt.Tx) error {
		origin, err := tx.CreateBucketIfNotExists([]byte("downstream origin"))
		if err != nil {
			return err
		}
		result, err := tx.CreateBucketIfNotExists([]byte("downstream result"))
		if err != nil {
			return err
		}

		// origin -> result
		if err = appendValueInList(d.Result, d.Origin, origin); err != nil {
			return err
		}

		// result -> origin
		return appendValueInList(d.Origin, d.Result, result)
	})
}

// AddDownstreams adds many downstream at once.
func AddDownstreams(ds []model.Downstream, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	err := DB.Batch(func(tx *bolt.Tx) error {
		origin, err := tx.CreateBucketIfNotExists([]byte("downstream origin"))
		if err != nil {
			return err
		}
		result, err := tx.CreateBucketIfNotExists([]byte("downstream result"))
		if err != nil {
			return err
		}
		for _, d := range ds {
			// origin -> result
			if err = appendValueInList(d.Result, d.Origin, origin); err != nil {
				return err
			}

			// result -> origin
			if err = appendValueInList(d.Origin, d.Result, result); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		errChan <- err
	}
}
