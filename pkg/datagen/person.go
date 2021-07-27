package gen

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Person struct {
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

//GeneratePerson takes a pointer to a patient.  It will randomize the
//patient name, DOB, and sex.  If any of those values are present in the
//input pointer, they will not be overwritten.
func GeneratePerson(p *Person) error {

	//Person age block
	minAge := 0
	maxAge := 105
	t := time.Now()
	if p.DateOfBirth.Before(t.AddDate(-maxAge, 0, 0)) || p.DateOfBirth.After(t.AddDate(-minAge, 0, 0)) {
		dob, err := randomDateByYear(minAge, maxAge)
		if err != nil {
			return err
		}
		spew.Dump(dob)
		p.DateOfBirth = dob
	}

	//Person code details
	file := dataPath + personCodeFile

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
			if len(p.Sex.ExtCode) < 1 {
				v, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				p.Sex = v.Code
			}
		case "race":
			if len(p.Race.ExtCode) < 1 {
				v, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				p.Race = v.Code
			}
		case "ethnicity":
			if len(p.Ethnicity.ExtCode) < 1 {
				v, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				p.Ethnicity = v.Code
			}
		case "language":
			if len(p.Language.ExtCode) < 1 {
				v, err := randomWeightedCode(&c.Codesets[i])
				if err != nil {
					return err
				}
				p.Language = v.Code
			}
		case "marital_status":
			if len(p.MaritalStatus.ExtCode) < 1 {
				t := time.Now()
				if p.DateOfBirth.After(t.AddDate(-18, 0, 0)) {
					num := randomInt(1, 4)
					switch num {
					case 1:
						v, err := returnPersonByInternalCode("marital_status", "other")
						if err != nil {
							return err
						}
						p.MaritalStatus = v
					case 2:
						v, err := returnPersonByInternalCode("marital_status", "single")
						if err != nil {
							return err
						}
						p.MaritalStatus = v
					case 3:
						v, err := returnPersonByInternalCode("marital_status", "unreported")
						if err != nil {
							return err
						}
						p.MaritalStatus = v
					default:
						v, err := returnPersonByInternalCode("marital_status", "unknown")
						if err != nil {
							return err
						}
						p.MaritalStatus = v
					}
				} else {
					v, err := randomWeightedCode(&c.Codesets[i])
					if err != nil {
						return err
					}
					p.MaritalStatus = v.Code
				}
			}
		}

	}

	//Person first name
	file = dataPath + firstNameFile

	j, err = os.Open(file)
	if err != nil {
		return err
	}
	defer j.Close()

	b, _ = ioutil.ReadAll(j)

	var fn FirstNames

	err = json.Unmarshal(b, &fn)
	if err != nil {
		return err
	}

	yob := p.DateOfBirth.Year()
	dec := strconv.Itoa(yob)
	spew.Dump(dec)
	dec = dec[0:3]
	unkSex := rand.Intn(1)

	for i := 0; i < len(fn.Names); i++ {
		if fn.Names[i].Decade[0:3] == dec {

			if p.Sex.ExtCode == "M" {
				first, err := randomValue(&fn.Names[i].Male)
				if err != nil {
					return err
				}
				p.FirstName = first
			} else if p.Sex.ExtCode == "F" {
				first, err := randomValue(&fn.Names[i].Female)
				if err != nil {
					return err
				}
				p.FirstName = first
			} else {
				if unkSex == 0 {
					first, err := randomValue(&fn.Names[i].Male)
					if err != nil {
						return err
					}
					p.FirstName = first
				} else {
					first, err := randomValue(&fn.Names[i].Female)
					if err != nil {
						return err
					}
					p.FirstName = first
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

	if p.Sex.ExtCode == "M" {
		mid, err := randomValue(&mn.Male)
		if err != nil {
			return err
		}
		p.MiddleName = mid
	} else if p.Sex.ExtCode == "F" {
		mid, err := randomValue(&mn.Female)
		if err != nil {
			return err
		}
		p.MiddleName = mid
	} else {
		if unkSex == 0 {
			mid, err := randomValue(&mn.Male)
			if err != nil {
				return err
			}
			p.MiddleName = mid
		} else {
			mid, err := randomValue(&mn.Female)
			if err != nil {
				return err
			}
			p.MiddleName = mid
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

	last, err := randomValue(&ln.Names)
	if err != nil {
		return err
	}
	p.LastName = last

	return nil
}

//returnPersonByInternalCode takes a codeset and an internal code.  It
//returns a Code struct that supplies the external values needed for
//message generation.
func returnPersonByInternalCode(cs, cv string) (Code, error) {
	var out Code

	//Loading code file
	file := dataPath + personCodeFile

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

	return out, nil
}
