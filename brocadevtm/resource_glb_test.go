package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api/glb"
	"github.com/sky-uk/go-rest-api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMGLBBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	glbName := fmt.Sprintf("acctest_brocadevtm_glb-%d", randomInt)
	glbResourceName := "brocadevtm_glb.acctest"
	fmt.Printf("\n\nGLB is %s.\n\n", glbName)

	domainsSetPattern := regexp.MustCompile(`domains\.[0-9]+`)
	lastResortResponsePattern := regexp.MustCompile(`last_resort_response\.[0-9]+`)
	locationDrainingPattern := regexp.MustCompile(`location_draining\.[0-9]+`)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMGLBCheckDestroy(state, glbName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMGLBNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMGLBInvalidAlgorithmTemplate(glbName),
				ExpectError: regexp.MustCompile(`must be one of chained, geo, hybrid, load, round_robin or weighted_random`),
			},
			{
				Config:      testAccBrocadeVTMGLBInvalidGeoEffectTemplate(glbName),
				ExpectError: regexp.MustCompile(`must be a whole number between 0 and 100 \(percentage\)`),
			},
			{
				Config: testAccBrocadeGLBCreateTemplate(glbName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMGLBExists(glbName, glbResourceName),
					resource.TestCheckResourceAttr(glbResourceName, "name", glbName),
					resource.TestCheckResourceAttr(glbResourceName, "algorithm", "weighted_random"),
					resource.TestCheckResourceAttr(glbResourceName, "all_monitors_needed", "true"),
					resource.TestCheckResourceAttr(glbResourceName, "auto_recovery", "true"),
					resource.TestCheckResourceAttr(glbResourceName, "chained_auto_failback", "true"),
					resource.TestCheckResourceAttr(glbResourceName, "disable_on_failure", "true"),
					resource.TestCheckResourceAttr(glbResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(glbResourceName, "return_ips_on_fail", "true"),
					resource.TestCheckResourceAttr(glbResourceName, "geo_effect", "10"),
					resource.TestCheckResourceAttr(glbResourceName, "ttl", "30"),
					resource.TestCheckResourceAttr(glbResourceName, "chained_location_order.#", "2"),
					resource.TestCheckResourceAttr(glbResourceName, "chained_location_order.0", "example-location-one"),
					resource.TestCheckResourceAttr(glbResourceName, "chained_location_order.1", "example-location-two"),
					resource.TestCheckResourceAttr(glbResourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(glbResourceName, "rules.0", "ruleOne"),
					resource.TestCheckResourceAttr(glbResourceName, "rules.1", "ruleTwo"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, domainsSetPattern, "example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, domainsSetPattern, "another-example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, lastResortResponsePattern, "192.168.120.10"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, lastResortResponsePattern, "192.168.12.10"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationDrainingPattern, "example-location-one"),
				),
			},
			{
				Config: testAccBrocadeGLBUpdateTemplate(glbName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMGLBExists(glbName, glbResourceName),
					resource.TestCheckResourceAttr(glbResourceName, "name", glbName),
					resource.TestCheckResourceAttr(glbResourceName, "algorithm", "geo"),
					resource.TestCheckResourceAttr(glbResourceName, "all_monitors_needed", "false"),
					resource.TestCheckResourceAttr(glbResourceName, "auto_recovery", "false"),
					resource.TestCheckResourceAttr(glbResourceName, "chained_auto_failback", "false"),
					resource.TestCheckResourceAttr(glbResourceName, "disable_on_failure", "false"),
					resource.TestCheckResourceAttr(glbResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(glbResourceName, "return_ips_on_fail", "false"),
					resource.TestCheckResourceAttr(glbResourceName, "geo_effect", "90"),
					resource.TestCheckResourceAttr(glbResourceName, "ttl", "60"),
					resource.TestCheckResourceAttr(glbResourceName, "chained_location_order.#", "1"),
					resource.TestCheckResourceAttr(glbResourceName, "chained_location_order.0", "example-location-one"),
					resource.TestCheckResourceAttr(glbResourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(glbResourceName, "rules.0", "ruleTwo"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, domainsSetPattern, "example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, lastResortResponsePattern, "192.168.120.10"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationDrainingPattern, "example-location-one"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationDrainingPattern, "example-location-two"),
				),
			},
		},
	})
}

func testAccBrocadeVTMGLBCheckDestroy(state *terraform.State, name string) error {

	vtmClient := testAccProvider.Meta().(*rest.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_glb" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		api := glb.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Brocade vTM GLB - error occurred while retrieving a list of all GLBs")
		}
		for _, glb := range api.ResponseObject().(*glb.GlobalLoadBalancers).Children {
			if glb.Name == name {
				return fmt.Errorf("Brocade vTM GLB %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMGLBExists(glbName, glbResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[glbResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM GLB %s wasn't found in resources", glbName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM GLB ID not set for %s in resources", glbName)
		}

		vtmClient := testAccProvider.Meta().(*rest.Client)
		api := glb.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		for _, glb := range api.ResponseObject().(*glb.GlobalLoadBalancers).Children {
			if glb.Name == glbName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM GLB %s not found on remote vTM", glbName)
	}
}

func testAccBrocadeVTMGLBNoNameTemplate() string {
	return fmt.Sprintf(`
resource "brocadevtm_glb" "acctest" {

}
`)
}

func testAccBrocadeVTMGLBInvalidAlgorithmTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_glb" "acctest" {
  name = "%s"
  algorithm = "INVALID_ALGO"
}
`, name)
}

func testAccBrocadeVTMGLBInvalidGeoEffectTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_glb" "acctest" {
  name = "%s"
  geo_effect = 101
}
`, name)
}

func testAccBrocadeGLBCreateTemplate(glbName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_glb" "acctest" {
  name = "%s"
  algorithm = "weighted_random"
  all_monitors_needed = true
  auto_recovery = true
  chained_auto_failback = true
  disable_on_failure = true
  enabled = true
  return_ips_on_fail = true
  geo_effect = 10
  ttl = 30
  chained_location_order = [ "example-location-one", "example-location-two" ]
  rules = [ "ruleOne", "ruleTwo" ]
  domains = [ "example.com", "another-example.com" ]
  last_resort_response = [ "192.168.12.10", "192.168.120.10" ]
  location_draining = [ "example-location-one" ]
}
`, glbName)
}

func testAccBrocadeGLBUpdateTemplate(glbName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_glb" "acctest" {
  name = "%s"
  algorithm = "geo"
  all_monitors_needed = false
  auto_recovery = false
  chained_auto_failback = false
  disable_on_failure = false
  enabled = false
  return_ips_on_fail = false
  geo_effect = 90
  ttl = 60
  chained_location_order = [ "example-location-one" ]
  rules = [ "ruleTwo" ]
  domains = [ "example.com" ]
  last_resort_response = [ "192.168.120.10" ]
  location_draining = [ "example-location-one", "example-location-two" ]
}
`, glbName)
}
