# Terraform Provider Github Actions (Terraform Plugin Framework)

A Terraform provider for managing Github Actions configuration through the Github REST API.
## Features

- **GitHub Actions Repository Variables**: Manage repository-level variables effortlessly.
- **Terraform Resource Lifecycle**: Full support for `Create`, `Read`, `Update`, and `Delete` (CRUD) operations.
- **Data Source**: Look up existing GitHub Actions repository variables.
- **Import Support**: Bring existing UI or API-created variables under Terraform management without recreation.
- **GitHub API Authentication**: Secure authentication using Personal Access Tokens (PAT) or GitHub App tokens.
- **Drift Detection**: Automatically detects out-of-band changes or deletions made outside of Terraform.

## Getting Started

### Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24
- A Github Personal Access Token (PAT) with `repo` or `actions:read/write` scope.

### Configuration

Configure the provider with your Github Personal Access Token via the `token` parameter or the `GITHUB_TOKEN` environment variable:

```hcl
provider "github-actions" {
    token = var.github_token # or set GITHUB_TOKEN env var
}
```

---

## Usage Examples

### Resource

Create and manage a repository variable:

```hcl
provider "github-actions" {}

resource "github_actions_variable" "environment" {
  repository = "natnael-meresa/demo"
  name       = "ENVIRONMENT"
  value      = "production"
}
```

### Data Source

Look up an existing repository variable:

```hcl
data "github_actions_variable" "environment" {
  repository = "natnael-meresa/demo"
  name       = "ENVIRONMENT"
}

output "variable_value" {
  value = data.github_actions_variable.environment.value
}
```

---

## Importing Existing Resources

To import an existing GitHub Actions variable into Terraform state:

```bash
terraform import github_actions_variable.environment natnael-meresa/demo:ENVIRONMENT
```

Or using an `import` block (Terraform 1.5+):

```hcl
import {
  to = github_actions_variable.environment
  id = "natnael-meresa/demo:ENVIRONMENT"
}
```

---

## Development

### Building the Provider

```bash
go build -o terraform-provider-github-actions
```

### Running Acceptance Tests

```bash
export GITHUB_TOKEN="your-github-token"
export TF_ACC=1
go test -v ./...
```

---

## License

This project is licensed under the MIT License.