package gen

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/bgrudt/radgen/config"
	"github.com/davecgh/go-spew/spew"
)

type Client struct {
	UQName             string `json:"uq_name"`
	Name               string `json:"name"`
	SendingApplication string `json:"sending_app"`
	SendingFacility    string `json:"sending_fac"`
}

type ClientList struct {
	ClinicalClients []ClientDetail `json:"clin_clients"`
}

type ClientDetail struct {
	UQName             string   `json:"uq_name"`
	Name               string   `json:"name"`
	SendingApplication string   `json:"sending_app"`
	SendingFacility    string   `json:"sending_fac"`
	Weight             int      `json:"weight"`
	Locales            []Locale `json:"locales"`
}

type Locale struct {
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
	Weight int    `json:"weight"`
}

func returnRandomClient(l *[]ClientDetail) (Client, error) {
	var out Client
	var count int
	var cur int

	for i := 0; i < len(*l); i++ {
		count = count + (*l)[i].Weight
	}

	rng := randomInt(1, count)

	for i := 0; i < len(*l); i++ {
		cur = cur + (*l)[i].Weight
		if cur >= rng {
			out.UQName = (*l)[i].UQName
			out.Name = (*l)[i].Name
			out.SendingApplication = (*l)[i].SendingApplication
			out.SendingFacility = (*l)[i].SendingFacility
			return out, nil
		}
	}

	return out, nil
}

func returnRandomLocale(l *[]Locale) (Locale, error) {
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

func GenerateClient(cfg *config.Cfg, m *Message) error {
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
		spew.Dump(err)
		return err
	}
	if len(m.Client.UQName) < 1 {
		f, err := returnRandomClient(&facs.ClinicalClients)
		if err != nil {
			return err
		}

		m.Client = f
	} else {
		for i := 0; i < len(facs.ClinicalClients); i++ {
			if m.Client.UQName == facs.ClinicalClients[i].UQName {
				m.Client.Name = facs.ClinicalClients[i].Name
				m.Client.SendingApplication = facs.ClinicalClients[i].SendingApplication
				m.Client.SendingFacility = facs.ClinicalClients[i].SendingFacility
				break
			}
		}
	}

	return nil
}
