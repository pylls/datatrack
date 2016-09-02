package facebook

import (
	"io"
	"strconv"
	"time"

	"github.com/pylls/datatrack/model"
	"golang.org/x/net/html"
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
				bTag, _ := z.TagName()
				tag := string(bTag)
				if state == adsScanning && tag == "h2" {
					state = adsHeader
				}
			case html.EndTagToken:
				bTag, _ := z.TagName()
				tag := string(bTag)
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

		var attributesS []string
		for _, a := range out.attributes {
			attributesS = append(attributesS, a.ID)
		}
		disclosed := model.Disclosed{
			Disclosure: disclosure.ID,
			Attribute:  attributesS,
		}
		out.discloseds = append(out.discloseds, disclosed)

		out.downstreams = append(out.downstreams, model.Downstream{
			Origin: <-disclosureChan,
			Result: disclosure.ID,
		})
		return out, nil
	}
}
