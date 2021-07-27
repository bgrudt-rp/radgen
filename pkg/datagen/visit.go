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

func GenerateVisit(v *Visit) error {
	//Loading code file
	file := dataPath + visitCodeFile

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
			if len(v.AdmitSource.ExtCode) < 1 {
				c, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				v.AdmitSource = c.Code
			}
		case "admit_type":
			if len(v.AdmitType.ExtCode) < 1 {
				c, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				v.AdmitType = c.Code
			}
		case "patient_class":
			if len(v.PatientClass.ExtCode) < 1 {
				c, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				v.PatientClass = c.Code
			}
		case "patient_type":
			if len(v.PatientType.ExtCode) < 1 {
				c, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				v.PatientType = c.Code
			}
		case "service":
			if len(v.Service.ExtCode) < 1 {
				c, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				v.Service = c.Code
			}

		}

	}

	return nil
}

//returnVisitByInternalCode takes a codeset and an internal code.  It
//returns a Code struct that supplies the external values needed for
//message generation.
func returnVisitByInternalCode(cs, cv string) (WeightedCode, error) {
	var out WeightedCode

	//Loading code file
	file := dataPath + visitCodeFile

	j, err := os.Open(file)
	if err != nil {
		return out, err
	}
	defer j.Close()

	b, _ := ioutil.ReadAll(j)

	var c CodeList

	err = json.Unmarshal(b, &c)
	if err != nil {
		return out, err
	}

	for i := 0; i < len(c.Codesets); i++ {
		if c.Codesets[i].Codeset == cs {
			for i2 := 0; i2 < len(c.Codesets); i2++ {
				if c.Codesets[i].Values[i2].IntCode == cv {
					out = c.Codesets[i].Values[i2]
				}
			}
		}
	}

	return out, nil
}
