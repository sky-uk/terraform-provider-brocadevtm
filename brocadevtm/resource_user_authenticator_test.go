package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMUserAuthenticatorBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	userAuthenticatorName := fmt.Sprintf("acctest_brocadevtm_user_authenticator-%d", randomInt)
	userAuthenticatorResourceName := "brocadevtm_user_authenticator.acctest"

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
				Config:      testAccBrocadeUserAuthenticatorNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeUserAuthenticatorInvalidType(),
				ExpectError: regexp.MustCompile(`Access level must be one of ldap, radius or tacacs_plus`),
			},
			{
				Config:      testAccBrocadeUserAuthenticatorInvalidDNMethod(),
				ExpectError: regexp.MustCompile(`Access level must be one of construct, none or search`),
			},
			{
				Config:      testAccBrocadeUserAuthenticatorInvalidTacacsAuthType(),
				ExpectError: regexp.MustCompile(`Access level must be one of ascii or pap`),
			},
			{
				Config:      testAccBrocadeUserAuthenticatorTooManyTacacs(),
				ExpectError: regexp.MustCompile(`tacacs_plus: attribute supports 1 item maximum, config has 2 declared`),
			},
			{
				Config:      testAccBrocadeUserAuthenticatorTooManyLDAP(),
				ExpectError: regexp.MustCompile(`ldap: attribute supports 1 item maximum, config has 2 declared`),
			},
			{
				Config:      testAccBrocadeUserAuthenticatorTooManyRadius(),
				ExpectError: regexp.MustCompile(`radius: attribute supports 1 item maximum, config has 2 declared`),
			},
			{

				Config: testAccBrocadeUserAuthenticatorCreate(userAuthenticatorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMUserAuthenticatorExists(userAuthenticatorName, userAuthenticatorResourceName),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "name", userAuthenticatorName),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "description", "Create user authenticator acceptance test"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "type", "ldap"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.#", "1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.base_dn", "test_dn"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.bind_dn", "test_bind_dn"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.dn_method", "search"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.fallback_group", "test_group"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.filter", "test_filter"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.group_attribute", "test_attribute"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.group_field", "test_group_field"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.group_filter", "test_group_filter"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.port", "180"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.search_dn", "test_search_dn"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.search_password", "password"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.server", "127.0.0.1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.timeout", "132"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.#", "1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.fallback_group", "test_group"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.group_attribute", "1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.group_vendor", "1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.nas_identifier", "test_nas_identifier"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.nas_ip_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.port", "180"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.secret", "secret"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.server", "127.0.0.1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.timeout", "132"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.#", "1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.auth_type", "ascii"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.fallback_group", "test_group"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.group_field", "test_group_field"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.group_service", "test_service"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.port", "180"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.secret", "secret"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.server", "127.0.0.1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.timeout", "132"),
				),
			},

			{
				Config: testAccBrocadeUserAuthenticatorUpdate(userAuthenticatorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMUserAuthenticatorExists(userAuthenticatorName, userAuthenticatorResourceName),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "name", userAuthenticatorName),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "description", "Create user authenticator acceptance test update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "type", "radius"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.#", "1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.base_dn", "test_dn_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.bind_dn", "test_bind_dn_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.dn_method", "search"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.fallback_group", "test_group_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.filter", "test_filter_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.group_attribute", "test_attribute_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.group_field", "test_group_field_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.group_filter", "test_group_filter_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.port", "360"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.search_dn", "test_search_dn_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.search_password", "password_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.server", "127.0.0.2"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "ldap.0.timeout", "264"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.#", "1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.fallback_group", "test_group_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.group_attribute", "2"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.group_vendor", "2"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.nas_identifier", "test_nas_identifier_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.nas_ip_address", "127.0.0.2"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.port", "360"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.secret", "secret_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.server", "127.0.0.2"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "radius.0.timeout", "264"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.#", "1"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.auth_type", "pap"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.fallback_group", "test_group_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.group_field", "test_group_field_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.group_service", "test_service_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.port", "360"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.secret", "secret_update"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.server", "127.0.0.2"),
					resource.TestCheckResourceAttr(userAuthenticatorResourceName, "tacacs_plus.0.timeout", "264"),
				),
			},
		},
	})
}

