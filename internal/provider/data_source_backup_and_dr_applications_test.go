package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBackupAndDRApplicationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDRApplicationsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_backup_and_dr_applications.all", "id"),
					resource.TestCheckResourceAttrSet("data.firefly_backup_and_dr_applications.all", "applications.#"),
				),
			},
		},
	})
}

func TestAccBackupAndDRApplicationsDataSource_FilterByProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDRApplicationsDataSourceFilterConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_backup_and_dr_applications.aws_only", "id"),
				),
			},
		},
	})
}

func testAccBackupAndDRApplicationsDataSourceConfig() string {
	return `
data "firefly_backup_and_dr_applications" "all" {}
`
}

func testAccBackupAndDRApplicationsDataSourceFilterConfig() string {
	return `
data "firefly_backup_and_dr_applications" "aws_only" {
  provider_type = "aws"
  status        = "Active"
}
`
}
