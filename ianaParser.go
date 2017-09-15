package ianaParser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// URLRootZoneDb - url of web page
	// The Root Zone Database represents the delegation details of top-level domains, including gTLDs such as .com, and country-code TLDs such as .uk. As the manager of the DNS root zone, we are responsible for coordinating these delegations in accordance with our policies and procedures.
	URLRootZoneDb = "https://www.iana.org/domains/root/db"
)

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

	z := html.NewTokenizer(bytes.NewReader(dat))
	insideTable := false
	/*
			Example of html part
		    <tr>
		        <td>

		            <span class="domain tld"><a href="/domains/root/db/aaa.html">.aaa</a></span></td>

		        <td>generic</td>
		        <td>American Automobile Association, Inc.</td>
		    </tr>
	*/
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			goto END
		case html.TextToken:
			/*
				short := strings.TrimSpace(string(z.Text()))
				if len(short) < 1 {
					continue
				}
				if short == "\n" {
					continue
				}*/
			if insideTable {
				fmt.Printf("text --> !%s!\n\n", z.Text())
			}
		case html.StartTagToken, html.EndTagToken:
			tag, hasAttr := z.TagName()
			if atom.Lookup(tag) == atom.Table {
				if insideTable {
					insideTable = false
					continue
				} else {
					insideTable = true
					r := z.Raw()
					if !strings.Contains(string(r), "iana-table") {
						insideTable = false
					}
				}
			}
			if insideTable {
				fmt.Printf("tag1 --> tag %s\n", tag)
				fmt.Printf("tag2 --> hasAttr %v\n", hasAttr)
				c, d, e := z.TagAttr()
				fmt.Printf("tag4 --> c %s\n", c)
				fmt.Printf("tag5 --> d %s\n", d)
				fmt.Printf("tag6 --> e %v\n", e)
				g := z.Text()
				fmt.Printf("tag7 --> g %v\n", g)
				r := z.Raw()
				fmt.Printf("tag8 --> raw %s\n", r)
			}
		}
	}
END:
	return rz, nil
}
