package ianaParser

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// URLRootZoneDb - url of web page
	// The Root Zone Database represents the delegation details of top-level domains, including gTLDs such as .com, and country-code TLDs such as .uk. As the manager of the DNS root zone, we are responsible for coordinating these delegations in accordance with our policies and procedures.
	URLRootZoneDb = "https://www.iana.org/domains/root/db"

	// URLIana - main url of IANA website
	URLIana = "https://www.iana.org"
)

var getBody *regexp.Regexp

func init() {
	getBody = regexp.MustCompile(`"(.*?)"`)
}

// RootZone - type of root zone
type RootZone struct {
	Domain                 string
	Type                   RootZoneType
	SponsoringOrganisation string
	URLorganization        string
}

// GetRootZone - get root zone database from website and parsing
func GetRootZone() (rz []RootZone, err error) {
	response, err := http.Get(URLRootZoneDb)
	if err != nil {
		return rz, err
	}

	defer func() {
		err2 := response.Body.Close()
		if err2 != nil {
			if err != nil {
				err = fmt.Errorf("Errors: %v\n%v", err, err2)
			} else {
				err = err2
			}
		}
	}()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, response.Body)
	if err != nil {
		return rz, err
	}
	dat := buffer.Bytes()

	// create a channel for output
	cBlock := make(chan [4]string)

	/* Example of html part
	   <tr>
	       <td><span class="domain tld"><a href="/domains/root/db/aaa.html">.aaa</a></span></td>
	       <td>generic</td>
	       <td>American Automobile Association, Inc.</td>
	   </tr>
	*/
	go func() {
		// close the channel
		defer close(cBlock)

		// analyze html page
		z := html.NewTokenizer(bytes.NewReader(dat))
		insideTable := false
		insideTableBody := false
		insideRow := false
		counter := 0
		var block [4]string
		for {
			tt := z.Next()
			switch tt {
			case html.ErrorToken:
				return
			case html.TextToken:
				if insideTable && insideTableBody && insideRow {
					short := strings.TrimSpace(string(z.Text()))
					if len(short) > 0 && short != "\n" {
						block[counter] = short
						counter++
					}
				}
			case html.StartTagToken, html.EndTagToken:
				tag, _ := z.TagName()
				if atom.Lookup(tag) == atom.Table {
					if insideTable {
						insideTable = false
						continue
					} else {
						insideTable = true
						if !strings.Contains(string(z.Raw()), "iana-table") {
							insideTable = false
						}
					}
				}
				if insideTable && atom.Lookup(tag) == atom.Tbody {
					if insideTableBody {
						insideTableBody = false
					} else {
						insideTableBody = true
					}
				}
				if insideTable && insideTableBody && atom.Lookup(tag) == atom.Tr {
					if insideRow {
						cBlock <- block
						counter = 0
						insideRow = false
					} else {
						insideRow = true
					}
				}

				if insideTable && insideTableBody && insideRow {
					s := strings.TrimSpace(string(z.Raw()))
					if strings.Contains(s, "href") {
						result := getBody.FindAllString(s, -1)
						// result : "/domains/root/db/xn--mgbayh7gpa.html"
						block[counter] = result[0][1 : len(result[0])-1]
						counter++
					}
				}

			}
		}
	}()

	for c := range cBlock {
		tz, err := convertToRootZoneType(c[2])
		if err != nil {
			return rz, err
		}
		rz = append(rz, RootZone{
			Domain: c[1],
			Type:   tz,
			SponsoringOrganisation: c[3],
			URLorganization:        URLIana + c[0],
		})
	}
	return rz, nil
}
