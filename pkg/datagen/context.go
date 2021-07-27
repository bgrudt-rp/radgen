package gen

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type ContextList struct {
	ClinicalContexts []ClinicalContext `json:"clin_contexts"`
}

type ClinicalContext struct {
	UQName    string            `json:"clin_client"`
	Locations []LocationContext `json:"locations"`
}

type LocationContext struct {
	UQName       string    `json:"uq_name"`
	LocationName string    `json:"loc_name"`
	LocationCode string    `json:"loc_code"`
	Weight       int       `json:"weight"`
	Contexts     []Context `json:"contexts"`
}

type Context struct {
	AdmitSource  string `json:"admit_source"`
	AdmitType    string `json:"admit_type"`
	Service      string `json:"service"`
	PatientClass string `json:"pt_class"`
	PatientType  string `json:"pt_type"`
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

//GenerateContext can only run once a client has been applied
//to the message supplied in the pointer.
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

	//Find my client and apply values
	for i := 0; i < len(cctx.ClinicalContexts); i++ {
		if m.Client.UQName == cctx.ClinicalContexts[i].UQName {
			//Generate context - visit info
			lctx, err := returnRandomLocation(&cctx.ClinicalContexts[i].Locations)
			if err != nil {
				spew.Dump(err)
				return err
			}
			m.Visit.Location.ExtCode = lctx.LocationCode
			m.Visit.Location.Description = lctx.LocationName
			lc, err := returnRandomContext(&lctx.Contexts)
			if err != nil {
				spew.Dump(err)
				return err
			}
			//Get visit code values
			spew.Dump(lc.AdmitSource)
			v, err := returnVisitByInternalCode("admit_source", lc.AdmitSource)
			spew.Dump(v)
			if err != nil {
				return err
			}
			m.Visit.AdmitSource.ExtCode = v.Code.ExtCode
			m.Visit.AdmitSource.Description = v.Code.Description

			v, err = returnVisitByInternalCode("admit_type", lc.AdmitType)
			if err != nil {
				return err
			}
			m.Visit.AdmitType.ExtCode = v.Code.ExtCode
			m.Visit.AdmitType.Description = v.Code.Description

			break
		}
	}
	return nil
}

func returnRandomContext(c *[]Context) (Context, error) {
	var out Context
	var count int
	var cur int

	for i := 0; i < len(*c); i++ {
		count = count + (*c)[i].Weight
	}

	rng := randomInt(1, count)

	for i := 0; i < len(*c); i++ {
		cur = cur + (*c)[i].Weight
		if cur >= rng {
			return (*c)[i], nil
		}
	}

	return out, nil
}

func returnRandomLocation(lc *[]LocationContext) (LocationContext, error) {
	var out LocationContext
	var count int
	var cur int

	for i := 0; i < len(*lc); i++ {
		count = count + (*lc)[i].Weight
	}

	rng := randomInt(1, count)

	for i := 0; i < len(*lc); i++ {
		cur = cur + (*lc)[i].Weight
		if cur >= rng {
			return (*lc)[i], nil
		}
	}

	return out, nil
}
