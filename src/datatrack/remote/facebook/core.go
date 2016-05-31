package facebook

import (
	"archive/zip"
	"datatrack/config"
	"datatrack/database"
	"datatrack/model"
	"errors"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var org = model.Organization{
	ID:          "Facebook",
	Name:        "Facebook",
	URL:         "https://www.facebook.com",
	Country:     "United States of America",
	Description: "Facebook Inc. is a technology company. Its main product is a social network with billions of users.",
}

const ProfilePicName string = "profile-facebook.jpg"

// parses a Facebook data file in zip format.
func ParseDataZip(reader io.Reader) (err error) {
	file, err := ioutil.TempFile("", "facebookdata")
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
	// TEMPORARY

	disclosure, err := model.MakeDisclosure(database.Self, org.ID, "",
		"", "", "", "")
	if err != nil {
		return err
	}
	err = database.AddDisclosure(disclosure)
	if err != nil {
		return err
	}
	var attributes []model.Attribute

	for _, f := range r.File {
		switch f.Name {
		case "photos/profile.jpg":
			reader, err := f.Open()
			if err != nil {
				return err
			}
			contents, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(config.StaticPath+"/img/"+ProfilePicName,
				contents, 0644)
			if err != nil {
				return err
			}

		case "index.htm":
			reader, err := f.Open()
			if err != nil {
				return err
			}
			attributes, err = ReadIndex(reader)
			if err != nil {
				return err
			}
		default:
			continue
		}
	}
	var attributes_s []string
	for _, a := range attributes {
		attributes_s = append(attributes_s, a.ID)
	}

	disclosed := model.Disclosed{
		Disclosure: disclosure.ID,
		Attribute:  attributes_s,
	}
	err = database.AddDisclosed(disclosed)
	if err != nil {
		return err
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	errChan := make(chan error, 1)
	go database.AddAttributes(attributes, wg, errChan)
	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}

	return nil
}

const (
	Scanning = iota
	SName    // S is for Scanning
	SHeader
	SData
)

func ReadIndex(reader io.Reader) (attributes []model.Attribute, err error) {
	z := html.NewTokenizer(reader)
	state := Scanning
	var header string
	var data = make([]string, 0)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			return attributes, nil
		case html.StartTagToken:
			if state == Scanning {
				tag_, _ := z.TagName()
				tag := string(tag_)
				switch tag {
				case "h1":
					state = SName
				case "th":
					state = SHeader
					header = ""
				case "td":
					state = SData
					data = data[:0]
				}
			}
		case html.EndTagToken:
			tag_, _ := z.TagName()
			tag := string(tag_)
			if state == SName && tag == "h1" {
				state = Scanning
			} else if state == SHeader && tag == "th" {
				state = Scanning
			} else if state == SData && tag == "td" {
				state = Scanning
				attrs, err := SaveScannedAttributes(header, data)
				if err != nil {
					return nil, err
				}
				attributes = append(attributes, attrs...)
			}

		case html.TextToken:
			if state == SName {
				err = database.SetUser(model.User{
					Name:    string(z.Text()),
					Picture: ProfilePicName,
				})
			} else if state == SHeader {
				if header != "" {
					return nil, errors.New("Found two data in a row")
				}
				header = string(z.Text())
			} else if state == SData {
				data = append(data, string(z.Text()))
			}
		}

		if err != nil {
			return nil, err
		}
	}
}

func Split(in []string, split string) (out []string) {
	for _, d := range in {
		a := strings.Split(d, split)
		out = append(out, a...)
	}
	return out
}

func SaveScannedAttributes(header string, data []string) (attrs []model.Attribute, err error) {
	var category string
	switch header {
	case "Profile":
		category = "user"
	case "Email address":
		category = "envelope"
	case "Registration Date":
		category = "calendar-o"
	case "Birthday":
		category = "birthday-cake"
	case "Gender":
		category = "venus-mars"
	case "Previous Names":
		category = "male"
		data = Split(data, " - ")
	case "Current location":
		category = "map-marker"
	case "Home Town":
		category = "globe"
	case "Family":
		category = "tree" // invented
	case "Education":
		category = "university"
	case "Spoken Languages":
		category = "language"
		data = Split(data, ", ")
	case "Activities":
		category = "child" // Taken from datatrack/remote/google
		data = Split(data, ", ")
	case "Interests":
		category = "smile-o" // invented
		data = Split(data, ", ")
	case "Music":
		category = "music" // invented
		data = Split(data, ", ")
	case "Books":
		category = "book" // invented
		data = Split(data, ", ")
	case "Movies":
		category = "film" // invented
		data = Split(data, ", ")
	case "Television":
		category = "television" // invented
		data = Split(data, ", ")
	case "Games":
		category = "gamepad" // invented
		data = Split(data, ", ")
	case "Other":
		category = "thumbs-o-up" // invented
		data = Split(data, ", ")
	case "Favourite sports":
		category = "soccer-ball-o" // invented
		data = Split(data, ", ")
	case "Favourite teams":
		category = "soccer-ball-o" // invented
		data = Split(data, ", ")
	case "Favourite athletes":
		category = "soccer-ball-o" // invented
		data = Split(data, ", ")
	case "Groups":
		category = "group" // invented
		data = Split(data, ", ")
	case "Networks":
		category = "share-alt" // invented
	case "Apps":
		category = "laptop" // invented
	case "Pages You Admin":
		category = "cog" // invented
	default:
		return nil, errors.New("Unexpected category: \"" + header + "\"")
	}
	for _, v := range data {
		attr, err := model.MakeAttribute(header, category, v)
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, attr)
	}
	return attrs, nil
}
