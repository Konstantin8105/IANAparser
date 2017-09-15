package ianaParser

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
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

	z := html.NewTokenizer(bytes.NewReader(buffer.Bytes()))
	insideTable := false
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			goto END
		case html.TextToken:
			if insideTable {
				fmt.Printf("text --> %s\n", z.Text())
			}
		case html.StartTagToken, html.EndTagToken:
			a, b := z.TagName()
			if string(a) != "table" {
				continue
			}
			if insideTable {
				fmt.Printf("tag --> a %s\n", a)
				fmt.Printf("tag --> b %v\n", b)
			}
			if insideTable {
				insideTable = false
			} else {
				insideTable = true
			}
		}
	}
END:
	return rz, nil
}
