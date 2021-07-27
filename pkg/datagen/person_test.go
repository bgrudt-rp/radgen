package gen

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNewDOB(t *testing.T) {
	n := 0
	for n < 10 {
		dob, err := randomDateByYear(1, 2)
		if err != nil {
			spew.Dump(err)
		}
		spew.Dump(dob)
		n++
	}
}

func TestRandomPerson(t *testing.T) {
	n := 0
	for n < 1 {
		var err error
		var msg Message

		err = GeneratePerson(&msg.Patient)
		if err != nil {
			fmt.Printf("There was an error randomizing the person")
		}
		err = GenerateAddress(&msg.Patient.Address)
		if err != nil {
			fmt.Printf("There was an error randomizing the person address")
		}

		spew.Dump(msg.Patient)

		n++
	}

}
