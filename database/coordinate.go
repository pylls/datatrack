package database

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sort"
	"strings"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/pylls/datatrack/ephemeral"
	"github.com/pylls/datatrack/model"
)

type sortCoordinates []model.Coordinate

func (a sortCoordinates) Len() int      { return len(a) }
func (a sortCoordinates) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortCoordinates) Less(i, j int) bool {
	return strings.Compare(a[i].ID, a[j].ID) == -1
}

// AddCoordinates adds many coordinates at once.
func AddCoordinates(cs []model.Coordinate, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	sort.Sort(sortCoordinates(cs))

	err := DB.Batch(func(tx *bolt.Tx) error {
		db, err := tx.CreateBucketIfNotExists([]byte("coordinate"))
		if err != nil {
			return err
		}
		dbTime, err := tx.CreateBucketIfNotExists([]byte("coordinate time"))
		if err != nil {
			return err
		}
		dbLat, err := tx.CreateBucketIfNotExists([]byte("coordinate latitude"))
		if err != nil {
			return err
		}

		for _, d := range cs {
			d.Latitude = PadCoordinate(d.Latitude)
			d.Longitude = PadCoordinate(d.Longitude)
			encoded := new(bytes.Buffer)
			enc := gob.NewEncoder(encoded)
			if err = enc.Encode(d); err != nil {
				return err
			}

			if err = db.Put([]byte(d.ID), ephemeral.Encrypt(encoded.Bytes())); err != nil {
				return err
			}
			if err = dbTime.Put([]byte(d.Timestamp), []byte(d.ID)); err != nil {
				return err
			}
			if err = appendValueInList(d.ID, d.Latitude, dbLat); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		errChan <- err
	}
}

// GetCoordinates gets the coordinates within (inclusive)
func GetCoordinates(neLat, neLng, swLat, swLng string) (reply []model.Coordinate, err error) {
	reply = make([]model.Coordinate, 0)
	minLng := []byte(PadCoordinate(swLng))
	maxLng := []byte(PadCoordinate(neLng))
	minLat := []byte(PadCoordinate(swLat))
	maxLat := []byte(PadCoordinate(neLat))

	err = DB.View(func(tx *bolt.Tx) error {
		db := tx.Bucket([]byte("coordinate"))
		c := tx.Bucket([]byte("coordinate latitude")).Cursor()

		// iterate over all ys from swLat to neLat
		for k, v := c.Seek(minLat); k != nil && bytes.Compare(k, maxLat) <= 0; k, v = c.Next() {
			// decode the list
			var list []string
			encoded := bytes.NewBuffer(ephemeral.Decrypt(v))
			dec := gob.NewDecoder(encoded)
			if err = dec.Decode(&list); err != nil {
				return err
			}

			// iterate over the list
			for _, value := range list {
				// decode Coordinate
				var cord model.Coordinate
				encoded := bytes.NewBuffer(ephemeral.Decrypt(db.Get([]byte(value))))
				dec := gob.NewDecoder(encoded)
				if err = dec.Decode(&cord); err != nil {
					return err
				}
				longitude := []byte(cord.Longitude)

				// filter based on within [swLng,neLng]
				if bytes.Compare(longitude, minLng) >= 0 && bytes.Compare(longitude, maxLng) <= 0 {
					reply = append(reply, cord)
				}
			}
		}
		return nil
	})
	return
}

// GetNextCoordinateChrono returns the next coordinate disclosed after the coordinate
// with the provided identifier.
func GetNextCoordinateChrono(id string) (exists bool, next model.Coordinate, err error) {
	exists = false
	err = DB.View(func(tx *bolt.Tx) error {
		db := tx.Bucket([]byte("coordinate"))

		// read target coordinate
		var cord model.Coordinate
		encoded := bytes.NewBuffer(ephemeral.Decrypt(db.Get([]byte(id))))
		dec := gob.NewDecoder(encoded)
		if err = dec.Decode(&cord); err != nil {
			return err
		}

		// attempt to find next
		c := tx.Bucket([]byte("coordinate time")).Cursor()
		c.Seek([]byte(cord.Timestamp))
		_, v := c.Next()
		if v == nil {
			return errors.New("no next coordinate")
		}
		encoded = bytes.NewBuffer(ephemeral.Decrypt(db.Get(v)))
		dec = gob.NewDecoder(encoded)
		return dec.Decode(&next)
	})
	if err == nil {
		exists = true
	} else if err.Error() == "no next coordinate" {
		err = nil
	}

	return
}

// GetPrevCoordinateChrono returns the previous coordinate disclosed before the coordinate
// with the provided identifier.
func GetPrevCoordinateChrono(id string) (exists bool, prev model.Coordinate, err error) {
	err = DB.View(func(tx *bolt.Tx) error {
		db := tx.Bucket([]byte("coordinate"))

		// read target coordinate
		var cord model.Coordinate
		encoded := bytes.NewBuffer(ephemeral.Decrypt(db.Get([]byte(id))))
		dec := gob.NewDecoder(encoded)
		if err = dec.Decode(&cord); err != nil {
			return err
		}

		// attempt to find next
		c := tx.Bucket([]byte("coordinate time")).Cursor()
		c.Seek([]byte(cord.Timestamp))
		_, v := c.Prev()
		if v == nil {
			return errors.New("no previous coordinate")
		}
		encoded = bytes.NewBuffer(ephemeral.Decrypt(db.Get(v)))
		dec = gob.NewDecoder(encoded)
		return dec.Decode(&prev)
	})
	if err == nil {
		exists = true
	} else if err.Error() == "no previous coordinate" {
		err = nil
	}
	return
}

// PadCoordinate pads the coordinate with prefix zeroes
func PadCoordinate(coord string) string {
	switch strings.Index(coord, ".") {
	case 0:
		return "000" + coord
	case 1:
		return "00" + coord
	case 2:
		return "0" + coord
	}
	return coord
}
