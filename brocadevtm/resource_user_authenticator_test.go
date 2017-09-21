package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api/user_groups"
	"github.com/sky-uk/go-rest-api"
	"net/http"
	"regexp"
	"testing"
	"github.com/sky-uk/go-brocade-vtm/api/user_authenticators"
)

func TestAccBrocadeVTMUserAuthenticatorBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	userAuthenticatorName := fmt.Sprintf("acctest_brocadevtm_user_authenticator-%d", randomInt)
	//userAuthenticatorResourceName := "brocadevtm_user_authenticator.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMUserAuthenticatorCheckDestroy(state, userAuthenticatorName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBrocadeUserAuthenticatorNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccBrocadeUserAuthenticatorInvalidType(),
				ExpectError: regexp.MustCompile(`Access level must be one of ldap, radius or tacas_plus`),
			},
			/*
			{

				Config: testAccBrocadeUserAuthenticatorCreate(userAuthenticatorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMUserAuthenticatorExists(userAuthenticatorName, userAuthenticatorResourceName),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "name", userAuthenticatorName),
				),
			},
			{
				Config: testAccBrocadeUserAuthenticatorUpdate(userAuthenticatorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMUserAuthenticatorExists(userAuthenticatorName, userAuthenticatorResourceName),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "name", userAuthenticatorName),
				),
			},
			*/
		},
	})
}

func testAccBrocadeVTMUserAuthenticatorCheckDestroy(state *terraform.State, name string) error {
	vtmClient := testAccProvider.Meta().(*rest.Client)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "infoblox_user_authenticator" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		api := userauthenticators.NewGet(name)
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Error: Brocade vTM error occurred while retrieving User Authenticator: %v", err)
		}
		if api.StatusCode() == http.StatusOK {
			return fmt.Errorf("Error: Brocade vTM User Authenticator %s still exists", name)
		}
	}
	return nil
}

func testAccBrocadeVTMUserAuthenticatorExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM User Group %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM User Group ID not set for %s in resources", name)
		}
		vtmClient := testAccProvider.Meta().(*rest.Client)
		api := usergroups.NewGet(name)
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Brocade vTM User Group - error while retrieving User Group: %v", err)
		}
		if api.StatusCode() == http.StatusOK {
			return nil
		}
		return fmt.Errorf("Brocade vTM User Group %s not found on remote vTM", name)
	}
}

func testAccBrocadeUserAuthenticatorCreate(name string) string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       name = "%s"
}
`, name)
}

func testAccBrocadeUserAuthenticatorUpdate(name string) string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       name = "%s"
}
`, name)
}

func testAccBrocadeUserAuthenticatorNoName() string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
	       description = "No Name Acceptance Test"
	       enabled = false
	       type = "ldap"

	       ldap = {
		  base_dn = "testupdated3333"
		  dn_method = "search"
		  port = 180
		  timeout = 132
	       }
	}
`)
}

func testAccBrocadeUserAuthenticatorInvalidType() string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       	       name = "invalidTypeUA"
	       description = "Invalid Type Acceptance Test"
	       enabled = false
	       type = "invalidtype"

	       ldap = {
		  base_dn = "testupdated3333"
		  dn_method = "search"
		  port = 180
		  timeout = 132
	       }
	}
`)
}