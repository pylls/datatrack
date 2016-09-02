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

// AddDisclosed adds that a disclosure disclosed many attributes.
func AddDisclosed(d model.Disclosed) (err error) {
	return DB.Batch(func(tx *bolt.Tx) error {
		db, err := tx.CreateBucketIfNotExists([]byte("disclosed"))
		if err != nil {
			return err
		}
		dbmap, err := tx.CreateBucketIfNotExists([]byte("disclosed map"))
		if err != nil {
			return err
		}

		encoded := new(bytes.Buffer)
		enc := gob.NewEncoder(encoded)
		err = enc.Encode(d)
		if err != nil {
			return err
		}
		err = db.Put([]byte(d.Disclosure), ephemeral.Encrypt(encoded.Bytes()))
		if err != nil {
			return err
		}

		// update attribute -> []disclosure
		for i := 0; i < len(d.Attribute); i++ {
			err = appendValueInList(d.Disclosure, d.Attribute[i], dbmap)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

type sortDisclosued []model.Disclosed

func (a sortDisclosued) Len() int      { return len(a) }
func (a sortDisclosued) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortDisclosued) Less(i, j int) bool {
	return strings.Compare(a[i].Disclosure, a[j].Disclosure) == -1
}

// AddDiscloseds adds many disclosed at once.
func AddDiscloseds(ds []model.Disclosed, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	sort.Sort(sortDisclosued(ds))

	err := DB.Batch(func(tx *bolt.Tx) error {
		db, err := tx.CreateBucketIfNotExists([]byte("disclosed"))
		if err != nil {
			return err
		}
		dbmap, err := tx.CreateBucketIfNotExists([]byte("disclosed map"))
		if err != nil {
			return err
		}

		for _, d := range ds {
			encoded := new(bytes.Buffer)
			enc := gob.NewEncoder(encoded)
			err = enc.Encode(d)
			if err != nil {
				return err
			}
			err = db.Put([]byte(d.Disclosure), ephemeral.Encrypt(encoded.Bytes()))
			if err != nil {
				return err
			}

			// update attribute -> []disclosure
			for i := 0; i < len(d.Attribute); i++ {
				err = appendValueInList(d.Disclosure, d.Attribute[i], dbmap)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		errChan <- err
	}
}

// GetDisclosed gets the disclosed for a specific disclosure id.
func GetDisclosed(id string) (d *model.Disclosed, err error) {
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("disclosed"))
		if b == nil {
			return errors.New("no disclosed bucket")
		}
		raw := ephemeral.Decrypt(b.Get([]byte(id)))
		if raw == nil {
			return errors.New("no such disclosure")
		}
		d = new(model.Disclosed)
		encoded := bytes.NewBuffer(raw)
		dec := gob.NewDecoder(encoded)
		err = dec.Decode(d)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// GetExplicitlyDisclosedAttributeIDs returns all explicitly disclosed attributes.
func GetExplicitlyDisclosedAttributeIDs() (IDs []string, err error) {
	IDs = make([]string, 0)
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("disclosed map"))
		if b == nil {
			return errors.New("no disclosed map bucket")
		}
		d := tx.Bucket([]byte("disclosure"))
		if d == nil {
			return errors.New("no disclosure bucket")
		}

		// each key has an attribute list
		seen := make(map[string]bool)
		err = b.ForEach(func(k, v []byte) error {
			list, err := getList(string(k), b)
			if err != nil {
				return err
			}

			// the list lists disclosure IDs the atrtribute has been disclosed to
			for i := 0; i < len(list); i++ {
				// get the disclosure and determine if explicit or not
				raw := ephemeral.Decrypt(d.Get([]byte(list[i])))
				if raw == nil {
					return errors.New("no such disclosure")
				}
				disc := new(model.Disclosure)
				encoded := bytes.NewBuffer(raw)
				dec := gob.NewDecoder(encoded)
				err = dec.Decode(disc)
				if err != nil {
					return err
				}

				// if the attribute has been explicitly disclosed once we're done
				if !strings.EqualFold(disc.Sender, disc.Recipient) {
					seen[string(k)] = true
					break
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		// assemble reply
		for id := range seen {
			result := make([]byte, len(id))
			copy(result, []byte(id))
			IDs = append(IDs, string(result))
		}
		return nil
	})
	return
}

// GetImplicitlyDisclosedAttributeIDs returns all implicitly disclosed attributes.
func GetImplicitlyDisclosedAttributeIDs() (IDs []string, err error) {
	IDs = make([]string, 0)
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("disclosed map"))
		if b == nil {
			return errors.New("no disclosed map bucket")
		}
		d := tx.Bucket([]byte("disclosure"))
		if d == nil {
			return errors.New("no disclosure bucket")
		}

		// each key has an attribute list
		return b.ForEach(func(k, v []byte) error {
			list, err := getList(string(k), b)
			if err != nil {
				return err
			}

			// the list lists disclosure IDs the atrtribute has been disclosed to
			for i := 0; i < len(list); i++ {
				// get the disclosure and determine if implicit or not
				raw := ephemeral.Decrypt(d.Get([]byte(list[i])))
				if raw == nil {
					return errors.New("no such disclosure")
				}
				disc := new(model.Disclosure)
				encoded := bytes.NewBuffer(raw)
				dec := gob.NewDecoder(encoded)
				err = dec.Decode(disc)
				if err != nil {
					return err
				}

				// in the model, sent to self -> implicit
				if strings.EqualFold(disc.Sender, disc.Recipient) {
					result := make([]byte, len(k))
					copy(result, []byte(k))
					IDs = append(IDs, string(result))
				}
			}

			return nil
		})

	})
	return
}
