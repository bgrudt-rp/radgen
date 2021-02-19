package gen

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Patient struct {
	FirstName   string
	MiddleName  string
	LastName    string
	Address     Address
	DateOfBirth time.Time
	Sex         Sex
	Race        Race
}

type Address struct {
	Address1 string
	Address2 string
	City     string
	State    string
	Zip      string
	County   string
}

type Sex struct {
	Code        string
	Description string
}

type Race struct {
	InternalCode         string
	Code                 string
	Description          string
	EthnicityCode        string
	EthnicityDescription string
}

type Addresses struct {
	Line1 []WeightedValue `json:"line_one,omitempty"`
	Line2 []WeightedValue `json:"line_two,omitempty"`
}

type CodeList struct {
	Codesets []Codeset `json:"codesets,omitempty"`
}

type Codeset struct {
	Codeset string         `json:"codeset,omitempty"`
	Values  []WeightedCode `json:"values,omitempty"`
}

type FirstNames struct {
	Names []DecadeName `json:"names,omitempty"`
}

type MiddleNames struct {
	Female []WeightedValue `json:"female,omitempty"`
	Male   []WeightedValue `json:"male,omitempty"`
}

type LastNames struct {
	Names []WeightedValue `json:"names,omitempty"`
}

type DecadeName struct {
	Decade string          `json:"decade,omitempty"`
	Female []WeightedValue `json:"female,omitempty"`
	Male   []WeightedValue `json:"male,omitempty"`
}

//RandomizePatient will randomize details about a patient that have not
//been explicitly set.  If the message that is referenced in the input params
//has non-nil patient fields, those fields will not be overwritten.
func RandomizePatient(m *Message) error {
	const (
		dataPath       = "/users/brandongrudt/projects/hl7cli/fixtures/datagen/"
		firstNameFile  = "first_name.json"
		lastNameFile   = "last_name.json"
		middleNameFile = "middle_name.json"
		addressFile    = "street_address.json"
	)

	err := GeneratePatientDateOfBirth(m)
	if err != nil {
		return err
	}

	err = GeneratePatientCodes(m)
	if err != nil {
		return err
	}

	err = GeneratePatientName(m)
	if err != nil {
		return err
	}

	err = GeneratePatientAddress(m)
	if err != nil {
		return err
	}
	return nil
}

func GeneratePatientAddress(m *Message) error {
	const (
		dataPath    = "/users/brandongrudt/projects/hl7cli/fixtures/datagen/"
		addressFile = "street_address.json"
	)

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

	l1, err := ReturnRandomValue(&addr.Line1)
	if err != nil {
		return err
	}

	if len(l1) > 0 {
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
		m.Patient.Address.Address1 = strconv.Itoa(sNum) + " " + l1

		l2, err := ReturnRandomValue(&addr.Line2)
		if err != nil {
			return err
		}
		m.Patient.Address.Address2 = l2
	} else {
		m.Patient.Address.City = ""
		m.Patient.Address.State = ""
		m.Patient.Address.Zip = ""
	}
	return nil
}

func GeneratePatientCodes(m *Message) error {
	const (
		dataPath = "/users/brandongrudt/projects/hl7cli/fixtures/datagen/"
		codeFile = "codes.json"
	)
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

	//Sex
	if len(m.Patient.Sex.Code) < 1 {
		for i := 0; i < len(c.Codesets); i++ {
			if c.Codesets[i].Codeset == "sex" {
				v, err := ReturnRandomCode(&c.Codesets[i].Values)
				if err != nil {
					return err
				}
				m.Patient.Sex.Code = v.ExtCode
				m.Patient.Sex.Description = v.Name
			}
		}
	}

	//Race
	if len(m.Patient.Race.Code) < 1 {
		for i := 0; i < len(c.Codesets); i++ {
			if c.Codesets[i].Codeset == "race" {
				v, err := ReturnRandomCode(&c.Codesets[i].Values)
				if err != nil {
					return err
				}
				m.Patient.Race.InternalCode = v.IntCode
				m.Patient.Race.Code = v.ExtCode
				m.Patient.Race.Description = v.Name
			}
		}
	}

	//Ethnicity
	if len(m.Patient.Race.EthnicityCode) < 1 {
		for i := 0; i < len(c.Codesets); i++ {
			if c.Codesets[i].Codeset == "ethnicity" {
				v, err := ReturnRandomCode(&c.Codesets[i].Values)
				if err != nil {
					return err
				}
				m.Patient.Race.EthnicityCode = v.ExtCode
				m.Patient.Race.EthnicityDescription = v.Name
			}
		}
	}
	return nil
}

