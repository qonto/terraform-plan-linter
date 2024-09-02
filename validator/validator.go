package validator

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/qonto/terraform-plan-linter/config"
)

const (
	RuleTypeTags = "tags"
	// Add more rule types here as needed
)

type ValidateFunc func(resourceContext ResourceContext, rule config.Rule) []ValidationIssue

var validateFuncs = map[string]ValidateFunc{
	RuleTypeTags: validateTags,
	// Add more validation functions here as needed
}

func ValidatePlan(planFile string, cfg *config.Config) ([]ValidationIssue, error) {
	plan, err := GetPlan(planFile)
	if err != nil {
		return nil, err
	}

	var allIssues []ValidationIssue

	for ruleName, rule := range cfg.Rules {
		// Create a set of AWS resources to validate for faster lookup
		resourceSet := make(map[string]bool)
		for _, resource := range rule.TargetAWSResources {
			resourceSet[resource] = true
		}

		if rule.FetchPossibleValues.URL != "" {
			fetchedValues, err := fetchPossibleValues(rule.FetchPossibleValues.URL)
			if err != nil {
				return nil, fmt.Errorf("error fetching possible values for rule '%s': %v", ruleName, err)
			}
			rule.PossibleValues = fetchedValues
		}

		validateFunc, ok := validateFuncs[rule.Type]
		if !ok {
			return nil, fmt.Errorf("unknown rule type for rule '%s': %s", ruleName, rule.Type)
		}

		issues := validateModule(plan.PlannedValues.RootModule, rule, validateFunc, resourceSet, "")
		allIssues = append(allIssues, issues...)
	}

	return allIssues, nil
}

func validateModule(module Module, rule config.Rule, validateFunc ValidateFunc, resourceSet map[string]bool, parentPath string) []ValidationIssue {
	var issues []ValidationIssue

	for _, resource := range module.Resources {
		if len(rule.TargetAWSResources) == 0 || resourceSet[resource.Type] {
			resourcePath := parentPath
			if resourcePath != "" {
				resourcePath += "."
			}
			resourcePath += fmt.Sprintf("%s.%s", resource.Type, resource.Name)

			resourceContext := ResourceContext{
				Resource: resource,
				Path:     resourcePath,
			}
			issues = append(issues, validateFunc(resourceContext, rule)...)
		}
	}

	for _, childModule := range module.ChildModules {
		modulePath := parentPath
		if modulePath != "" {
			modulePath += "."
		}
		modulePath += childModule.Address
		issues = append(issues, validateModule(childModule, rule, validateFunc, resourceSet, modulePath)...)
	}

	return issues
}

func fetchPossibleValues(url string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON array: %v", err)
	}
	return result, nil

}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
