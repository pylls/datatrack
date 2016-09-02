package util

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteRead(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test.insynd.rw")
	data := []byte("some data to write")

	err := WriteToLocation(data, path)
	if err != nil {
		t.Error("failed to write testdata: " + err.Error())
	}

	buffer, err := ReadFromLocation(path)
	if err != nil {
		t.Error("failed to read testdata: " + err.Error())
	}

	if !bytes.Equal(data, buffer) {
		t.Error("read data differs from written data")
	}
	os.Remove(path)
}
