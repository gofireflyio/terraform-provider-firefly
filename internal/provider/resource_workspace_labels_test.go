package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceLabelsResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceLabelsResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "workspace_id", "test-workspace-id"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.#", "3"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.0", "production"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.1", "critical"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.2", "managed"),
					resource.TestCheckResourceAttrSet("firefly_workspace_labels.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_workspace_labels.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccWorkspaceLabelsResourceUpdatedConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.#", "4"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.0", "production"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.1", "critical"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.2", "managed"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.test", "labels.3", "terraform"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccWorkspaceLabelsResource_singleLabel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceLabelsResourceSingleConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_workspace_labels.single", "workspace_id", "single-workspace-id"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.single", "labels.#", "1"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.single", "labels.0", "test"),
				),
			},
		},
	})
}

func TestAccWorkspaceLabelsResource_emptyLabels(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceLabelsResourceEmptyConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_workspace_labels.empty", "workspace_id", "empty-workspace-id"),
					resource.TestCheckResourceAttr("firefly_workspace_labels.empty", "labels.#", "0"),
				),
			},
		},
	})
}

func testAccWorkspaceLabelsResourceConfig() string {
	return `
resource "firefly_workspace_labels" "test" {
  workspace_id = "test-workspace-id"
  labels       = ["production", "critical", "managed"]
}
`
}

func testAccWorkspaceLabelsResourceUpdatedConfig() string {
	return `
resource "firefly_workspace_labels" "test" {
  workspace_id = "test-workspace-id"
  labels       = ["production", "critical", "managed", "terraform"]
}
`
}

func testAccWorkspaceLabelsResourceSingleConfig() string {
	return `
resource "firefly_workspace_labels" "single" {
  workspace_id = "single-workspace-id"
  labels       = ["test"]
}
`
}

func testAccWorkspaceLabelsResourceEmptyConfig() string {
	return `
resource "firefly_workspace_labels" "empty" {
  workspace_id = "empty-workspace-id"
  labels       = []
}
`
}