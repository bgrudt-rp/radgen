package gen

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRandomEncounter(t *testing.T) {
	n := 0
	for n < 1 {
		var err error
		var msg Message

		err = GenerateClient(&msg.Client)
		if err != nil {
			fmt.Printf("There was an error generating the client")
		}

		err = GeneratePerson(&msg.Patient)
		if err != nil {
			fmt.Printf("There was an error randomizing the patient")
		}

		err = GenerateAddress(&msg.Patient.Address)
		if err != nil {
			fmt.Printf("There was an error randomizing the patient address")
		}

		err = GenerateVisit(&msg.Visit)
		if err != nil {
			fmt.Printf("There was an error randomizing the visit")
		}

		spew.Dump(msg)

		n++
	}

}

func TestContextualEncounter(t *testing.T) {
	n := 0
	for n < 7 {
		var err error
		var msg Message

		err = GenerateClient(&msg.Client)
		if err != nil {
			fmt.Printf("There was an error generating the client")
		}

		err = GenerateContext(&msg)
		if err != nil {
			fmt.Printf("There was an error generating the context")
		}

		err = GeneratePerson(&msg.Patient)
		if err != nil {
			fmt.Printf("There was an error randomizing the patient")
		}

		err = GenerateAddress(&msg.Patient.Address)
		if err != nil {
			fmt.Printf("There was an error randomizing the patient address")
		}

		err = GenerateVisit(&msg.Visit)
		if err != nil {
			fmt.Printf("There was an error randomizing the visit")
		}

		spew.Dump(msg)

		n++
	}
}
