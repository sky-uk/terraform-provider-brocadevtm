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

func TestAccBrocadeVTMDNSZoneBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	dnsZoneName := fmt.Sprintf("acctest_brocadevtm_dns_zone-%d", randomInt)
	dnsZoneFileName := fmt.Sprintf("acctest_brocadevtm_dns_zone_file-%d", randomInt)
	dnsZoneFileNameUpdate := fmt.Sprintf("acctest_brocadevtm_dns_zone_file-%d", randomInt)
	dnsZoneResourceName := "brocadevtm_dns_zone.acctest"
	fmt.Printf("\n\nDNS zone is %s.\n\n", dnsZoneName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMDNSZoneCheckDestroy(state, dnsZoneName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMDNSZoneNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccBrocadeDNSZoneCreateTemplate(dnsZoneName, dnsZoneFileName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMDNSZoneExists(dnsZoneName, dnsZoneResourceName),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "name", dnsZoneName),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "origin", "example.com"),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "zone_file", dnsZoneFileName),
				),
			},
			{
				Config: testAccBrocadeDNSZoneUpdateTemplate(dnsZoneName, dnsZoneFileNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMDNSZoneExists(dnsZoneName, dnsZoneResourceName),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "name", dnsZoneName),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "origin", "updated-example.com"),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "zone_file", dnsZoneFileNameUpdate),
				),
			},
		},
	})
}

func testAccBrocadeVTMDNSZoneCheckDestroy(state *terraform.State, name string) error {

	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_dns_zone" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		zones, err := client.GetAllResources("dns_server/zones")
		if err != nil {
			return fmt.Errorf("Brocade vTM DNS zone - error occurred whilst retrieving a list of all DNS zones: %+v", err)
		}
		for _, dnsZone := range zones {
			if dnsZone["name"] == name {
				return fmt.Errorf("Brocade vTM DNS zone %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMDNSZoneExists(dnsZoneName, dnsZoneResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[dnsZoneResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM DNS zone %s wasn't found in resources", dnsZoneName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM DNS zone ID not set for %s in resources", dnsZoneName)
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		zones, err := client.GetAllResources("dns_server/zones")
		if err != nil {
			return fmt.Errorf("Error getting all dns zones: %+v", err)
		}
		for _, dnsZone := range zones {
			if dnsZone["name"] == dnsZoneName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM DNS zone %s not found on remote vTM", dnsZoneName)
	}
}

func testAccBrocadeVTMDNSZonePrepare(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone_file" "acctest" {
  name = "%s"
  dns_zone_config = <<DNS_ZONE_CONFIG
$TTL 3600
@                       30  IN  SOA example.com. hostmaster.isp.example.com. (
                                    2017092901 ; serial
                                    3600       ; refresh after 1 hour
                                    300        ; retry after 5 minutes
                                    1209600    ; expire after 2 weeks
                                    30 )       ; minimum TTL of 30 seconds
                                    IN  NS  ns1.example.com.
ns1				60  IN  A   10.0.0.1
;
; Services - Each service in a location has a unique IP address. Two locations = two IPs.
;
example-service         	60  IN  A   10.100.10.5
                        	60  IN  A   10.100.20.5
another-example-service         60  IN  A   10.100.10.6
                        	60  IN  A   10.100.20.6
DNS_ZONE_CONFIG
}
`, name)
}

func testAccBrocadeVTMDNSZoneNoNameTemplate() string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone" "acctest" {
  origin = "example.com"
  zone_file = "example.com.db"
}
`)
}

func testAccBrocadeDNSZoneCreateTemplate(name, dnsZoneFileName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone" "acctest" {
  name = "%s"
  origin = "example.com"
  zone_file = "${brocadevtm_dns_zone_file.acctest.name}"
}
%s
`, name, testAccBrocadeVTMDNSZonePrepare(dnsZoneFileName))
}

func testAccBrocadeDNSZoneUpdateTemplate(name, dnsZoneFileName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone" "acctest" {
  name = "%s"
  origin = "updated-example.com"
  zone_file = "${brocadevtm_dns_zone_file.acctest.name}"
}
%s
`, name, testAccBrocadeVTMDNSZonePrepare(dnsZoneFileName))
}
