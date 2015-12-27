package database

import (
	"datatrack/model"
	"sync"

	"github.com/boltdb/bolt"
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
		err = appendValueInList(d.Result, d.Origin, origin)
		if err != nil {
			return err
		}

		// result -> origin
		err = appendValueInList(d.Origin, d.Result, result)
		if err != nil {
			return err
		}
		return nil
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
			err = appendValueInList(d.Result, d.Origin, origin)
			if err != nil {
				return err
			}

			// result -> origin
			err = appendValueInList(d.Origin, d.Result, result)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		errChan <- err
	}
}
