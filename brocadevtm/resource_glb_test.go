package brocadevtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func TestAccBrocadeVTMGLBBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	glbName := fmt.Sprintf("acctest_brocadevtm_glb-%d", randomInt)
	glbResourceName := "brocadevtm_glb.acctest"
	fmt.Printf("\n\nGLB is %s.\n\n", glbName)

	domainsSetPattern := regexp.MustCompile(`domains\.[0-9]+`)
	lastResortResponsePattern := regexp.MustCompile(`last_resort_response\.[0-9]+`)
	locationDrainingPattern := regexp.MustCompile(`location_draining\.[0-9]+`)
	locationSettingsIPPattern := regexp.MustCompile(`location_settings\.[0-9]+\.ip_addresses\.[0-9]+`)
	locationSettingsLocationPattern := regexp.MustCompile(`location_settings\.[0-9]+\.location`)
	locationSettingsWeightPattern := regexp.MustCompile(`location_settings\.[0-9]+\.weight`)
	locationSettingsMonitorPattern := regexp.MustCompile(`location_settings\.[0-9]+\.monitors\.[0-9]+`)
	dnsSecDomainPattern := regexp.MustCompile(`dns_sec_keys\.[0-9]+\.domain`)
	dnsSecSSLKeysPattern := regexp.MustCompile(`dns_sec_keys\.[0-9]+\.ssl_keys\.[0-9]+`)

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
				Config:      testAccBrocadeVTMGLBInvalidLocationWeightTemplate(glbName),
				ExpectError: regexp.MustCompile(`must be a whole number between 1 and 100`),
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
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsIPPattern, "192.168.234.56"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsIPPattern, "192.0.2.2"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsLocationPattern, "example-location-one"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsWeightPattern, "34"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsMonitorPattern, "glb-example-monitor"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsMonitorPattern, "glb-example-monitor2"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsIPPattern, "192.168.17.56"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsIPPattern, "192.168.8.22"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsLocationPattern, "example-location-two"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsWeightPattern, "66"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsMonitorPattern, "glb-example-monitor"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecDomainPattern, "example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecSSLKeysPattern, "another-example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecSSLKeysPattern, "example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecDomainPattern, "another-example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecSSLKeysPattern, "example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecSSLKeysPattern, "another-example.com"),
					resource.TestCheckResourceAttr(glbResourceName, "logging_enabled", "true"),
					resource.TestCheckResourceAttr(glbResourceName, "log_file_name", "/var/log/brocadevtm/test.log"),
					//resource.TestCheckResourceAttr(glbResourceName, "log_format", "%g, %n, %d, %a, %t, %s, %l, %q"), Test gets data which looks like it has attempted to interpret % symbol. Commenting out until we have a solution.
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
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationDrainingPattern, "example-location-one"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsIPPattern, "10.56.78.34"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsIPPattern, "10.23.189.47"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsLocationPattern, "example-location-two"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsWeightPattern, "50"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsMonitorPattern, "glb-example-monitor"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsIPPattern, "192.168.6.12"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsIPPattern, "192.168.89.11"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsLocationPattern, "example-location-three"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsWeightPattern, "78"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsMonitorPattern, "glb-example-monitor2"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, locationSettingsMonitorPattern, "glb-example-monitor3"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecDomainPattern, "another-example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecSSLKeysPattern, "another-example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecDomainPattern, "example.com"),
					util.AccTestCheckValueInKeyPattern(glbResourceName, dnsSecSSLKeysPattern, "example.com"),
					resource.TestCheckResourceAttr(glbResourceName, "logging_enabled", "false"),
					resource.TestCheckResourceAttr(glbResourceName, "log_file_name", "/var/log/brocadevtm/updated-test.log"),
					//resource.TestCheckResourceAttr(glbResourceName, "log_format", "%a, %t, %s, %l, %q"),
				),
			},
		},
	})
}

func testAccBrocadeVTMGLBCheckDestroy(state *terraform.State, glbName string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_glb" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		glbServices, err := client.GetAllResources("glb_services")
		if err != nil {
			return fmt.Errorf("Brocade vTM GLB - error while retrieving GLB: %v", err)
		}
		for _, glb := range glbServices {
			if glb["name"] == glbName {
				return fmt.Errorf("Brocade vTM GLB %s still exists", glbName)
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
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		glbServices, err := client.GetAllResources("glb_services")
		if err != nil {
			return fmt.Errorf("Brocade vTM GLB - error while retrieving GLB: %v", err)
		}
		for _, glb := range glbServices {
			if glb["name"] == glbName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM  GLB %s not found on remote vTM", glbName)
	}
}

func testAccBrocadeVTMGLBNoNameTemplate() string {
	return `
resource "brocadevtm_glb" "acctest" {

}
`
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

func testAccBrocadeVTMGLBInvalidLocationWeightTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_glb" "acctest" {
  name = "%s"
  location_settings = [
    {
      ip_addresses = [ "192.168.234.56", "192.0.2.2" ]
      location = "example-location-one"
      weight = 101
      monitors = [ "glb-example-monitor", "glb-example-monitor2" ]
    },
  ]
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
  location_settings = [
    {
      ip_addresses = [ "192.168.234.56", "192.0.2.2" ]
      location = "example-location-one"
      weight = 34
      monitors = [ "glb-example-monitor", "glb-example-monitor2" ]
    },
    {
      ip_addresses = [ "192.168.17.56", "192.168.8.22" ]
      location = "example-location-two"
      weight = 66
      monitors = [ "glb-example-monitor" ]
    },
  ]
  dns_sec_keys = [
    {
      domain = "example.com"
      ssl_keys = [ "another-example.com", "example.com" ]
    },
    {
      domain = "another-example.com"
      ssl_keys = [ "example.com", "another-example.com" ]
    },
  ]
  logging_enabled = true
  log_file_name = "/var/log/brocadevtm/test.log"
  //log_format = ""
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
  location_settings = [
    {
      ip_addresses = [ "10.56.78.34", "10.23.189.47" ]
      location = "example-location-two"
      weight = 50
      monitors = [ "glb-example-monitor" ]
    },
    {
      ip_addresses = [ "192.168.6.12", "192.168.89.11" ]
      location = "example-location-three"
      weight = 78
      monitors = [ "glb-example-monitor2", "glb-example-monitor3" ]
    },
  ]
  dns_sec_keys = [
    {
      domain = "another-example.com"
      ssl_keys = [ "another-example.com" ]
    },
    {
      domain = "example.com"
      ssl_keys = [ "example.com" ]
    },
  ]
  logging_enabled = false
  log_file_name = "/var/log/brocadevtm/updated-test.log"
  //log_format = ""
}
`, glbName)
}
