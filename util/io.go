package util

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func WriteToLocation(data []byte, location string) (err error) {
	// make sure we can write to the directory
	err = os.MkdirAll(filepath.Dir(location), 0700)
	if err != nil {
		return errors.New("failed to create directories to location: " + err.Error())
	}

	// write tmp file
	tmpLocation := location + "TMP"
	err = ioutil.WriteFile(tmpLocation, data, 0600)
	if err != nil {
		err = errors.New("Failed to write to temporary file location: " + err.Error())
	}

	// remove the old file at the location
	err = os.Remove(location)
	if err != nil {
		err = errors.New("Failed to remove the old file at the given location: " + err.Error())
	}

	// rename
	err = os.Rename(tmpLocation, location)
	if err != nil {
		err = errors.New("Failed to rename the temporary file to the proper location: " + err.Error())
	}

	return
}

func ReadFromLocation(location string) (data []byte, err error) {
	data, err = ioutil.ReadFile(location)
	if err != nil {
		return nil, errors.New("Failed to read the file at " + location + ": " + err.Error())
	}

	return
}

func ReadFileOrDie(location string) []byte {
	encoded, err := ReadFromLocation(location)
	if err != nil {
		log.Printf("failed to read from location: %s", location)
		log.Fatal(err)
	}
	return encoded
}

func Link(oldname, newname string) (err error) {
	os.MkdirAll(filepath.Dir(newname), 0700)
	e := os.Link(oldname, newname)
	if e == nil {
		return
	}
	cmd := exec.Command("cmd", "/c", "mklink", "/H", newname, oldname)
	output, er := cmd.CombinedOutput()
	if er != nil {
		return errors.New("failed to create symlink, golang os error (" + e.Error() + "), windows-specific error (" + er.Error() + string(output) + oldname + newname + ")")
	}
	return
}
