terraform {
    required_providers {
        github_actions = {
            source = "github.com/natnael-meresa/github-actions"
        }
    }
}

provider "github_actions" {}

resource "github_actions_variable" "environment" {
    repository = "natnael-meresa/terraform-provider-github-actions"
    name = "environment"
    value = "dev"
}