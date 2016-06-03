package facebook

import (
	"datatrack/database"
	"datatrack/model"
	"golang.org/x/net/html"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	msgScan = iota
	msgThread
	msgTimestamp
)

// Save the earliest and latest time you spoke with a user
func createReadMessages(usernameChan chan string) readFun {
	return func(reader io.Reader) (out ReadFunOutput, err error) {
		z := html.NewTokenizer(reader)
		state := msgScan
		user_name := <-usernameChan

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
				if state == msgScan {
					tag_, hasattrs := z.TagName()
					tag := string(tag_)
					if hasattrs {
						if tag == "div" {
							if CheckAttrPresent(z, "class", "thread") {
								state = msgThread
							}
						} else if tag == "span" {
							if CheckAttrPresent(z, "class", "meta") {
								state = msgTimestamp
							}
						}
					}
				}
			case html.TextToken:
				switch state {
				case msgThread:
					state = msgScan
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
				case msgTimestamp:
					state = msgScan
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
}
