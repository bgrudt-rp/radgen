package gen

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type ContextList struct {
	ClinicalContexts []ClinicalContext `json:"clin_clients"`
}

type ClinicalContext struct {
	UQName   string    `json:"clin_client"`
	Contexts []Context `json:"contexts"`
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

func GenerateContext(m *Message) error {
	file := dataPath + contextFile

	j, err := os.Open(file)
	if err != nil {
		return err
	}
	defer j.Close()

	b, _ := ioutil.ReadAll(j)

	var cctx ContextList

	err = json.Unmarshal(b, &cctx)
	if err != nil {
		spew.Dump(err)
		return err
	}

	//Find my facility and apply values
	for i := 0; i < len(cctx.ClinicalContexts); i++ {
		if m.Client.UQName == cctx.ClinicalContexts[i].UQName {
			// 	//Generate context - visit info
			ct, err := ReturnRandomContext(&cctx.ClinicalContexts[i].Contexts)
			if err != nil {
				spew.Dump(err)
				return err
			}

			m.Visit.Location.Description = ct.LocationName
			m.Visit.Location.ExtCode = ct.LocationCode
			m.Visit.Service.ExtCode = ct.Service
			m.Visit.PatientClass.ExtCode = ct.PatientClass
			m.Visit.PatientType.ExtCode = ct.PatientType
			m.Visit.AdmitType.ExtCode = ct.AdmitType
			m.Visit.AdmitSource.ExtCode = ct.AdmitSource

			break
		}
	}
	return nil
}
