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

func TestAccBrocadeVTMBandwidthBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	bandwidthName := fmt.Sprintf("acctest_brocadevtm_bandwidth-%d", randomInt)
	bandwidthResourceName := fmt.Sprintf("brocadevtm_bandwidth.acctest")
	fmt.Printf("\n\nBandwidth is %s.\n\n", bandwidthName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMBandwidthCheckDestroy(state, bandwidthName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMBandwidthNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMBandwidthInvalidSharingOption(bandwidthName),
				ExpectError: regexp.MustCompile(`must be one of cluster, connection, machine`),
			},
			{
				Config: testAccBrocadeBandwidthCreateTemplate(bandwidthName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMBandwidthExists(bandwidthName, bandwidthResourceName),
					resource.TestCheckResourceAttr(bandwidthResourceName, "name", bandwidthName),
					resource.TestCheckResourceAttr(bandwidthResourceName, "maximum", "13456"),
					resource.TestCheckResourceAttr(bandwidthResourceName, "note", "Acceptance test"),
					resource.TestCheckResourceAttr(bandwidthResourceName, "sharing", "cluster"),
				),
			},
			{
				Config: testAccBrocadeBandwidthUpdateTemplate(bandwidthName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMBandwidthExists(bandwidthName, bandwidthResourceName),
					resource.TestCheckResourceAttr(bandwidthResourceName, "name", bandwidthName),
					resource.TestCheckResourceAttr(bandwidthResourceName, "maximum", "65432"),
					resource.TestCheckResourceAttr(bandwidthResourceName, "note", "Acceptance test - updated"),
					resource.TestCheckResourceAttr(bandwidthResourceName, "sharing", "connection"),
				),
			},
		},
	})
}

func testAccBrocadeVTMBandwidthCheckDestroy(state *terraform.State, name string) error {

	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_bandwidth" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		bandwidthClasses, err := client.GetAllResources("bandwidth")
		if err != nil {
			return fmt.Errorf("Brocade vTM error whilst retrieving bandwidth classes: %+v", err)
		}
		for _, bandwidthClass := range bandwidthClasses {
			if bandwidthClass["name"] == name {
				return fmt.Errorf("Brocade vTM Bandwidth Class %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMBandwidthExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Bandwidth Class %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Bandwidth Class ID not set for %s in resources", name)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		bandwidthClasses, err := client.GetAllResources("bandwidth")
		if err != nil {
			return fmt.Errorf("Brocade vTM error whilse retriving bandwidth classes: %v", err)
		}
		for _, bandwidthClass := range bandwidthClasses {
			if bandwidthClass["name"] == name {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM Bandwidth Class %s not found on remote vTM", name)
	}
}

func testAccBrocadeVTMBandwidthNoNameTemplate() string {
	return fmt.Sprintf(`
resource "brocadevtm_bandwidth" "acctest" {
}
`)
}

func testAccBrocadeVTMBandwidthInvalidSharingOption(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_bandwidth" "acctest" {
  name = "%s"
  sharing = "INVALID OPTION"
}
`, name)
}

func testAccBrocadeBandwidthCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_bandwidth" "acctest" {
  name = "%s"
  maximum = 13456
  note = "Acceptance test"
  sharing = "cluster"
}
`, name)
}

func testAccBrocadeBandwidthUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_bandwidth" "acctest" {
  name = "%s"
  maximum = 65432
  note = "Acceptance test - updated"
  sharing = "connection"
}
`, name)
}
