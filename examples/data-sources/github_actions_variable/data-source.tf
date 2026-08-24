data "github_actions_variable" "environment" {
  repository: "natnael-meresa/test-repo"
  name = "ENVIRONMENT"
}

output "environment" {
  value = data.github_actions_variable.environment.value
}
