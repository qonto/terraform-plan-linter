package validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TerraformPlan struct {
	PlannedValues struct {
		RootModule Module `json:"root_module"`
	} `json:"planned_values"`
}

type Module struct {
	Address      string     `json:"address"`
	Resources    []Resource `json:"resources"`
	ChildModules []Module   `json:"child_modules"`
}

type Resource struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Values struct {
		Tags map[string]string `json:"tags"`
	} `json:"values"`
}

type ResourceContext struct {
	Resource Resource
	Path     string
}

func GetPlan(planFile string) (*TerraformPlan, error) {
	planData, err := ioutil.ReadFile(planFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var plan TerraformPlan
	if err := json.Unmarshal(planData, &plan); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &plan, nil
}
