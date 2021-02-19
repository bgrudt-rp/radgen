package gen

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type Facility struct {
	UQName             string `json:"uq_name"`
	Name               string `json:"name"`
	SendingApplication string `json:"sending_app"`
	SendingFacility    string `json:"sending_fac"`
}

type Facilities struct {
	Facilities []FacilityDetails `json:"facilities"`
}

type FacilityDetails struct {
	UQName             string    `json:"uq_name"`
	Name               string    `json:"name"`
	SendingApplication string    `json:"sending_app"`
	SendingFacility    string    `json:"sending_fac"`
	Contexts           []Context `json:"contexts"`
	Locales            []Locale  `json:"locales"`
}

type Context struct {
	UQName       string `json:"uq_name"`
	LocationName string `json:"loc_name"`
	LocationCode string `json:"loc_code"`
	Service      string `json:"service"`
	PatientClass string `json:"pt_class"`
	PatientType  string `json:"pt_type"`
	AdmitType    string `json:"admit_type"`
	AdmitSource  string `json:"admit_source"`
	Weight       int    `json:"weight"`
}

type Locale struct {
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
	Weight int    `json:"weight"`
}

func ReturnRandomContext(l *[]Context) (Context, error) {
	var out Context
	var count int
	var cur int

	for i := 0; i < len(*l); i++ {
		count = count + (*l)[i].Weight
	}

	rng := randomInt(1, count)

	for i := 0; i < len(*l); i++ {
		cur = cur + (*l)[i].Weight
		if cur >= rng {
			out = (*l)[i]
			return out, nil
		}
	}

	return out, nil
}

func ReturnRandomLocale(l *[]Locale) (Locale, error) {
	var out Locale
	var count int
	var cur int

	for i := 0; i < len(*l); i++ {
		count = count + (*l)[i].Weight
	}

	rng := randomInt(1, count)

	for i := 0; i < len(*l); i++ {
		cur = cur + (*l)[i].Weight
		if cur >= rng {
			out = (*l)[i]
			return out, nil
		}
	}

	return out, nil
}

func GenerateContext(m *Message) error {
	facFile := "/users/brandongrudt/projects/hl7cli/fixtures/datagen/facilities.json"

	j, err := os.Open(facFile)
	if err != nil {
		return err
	}
	defer j.Close()

	b, _ := ioutil.ReadAll(j)

	var facs Facilities

	err = json.Unmarshal(b, &facs)
	if err != nil {
		spew.Dump(err)
		return err
	}

	for i := 0; i < len(facs.Facilities); i++ {
		if m.Facility.UQName == facs.Facilities[i].UQName {
			m.Facility.Name = facs.Facilities[i].Name
			m.Facility.SendingApplication = facs.Facilities[i].SendingApplication
			m.Facility.SendingFacility = facs.Facilities[i].SendingFacility

			//Generate patient city, state, zip
			rng := randomInt(1, 100)
			count := 0
			for ii := 0; ii < len(facs.Facilities[i].Locales); ii++ {
				up := facs.Facilities[i].Locales[ii].Weight
				count = count + up
				if count >= rng {
					m.Patient.Address.City = facs.Facilities[i].Locales[ii].City
					m.Patient.Address.State = facs.Facilities[i].Locales[ii].State
					m.Patient.Address.Zip = facs.Facilities[i].Locales[ii].Zip
					break
				}
			}

			// 	//Generate context - visit info
			rng = randomInt(1, 100)
			count = 0
			for ii := 0; ii < len(facs.Facilities[i].Contexts); ii++ {
				up := facs.Facilities[i].Contexts[ii].Weight
				count = count + up
				if count >= rng {
					m.Visit.LocationName = facs.Facilities[i].Contexts[ii].LocationName
					m.Visit.LocationCode = facs.Facilities[i].Contexts[ii].LocationCode
					m.Visit.Service = facs.Facilities[i].Contexts[ii].Service
					m.Visit.PatientClass = facs.Facilities[i].Contexts[ii].PatientClass
					m.Visit.PatientType = facs.Facilities[i].Contexts[ii].PatientType
					m.Visit.AdmitType = facs.Facilities[i].Contexts[ii].AdmitType
					m.Visit.AdmitSource = facs.Facilities[i].Contexts[ii].AdmitSource
					break
				}
			}
			break
		}
	}

	return nil
}
