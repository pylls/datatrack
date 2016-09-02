package facebook

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/marcelfarres/datatrack/config"
	"github.com/marcelfarres/datatrack/database"
	"github.com/marcelfarres/datatrack/model"
	"golang.org/x/net/html"
)

var org = model.Organization{
	ID:          "Facebook",
	Name:        "Facebook",
	URL:         "https://www.facebook.com",
	Country:     "United States of America",
	Description: "Facebook Inc. is a technology company. Its main product is a social network with billions of users.",
}

const ProfilePicName string = "profile-facebook.jpg"

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

// ParseDataZip parses a Facebook data file in zip format.
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
	// Send the user's name from the profile parser to the messages parser
	usernameChan := make(chan string, 1)
	// Send the main disclosure ID to the ads parser
	disclosureChan := make(chan string, 1)

	for _, f := range r.File {
		switch f.Name {
		case "photos/profile.jpg":
			wg.Add(1)
			go ReadFile(f, SaveProfilePic, wg, errChan)
		case "index.htm":
			wg.Add(1)
			go ReadFile(f, createReadIndex(usernameChan, disclosureChan), wg, errChan)
		case "html/messages.htm":
			wg.Add(1)
			go ReadFile(f, createReadMessages(usernameChan), wg, errChan)
		case "html/security.htm":
			wg.Add(1)
			go ReadFile(f, ReadSecurity, wg, errChan)
		case "html/ads.htm":
			wg.Add(1)
			go ReadFile(f, createReadAds(disclosureChan), wg, errChan)
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

type readFun func(reader io.Reader) (out ReadFunOutput, err error)

func ReadFile(f *zip.File, fun readFun, wg *sync.WaitGroup, errChan chan error) {
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
	if len(out.coordinates) > 0 {
		wg.Add(1)
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

	return out, err
}

func CheckAttrPresent(z *html.Tokenizer, key string, value string) bool {
	k, v, hasNext := z.TagAttr()
	for hasNext && string(k) != key {
		k, v, hasNext = z.TagAttr()
	}
	return string(k) == key && string(v) == value
}

// Save the locations and IPs you accessed Facebook from and when you did so
func ReadSecurity(reader io.Reader) (out ReadFunOutput, err error) {
	var (
		curAct     ActDate
		curIP      IPAddr
		curLatLong LatLong
	)

	z := html.NewTokenizer(reader)
	state := SecScan

	ipCreation := make(map[time.Time]IPAddr)
	locationCreation := make(map[time.Time]LatLong)
	activities := make(map[IPAddr][]ActDate) // map IP to activities
outer:
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			break outer
		case html.StartTagToken:
			bTag, hasAttr := z.TagName()
			tag := string(bTag)
			if state == SecScan && tag == "h2" {
				state = SecHeader
			} else if state == SecAccountActivity && tag == "li" {
				state = SecItem
			} else if state == SecLoginProtection && tag == "li" {
				state = SecProtectionItem
			} else if hasAttr && tag == "p" &&
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
			bTag, _ := z.TagName()
			tag := string(bTag)
			if tag == "p" {
				if state == SecItemMeta {
					state = SecAccountActivity
					activities[curIP] = append(activities[curIP], curAct)
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
				curAct.Act = text
			case SecItemMeta:
				date, err := time.Parse(
					"Monday, 2 January 2006 at 15:04 UTC-07", text)
				if err == nil {
					curAct.Date = date
					break
				}
				if strings.HasPrefix(text, "IP Address: ") {
					curIP = IPAddr(strings.TrimPrefix(string(text), "IP Address: "))
				} else if strings.HasPrefix(text, "Browser: ") {
					curAct.UserAgent = strings.TrimPrefix(text, "Browser: ")
				}
			case SecProtectionItem:
				if strings.HasPrefix(text, "IP Address: ") {
					curIP = IPAddr(strings.TrimPrefix(string(text), "IP Address: "))
					state = SecIP
				} else if strings.HasPrefix(text, "Estimated location inferred from IP: ") {
					c := strings.Split(strings.TrimPrefix(text,
						"Estimated location inferred from IP: "), ", ")
					curLatLong.Latitude, err = strconv.ParseFloat(c[0], 64)
					if err != nil {
						return out, err
					}
					curLatLong.Longitude, err = strconv.ParseFloat(c[1], 64)
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
						locationCreation[date] = curLatLong
					} else {
						ipCreation[date] = curIP
					}
				}
			}
		}
	}
	for date, ip := range ipCreation {
		var attrs []model.Attribute
		a, err := model.MakeAttribute("IP Address", "flag-checkered", string(ip))
		if err != nil {
			return out, err
		}
		attrs = append(attrs, a)

		if _, actsOk := activities[ip]; actsOk {
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
			strconv.FormatInt(unixTimestamp*1000, 10), "", "", "", "") // ignored error !

		out.disclosures = append(out.disclosures, disclosure)

		if latLong, locOk := locationCreation[date]; locOk {
			self_disclosure, err := model.MakeDisclosure(org.ID, org.ID,
				// unixTimestamp is in seconds, we need milliseconds
				strconv.FormatInt(unixTimestamp*1000, 10), "", "", "", "")
			out.disclosures = append(out.disclosures, self_disclosure)

			coordAttr, err := model.MakeAttribute("Coordinates", "map-marker",
				fmt.Sprintf("%f, %f", latLong.Latitude, latLong.Longitude))
			if err != nil {
				return out, err
			}
			out.attributes = append(out.attributes, coordAttr)

			out.coordinates = append(out.coordinates,
				model.MakeCoordinate(fmt.Sprintf("%f", latLong.Latitude),
					fmt.Sprintf("%f", latLong.Longitude),
					self_disclosure.ID, self_disclosure.Timestamp))
			disclosedDownstream := model.Disclosed{
				Disclosure: self_disclosure.ID,
				Attribute:  []string{coordAttr.ID},
			}
			out.discloseds = append(out.discloseds, disclosedDownstream)

			out.downstreams = append(out.downstreams, model.Downstream{
				Origin: disclosure.ID,
				Result: self_disclosure.ID,
			})

		}

		var attributesS []string
		for _, a := range attrs {
			attributesS = append(attributesS, a.ID)
		}
		disclosed := model.Disclosed{
			Disclosure: disclosure.ID,
			Attribute:  attributesS,
		}
		out.discloseds = append(out.discloseds, disclosed)
	}
	return out, nil
}
