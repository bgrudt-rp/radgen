package gen

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Code struct {
	ExtCode     string `json:"ext_code"`
	Description string `json:"description"`
}

type CodeList struct {
	Codesets []Codeset `json:"codesets"`
}

type Codeset struct {
	Codeset  string         `json:"codeset"`
	FillRate int            `json:"fill_rate"`
	Values   []WeightedCode `json:"values"`
}

type WeightedCode struct {
	IntCode string `json:"int_code"`
	Weight  int    `json:"weight"`
	Code    Code   `json:"code"`
}

type WeightedValue struct {
	Value  string `json:"value"`
	Weight int    `json:"weight"`
}

func randomDate(start, end time.Time) time.Time {
	min := time.Time(start).Unix()
	max := time.Time(end).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func randomInt(min, max int) int {
	var out int
	rand.Seed(time.Now().UnixNano())
	out = rand.Intn(max-min+1) + min
	return out
}

//ReturnByInternalCode takes a codeset and an internal code.  It
//returns a Code struct that supplies the external values needed for
//message generation.
func ReturnByInternalCode(cs, cv string) (Code, error) {
	var out Code

	//Loading code file
	file := dataPath + codeFile

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

//ReturnRandomCode takes a list of weighted codes and returns
//one of the corresponding codes at random.  The Weight values of the
//referenced list must be > 0 to avoid processing error.
func ReturnRandomWeightedCode(l *Codeset) (WeightedCode, error) {
	var out WeightedCode
	var count int
	var cur int

	bng := randomInt(1, 1000)
	if bng <= l.FillRate {

		for i := 0; i < len(l.Values); i++ {
			count = count + (l.Values)[i].Weight
		}

		rng := randomInt(1, count)

		for i := 0; i < len(l.Values); i++ {
			cur = cur + (l.Values)[i].Weight
			if cur >= rng {
				out = (l.Values)[i]
				return out, nil
			}
		}
	}

	return out, nil
}

//ReturnRandomValues takes a list of weighted values and returns
//one of the corresponding values at random.  The Weight values of the
//referenced list must be > 0 to avoid processing error.
func ReturnRandomValue(l *[]WeightedValue) (string, error) {
	var n string
	var count int
	var cur int

	for i := 0; i < len(*l); i++ {
		count = count + (*l)[i].Weight
	}

	rng := randomInt(1, count)

	for i := 0; i < len(*l); i++ {
		cur = cur + (*l)[i].Weight
		if cur >= rng {
			return (*l)[i].Value, nil
		}
	}

	return n, nil
}
