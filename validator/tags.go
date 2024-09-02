package validator

import (
	"fmt"
	"strings"

	"github.com/qonto/terraform-plan-linter/config"
)

func validateTags(rc ResourceContext, rule config.Rule) []ValidationIssue {
	var issues []ValidationIssue

	tagValue, ok := rc.Resource.Values.Tags[rule.Key]
	if !ok {
		return []ValidationIssue{
			NewError(
				rc.Path,
				fmt.Sprintf("has missing tag value %s", rule.Key),
				"",
			),
		}
	}
	if !contains(rule.PossibleValues, tagValue) {
		return []ValidationIssue{
			NewError(
				rc.Path,
				fmt.Sprintf("has invalid '%s' tag value", rule.Key),
				fmt.Sprintf("%s \nShould be one of %v",
					tagValue, strings.Join(rule.PossibleValues, ", "),
				),
			),
		}
	}

	return issues
}
