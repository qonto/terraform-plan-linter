# Terraform plan linter

Basic linter for terraform plan output. Currently in development.


## Usage

```bash
# Build project
go build

# Get terraform plan
terraform plan -out=plan.out
terraform show -json plan.out > plan.json

./terraform-plan-lint plan.json
```


Configuration file example

```yaml
rules:
  exampleRule:
    type: tags
    key: Owner
    possible_values: ['team1', 'team2']
    target_aws_resources: # Optional, if not provided it will check all resources
      - aws_s3_bucket
    fetch_possible_values_from:
      url: https://myapi.com/possible_values # Optional, otherwise it will use possible_values
```

