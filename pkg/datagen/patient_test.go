package gen

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRandomizer(t *testing.T) {
	n := 0
	for n < 20 {
		var err error
		var msg Message
		msg.Facility.UQName = "MemHospUC"

		err = GenerateContext(&msg)
		if err != nil {
			fmt.Printf("There was an error generating the context")
		}

		err = RandomizePatient(&msg)
		if err != nil {
			fmt.Printf("There was an error randomizing the patient.")
		}
		spew.Dump(msg)

		n++
	}

}
