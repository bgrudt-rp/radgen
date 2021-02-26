package gen

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Visit struct {
	Location     Code `json:"loc_name"`
	Service      Code `json:"service"`
	PatientClass Code `json:"pt_class"`
	PatientType  Code `json:"pt_type"`
	AdmitType    Code `json:"admit_type"`
	AdmitSource  Code `json:"admit_source"`
}

func RandomizeVisit(m *Message) error {
	err := GenerateVisitCodes(m)
	if err != nil {
		return err
	}

	return nil
}

func GenerateVisitCodes(m *Message) error {
	//Loading code file
	file := dataPath + codeFile

	j, err := os.Open(file)
	if err != nil {
		return err
	}
	defer j.Close()

	b, _ := ioutil.ReadAll(j)

	var c CodeList

	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	for i := 0; i < len(c.Codesets); i++ {
		switch c.Codesets[i].Codeset {
		case "admit_source":
			if len(m.Visit.AdmitSource.ExtCode) < 1 {
				v, err := ReturnRandomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				m.Visit.AdmitSource = v.Code
			}

		}

	}

	return nil
}
