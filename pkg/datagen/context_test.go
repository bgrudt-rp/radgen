package gen

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRandomContext(t *testing.T) {
	n := 0
	for n < 1 {
		var err error
		var msg Message

		err = GenerateClient(&msg.Client)
		if err != nil {
			fmt.Printf("There was an error generating the client")
		}

		err = GenerateContext(&msg)
		if err != nil {
			fmt.Printf("There was an error randomizing the context")
		}

		spew.Dump(msg)

		n++
	}
}
