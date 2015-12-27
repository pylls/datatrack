package util

import (
	"bytes"
	"testing"
)

func TestDestroyByteSlice(t *testing.T) {
	data := []byte("DATADATADATA")
	DestroyByteSlice(data)
	if bytes.Equal(data, []byte("DATADATADATA")) {
		t.Fail()
	} // TODO: check with a real memory dump
}

func TestDestroyByteSliceShallow(t *testing.T) {
	data := []byte("DATADATADATA")
	DestroyByteSliceShallow(data)
	if bytes.Equal(data, []byte("DATADATADATA")) {
		t.Fail()
	} // TODO: check with a real memory dump
}
