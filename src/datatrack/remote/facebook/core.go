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
	"strconv"
	"strings"
	"sync"
	"time"
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

	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)

	for _, f := range r.File {
		switch f.Name {
		case "photos/profile.jpg":
			wg.Add(1)
			go ReadFile(f, SaveProfilePic, wg, errChan)
		case "index.htm":
			wg.Add(1)
			go ReadFile(f, ReadIndex, wg, errChan)
		case "html/messages.htm":
			wg.Add(1)
			go ReadFile(f, ReadMessages, wg, errChan)
		default:
			continue
		}
	}
	if err := database.AddOrganization(org); err != nil {
		return err
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}
	return nil
}

type ReadFunOutput struct {
	attributes  []model.Attribute
	disclosures []model.Disclosure
	discloseds  []model.Disclosed
	downstreams []model.Downstream
}

type read_func func(reader io.Reader) (out ReadFunOutput, err error)

func ReadFile(f *zip.File, fun read_func, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	reader, err := f.Open()
	if err != nil {
		errChan <- err
		return
	}
	out, err := fun(reader)
	if err != nil {
		errChan <- err
		return
	}
	if len(out.attributes) > 0 {
		wg.Add(1)
		go database.AddAttributes(out.attributes, wg, errChan)
	}
	if len(out.disclosures) > 0 {
		wg.Add(1)
		go database.AddDisclosures(out.disclosures, wg, errChan)
	}
	if len(out.discloseds) > 0 {
		wg.Add(1)
		go database.AddDiscloseds(out.discloseds, wg, errChan)
	}
	if len(out.downstreams) > 0 {
		wg.Add(1)
		go database.AddDownstreams(out.downstreams, wg, errChan)
	}
}

func SaveProfilePic(reader io.Reader) (out ReadFunOutput, err error) {
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return out, err
	}
	err = ioutil.WriteFile(config.StaticPath+"/img/"+ProfilePicName,
		contents, 0644)
	if err != nil {
		return out, err
	}
	return out, nil
}

const (
	Scanning = iota
	SName    // S is for Scanning
	SHeader
	SData
	SParseDate
)

func ReadIndex(reader io.Reader) (out ReadFunOutput, err error) {
	z := html.NewTokenizer(reader)
	state := Scanning
	var header string
	var data = make([]string, 0)
	var unixTimestamp int64
outer:
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			break outer
		case html.StartTagToken:
			if state == Scanning {
				tag_, hasattrs := z.TagName()
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
				case "div":
					if hasattrs {
						key, value, _ := z.TagAttr()
						if string(key) == "class" && string(value) == "footer" {
							state = SParseDate
						}
					}
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
					return out, err
				}
				out.attributes = append(out.attributes, attrs...)
			} else if state == SParseDate && tag == "div" {
				state = Scanning
			}

		case html.TextToken:
			switch state {
			case SName:
				err = database.SetUser(model.User{
					Name:    string(z.Text()),
					Picture: ProfilePicName,
				})
			case SHeader:
				if header != "" {
					return out, errors.New("Found two data in a row")
				}
				header = string(z.Text())
			case SData:
				data = append(data, string(z.Text()))
			case SParseDate:
				stringDate := strings.Split(string(z.Text()), ", ")
				var date time.Time
				date, err = time.Parse("2 Jan 2006 at 15:04 UTC-07",
					stringDate[len(stringDate)-1])
				// The user interface is not timezone-aware, so we strip the timezone
				_, offset := date.Zone()
				unixTimestamp = date.Unix() + int64(offset)
			}
		}

		if err != nil {
			return out, err
		}
	}

	disclosure, err := model.MakeDisclosure(database.Self, org.ID,
		// unixTimestamp is in seconds, we need milliseconds
		strconv.FormatInt(unixTimestamp*1000, 10), "", "", "", "")
	if err != nil {
		return out, err
	}
	out.disclosures = append(out.disclosures, disclosure)

	var attributes_s []string
	for _, a := range out.attributes {
		attributes_s = append(attributes_s, a.ID)
	}
	disclosed := model.Disclosed{
		Disclosure: disclosure.ID,
		Attribute:  attributes_s,
	}
	out.discloseds = append(out.discloseds, disclosed)
	return out, nil
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

