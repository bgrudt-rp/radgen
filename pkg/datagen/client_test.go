package gen

import (
	"fmt"
	"testing"

	"github.com/bgrudt/radgen/config"
	"github.com/davecgh/go-spew/spew"
)

func TestRandomCient(t *testing.T) {
	n := 0
	for n < 3 {
		var err error
		var msg Message

		err = GenerateClient(config.Cfg, &msg)
		if err != nil {
			fmt.Printf("There was an error generating the client")
		}

		spew.Dump(msg)

		n++
	}
}

func TestStaticFacility(t *testing.T) {
	n := 0
	for n < 10 {
		var err error
		var msg Message
		msg.Client.UQName = "memhosp_uc"

		err = GenerateClient(config.Cfg, &msg)
		if err != nil {
			fmt.Printf("There was an error generating the client.")
		}

		spew.Dump(msg)

		n++
	}
}

func TestBlank(t *testing.T) {
	var test string
	if test == "potato" {
		fmt.Printf("potato")
	} else {
		fmt.Printf("not potato")
	}
}
