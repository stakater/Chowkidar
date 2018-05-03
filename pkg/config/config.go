package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Controllers []Controller `json:controllers`
}

type Controller struct {
	Type           string
	WatchCriterion []string
	Actions        []Action
	Operator       string
}

type Action struct {
	Name   string
	Params map[string]string
}

func ReadConfig(filePath string) Config {
	var config Config
	// Read YML
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	// Unmarshall
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Panic(err)
	}

	return config
}
func WriteConfig(config Config, path string) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
