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

func TestAccPulseVTMBandwidthBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	bandwidthName := fmt.Sprintf("acctest_pulsevtm_bandwidth-%d", randomInt)
	bandwidthResourceName := fmt.Sprintf("pulsevtm_bandwidth.acctest")
	fmt.Printf("\n\nBandwidth is %s.\n\n", bandwidthName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccPulseVTMBandwidthCheckDestroy(state, bandwidthName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccPulseVTMBandwidthNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPulseVTMBandwidthInvalidSharingOption(bandwidthName),
				ExpectError: regexp.MustCompile(`must be one of cluster, connection, machine`),
			},
			{
				Config: testAccPulseBandwidthCreateTemplate(bandwidthName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMBandwidthExists(bandwidthName, bandwidthResourceName),
					resource.TestCheckResourceAttr(bandwidthResourceName, "name", bandwidthName),
					resource.TestCheckResourceAttr(bandwidthResourceName, "maximum", "13456"),
					resource.TestCheckResourceAttr(bandwidthResourceName, "note", "Acceptance test"),
					resource.TestCheckResourceAttr(bandwidthResourceName, "sharing", "cluster"),
				),
			},
			{
				Config: testAccPulseBandwidthUpdateTemplate(bandwidthName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMBandwidthExists(bandwidthName, bandwidthResourceName),
					resource.TestCheckResourceAttr(bandwidthResourceName, "name", bandwidthName),
					resource.TestCheckResourceAttr(bandwidthResourceName, "maximum", "65432"),
					resource.TestCheckResourceAttr(bandwidthResourceName, "note", "Acceptance test - updated"),
					resource.TestCheckResourceAttr(bandwidthResourceName, "sharing", "connection"),
				),
			},
		},
	})
}

func testAccPulseVTMBandwidthCheckDestroy(state *terraform.State, name string) error {

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "pulsevtm_bandwidth" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()

		bandwidthClasses, err := client.GetAllResources("bandwidth")
		if err != nil {
			return fmt.Errorf("[ERROR] Pulse vTM error whilst retrieving bandwidth classes: %+v", err)
		}
		for _, bandwidthClass := range bandwidthClasses {
			if bandwidthClass["name"] == name {
				return fmt.Errorf("[ERROR] Pulse vTM Bandwidth Class %s still exists", name)
			}
		}
	}
	return nil
}

func testAccPulseVTMBandwidthExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Pulse vTM Bandwidth Class %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Pulse vTM Bandwidth Class ID not set for %s in resources", name)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()

		bandwidthClasses, err := client.GetAllResources("bandwidth")
		if err != nil {
			return fmt.Errorf("[ERROR] Pulse vTM error whilst retriving bandwidth classes: %v", err)
		}
		for _, bandwidthClass := range bandwidthClasses {
			if bandwidthClass["name"] == name {
				return nil
			}
		}
		return fmt.Errorf("[ERROR] Pulse vTM Bandwidth Class %s not found on remote vTM", name)
	}
}

func testAccPulseVTMBandwidthNoNameTemplate() string {
	return fmt.Sprintf(`
resource "pulsevtm_bandwidth" "acctest" {
}
`)
}

func testAccPulseVTMBandwidthInvalidSharingOption(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_bandwidth" "acctest" {
  name = "%s"
  sharing = "INVALID OPTION"
}
`, name)
}

func testAccPulseBandwidthCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_bandwidth" "acctest" {
  name = "%s"
  maximum = 13456
  note = "Acceptance test"
  sharing = "cluster"
}
`, name)
}

func testAccPulseBandwidthUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_bandwidth" "acctest" {
  name = "%s"
  maximum = 65432
  note = "Acceptance test - updated"
  sharing = "connection"
}
`, name)
}
