package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSCustomDBRoles_basic(t *testing.T) {
	resourceName := "mongodbatlas_custom_db_role.test"
	dataSourceName := "data.mongodbatlas_custom_db_roles.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	roleName := fmt.Sprintf("test-acc-custom_role-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasCustomDBRolesConfig(projectID, roleName, "INSERT", fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
				Check: resource.ComposeTestCheckFunc(
					// Test for Resource
					testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_name"),
					resource.TestCheckResourceAttrSet(resourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "INSERT"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),

					// Test for Data source
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasCustomDBRolesConfig(projectID, roleName, action, databaseName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_custom_db_role" "test" {
			project_id = "%s"
			role_name  = "%s"

			actions {
				action = "%s"
				resources {
					collection_name = ""
					database_name   = "%s"
				}
			}
		}

		data "mongodbatlas_custom_db_roles" "test" {
			project_id = "${mongodbatlas_custom_db_role.test.project_id}"
		}
	`, projectID, roleName, action, databaseName)
}
