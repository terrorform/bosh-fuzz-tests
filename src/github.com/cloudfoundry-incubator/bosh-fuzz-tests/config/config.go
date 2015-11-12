package config

import (
	"encoding/json"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Config struct {
	fs boshsys.FileSystem

	Parameters                Parameters `json:"parameters"`
	NumberOfConsequentDeploys int        `json:"number_of_consequent_deploys"`
}

type Parameters struct {
	NameLength        []int      `json:"name_length"`
	Instances         []int      `json:"instances"`
	AvailabilityZones [][]string `json:"availability_zones"`
}

func NewConfig(fs boshsys.FileSystem) *Config {
	return &Config{
		fs: fs,
	}
}

func (c *Config) Load(configPath string) error {
	contents, err := c.fs.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(contents), &c)
	if err != nil {
		return err
	}

	return nil
}