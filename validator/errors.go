package validator

import (
	"fmt"

	"github.com/fatih/color"
)

type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
)

type ValidationIssue struct {
	Severity Severity
	Message  string
	KeyError string
	Path     string
}

func init() {
	color.NoColor = false // Force colorization
}

var (
	redBold    = color.New(color.FgRed, color.Bold).SprintFunc()
	yellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	cyanBold   = color.New(color.FgCyan, color.Bold).SprintFunc()
)

func NewError(path string, keyError string, message string) ValidationIssue {
	return ValidationIssue{
		Severity: SeverityError,
		Path:     path,
		Message:  message,
		KeyError: keyError,
	}
}

func NewWarning(path string, keyError string, message string) ValidationIssue {
	return ValidationIssue{
		Severity: SeverityWarning,
		Path:     path,
		Message:  message,
		KeyError: keyError,
	}
}

func (vi ValidationIssue) String() string {
	var severityColor func(...interface{}) string
	switch vi.Severity {
	case SeverityError:
		severityColor = redBold
	case SeverityWarning:
		severityColor = yellowBold
	default:
		severityColor = color.New().SprintFunc()
	}

	return fmt.Sprintf(
		"%s resource %s %s: %s",
		severityColor(vi.Severity),
		cyanBold(vi.Path),
		vi.KeyError,
		vi.Message,
	)
}
