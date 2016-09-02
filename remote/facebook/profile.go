package facebook

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pylls/datatrack/database"
	"github.com/pylls/datatrack/model"
	"golang.org/x/net/html"
)

const (
	Scanning = iota
	SName    // S is for Scanning
	SHeader
	SData
	SParseDate
)

func createReadIndex(usernameChan chan string, disclosureChan chan string) readFun {
	return func(reader io.Reader) (out ReadFunOutput, err error) {
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
					attrs, err := saveScannedAttributes(header, data)
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
					name := string(z.Text())
					err = database.SetUser(model.User{
						Name:    name,
						Picture: ProfilePicName,
					})
					usernameChan <- name
					close(usernameChan)
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
		disclosureChan <- disclosure.ID
		close(disclosureChan)

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
}

func saveScannedAttributes(header string, data []string) (attrs []model.Attribute, err error) {
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
		data = split(data, " - ")
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
		data = split(data, ", ")
	case "Activities":
		category = "child" // Taken from datatrack/remote/google
		data = split(data, ", ")
	case "Interests":
		category = "smile-o" // invented
		data = split(data, ", ")
	case "Music":
		category = "music" // invented
		data = split(data, ", ")
	case "Books":
		category = "book" // invented
		data = split(data, ", ")
	case "Movies":
		category = "film" // invented
		data = split(data, ", ")
	case "Television":
		category = "television" // invented
		data = split(data, ", ")
	case "Games":
		category = "gamepad" // invented
		data = split(data, ", ")
	case "Other":
		category = "thumbs-o-up" // invented
		data = split(data, ", ")
	case "Favourite sports":
		category = "soccer-ball-o" // invented
		data = split(data, ", ")
	case "Favourite teams":
		category = "soccer-ball-o" // invented
		data = split(data, ", ")
	case "Favourite athletes":
		category = "soccer-ball-o" // invented
		data = split(data, ", ")
	case "Groups":
		category = "group" // invented
		data = split(data, ", ")
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

func split(in []string, split string) (out []string) {
	for _, d := range in {
		a := strings.Split(d, split)
		out = append(out, a...)
	}
	return out
}
