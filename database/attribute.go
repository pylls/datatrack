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

// AddAttribute adds an attribute to the database.
func AddAttribute(a model.Attribute) (err error) {
	return DB.Batch(func(tx *bolt.Tx) error {
		ab, err := tx.CreateBucketIfNotExists([]byte("attribute"))
		if err != nil {
			return err
		}
		abmap, err := tx.CreateBucketIfNotExists([]byte("attribute map"))
		if err != nil {
			return err
		}
		encoded := new(bytes.Buffer)
		enc := gob.NewEncoder(encoded)
		if err = enc.Encode(a); err != nil {
			return err
		}
		if err = ab.Put([]byte(a.ID), ephemeral.Encrypt(encoded.Bytes())); err != nil {
			return err
		}

		// update type->id map
		if err = appendValueInMap(a.Type, a.ID, "type2id", abmap); err != nil {
			return err
		}

		// update type->value map
		return appendValueInMap(a.Type, a.Value, "type2value", abmap)
	})
}

type sortAttribute []model.Attribute

func (a sortAttribute) Len() int      { return len(a) }
func (a sortAttribute) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortAttribute) Less(i, j int) bool {
	return strings.Compare(a[i].ID, a[j].ID) == -1
}

// AddAttributes adds many attributes at once.
func AddAttributes(as []model.Attribute, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	sort.Sort(sortAttribute(as))

	err := DB.Batch(func(tx *bolt.Tx) error {
		ab, err := tx.CreateBucketIfNotExists([]byte("attribute"))
		if err != nil {
			return err
		}
		abmap, err := tx.CreateBucketIfNotExists([]byte("attribute map"))
		if err != nil {
			return err
		}

		// read maps
		type2id, err := getMap("type2id", abmap)
		if err != nil {
			return err
		}
		type2value, err := getMap("type2value", abmap)
		if err != nil {
			return err
		}

		// write all attributes
		for _, a := range as {
			encoded := new(bytes.Buffer)
			enc := gob.NewEncoder(encoded)
			err = enc.Encode(a)
			if err != nil {
				return err
			}
			err = ab.Put([]byte(a.ID), ephemeral.Encrypt(encoded.Bytes()))
			if err != nil {
				return err
			}

			// update type->id map
			data, exists := type2id[a.Type]
			if !exists {
				data = make([]string, 0, 1)
			}
			type2id[a.Type] = append(data, a.ID)

			// update type->value map
			if data, exists = type2value[a.Type]; !exists {
				data = make([]string, 0, 1)
			}
			type2value[a.Type] = append(data, a.Value)
		}

		// write maps
		if err = writeMap("type2id", type2id, abmap); err != nil {
			return err
		}

		return writeMap("type2value", type2value, abmap)
	})
	if err != nil {
		errChan <- err
	}
}

// GetAttribute returns the attribute with the provided id.
func GetAttribute(id string) (a *model.Attribute, err error) {
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("attribute"))
		if b == nil {
			return errors.New("no attribute bucket")
		}
		raw := ephemeral.Decrypt(b.Get([]byte(id)))
		if raw == nil {
			return errors.New("no such attribute")
		}
		a = new(model.Attribute)
		encoded := bytes.NewBuffer(raw)
		dec := gob.NewDecoder(encoded)
		err = dec.Decode(a)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// GetAttributeIDs returns all attribute identifiers.
