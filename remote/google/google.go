package google

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"github.com/pylls/datatrack/database"
	"github.com/pylls/datatrack/model"
)

var org = model.Organization{
	ID:          "Google",
	Name:        "Google",
	URL:         "https://www.google.com",
	Description: "Google Inc. is an American multinational technology company specializing in Internet-related services and products.",
}

// ParseTakeoutGzip parses a Google takeout file in tar.gz (.tgz) format.
func ParseTakeoutGZIP(reader io.Reader) (err error) {
	greader, err := gzip.NewReader(reader)
	if err != nil {
		return
	}
	if err := database.AddOrganization(org); err != nil {
		return err
	}
	treader := tar.NewReader(greader)
	for h, err := treader.Next(); err == nil; h, err = treader.Next() {
		switch h.Name {
		case "Takeout/Location History/LocationHistory.json":
			if err := LFromTakeout(treader); err != nil {
				return err
			}
		case "Takeout/YouTube/history/watch-history.json":
			if err := WFromTakeout(treader); err != nil {
				return err
			}
		default:
			continue
		}
	}

	return
}

// ParseTakeoutZip parses a Google takeout file in zip format.
func ParseTakeoutZip(reader io.Reader) (err error) {
	file, err := ioutil.TempFile("", "googletakeout")
	if err != nil {
		return
	}
	defer file.Close()
	defer os.Remove(file.Name())

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}

	r, err := zip.OpenReader(file.Name())
	if err != nil {
		return err
	}
	defer r.Close()
	if err := database.AddOrganization(org); err != nil {
		return err
	}
	for _, f := range r.File {
		switch f.Name {
		case "Takeout/Location History/LocationHistory.json":
			reader, err := f.Open()
			if err != nil {
				return err
			}
			if err := LFromTakeout(reader); err != nil {
				return err
			}
			continue
		case "Takeout/YouTube/history/watch-history.json":
			reader, err := f.Open()
			if err != nil {
				return err
			}
			if err := WFromTakeout(reader); err != nil {
				return err
			}
			continue
		default:
			continue
		}

	}
	return nil
}
