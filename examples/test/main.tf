terraform {
  required_providers {
    github-actions = {
      source = "github.com/natnael-meresa/github-actions"
    }
  }
}

provider "github-actions" {}

resource "github_actions_variable" "environment" {
 repository = "natnael-meresa/terraform-provider-github-actions"
  name       = "ENVIRONMENT"
  value      = "development"
}