func testAccBrocadeVTMUserAuthenticatorCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_user_authenticator" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		authenticators, err := client.GetAllResources("user_authenticators")
		if err != nil {
			return fmt.Errorf("Error getting all User Authenticators: %+v", err)
		}
		for _, authenticator := range authenticators {
			if authenticator["name"] == name {
				return fmt.Errorf("Brocade vTM User Authenticator %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMUserAuthenticatorExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM User Authenticator %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM User Authenticator ID not set for %s in resources", name)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		authenticators, err := client.GetAllResources("user_authenticators")
		if err != nil {
			return fmt.Errorf("Error getting all User Authenticators: %+v", err)
		}
		for _, authenticator := range authenticators {
			if authenticator["name"] == name {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM User Authenticator %s not found on remote vTM", name)
	}
}

func testAccBrocadeUserAuthenticatorCreate(name string) string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       name = "%s"
       description = "Create user authenticator acceptance test"
       enabled = false
       type = "ldap"
       ldap = {
		  base_dn = "test_dn"
		  bind_dn = "test_bind_dn"
		  dn_method = "search"
		  fallback_group = "test_group"
		  filter = "test_filter"
		  group_attribute = "test_attribute"
		  group_field = "test_group_field"
		  group_filter = "test_group_filter"
		  port = 180
		  search_dn = "test_search_dn"
		  search_password = "password"
		  server = "127.0.0.1"
		  timeout = 132
	       }
       radius = {
		  fallback_group = "test_group"
		  group_attribute = 1
		  group_vendor = 1
		  nas_identifier = "test_nas_identifier"
		  nas_ip_address = "127.0.0.1"
		  port = 180
		  secret = "secret"
		  server = "127.0.0.1"
		  timeout = 132
	       }
       tacacs_plus = {
		  auth_type = "ascii"
		  fallback_group = "test_group"
		  group_field = "test_group_field"
		  group_service = "test_service"
		  port = 180
		  secret = "secret"
		  server = "127.0.0.1"
		  timeout = 132
       }

}
`, name)
}

func testAccBrocadeUserAuthenticatorUpdate(name string) string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       name = "%s"
       description = "Create user authenticator acceptance test update"
       enabled = true
       type = "radius"
       ldap = {
		  base_dn = "test_dn_update"
		  bind_dn = "test_bind_dn_update"
		  dn_method = "search"
		  fallback_group = "test_group_update"
		  filter = "test_filter_update"
		  group_attribute = "test_attribute_update"
		  group_field = "test_group_field_update"
		  group_filter = "test_group_filter_update"
		  port = 360
		  search_dn = "test_search_dn_update"
		  search_password = "password_update"
		  server = "127.0.0.2"
		  timeout = 264
	       }
       radius = {
		  fallback_group = "test_group_update"
		  group_attribute = 2
		  group_vendor = 2
		  nas_identifier = "test_nas_identifier_update"
		  nas_ip_address = "127.0.0.2"
		  port = 360
		  secret = "secret_update"
		  server = "127.0.0.2"
		  timeout = 264
	       }
       tacacs_plus = {
		  auth_type = "pap"
		  fallback_group = "test_group_update"
		  group_field = "test_group_field_update"
		  group_service = "test_service_update"
		  port = 360
		  secret = "secret_update"
		  server = "127.0.0.2"
		  timeout = 264
       }

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
		  base_dn = "testdn"
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
		  base_dn = "testdn"
		  dn_method = "search"
		  port = 180
		  timeout = 132
	       }
	}
`)
}

