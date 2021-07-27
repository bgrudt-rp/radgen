package gen

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRandomCient(t *testing.T) {
	n := 0
	for n < 10 {
		var err error
		var msg Message

		err = GenerateClient(&msg.Client)
		if err != nil {
			fmt.Printf("There was an error generating the client")
		}

		spew.Dump(msg.Client)

		n++
	}
}

func TestStaticFacility(t *testing.T) {
	n := 0
	for n < 10 {
		var err error
		var msg Message
		msg.Client.UQName = "memhosp_uc"

		err = GenerateClient(&msg.Client)
		if err != nil {
			fmt.Printf("There was an error generating the client.")
		}

		spew.Dump(msg.Client)

		n++
	}
}
