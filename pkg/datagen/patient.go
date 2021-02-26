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
	FirstName     string
	MiddleName    string
	LastName      string
	Address       Address
	DateOfBirth   time.Time
	Sex           Code
	Race          Code
	Ethnicity     Code
	Language      Code
	MaritalStatus Code
}

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

//GeneratePatientAddress uses the locale list of a client
//to assign a random street address for a patient.
func GeneratePatientAddress(m *Message) error {
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

		file2 := dataPath + clientFile
		j2, err := os.Open(file2)
		if err != nil {
			return err
		}
		defer j2.Close()

		b2, _ := ioutil.ReadAll(j2)

		var f2 ClientList

		err = json.Unmarshal(b2, &f2)
		if err != nil {
			return err
		}
		for i := 0; i < len(f2.ClinicalClients); i++ {
			if m.Client.UQName == f2.ClinicalClients[i].UQName {
				lo, err := returnRandomLocale(&f2.ClinicalClients[i].Locales)
				if err != nil {
					return err
				}
				m.Patient.Address.City = lo.City
				m.Patient.Address.State = lo.State
				m.Patient.Address.Zip = lo.Zip

				break
			}
		}

	} else {
		m.Patient.Address.City = ""
		m.Patient.Address.State = ""
		m.Patient.Address.Zip = ""
	}
	return nil
}

//GeneratePatientCodes takes an inbound message pointer and assigns
//coded values to the patient level of the message struct.
func GeneratePatientCodes(m *Message) error {
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
		case "sex":
			if len(m.Patient.Sex.ExtCode) < 1 {
				v, err := ReturnRandomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				m.Patient.Sex = v.Code
			}
		case "race":
			if len(m.Patient.Race.ExtCode) < 1 {
				v, err := ReturnRandomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				m.Patient.Race = v.Code
			}
		case "ethnicity":
			if len(m.Patient.Ethnicity.ExtCode) < 1 {
				v, err := ReturnRandomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				m.Patient.Ethnicity = v.Code
			}
		case "language":
			if len(m.Patient.Language.ExtCode) < 1 {
				v, err := ReturnRandomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				m.Patient.Language = v.Code
			}
		case "marital_status":
			if len(m.Patient.Sex.ExtCode) < 1 {
				t := time.Now()
				if m.Patient.DateOfBirth.After(t.AddDate(-18, 0, 0)) {
					num := randomInt(1, 4)
					switch num {
					case 1:
						v, err := ReturnByInternalCode("marital_status", "other")
						if err != nil {
							return err
						}
						m.Patient.MaritalStatus = v
					case 2:
						v, err := ReturnByInternalCode("marital_status", "single")
						if err != nil {
							return err
						}
						m.Patient.MaritalStatus = v
					case 3:
						v, err := ReturnByInternalCode("marital_status", "unreported")
						if err != nil {
							return err
						}
						m.Patient.MaritalStatus = v
					default:
						v, err := ReturnByInternalCode("marital_status", "unknown")
						if err != nil {
							return err
						}
						m.Patient.MaritalStatus = v
					}
				}

				v, err := ReturnRandomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				m.Patient.Sex = v.Code
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

//GeneratePatientName should be run after the patient age and sex have been assigned
//to a message.  This function will generate a first, middle, and last name
//based on their sex and decade of birth.  It will generate a random DOB if
//the input pointer does not contain a viable DOB.
func GeneratePatientName(m *Message) error {
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

			if m.Patient.Sex.ExtCode == "M" {
				first, err := ReturnRandomValue(&fn.Names[i].Male)
				if err != nil {
					return err
				}
				m.Patient.FirstName = first
			} else if m.Patient.Sex.ExtCode == "F" {
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

	if m.Patient.Sex.ExtCode == "M" {
		mid, err := ReturnRandomValue(&mn.Male)
		if err != nil {
			return err
		}
		m.Patient.MiddleName = mid
	} else if m.Patient.Sex.ExtCode == "F" {
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