func GetAttributeIDs() (IDs []string, err error) {
	IDs = make([]string, 0)
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("attribute"))
		if b == nil {
			return errors.New("no attribute bucket")
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

// GetAttributeIDsToOrg returns all attribute IDs disclosed to an organization.
func GetAttributeIDsToOrg(org string) (IDs []string, err error) {
	IDs = make([]string, 0)

	// disclosures to org
	workIDs, err := GetDisclosureIDsToOrg(org)
	if err != nil {
		return
	}

	distinctAttr := make(map[string]bool)
	for i := 0; i < len(workIDs); i++ {
		// disclosed attributes in each disclosure
		disclosed, err := GetDisclosed(workIDs[i])
		if err != nil {
			return nil, err
		}
		// collect distinct attribute IDs, abusing map
		for _, id := range disclosed.Attribute {
			distinctAttr[id] = true
		}
	}
	// assemble reply
	for id := range distinctAttr {
		IDs = append(IDs, id)
	}

	return
}

// GetExplicitlyDisclosedAttributeIDsToOrg is like GetAttributeIDsToOrg but for
// explicitly disclosed attributes only.
func GetExplicitlyDisclosedAttributeIDsToOrg(org string) (IDs []string, err error) {
	IDs = make([]string, 0)

	// explicit disclosures to org
	workIDs, err := GetExplicitDisclosureIDsToOrg(org)
	if err != nil {
		return
	}

	distinctAttr := make(map[string]bool)
	for i := 0; i < len(workIDs); i++ {
		// disclosed attributes in each disclosure
		disclosed, err := GetDisclosed(workIDs[i])
		if err != nil {
			return nil, err
		}
		// collect distinct attribute IDs, abusing map
		for _, id := range disclosed.Attribute {
			distinctAttr[id] = true
		}
	}
	// assemble reply
	for id := range distinctAttr {
		IDs = append(IDs, id)
	}

	return
}

// GetImplicitlyDisclosedAttributeIDsToOrg is like GetAttributeIDsToOrg but for
// implicitly disclosed attributes only.
func GetImplicitlyDisclosedAttributeIDsToOrg(org string) (IDs []string, err error) {
	IDs = make([]string, 0)

	// explicit disclosures to org
	workIDs, err := GetImplictDisclosureIDsToOrg(org)
	if err != nil {
		return
	}

	distinctAttr := make(map[string]bool)
	for i := 0; i < len(workIDs); i++ {
		// disclosed attributes in each disclosure
		disclosed, err := GetDisclosed(workIDs[i])
		if err != nil {
			return nil, err
		}
		// collect distinct attribute IDs, abusing map
		for _, id := range disclosed.Attribute {
			distinctAttr[id] = true
		}
	}
	// assemble reply
	for id := range distinctAttr {
		IDs = append(IDs, id)
	}

	return
}

// GetAttributeTypeIDs returns all attribute type IDs.
func GetAttributeTypeIDs() (IDs []string, err error) {
	IDs = make([]string, 0)
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("attribute map"))
		if b == nil {
			return errors.New("no attribute map bucket")
		}
		themap, err := getMap("type2id", b)
		if err != nil {
			return err
		}

		for id := range themap {
			result := make([]byte, len(id))
			copy(result, []byte(id))
			IDs = append(IDs, string(result))
		}

		return nil
	})
	return
}

// GetAttributeTypeValues returns all attribute values for a type.
func GetAttributeTypeValues(thetype string) (values []string, err error) {
	values = make([]string, 0)
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("attribute map"))
		if b == nil {
			return errors.New("no attribute map bucket")
		}
		themap, err := getMap("type2value", b)
		if err != nil {
			return err
		}

		val, exists := themap[thetype]
		if !exists {
			return errors.New("no such values")
		}

		for i := 0; i < len(val); i++ {
			result := make([]byte, len(val[i]))
			copy(result, []byte(val[i]))
			values = append(values, string(result))
		}

		return nil
	})
	return
}

// GetExplicitlyDisclosedAttributeTypeIDs returns all attribute type IDs for
// explicitly disclosed disclosures.
func GetExplicitlyDisclosedAttributeTypeIDs() (IDs []string, err error) {
	IDs = make([]string, 0)

	// find out which attributes were explicitly disclosed
	workIDs, err := GetExplicitlyDisclosedAttributeIDs()
	if err != nil {
		return
	}

	// find the type, abuse map for tracking unique IDs
	types := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return IDs, err
		}
		types[attr.Type] = true
	}

	for typeID := range types {
		IDs = append(IDs, typeID)
	}
	return
}

// GetExplicitlyDisclosedAttributeTypeValues returns all attribute values for a type for
// explicitly disclosed disclosures.
func GetExplicitlyDisclosedAttributeTypeValues(thetype string) (values []string, err error) {
	values = make([]string, 0)
	// find out which attributes were explicitly disclosed
	workIDs, err := GetExplicitlyDisclosedAttributeIDs()
	if err != nil {
		return
	}

	// find the type, abuse map for tracking unique IDs
	valuesmap := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return values, err
		}
		valuesmap[attr.Value] = true
	}

	for v := range valuesmap {
		values = append(values, v)
	}
	return
}

// GetImplicitlyDisclosedAttributeTypeIDs returns all attribute type IDs for
// implicity disclosed disclosures.
func GetImplicitlyDisclosedAttributeTypeIDs() (IDs []string, err error) {
	IDs = make([]string, 0)

	// find out which attributes were implicity disclosed
	workIDs, err := GetImplicitlyDisclosedAttributeIDs()
	if err != nil {
		return
	}

	// find the type, abuse map for tracking unique IDs
	types := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return IDs, err
		}
		types[attr.Type] = true
	}

	for typeID := range types {
		IDs = append(IDs, typeID)
	}
	return
}

