package pulsevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-pulse-vtm/api"
	"regexp"
	"testing"
)

func TestAccPulseVTMPersistenceBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	persistanceName := fmt.Sprintf("acctest_pulsevtm_persistence-%d", randomInt)
	persistenceResourceName := fmt.Sprintf("pulsevtm_persistence.acctest")
	fmt.Printf("\nPersistance is %s.\n\n", persistanceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccPulseVTMPersistenceCheckDestroy(state, persistanceName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccPulseVTMPersistanceNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPulseVTMPersistenceInvalidFailureModeTemplate(persistanceName),
				ExpectError: regexp.MustCompile(`must be one of close, new_node or url`),
			},
			{
				Config:      testAccPulseVTMPersistenceInvalidTypeTemplate(persistanceName),
				ExpectError: regexp.MustCompile(`must be one of asp, cookie, ip, j2ee, named, ssl, transparent, universal or x_zeus`),
			},
			{
				Config: testAccPulsePersistenceCreateTemplate(persistanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMPersistenceExists(persistanceName, persistenceResourceName),
					resource.TestCheckResourceAttr(persistenceResourceName, "name", persistanceName),
					resource.TestCheckResourceAttr(persistenceResourceName, "cookie", "example-cookie"),
					resource.TestCheckResourceAttr(persistenceResourceName, "delete", "true"),
					resource.TestCheckResourceAttr(persistenceResourceName, "failure_mode", "url"),
					resource.TestCheckResourceAttr(persistenceResourceName, "note", "Acceptance test"),
					resource.TestCheckResourceAttr(persistenceResourceName, "subnet_prefix_length_v4", "24"),
					resource.TestCheckResourceAttr(persistenceResourceName, "subnet_prefix_length_v6", "64"),
					resource.TestCheckResourceAttr(persistenceResourceName, "type", "cookie"),
					resource.TestCheckResourceAttr(persistenceResourceName, "url", "http://www.example.com/"),
				),
			},
			{
				Config: testAccPulsePersistenceUpdateTemplate(persistanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMPersistenceExists(persistanceName, persistenceResourceName),
					resource.TestCheckResourceAttr(persistenceResourceName, "name", persistanceName),
					resource.TestCheckResourceAttr(persistenceResourceName, "cookie", "another-example-cookie"),
					resource.TestCheckResourceAttr(persistenceResourceName, "delete", "false"),
					resource.TestCheckResourceAttr(persistenceResourceName, "failure_mode", "new_node"),
					resource.TestCheckResourceAttr(persistenceResourceName, "note", "Acceptance test - updated"),
					resource.TestCheckResourceAttr(persistenceResourceName, "subnet_prefix_length_v4", "16"),
					resource.TestCheckResourceAttr(persistenceResourceName, "subnet_prefix_length_v6", "32"),
					resource.TestCheckResourceAttr(persistenceResourceName, "type", "j2ee"),
					resource.TestCheckResourceAttr(persistenceResourceName, "url", "http://www.another-example.com/"),
				),
			},
		},
	})
}

func testAccPulseVTMPersistenceCheckDestroy(state *terraform.State, name string) error {

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "pulsevtm_persistence" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()

		persistenceClasses, err := client.GetAllResources("persistence")
		if err != nil {
			return fmt.Errorf("[ERROR] Pulse vTM error whilst retrieving list of persistence classes: %+v", err)
		}
		for _, persistenceClass := range persistenceClasses {
			if persistenceClass["name"] == name {
				return fmt.Errorf("[ERROR] Pulse vTM Persistance %s still exists", name)
			}
		}
	}
	return nil
}

func testAccPulseVTMPersistenceExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Pulse vTM Persistence %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Pulse vTM Persistance ID not set for %s in resources", name)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()

		persistenceClasses, err := client.GetAllResources("persistence")
		if err != nil {
			return fmt.Errorf("[ERROR] Pulse vTM error whilst retriving Persistance classes: %v", err)
		}
		for _, persistenceClass := range persistenceClasses {
			if persistenceClass["name"] == name {
				return nil
			}
		}
		return fmt.Errorf("[ERROR] Pulse vTM Perstistence %s not found on remote vTM", name)
	}
}

func testAccPulseVTMPersistanceNoNameTemplate() string {
	return fmt.Sprintf(`
resource "pulsevtm_persistence" "acctest" {
}
`)
}

func testAccPulseVTMPersistenceInvalidFailureModeTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_persistence" "acctest" {
  name = "%s"
  failure_mode = "INVALID"
}
`, name)
}

func testAccPulseVTMPersistenceInvalidTypeTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_persistence" "acctest" {
  name = "%s"
  type = "INVALID"
}
`, name)
}

func testAccPulsePersistenceCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_persistence" "acctest" {
  name = "%s"
  cookie = "example-cookie"
  delete = true
  failure_mode = "url"
  note = "Acceptance test"
  subnet_prefix_length_v4 = 24
  subnet_prefix_length_v6 = 64
  type = "cookie"
  url = "http://www.example.com/"
}
`, name)
}

func testAccPulsePersistenceUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_persistence" "acctest" {
  name = "%s"
  cookie = "another-example-cookie"
  delete = false
  failure_mode = "new_node"
  note = "Acceptance test - updated"
  subnet_prefix_length_v4 = 16
  subnet_prefix_length_v6 = 32
  type = "j2ee"
  url = "http://www.another-example.com/"
}
`, name)
}
