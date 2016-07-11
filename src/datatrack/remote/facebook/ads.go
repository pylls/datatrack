package facebook

import (
	"datatrack/model"
	"golang.org/x/net/html"
	"io"
	"strconv"
	"time"
)

const (
	adsScanning = iota
	adsHeader
	adsTopics
)

func createReadAds(disclosureChan chan string) readFun {
	return func(reader io.Reader) (out ReadFunOutput, err error) {
		z := html.NewTokenizer(reader)
		state := adsScanning
	outer:
		for {
			tt := z.Next()

			switch tt {
			case html.ErrorToken:
				// End of the document, we're done
				break outer
			case html.StartTagToken:
				tag_, _ := z.TagName()
				tag := string(tag_)
				if state == adsScanning && tag == "h2" {
					state = adsHeader
				}
			case html.EndTagToken:
				tag_, _ := z.TagName()
				tag := string(tag_)
				if tag == "ul" {
					state = adsScanning
				}
			case html.TextToken:
				text := string(z.Text())
				if state == adsHeader {
					if text == "Ads Topics" {
						state = adsTopics
					} else {
						state = adsScanning
					}
				} else if state == adsTopics {
					attr, err := model.MakeAttribute("Ad topic", "cube", text)
					if err != nil {
						return out, err
					}
					out.attributes = append(out.attributes, attr)
				}

				if err != nil {
					return out, err
				}
			}
		}

		disclosure, err := model.MakeDisclosure(org.ID, org.ID,
			strconv.FormatInt(time.Now().Unix()*1000, 10), "", "", "", "")
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

		out.downstreams = append(out.downstreams, model.Downstream{
			Origin: <-disclosureChan,
			Result: disclosure.ID,
		})
		return out, nil
	}
}
