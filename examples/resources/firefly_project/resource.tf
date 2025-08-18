resource "firefly_project" "example" {
  name        = "Example Project"
  description = "An example Firefly project"
  labels      = ["example", "terraform"]

  variables {
    key         = "ENVIRONMENT"
    value       = "development"
    sensitivity = "string"
    destination = "env"
  }
}