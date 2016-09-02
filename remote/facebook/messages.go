package facebook

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/marcelfarres/datatrack/database"
	"github.com/marcelfarres/datatrack/model"
	"golang.org/x/net/html"
)

const (
	msgScan = iota
	msgThread
	msgTimestamp
)

// Save the earliest and latest time you spoke with a user
func createReadMessages(usernameChan chan string) readFun {
	return func(reader io.Reader) (out ReadFunOutput, err error) {
		var currentNames []string

		z := html.NewTokenizer(reader)
		state := msgScan
		userName := <-usernameChan

		firstMsg := make(map[string]time.Time)
		lastMsg := make(map[string]time.Time)
		numMsg := make(map[string]int64)
	outer:
		for {
			tt := z.Next()

			switch tt {
			case html.ErrorToken:
				// End of the document, we're done
				break outer
			case html.StartTagToken:
				if state == msgScan {
					bTag, hasAttr := z.TagName()
					tag := string(bTag)
					if hasAttr {
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
					currentNames = nil
					containsUserName := false
					for _, n := range names {
						if n == userName {
							containsUserName = true
						} else {
							currentNames = append(currentNames, n)
						}
					}
					if !containsUserName {
						currentNames = nil
					}
					for _, n := range currentNames {
						if _, ok := firstMsg[n]; !ok {
							// last possible time
							firstMsg[n] = time.Now()
							// first possible time
							lastMsg[n] = time.Time{}
							numMsg[n] = 0
						}
					}
				case msgTimestamp:
					state = msgScan
					date, err := time.Parse(
						"Monday, 2 January 2006 at 15:04 UTC-07", string(z.Text()))
					if err != nil {
						return out, err
					}
					for _, n := range currentNames {
						if date.Before(firstMsg[n]) {
							firstMsg[n] = date
						}
						if date.After(lastMsg[n]) {
							lastMsg[n] = date
						}
						numMsg[n]++
					}
				}
			}

			if err != nil {
				return out, err
			}
		}

		for name, date := range lastMsg {
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
				firstMsg[name].Format("2 Jan 2006 at 15:04 UTC-07"))
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
			nMsg, err := model.MakeAttribute("Number of Messages", "comments-o",
				strconv.FormatInt(numMsg[name], 10))
			if err != nil {
				return out, err
			}
			out.attributes = append(out.attributes, first, last, user, nMsg)
			disclosed := model.Disclosed{
				Disclosure: disclosure.ID,
				Attribute:  []string{user.ID, nMsg.ID, first.ID, last.ID},
			}
			out.discloseds = append(out.discloseds, disclosed)
		}
		return out, nil
	}
}
