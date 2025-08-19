package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRunnersWorkspaceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRunnersWorkspaceResourceConfig("test-workspace"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "name", "test-workspace"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "description", "Test runners workspace"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "repository", "myorg/infrastructure"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "vcs_type", "github"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "default_branch", "main"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "iac_type", "terraform"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "terraform_version", "1.6.0"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "apply_rule", "manual"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "triggers.#", "1"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "triggers.0", "merge"),
					resource.TestCheckResourceAttrSet("firefly_workflows_runners_workspace.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_workflows_runners_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccRunnersWorkspaceResourceConfig("test-workspace-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "name", "test-workspace-updated"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "apply_rule", "manual"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRunnersWorkspaceResource_withProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRunnersWorkspaceResourceWithProjectConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project.test", "name", "test-project-for-workspace"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "name", "workspace-with-project"),
					resource.TestCheckResourceAttrPair("firefly_workflows_runners_workspace.test", "project_id", "firefly_project.test", "id"),
				),
			},
		},
	})
}

func TestAccRunnersWorkspaceResource_withVariables(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRunnersWorkspaceResourceWithVariablesConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "name", "workspace-with-vars"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "variables.#", "2"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "variables.0.key", "ENVIRONMENT"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "variables.0.value", "test"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "variables.0.sensitivity", "string"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "variables.0.destination", "env"),
				),
			},
		},
	})
}

func TestAccRunnersWorkspaceResource_withVariableSets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRunnersWorkspaceResourceWithVariableSetsConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_variable_set.test", "name", "test-varset-for-workspace"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "name", "workspace-with-varsets"),
					resource.TestCheckResourceAttr("firefly_workflows_runners_workspace.test", "consumed_variable_sets.#", "1"),
					resource.TestCheckResourceAttrPair("firefly_workflows_runners_workspace.test", "consumed_variable_sets.0", "firefly_variable_set.test", "id"),
				),
			},
		},
	})
}

func testAccRunnersWorkspaceResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "firefly_workflows_runners_workspace" "test" {
  name              = %[1]q
  description       = "Test runners workspace"
  repository        = "myorg/infrastructure"
  vcs_integration_id = "test-vcs-integration-id"
  vcs_type         = "github"
  default_branch   = "main"
  working_directory = "environments/test"
  iac_type         = "terraform"
  terraform_version = "1.6.0"
  apply_rule       = "manual"
  triggers         = ["merge"]
  labels           = ["test", "terraform"]
}
`, name)
}

func testAccRunnersWorkspaceResourceWithProjectConfig() string {
	return `
resource "firefly_project" "test" {
  name        = "test-project-for-workspace"
  description = "Test project for workspace"
  labels      = ["test", "project"]
}

resource "firefly_workflows_runners_workspace" "test" {
  name              = "workspace-with-project"
  description       = "Workspace linked to project"
  project_id        = firefly_project.test.id
  repository        = "myorg/infrastructure"
  vcs_integration_id = "test-vcs-integration-id"
  vcs_type         = "github"
  default_branch   = "main"
  iac_type         = "terraform"
  terraform_version = "1.6.0"
  apply_rule       = "manual"
  triggers         = ["merge"]
  labels           = ["test", "with-project"]
}
`
}

func testAccRunnersWorkspaceResourceWithVariablesConfig() string {
	return `
resource "firefly_workflows_runners_workspace" "test" {
  name              = "workspace-with-vars"
  description       = "Workspace with variables"
  repository        = "myorg/infrastructure"
  vcs_integration_id = "test-vcs-integration-id"
  vcs_type         = "github"
  default_branch   = "main"
  iac_type         = "terraform"
  terraform_version = "1.6.0"
  apply_rule       = "auto"
  triggers         = ["merge", "push"]
  labels           = ["test", "variables"]
  
  variables {
    key         = "ENVIRONMENT"
    value       = "test"
    sensitivity = "string"
    destination = "env"
  }
  
  variables {
    key         = "TF_VAR_region"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "iac"
  }
}
`
}

func testAccRunnersWorkspaceResourceWithVariableSetsConfig() string {
	return `
resource "firefly_variable_set" "test" {
  name        = "test-varset-for-workspace"
  description = "Variable set for workspace testing"
  labels      = ["test", "for-workspace"]
  
  variables {
    key         = "SHARED_CONFIG"
    value       = "shared-value"
    sensitivity = "string"
    destination = "env"
  }
}

resource "firefly_workflows_runners_workspace" "test" {
  name                   = "workspace-with-varsets"
  description            = "Workspace consuming variable sets"
  repository             = "myorg/infrastructure"
  vcs_integration_id     = "test-vcs-integration-id"
  vcs_type              = "github"
  default_branch        = "main"
  iac_type              = "terraform"
  terraform_version     = "1.6.0"
  apply_rule            = "manual"
  triggers              = ["merge"]
  labels                = ["test", "with-varsets"]
  consumed_variable_sets = [firefly_variable_set.test.id]
}
`
}