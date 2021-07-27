package gen

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

type Address struct {
	Address1 string
	Address2 string
	City     string
	State    string
	Zip      string
	County   string
}

type Addresses struct {
	Line1 []WeightedValue `json:"line_one,omitempty"`
	Line2 []WeightedValue `json:"line_two,omitempty"`
}

//GeneratePatientAddress uses the locale list of a client
//to assign a random street address for a patient.
func GenerateAddress(a *Address) error {
	//Street address
	file := dataPath + addressFile

	j, err := os.Open(file)
	if err != nil {
		return err
	}
	defer j.Close()

	b, _ := ioutil.ReadAll(j)

	var addr Addresses

	err = json.Unmarshal(b, &addr)
	if err != nil {
		return err
	}

	//don't overwrite valid address
	if len(a.Address1) < 1 {
		l1, err := randomValue(&addr.Line1)
		if err != nil {
			return err
		}
		a.Address1 = l1
	}

	//avoid filling for blank address line 1
	if len(a.Address1) > 0 {
		var sNum int

		num := randomInt(1, 3)
		switch num {
		case 1:
			sNum = randomInt(100, 999)
		case 2:
			sNum = randomInt(100, 1000)
			sNum = sNum * 5
		default:
			sNum = randomInt(100, 500)
			sNum = sNum * 100
		}
		a.Address1 = strconv.Itoa(sNum) + " " + a.Address1

		l2, err := randomValue(&addr.Line2)
		if err != nil {
			return err
		}
		a.Address2 = l2

		//locale is randomly selected across all
		//valid client locales
		file := dataPath + clientFile
		j, err := os.Open(file)
		if err != nil {
			return err
		}
		defer j.Close()

		b, _ := ioutil.ReadAll(j)

		var facs ClientList

		err = json.Unmarshal(b, &facs)
		if err != nil {
			return err
		}
		var count int
		var cur int
		for i := 0; i < len(facs.ClinicalClients); i++ {
			count = count + (facs.ClinicalClients)[i].Weight
		}

		rng := randomInt(1, count)

		for i := 0; i < len(facs.ClinicalClients); i++ {
			cur = cur + (facs.ClinicalClients)[i].Weight
			if cur >= rng {
				l, err := returnRandomLocale(&facs.ClinicalClients[i].Locales)
				if err != nil {
					return err
				}
				a.City = l.City
				a.State = l.State
				a.Zip = l.Zip

				break
			}
		}
	}

	return nil
}
