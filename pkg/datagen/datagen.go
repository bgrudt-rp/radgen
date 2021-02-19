package gen

import (
	"math/rand"
	"time"
)

type Message struct {
	Facility Facility
	Patient  Patient
	Visit    Visit
}

type WeightedCode struct {
	Name    string `json:"name,omitempty"`
	IntCode string `json:"int_code,omitempty"`
	ExtCode string `json:"ext_code,omitempty"`
	Weight  int    `json:"weight,omitempty"`
}

type WeightedValue struct {
	Value  string `json:"value,omitempty"`
	Weight int    `json:"weight,omitempty"`
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

//ReturnRandomCode takes a list of weighted codes and returns
//one of the corresponding codes at random.  The Weight values of the
//referenced list must be > 0 to avoid processing error.
func ReturnRandomCode(l *[]WeightedCode) (WeightedCode, error) {
	var out WeightedCode
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
