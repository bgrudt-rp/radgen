package config

import "github.com/davecgh/go-spew/spew"

type Config struct {
	Filestore Filestore
}

type Filestore struct {
	DataPath       string
	FirstNameFile  string
	LastNameFile   string
	MiddleNameFile string
	AddressFile    string
	ClientFile     string
	CodeFile       string
	ContextFile    string
}

func init() {
	Cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	spew.Dump(Cfg)
}

func NewConfig() (*Config, error) {
	cfg := Config{
		Filestore: Filestore{
			DataPath:       "/users/brandongrudt/projects/radgen/fixtures/",
			FirstNameFile:  "first_name.json",
			LastNameFile:   "last_name.json",
			MiddleNameFile: "middle_name.json",
			AddressFile:    "street_address.json",
			CodeFile:       "codes.json",
			ClientFile:     "clin_clients.json",
			ContextFile:    "clin_contexts.json",
		},
	}
	return &cfg, nil
}
