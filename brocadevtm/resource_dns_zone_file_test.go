package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"log"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMDNSZoneFileBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	dnsZoneFileName := fmt.Sprintf("acctest_brocadevtm_dns_zone_file-%d", randomInt)
	dnsZoneFileResourceName := "brocadevtm_dns_zone_file.acctest"
	fmt.Printf("\n\nDNS zone file is %s.\n\n", dnsZoneFileName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMDNSZoneFileCheckDestroy(state, dnsZoneFileName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMDNSZoneFileNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccBrocadeDNSZoneFileCreateTemplate(dnsZoneFileName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMDNSZoneFileExists(dnsZoneFileName, dnsZoneFileResourceName),
					resource.TestCheckResourceAttr(dnsZoneFileResourceName, "name", dnsZoneFileName),
					resource.TestMatchResourceAttr(dnsZoneFileResourceName, "dns_zone_config", regexp.MustCompile(`example-service`)),
				),
			},
			{
				Config: testAccBrocadeDNSZoneFileUpdateTemplate(dnsZoneFileName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMDNSZoneFileExists(dnsZoneFileName, dnsZoneFileResourceName),
					resource.TestCheckResourceAttr(dnsZoneFileResourceName, "name", dnsZoneFileName),
					resource.TestMatchResourceAttr(dnsZoneFileResourceName, "dns_zone_config", regexp.MustCompile(``)),
					resource.TestMatchResourceAttr(dnsZoneFileResourceName, "dns_zone_config", regexp.MustCompile(`updated-example-service`)),
				),
			},
		},
	})
}

func testAccBrocadeVTMDNSZoneFileCheckDestroy(state *terraform.State, name string) error {

	log.Println("Checking DESTROY")
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	client.WorkWithConfigurationResources()
	zoneConfig := new([]byte)
	err := client.GetByName("dns_server/zone_files", name, zoneConfig)
	if err != nil {
		return nil
	}
	return fmt.Errorf("Error: resource %s still exists", name)
}

func testAccBrocadeVTMDNSZoneFileExists(dnsZoneFileName, dnsZoneResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		log.Println("Checking EXISTS")
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)

		client.WorkWithConfigurationResources()
		zoneConfig := new([]byte)
		err := client.GetByName("dns_server/zone_files", dnsZoneFileName, zoneConfig)
		if err != nil {
			return fmt.Errorf("Error: resource %s doesn't exists", dnsZoneFileName)
		}
		return nil
	}
}

func testAccBrocadeVTMDNSZoneFileNoNameTemplate() string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone_file" "acctest" {

}
`)
}

func testAccBrocadeDNSZoneFileCreateTemplate(name string) string {
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
;
; Services - Each service in a location has a unique IP address. Two locations = two IPs.
;
example-service         60  IN  A   10.100.10.5
                        60  IN  A   10.100.20.5
another-example-service             60  IN  A   10.100.10.6
                        60  IN  A   10.100.20.6
DNS_ZONE_CONFIG
}
`, name)
}

func testAccBrocadeDNSZoneFileUpdateTemplate(name string) string {
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
;
; Services - Each service in a location has a unique IP address. Two locations = two IPs.
;
updated-example-service 60  IN  A   10.100.10.5
                        60  IN  A   10.100.20.5
another-example-service             60  IN  A   10.100.10.6
                        60  IN  A   10.100.20.6
DNS_ZONE_CONFIG
}
`, name)
}
