package brocadevtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
)

func TestAccBrocadeVTMUserGroupBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	userGroupName := fmt.Sprintf("acctest_brocadevtm_user_group-%d", randomInt)
	userGroupResourceName := "brocadevtm_user_group.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMUserGroupCheckDestroy(state, userGroupName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMUserGroupNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccBrocadeUserGroupCreate(userGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMUserGroupExists(userGroupName, userGroupResourceName),
					resource.TestCheckResourceAttr(userGroupResourceName, "name", userGroupName),
					resource.TestCheckResourceAttr(userGroupResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(userGroupResourceName, "password_expire_time", "300"),
					resource.TestCheckResourceAttr(userGroupResourceName, "timeout", "300"),
					resource.TestCheckResourceAttr(userGroupResourceName, "permissions.#", "1"),
				),
			},
			{
				Config: testAccBrocadeUserGroupUpdate(userGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMUserGroupExists(userGroupName, userGroupResourceName),
					resource.TestCheckResourceAttr(userGroupResourceName, "name", userGroupName),
					resource.TestCheckResourceAttr(userGroupResourceName, "description", "test description - updated"),
					resource.TestCheckResourceAttr(userGroupResourceName, "password_expire_time", "600"),
					resource.TestCheckResourceAttr(userGroupResourceName, "timeout", "600"),
					resource.TestCheckResourceAttr(userGroupResourceName, "permissions.#", "2"),
				),
			},
		},
	})
}

func testAccBrocadeVTMUserGroupCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "infoblox_user_group" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		userGroups, err := client.GetAllResources("user_groups")
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM User Group error whilst retrieving %s: %v", name, err)
		}
		for _, individualUserGroup := range userGroups {
			if individualUserGroup["name"] == name {
				return fmt.Errorf("[ERROR] Brocade vTM User Group %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMUserGroupExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Brocade vTM User Group %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Brocade vTM User Group ID not set for %s in resources", name)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		userGroups, err := client.GetAllResources("user_groups")
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM User Group error whilst retrieving %s: %v", name, err)
		}
		for _, individualUserGroup := range userGroups {
			if individualUserGroup["name"] == name {
				return nil
			}
		}
		return fmt.Errorf("[ERROR] Brocade vTM User Group %s not found on remote vTM", name)
	}
}

func testAccBrocadeVTMUserGroupNoName() string {
	return fmt.Sprintf(`
resource "brocadevtm_user_group" "acctest" {
       description = "test description"
       password_expire_time = 300
       timeout = 300
       permissions = {
          name =  "Web_Cache"
          access_level = "FULL"
       }
}
`)
}

func testAccBrocadeUserGroupCreate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_user_group" "acctest" {
       name = "%s"
       description = "test description"
       password_expire_time = 300
       timeout = 300
       permissions = {
          name =  "TestPermissionOne"
          access_level = "full"
       }
}
`, name)
}

func testAccBrocadeUserGroupUpdate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_user_group" "acctest" {
       name = "%s"
       description = "test description - updated"
       password_expire_time = 600
       timeout = 600
       permissions = {
          name =  "TestPermissionOne"
          access_level = "ro"
       }
        permissions = {
          name =  "TestPermissionTwo"
          access_level = "ro"
       }
}
`, name)
}