func testAccBrocadeUserAuthenticatorInvalidDNMethod() string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       	       name = "invalidTypeUA"
	       description = "Invalid Type Acceptance Test"
	       enabled = false
	       type = "ldap"

	       ldap = {
		  base_dn = "testdn"
		  dn_method = "kmkm"
		  port = 180
		  timeout = 132
	       }
	}
`)
}

func testAccBrocadeUserAuthenticatorInvalidTacacsAuthType() string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       	       name = "invalidTypeUA"
	       description = "Invalid Type Acceptance Test"
	       enabled = false
	       type = "tacacs_plus"

	       tacacs_plus = {
		  auth_type = "invalid"
		  fallback_group = "test-fallback-group"
		  group_field = "test-group"
		  group_service = "zeus"
		  port = 80
		  secret = "testsecret"
		  server = "127.0.0.1"
		  timeout = 120
	       }
	}
`)
}

func testAccBrocadeUserAuthenticatorTooManyTacacs() string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       	       name = "invalidTypeUA"
	       description = "Invalid Type Acceptance Test"
	       enabled = false
	       type = "tacacs_plus"

	       tacacs_plus = {
		  auth_type = "asciii"
		  fallback_group = "test-fallback-group"
		  group_field = "test-group"
		  group_service = "zeus"
		  port = 80
		  secret = "testsecret"
		  server = "127.0.0.1"
		  timeout = 120
	       }
	       tacacs_plus = {
		  auth_type = "asciii"
		  fallback_group = "test-fallback-group"
		  group_field = "test-group"
		  group_service = "zeus"
		  port = 80
		  secret = "testsecret"
		  server = "127.0.0.1"
		  timeout = 120
	       }
	}
`)
}

func testAccBrocadeUserAuthenticatorTooManyLDAP() string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       	       name = "invalidTypeUA"
	       description = "Invalid Type Acceptance Test"
	       enabled = false
	       type = "tacacs_plus"

     	       ldap = {
		  base_dn = "test_dn_update"
		  bind_dn = "test_bind_dn_update"
		  dn_method = "search"
		  fallback_group = "test_group_update"
		  filter = "test_filter_update"
		  group_attribute = "test_attribute_update"
		  group_field = "test_group_field_update"
		  group_filter = "test_group_filter_update"
		  port = 360
		  search_dn = "test_search_dn_update"
		  search_password = "password_update"
		  server = "127.0.0.2"
		  timeout = 264
	       }
     	       ldap = {
		  base_dn = "test_dn_update"
		  bind_dn = "test_bind_dn_update"
		  dn_method = "search"
		  fallback_group = "test_group_update"
		  filter = "test_filter_update"
		  group_attribute = "test_attribute_update"
		  group_field = "test_group_field_update"
		  group_filter = "test_group_filter_update"
		  port = 360
		  search_dn = "test_search_dn_update"
		  search_password = "password_update"
		  server = "127.0.0.2"
		  timeout = 264
	       }
	}
`)
}

func testAccBrocadeUserAuthenticatorTooManyRadius() string {
	return fmt.Sprintf(`
       resource "brocadevtm_user_authenticator" "acctest" {
       	       name = "invalidTypeUA"
	       description = "Invalid Type Acceptance Test"
	       enabled = false
	       type = "tacacs_plus"

       radius = {
		  fallback_group = "test_group_update"
		  group_attribute = 2
		  group_vendor = 2
		  nas_identifier = "test_nas_identifier_update"
		  nas_ip_address = "127.0.0.2"
		  port = 360
		  secret = "secret_update"
		  server = "127.0.0.2"
		  timeout = 264
	       }
       radius = {
		  fallback_group = "test_group_update"
		  group_attribute = 2
		  group_vendor = 2
		  nas_identifier = "test_nas_identifier_update"
		  nas_ip_address = "127.0.0.2"
		  port = 360
		  secret = "secret_update"
		  server = "127.0.0.2"
		  timeout = 264
	       }
       }
`)
}
