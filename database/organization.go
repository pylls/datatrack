package database

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/boltdb/bolt"

	"github.com/pylls/datatrack/ephemeral"
	"github.com/pylls/datatrack/model"
)

// NoSuchOrgError is the error message on no such organization
const NoSuchOrgError = "no such organization"

// AddOrganization adds an organization to the database.
func AddOrganization(o model.Organization) (err error) {
	return DB.Batch(func(tx *bolt.Tx) error {
		db, err := tx.CreateBucketIfNotExists([]byte("organization"))
		if err != nil {
			return err
		}
		encoded := new(bytes.Buffer)
		enc := gob.NewEncoder(encoded)
		err = enc.Encode(o)
		if err != nil {
			return err
		}
		err = db.Put([]byte(o.ID), ephemeral.Encrypt(encoded.Bytes()))
		if err != nil {
			return err
		}

		return nil
	})
}

// GetOrganization returns the organization with the provided identifier. Returns
// NoSuchOrgError on no organization with the provided identifier.
func GetOrganization(id string) (org *model.Organization, err error) {
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("organization"))
		if b == nil {
			return errors.New("no organization bucket")
		}
		raw := ephemeral.Decrypt(b.Get([]byte(id)))
		if raw == nil {
			return errors.New(NoSuchOrgError)
		}
		org = new(model.Organization)
		encoded := bytes.NewBuffer(raw)
		dec := gob.NewDecoder(encoded)
		err = dec.Decode(org)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// GetOrganizationIDs returns all the organization IDs.
func GetOrganizationIDs() (IDs []string, err error) {
	IDs = make([]string, 0)
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("organization"))
		if b == nil {
			return errors.New("no organization bucket")
		}
		return b.ForEach(func(k, v []byte) error {
			result := make([]byte, len(k))
			copy(result, k)
			IDs = append(IDs, string(result))
			return nil
		})
	})
	return
}

// GetReceivingOrgIDs returns the organization IDs of organizations that
// have received a particular attribute.
func GetReceivingOrgIDs(attribute string) (IDs []string, err error) {
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

		list, err := getList(attribute, b)
		if err != nil {
			return err
		}

		// the list lists disclosure IDs the atttribute has been disclosed to
		seen := make(map[string]bool)
		for i := 0; i < len(list); i++ {
			// get the disclosure
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
			seen[disc.Recipient] = true
		}

		for id := range seen {
			result := make([]byte, len(id))
			copy(result, []byte(id))
			IDs = append(IDs, string(result))
		}
		return nil
	})
	return
}
