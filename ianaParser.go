package ianaParser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
)

var getBody *regexp.Regexp

func init() {
	getBody = regexp.MustCompile(`"(.*?)"`)
}

// RootZoneType - type of root zone
type RootZoneType int

// Type of root zone
const (
	Generic RootZoneType = iota
	CountryCode
	Sponsored
)

// RootZone - type of root zone
type RootZone struct {
	Domain                 string
	Type                   RootZoneType
	SponsoringOrganisation string
}

// GetRootZone - get root zone database
func GetRootZone() (rz RootZone, err error) {
	dat, err := ioutil.ReadFile("page")
	if err != nil {

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

		err = ioutil.WriteFile("page", buffer.Bytes(), 0666)
		if err != nil {
			panic(err)
		}
		dat = buffer.Bytes()
	}

	// create a channel for output
	cBlock := make(chan [4]string)

	go func() {
		defer close(cBlock)
		z := html.NewTokenizer(bytes.NewReader(dat))
		insideTable := false
		/* Example of html part
		   <tr>
		       <td>

		           <span class="domain tld"><a href="/domains/root/db/aaa.html">.aaa</a></span></td>

		       <td>generic</td>
		       <td>American Automobile Association, Inc.</td>
		   </tr>
		*/
		insideRow := false
		counter := 0
		var block [4]string
		for {
			tt := z.Next()
			switch tt {
			case html.ErrorToken:
				return
			case html.TextToken:
				if insideTable && insideRow {
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
				if insideTable && atom.Lookup(tag) == atom.Tr {
					if insideRow {
						cBlock <- block
						counter = 0
						insideRow = false
					} else {
						insideRow = true
					}
				}

				if insideTable && insideRow {
					s := strings.TrimSpace(string(z.Raw()))
					if strings.Contains(s, "href") {
						//r := regexp.MustCompile(`"(.*?)"`)
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
		for i := 0; i < 4; i++ {
			for j := 0; j < i; j++ {
				fmt.Printf("\t")
			}
			fmt.Println(c[i])
		}
	}
	return rz, nil
}
