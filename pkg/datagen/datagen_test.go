package gen

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRandomizer(t *testing.T) {
	n := 0
	for n < 11 {
		var err error
		var msg Message

		err = GenerateClient(config.Cfg, &msg)
		if err != nil {
			fmt.Printf("There was an error generating the facility")
		}

		err = GenerateContext(&msg)
		if err != nil {
			fmt.Printf("There was an error generating the context")
		}

		err = RandomizePatient(&msg)
		if err != nil {
			fmt.Printf("There was an error randomizing the patient")
		}

		err = RandomizeVisit(&msg)
		if err != nil {
			fmt.Printf("There was an error randomizing the visit")
		}

		spew.Dump(msg)

		n++
	}

}
