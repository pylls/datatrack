package facebook

import (
	"archive/zip"
	"datatrack/config"
	"datatrack/database"
	"datatrack/model"
	"errors"
	"fmt"
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

			err = database.SetUser(model.User{
				Name:    "Adria Garriga",
				Picture: ProfilePicName,
			})
		case "index.htm":
			wg.Add(1)
			go ReadFile(f, ReadIndex, wg, errChan)
		case "html/messages.htm":
			wg.Add(1)
			go ReadFile(f, ReadMessages, wg, errChan)
		case "html/security.htm":
			wg.Add(1)
			go ReadFile(f, ReadSecurity, wg, errChan)
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
	coordinates []model.Coordinate
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
		fmt.Printf("For func %v: %d attributes\n", fun, len(out.attributes))
		go database.AddAttributes(out.attributes, wg, errChan)
	}
	if len(out.disclosures) > 0 {
		wg.Add(1)
		fmt.Printf("For func %v: %d disclosures\n", fun, len(out.attributes))
		go database.AddDisclosures(out.disclosures, wg, errChan)
	}
	if len(out.discloseds) > 0 {
		wg.Add(1)
		fmt.Printf("For func %v: %d discloseds\n", fun, len(out.attributes))
		go database.AddDiscloseds(out.discloseds, wg, errChan)
	}
	if len(out.downstreams) > 0 {
		wg.Add(1)
		fmt.Printf("For func %v: %d downstreams\n", fun, len(out.attributes))
		go database.AddDownstreams(out.downstreams, wg, errChan)
	}
	if len(out.coordinates) > 0 {
		wg.Add(1)
		fmt.Printf("For func %v: %d coordinates\n", fun, len(out.coordinates))
		go database.AddCoordinates(out.coordinates, wg, errChan)
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

func CheckAttrPresent(z *html.Tokenizer, key string, value string) bool {
	k, v, hasnext := z.TagAttr()
	for hasnext && string(k) != key {
		k, v, hasnext = z.TagAttr()
	}
	return string(k) == key && string(v) == value
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
						if CheckAttrPresent(z, "class", "thread") {
							state = MsgThread
						}
					} else if tag == "span" {
						if CheckAttrPresent(z, "class", "meta") {
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
				date, err := time.Parse(
					"Monday, 2 January 2006 at 15:04 UTC-07", string(z.Text()))
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

const (
	SecScan = iota
	SecHeader
	SecAccountActivity
	SecItem
	SecItemMeta
	SecLoginProtection
	SecProtectionItem
	SecIP
	SecLocation
	SecIPMeta
	SecLocationMeta
)

type LatLong struct {
	Latitude  float64
	Longitude float64
}
type ActDate struct {
	Act       string
	Date      time.Time
	UserAgent string
}
type IPAddr string

// Save the locations and IPs you accessed Facebook from and when you did so
func ReadSecurity(reader io.Reader) (out ReadFunOutput, err error) {
	z := html.NewTokenizer(reader)
	state := SecScan

	ip_creation := make(map[time.Time]IPAddr)
	location_creation := make(map[time.Time]LatLong)
	activities := make(map[IPAddr][]ActDate) // map IP to activities
	var cur_act ActDate
	var cur_ip IPAddr
	var cur_latlong LatLong
outer:
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			break outer
		case html.StartTagToken:
			tag_, hasattrs := z.TagName()
			tag := string(tag_)
			if state == SecScan && tag == "h2" {
				state = SecHeader
			} else if state == SecAccountActivity && tag == "li" {
				state = SecItem
			} else if state == SecLoginProtection && tag == "li" {
				state = SecProtectionItem
			} else if hasattrs && tag == "p" &&
				CheckAttrPresent(z, "class", "meta") {
				switch state {
				case SecItem:
					state = SecItemMeta
				case SecIP:
					state = SecIPMeta
				case SecLocation:
					state = SecLocationMeta
				}
			}
		case html.EndTagToken:
			tag_, _ := z.TagName()
			tag := string(tag_)
			if tag == "p" {
				if state == SecItemMeta {
					state = SecAccountActivity
					activities[cur_ip] = append(activities[cur_ip], cur_act)
				} else if state == SecIPMeta || state == SecLocationMeta {
					state = SecProtectionItem
				}
			} else if tag == "ul" {
				state = SecScan
			}
		case html.TextToken:
			text := string(z.Text())
			switch state {
			case SecHeader:
				if text == "Account activity" {
					state = SecAccountActivity
				} else if text == "Login Protection Data" {
					state = SecLoginProtection
				}
			case SecItem:
				cur_act.Act = text
			case SecItemMeta:
				date, err := time.Parse(
					"Monday, 2 January 2006 at 15:04 UTC-07", text)
				if err == nil {
					cur_act.Date = date
					break
				}
				if strings.HasPrefix(text, "IP Address: ") {
					cur_ip = IPAddr(strings.TrimPrefix(string(text), "IP Address: "))
				} else if strings.HasPrefix(text, "Browser: ") {
					cur_act.UserAgent = strings.TrimPrefix(text, "Browser: ")
				}
			case SecProtectionItem:
				if strings.HasPrefix(text, "IP Address: ") {
					cur_ip = IPAddr(strings.TrimPrefix(string(text), "IP Address: "))
					state = SecIP
				} else if strings.HasPrefix(text, "Estimated location inferred from IP: ") {
					c := strings.Split(strings.TrimPrefix(text,
						"Estimated location inferred from IP: "), ", ")
					cur_latlong.Latitude, err = strconv.ParseFloat(c[0], 64)
					if err != nil {
						return out, err
					}
					cur_latlong.Longitude, err = strconv.ParseFloat(c[1], 64)
					if err != nil {
						return out, err
					}
					state = SecLocation
				} else {
					state = SecLoginProtection
				}
			case SecIPMeta, SecLocationMeta:
				var date time.Time
				filled := false
				if strings.HasPrefix(text, "Created: ") {
					s := strings.TrimPrefix(text, "Created: ")
					date, err = time.Parse(
						"Monday, 2 January 2006 at 15:04 UTC-07", s)
					filled = true
				} else if strings.HasPrefix(text, "Updated: ") {
					s := strings.TrimPrefix(text, "Updated: ")
					date, err = time.Parse(
						"Monday, 2 January 2006 at 15:04 UTC-07", s)
					filled = true
				}
				if filled {
					if err != nil {
						return out, err
					}
					if state == SecLocationMeta {
						location_creation[date] = cur_latlong
					} else {
						ip_creation[date] = cur_ip
					}
				}
			}
		}
	}
	for date, ip := range ip_creation {
		var attrs []model.Attribute
		a, err := model.MakeAttribute("IP Address", "flag-checkered", string(ip))
		if err != nil {
			return out, err
		}
		attrs = append(attrs, a)
		_, acts_ok := activities[ip]
		if acts_ok {
			for _, action := range activities[ip] {
				a, err := model.MakeAttribute(action.Act, "gear",
					action.Date.Format("2 Jan 2006 at 15:04 UTC-07")+"; "+action.UserAgent)
				if err != nil {
					return out, err
				}
				attrs = append(attrs, a)
			}
		}
		out.attributes = append(out.attributes, attrs...)

		_, offset := date.Zone()
		unixTimestamp := date.Unix() + int64(offset)
		disclosure, err := model.MakeDisclosure(database.Self, org.ID,
			// unixTimestamp is in seconds, we need milliseconds
			strconv.FormatInt(unixTimestamp*1000, 10), "", "", "", "")

		out.disclosures = append(out.disclosures, disclosure)
		latlong, loc_ok := location_creation[date]
		if loc_ok {
			//			self_disclosure, err := model.MakeDisclosure(org.ID, org.ID,
			// unixTimestamp is in seconds, we need milliseconds
			//				strconv.FormatInt(unixTimestamp*1000, 10), "", "", "", "")
			//			out.disclosures = append(out.disclosures, self_disclosure)
			self_disclosure := disclosure

			coord_attr, err := model.MakeAttribute("Coordinates", "map-marker",
				fmt.Sprintf("%f, %f", latlong.Latitude, latlong.Longitude))
			if err != nil {
				return out, err
			}
			out.attributes = append(out.attributes, coord_attr)

			out.coordinates = append(out.coordinates,
				model.MakeCoordinate(fmt.Sprintf("%f", latlong.Latitude),
					fmt.Sprintf("%f", latlong.Longitude),
					self_disclosure.ID, self_disclosure.Timestamp))
			disclosed_downstream := model.Disclosed{
				Disclosure: self_disclosure.ID,
				Attribute:  []string{coord_attr.ID},
			}
			out.discloseds = append(out.discloseds, disclosed_downstream)

			/*			out.downstreams = append(out.downstreams, model.Downstream{
						Origin: disclosure.ID,
						Result: self_disclosure.ID,
					})*/

		}

		var attributes_s []string
		for _, a := range attrs {
			attributes_s = append(attributes_s, a.ID)
		}
		disclosed := model.Disclosed{
			Disclosure: disclosure.ID,
			Attribute:  attributes_s,
		}
		out.discloseds = append(out.discloseds, disclosed)
	}

	fmt.Printf("len(ip_creation) = %d\n", len(ip_creation))
	fmt.Printf("len(location_creation) = %d\n", len(location_creation))
	fmt.Printf("len(activities) = %d\n", len(activities))
	s := 0
	for _, a := range activities {
		s += len(a)
	}
	fmt.Printf("len(len(activities)) = %d\n", s)
	fmt.Printf("len(out.attributes) = %d\n", len(out.attributes))
	fmt.Printf("len(out.coordinates) = %d\n", len(out.coordinates))
	fmt.Printf("len(out.discloseds) = %d\n", len(out.discloseds))
	fmt.Printf("len(out.disclosures) = %d\n", len(out.disclosures))
	fmt.Printf("len(out.downstreams) = %d\n", len(out.downstreams))
	return out, nil
}
