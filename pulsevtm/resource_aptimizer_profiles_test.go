package pulsevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-pulse-vtm/api"
	"net/http"
	"regexp"
	"testing"
)

func TestAccPulseVTMAptimizerProfilesBasic(t *testing.T) {

	aptimizerProfileName := acctest.RandomWithPrefix("acctest_pulsevtm_aptimizer_profiles")
	resourceName := "pulsevtm_aptimizer_profile.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccPulseVTMAptimizerProfilesCheckDestroy(state, aptimizerProfileName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccPulseVTMAptimizerProfilesNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPulseVTMAptimizerProfilesInvalidModeTemplate(),
				ExpectError: regexp.MustCompile(`must be one of active, idle or stealth`),
			},
			{
				Config: testAccPulseAptimizerProfilesCreateTemplate(aptimizerProfileName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMAptimizerProfilesExists(aptimizerProfileName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", aptimizerProfileName),
				),
			},
			{
				Config: testAccPulseAptimizerProfilesUpdateTemplate(aptimizerProfileName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMAptimizerProfilesExists(aptimizerProfileName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", aptimizerProfileName),
				),
			},
		},
	})
}

func testAccPulseVTMAptimizerProfilesCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "pulsevtm_aptimizer_profile" {
			continue
		}
		aptimizerProfileConfig := make(map[string]interface{})

		err := client.GetByName("aptimizer/profiles", rs.Primary.ID, &aptimizerProfileConfig)
		if client.StatusCode == http.StatusOK {
			return fmt.Errorf("[ERROR] Pulse vTM Check Destroy Error: Aptimizer Profile %s still exists", name)
		}
		if client.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("[ERROR] Pulse vTM Check Destroy Error: Aptimizer Profile %+v ", err)
	}
	return nil
}

func testAccPulseVTMAptimizerProfilesExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Pulse vTM Aptimizer Profile ID not set for %s in resources", name)
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()
		aptimizerProfileConfig := make(map[string]interface{})
		err := client.GetByName("aptimizer/profiles", name, &aptimizerProfileConfig)
		if client.StatusCode != http.StatusOK {
			return fmt.Errorf("[ERROR] Pulse vTM error whilst retrieving VTM Aptimizer Profile: %+v", err)
		}
		return nil
	}
}

func testAccPulseVTMAptimizerProfilesNoNameTemplate() string {
	return `
resource "pulsevtm_aptimizer_profile" "acctest" {
  	background_after = 50
  	background_on_additional_resources = true
  	show_info_bar = true
  	mode = "active"
}
`
}

func testAccPulseVTMAptimizerProfilesInvalidModeTemplate() string {
	return `
resource "pulsevtm_aptimizer_profile" "acctest" {
	name = "invalidModeTest"
  	background_after = 50
  	background_on_additional_resources = true
  	show_info_bar = true
  	mode = "INVALID"
}
`
}

func testAccPulseAptimizerProfilesCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_aptimizer_profile" "acctest" {
  	name = "%s"
	background_after = 50
	background_on_additional_resources = true
	show_info_bar = true
	mode = "active"
}
`, name)
}

func testAccPulseAptimizerProfilesUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_aptimizer_profile" "acctest" {
	name = "%s"
	background_after = 100
	background_on_additional_resources = false
	show_info_bar = false
	mode = "stealth"
}
`, name)
}
