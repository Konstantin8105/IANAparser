package ianaParser_test

import (
	"fmt"
	"testing"

	"github.com/Konstantin8105/IANAparser"
)

func TestSuccess(t *testing.T) {
	rz, err := ianaParser.GetRootZone()

	fmt.Println("rz  = ", rz)
	fmt.Println("err = ", err)
}
