package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMAptimizerProfilesBasic(t *testing.T) {

	aptimizerProfileName := acctest.RandomWithPrefix("acctest_brocadevtm_aptimizer_profiles")
	resourceName := "brocadevtm_aptimizer_profile.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMAptimizerProfilesCheckDestroy(state, aptimizerProfileName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMAptimizerProfilesNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMAptimizerProfilesInvalidModeTemplate(),
				ExpectError: regexp.MustCompile(`must be one of active, idle or stealth`),
			},
			{
				Config: testAccBrocadeAptimizerProfilesCreateTemplate(aptimizerProfileName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMAptimizerProfilesExists(aptimizerProfileName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", aptimizerProfileName),
				),
			},
			{
				Config: testAccBrocadeAptimizerProfilesUpdateTemplate(aptimizerProfileName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMAptimizerProfilesExists(aptimizerProfileName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", aptimizerProfileName),
				),
			},
		},
	})
}

func testAccBrocadeVTMAptimizerProfilesCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_aptimizer_profile" {
			continue
		}
		aptimizerProfileConfig := make(map[string]interface{})

		err := client.GetByName("aptimizer/profiles", rs.Primary.ID, &aptimizerProfileConfig)
		if client.StatusCode == http.StatusOK {
			return fmt.Errorf("Brocade vTM Check Destroy Error: Aptimizer Profile %s still exists", name)
		}
		if client.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Brocade vTM Check Destroy Error: Aptimizer Profile %+v ", err)
	}
	return nil
}

func testAccBrocadeVTMAptimizerProfilesExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Aptimizer Profile ID not set for %s in resources", name)
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()
		aptimizerProfileConfig := make(map[string]interface{})
		err := client.GetByName("aptimizer/profiles", name, &aptimizerProfileConfig)
		if client.StatusCode != http.StatusOK {
			return fmt.Errorf("Brocade vTM error whilst retrieving VTM Aptimizer Profile: %+v", err)
		}
		return nil
	}
}

func testAccBrocadeVTMAptimizerProfilesNoNameTemplate() string {
	return `
resource "brocadevtm_aptimizer_profile" "acctest" {
  	background_after = 50
  	background_on_additional_resources = true
  	show_info_bar = true
  	mode = "active"
}
`
}


func testAccBrocadeVTMAptimizerProfilesInvalidModeTemplate() string {
	return `
resource "brocadevtm_aptimizer_profile" "acctest" {
	name = "invalidModeTest"
  	background_after = 50
  	background_on_additional_resources = true
  	show_info_bar = true
  	mode = "INVALID"
}
`
}

func testAccBrocadeAptimizerProfilesCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_aptimizer_profile" "acctest" {
  	name = "%s"
	background_after = 50
	background_on_additional_resources = true
	show_info_bar = true
	mode = "active"
}
`, name)
}

func testAccBrocadeAptimizerProfilesUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_aptimizer_profile" "acctest" {
	name = "%s"
	background_after = 100
	background_on_additional_resources = false
	show_info_bar = false
	mode = "stealth"
}
`, name)
}