const (
	MsgScan = iota
	MsgThread
	MsgTimestamp
)

// Save the earliest and latest time you spoke with a user
func ReadMessages(reader io.Reader) (out ReadFunOutput, err error) {
	z := html.NewTokenizer(reader)
	state := MsgScan
	var user_name string
	// Busy wait, but it's clean code
	for user_name == "" {
		user, err := database.GetUser()
		if err != nil {
			time.Sleep(20 * time.Millisecond)
		} else {
			user_name = user.Name
		}
	}

	first_msg := make(map[string]time.Time)
	last_msg := make(map[string]time.Time)
	num_msg := make(map[string]int64)
	var current_names []string
outer:
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			break outer
		case html.StartTagToken:
			if state == MsgScan {
				tag_, hasattrs := z.TagName()
				tag := string(tag_)
				if hasattrs {
					if tag == "div" {
						key, value, hasnext := z.TagAttr()
						for hasnext && string(key) != "class" {
							key, value, hasnext = z.TagAttr()
						}
						if string(key) == "class" && string(value) == "thread" {
							state = MsgThread
						}
					} else if tag == "span" {
						key, value, hasnext := z.TagAttr()
						for hasnext && string(key) != "class" {
							key, value, hasnext = z.TagAttr()
						}
						if string(key) == "class" && string(value) == "meta" {
							state = MsgTimestamp
						}
					}
				}
			}
		case html.TextToken:
			switch state {
			case MsgThread:
				state = MsgScan
				names := strings.Split(string(z.Text()), ", ")
				current_names = nil
				contains_user_name := false
				for _, n := range names {
					if n == user_name {
						contains_user_name = true
					} else {
						current_names = append(current_names, n)
					}
				}
				if !contains_user_name {
					current_names = nil
				}
				for _, n := range current_names {
					if _, ok := first_msg[n]; !ok {
						// last possible time
						first_msg[n] = time.Now()
						// first possible time
						last_msg[n] = time.Time{}
						num_msg[n] = 0
					}
				}
			case MsgTimestamp:
				state = MsgScan
				date, err := time.Parse("Monday, 2 January 2006 at 15:04 UTC-07",
					string(z.Text()))
				if err != nil {
					return out, err
				}
				for _, n := range current_names {
					if date.Before(first_msg[n]) {
						first_msg[n] = date
					}
					if date.After(last_msg[n]) {
						last_msg[n] = date
					}
					num_msg[n]++
				}
			}
		}

		if err != nil {
			return out, err
		}
	}

	for name, date := range last_msg {
		_, offset := date.Zone()
		unixTimestamp := date.Unix() + int64(offset)
		disclosure, err := model.MakeDisclosure(database.Self, org.ID,
			// unixTimestamp is in seconds, we need milliseconds
			strconv.FormatInt(unixTimestamp*1000, 10), "", "", "", "")
		if err != nil {
			return out, err
		}
		out.disclosures = append(out.disclosures, disclosure)

		first, err := model.MakeAttribute("First Message", "comments-o",
			first_msg[name].Format("2 Jan 2006 at 15:04 UTC-07"))
		if err != nil {
			return out, err
		}
		last, err := model.MakeAttribute("Last Message", "comments-o",
			date.Format("2 Jan 2006 at 15:04 UTC-07"))
		if err != nil {
			return out, err
		}
		user, err := model.MakeAttribute("Recipient", "comments-o", name)
		if err != nil {
			return out, err
		}
		nmsg, err := model.MakeAttribute("Number of Messages", "comments-o",
			strconv.FormatInt(num_msg[name], 10))
		if err != nil {
			return out, err
		}
		out.attributes = append(out.attributes, first, last, user, nmsg)
		disclosed := model.Disclosed{
			Disclosure: disclosure.ID,
			Attribute:  []string{user.ID, nmsg.ID, first.ID, last.ID},
		}
		out.discloseds = append(out.discloseds, disclosed)
	}
	return out, nil
}