func GeneratePatientDateOfBirth(m *Message) error {
	//DOB - will return random DOB from last 100 years
	t := time.Now()
	if m.Patient.DateOfBirth.Before(t.AddDate(-120, 0, 0)) {
		st := t.AddDate(-100, 0, 0)
		dob := randomDate(st, t)
		m.Patient.DateOfBirth = dob
	}

	return nil
}

func GeneratePatientName(m *Message) error {
	const (
		dataPath       = "/users/brandongrudt/projects/hl7cli/fixtures/datagen/"
		firstNameFile  = "first_name.json"
		lastNameFile   = "last_name.json"
		middleNameFile = "middle_name.json"
	)

	//Generate a DOB if there is none - needed for name
	err := GeneratePatientDateOfBirth(m)
	if err != nil {
		return err
	}

	//Patient first name
	file := dataPath + firstNameFile

	j, err := os.Open(file)
	if err != nil {
		return err
	}
	defer j.Close()

	b, _ := ioutil.ReadAll(j)

	var fn FirstNames

	err = json.Unmarshal(b, &fn)
	if err != nil {
		return err
	}

	yob := m.Patient.DateOfBirth.Year()
	dec := strconv.Itoa(yob)
	dec = dec[0:3]
	unkSex := rand.Intn(1)

	for i := 0; i < len(fn.Names); i++ {
		if fn.Names[i].Decade[0:3] == dec {

			if m.Patient.Sex.Code == "M" {
				first, err := ReturnRandomValue(&fn.Names[i].Male)
				if err != nil {
					return err
				}
				m.Patient.FirstName = first
			} else if m.Patient.Sex.Code == "F" {
				first, err := ReturnRandomValue(&fn.Names[i].Female)
				if err != nil {
					return err
				}
				m.Patient.FirstName = first
			} else {
				if unkSex == 0 {
					first, err := ReturnRandomValue(&fn.Names[i].Male)
					if err != nil {
						return err
					}
					m.Patient.FirstName = first
				} else {
					first, err := ReturnRandomValue(&fn.Names[i].Female)
					if err != nil {
						return err
					}
					m.Patient.FirstName = first
				}

				break
			}
		}
	}

	//Middle name
	file = dataPath + middleNameFile

	j, err = os.Open(file)
	if err != nil {
		return err
	}
	defer j.Close()

	b, _ = ioutil.ReadAll(j)

	var mn MiddleNames

	err = json.Unmarshal(b, &mn)
	if err != nil {
		return err
	}

	if m.Patient.Sex.Code == "M" {
		mid, err := ReturnRandomValue(&mn.Male)
		if err != nil {
			return err
		}
		m.Patient.MiddleName = mid
	} else if m.Patient.Sex.Code == "F" {
		mid, err := ReturnRandomValue(&mn.Female)
		if err != nil {
			return err
		}
		m.Patient.MiddleName = mid
	} else {
		if unkSex == 0 {
			mid, err := ReturnRandomValue(&mn.Male)
			if err != nil {
				return err
			}
			m.Patient.MiddleName = mid
		} else {
			mid, err := ReturnRandomValue(&mn.Female)
			if err != nil {
				return err
			}
			m.Patient.MiddleName = mid
		}
	}

	//Last name
	file = dataPath + lastNameFile

	j, err = os.Open(file)
	if err != nil {
		return err
	}
	defer j.Close()

	b, _ = ioutil.ReadAll(j)

	var ln LastNames

	err = json.Unmarshal(b, &ln)
	if err != nil {
		return err
	}

	last, err := ReturnRandomValue(&ln.Names)
	if err != nil {
		return err
	}
	m.Patient.LastName = last

	return nil
}
