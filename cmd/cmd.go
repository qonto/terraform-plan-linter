package cmd

import (
	"fmt"
	"github.com/qonto/terraform-plan-linter/config"
	"github.com/qonto/terraform-plan-linter/validator"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "terraform-plan-linter [terraform plan file]",
	Short: "Validate Terraform plan against defined rules",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		planFile := args[0]
		issues, err := validator.ValidatePlan(planFile, cfg)
		if err != nil {
			fmt.Printf("Error validating plan: %v\n", err)
			os.Exit(1)
		}

		if len(issues) > 0 {
			fmt.Println("Validation completed. The following issues were found:")
			hasErrors := false
			for _, issue := range issues {
				fmt.Println(issue.String())

				if issue.Severity == validator.SeverityError {
					hasErrors = true
				}
			}
			if hasErrors {
				os.Exit(1)
			}
		} else {
			fmt.Println("Validation passed")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".terraform-plan-lint.yaml", "config file (default is .terraform-plan-lint.yaml)")
}
