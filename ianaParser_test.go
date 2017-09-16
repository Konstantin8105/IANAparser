package ianaParser_test

import (
	"testing"

	"github.com/Konstantin8105/IANAparser"
)

func TestSuccess(t *testing.T) {
	rz, err := ianaParser.GetRootZone()
	if err != nil {
		t.Error(err)
	}
	for _, r := range rz {
		if r.Domain[0] != '.' && ([]rune(r.Domain))[1] != '.' {
			t.Errorf("All domain must start from point. See %v\nRootZone: %v", r.Domain[0], r)
		}
		if len(r.URLorganization) == 0 {
			t.Errorf("Each root domain must have orgnization. See %v\nRootZone: %v", r.URLorganization, r)
		}
	}
}
