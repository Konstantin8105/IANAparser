package ianaParser

import (
	"fmt"
	"strings"
)

// RootZoneType - type of root zone
type RootZoneType int

// Type of root zone
const (
	Generic RootZoneType = iota
	CountryCode
	Sponsored
	Infrastructure
	GenericRestricted
	Test
)

var convert = [...]struct {
	name string
	tz   RootZoneType
}{
	{"generic", Generic},
	{"country-code", CountryCode},
	{"sponsored", Sponsored},
	{"infrastructure", Infrastructure},
	{"generic-restricted", GenericRestricted},
	{"test", Test},
}

func convertToRootZoneType(t string) (tr RootZoneType, err error) {
	t = strings.ToLower(t)
	for c := range convert {
		if t == convert[c].name {
			return convert[c].tz, nil
		}
	}
	return tr, fmt.Errorf("Cannot convert root zone type: %v", t)
}

func (t RootZoneType) String() string {
	for c := range convert {
		if t == convert[c].tz {
			return convert[c].name
		}
	}
	return ""
}
