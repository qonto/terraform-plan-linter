package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Rules map[string]Rule `yaml:"rules"`
}

type Rule struct {
	Type                string   `yaml:"type"`
	Key                 string   `yaml:"key"`
	PossibleValues      []string `yaml:"possible_values"`
	TargetAWSResources  []string `yaml:"target_aws_resources"`
	FetchPossibleValues struct {
		URL string `yaml:"url"`
	} `yaml:"fetch_possible_values_from"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config %s: %v", filename, err)
	}

	return &config, nil
}
