package main

import (
	"io/ioutil"
	"launchpad.net/goyaml"
)

const (
	config_loc = "config/"
)

var Config_Environment string

func init() {
	Config_Environment = os.Getenv("GO_ENV")
	if Config_Environment == "" {
		Config_Environment = "development"
	}
}

// call with ("filename.yml", &structure)
func Config_GetConfig(name string, out interface{}) error {
	bytes, err := ioutil.ReadFile(config_loc + name + ".yml")
	if err != nil {
		return err
	}
	// needed &out, or just out?
	return goyaml.Unmarshal(bytes, &out)
}