// GetImplicitlyDisclosedAttributeTypeValues returns all attribute values for a type for
// implicty disclosed disclosures.
func GetImplicitlyDisclosedAttributeTypeValues(thetype string) (values []string, err error) {
	values = make([]string, 0)
	// find out which attributes were explicitly disclosed
	workIDs, err := GetImplicitlyDisclosedAttributeIDs()
	if err != nil {
		return
	}

	// find the values, abuse map for tracking unique IDs
	valuesmap := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return values, err
		}
		if attr.Type == thetype {
			valuesmap[attr.Value] = true
		}
	}

	for v := range valuesmap {
		values = append(values, v)
	}
	return
}

// GetAttributeTypeIDsToOrg returns all attribute type IDs disclosed to an organization.
func GetAttributeTypeIDsToOrg(org string) (IDs []string, err error) {
	IDs = make([]string, 0)

	// find which attributes we disclosed
	workIDs, err := GetAttributeIDsToOrg(org)
	if err != nil {
		return
	}

	// find the types, abuse map for tracking unique IDs
	types := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return IDs, err
		}
		types[attr.Type] = true
	}

	for typeID := range types {
		IDs = append(IDs, typeID)
	}
	return
}

// GetAttributeTypeValuesToOrg returns all attribute type values disclosed to an organization.
func GetAttributeTypeValuesToOrg(org string, thetype string) (values []string, err error) {
	values = make([]string, 0)

	// find which attributes we disclosed
	workIDs, err := GetAttributeIDsToOrg(org)
	if err != nil {
		return
	}

	// find the type, abuse map for tracking unique IDs
	valuesmap := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return values, err
		}
		if strings.EqualFold(thetype, attr.Type) {
			valuesmap[attr.Value] = true
		}
	}

	for v := range valuesmap {
		values = append(values, v)
	}
	return
}

// GetExplicitlyDisclosedAttributeTypeIDsToOrg returns all attribute type IDs for
// explicitly disclosed disclosures to a particular organization.
func GetExplicitlyDisclosedAttributeTypeIDsToOrg(org string) (IDs []string, err error) {
	IDs = make([]string, 0)

	// find out which attributes were explicitly disclosed
	workIDs, err := GetExplicitlyDisclosedAttributeIDsToOrg(org)
	if err != nil {
		return
	}

	// find the type, abuse map for tracking unique IDs
	types := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return IDs, err
		}
		types[attr.Type] = true
	}

	for typeID := range types {
		IDs = append(IDs, typeID)
	}
	return
}

// GetExplicitlyDisclosedAttributeTypeValuesToOrg returns all attribute values for a type for
// explicitly disclosed disclosures to a particular organization.
func GetExplicitlyDisclosedAttributeTypeValuesToOrg(org string, thetype string) (values []string, err error) {
	values = make([]string, 0)
	// find out which attributes were explicitly disclosed
	workIDs, err := GetExplicitlyDisclosedAttributeIDsToOrg(org)
	if err != nil {
		return
	}

	// find the value, abuse map for tracking unique IDs
	valuesmap := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return values, err
		}
		valuesmap[attr.Value] = true
	}

	for v := range valuesmap {
		values = append(values, v)
	}
	return
}

// GetImplicitlyDisclosedAttributeTypeIDsToOrg returns all attribute type IDs for
// implicity disclosed disclosures to a particular organization.
func GetImplicitlyDisclosedAttributeTypeIDsToOrg(org string) (IDs []string, err error) {
	IDs = make([]string, 0)

	// find out which attributes were explicitly disclosed
	workIDs, err := GetImplicitlyDisclosedAttributeIDsToOrg(org)
	if err != nil {
		return
	}

	// find the type, abuse map for tracking unique IDs
	types := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return IDs, err
		}
		types[attr.Type] = true
	}

	for typeID := range types {
		IDs = append(IDs, typeID)
	}
	return
}

// GetImplicitlyDisclosedAttributeTypeValuesToOrg returns all attribute values for a type for
// implicty disclosed disclosures to a particular organization.
func GetImplicitlyDisclosedAttributeTypeValuesToOrg(org, thetype string) (values []string, err error) {
	values = make([]string, 0)
	// find out which attributes were explicitly disclosed
	workIDs, err := GetImplicitlyDisclosedAttributeIDsToOrg(org)
	if err != nil {
		return
	}

	// find the type, abuse map for tracking unique IDs
	valuesmap := make(map[string]bool)
	for _, id := range workIDs {
		attr, err := GetAttribute(id)
		if err != nil {
			return values, err
		}
		if attr.Type == thetype {
			valuesmap[attr.Value] = true
		}
	}

	for v := range valuesmap {
		values = append(values, v)
	}
	return
}
