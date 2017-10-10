package brocadevtm

/*
import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api/dns_zone_file"
	"github.com/sky-uk/go-rest-api"
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
					resource.TestMatchResourceAttr(dnsZoneFileResourceName, "dns_zone_file", regexp.MustCompile(`example-service`)),
				),
			},
			{
				Config: testAccBrocadeDNSZoneFileUpdateTemplate(dnsZoneFileName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMDNSZoneFileExists(dnsZoneFileName, dnsZoneFileResourceName),
					resource.TestCheckResourceAttr(dnsZoneFileResourceName, "name", dnsZoneFileName),
					resource.TestMatchResourceAttr(dnsZoneFileResourceName, "dns_zone_file", regexp.MustCompile(``)),
					resource.TestMatchResourceAttr(dnsZoneFileResourceName, "dns_zone_file", regexp.MustCompile(`updated-example-service`)),
				),
			},
		},
	})
}

func testAccBrocadeVTMDNSZoneFileCheckDestroy(state *terraform.State, name string) error {

	vtmClient := testAccProvider.Meta().(*rest.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_dns_zone_file" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		api := dnsZoneFile.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Brocade vTM DNS Zone File - error occurred whilst retrieving a list of all DNS zone files")
		}
		for _, dnsZoneFile := range api.ResponseObject().(*dnsZoneFile.DNSZoneFiles).Children {
			if dnsZoneFile.Name == name {
				return fmt.Errorf("Brocade vTM DNS zone file %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMDNSZoneFileExists(dnsZoneFileName, dnsZoneResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[dnsZoneResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM DNS zone file %s wasn't found in resources", dnsZoneFileName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM DNS zone file ID not set for %s in resources", dnsZoneFileName)
		}

		vtmClient := testAccProvider.Meta().(*rest.Client)
		api := dnsZoneFile.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		for _, dnsZoneFile := range api.ResponseObject().(*dnsZoneFile.DNSZoneFiles).Children {
			if dnsZoneFile.Name == dnsZoneFileName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM DNS zone file %s not found on remote vTM", dnsZoneFileName)
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
  dns_zone_file = <<DNS_ZONE_FILE
$TTL 3600
@				30	IN	SOA	ns1.example.com. hostmaster.isp.sky.com. (
							01	; serial
							3600	; refresh after 1 hour
							300	; retry after 5 minutes
							1209600	; expire after 2 weeks
							30)	; minimum TTL of 30 seconds
@				30	IN	NS	ns1.example.com.
ns1				30	IN	A	10.0.0.2
example-service			60	IN	A	10.1.0.2
				60	IN	A	10.1.1.2
another-example-service		60	IN	A	10.2.0.2
				60	IN	A	10.2.1.2
DNS_ZONE_FILE
}
`, name)
}

func testAccBrocadeDNSZoneFileUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone_file" "acctest" {
  name = "%s"
  dns_zone_file = <<DNS_ZONE_FILE
$TTL 3600
@ 				30	IN 	SOA 	ns2.example.com. hostmaster.isp.sky.com. (
							02	; serial
							3600	; refresh after 1 hour
							300	; retry after 5 minutes
							1209600	; expire after 2 weeks
							30)	; minimum TTL of 30 seconds
@				30	IN	NS	ns2.example.com.
ns1				30	IN	A	10.100.0.2
updated-example-service		30	IN	A	10.110.0.2
				30	IN	A	10.110.1.2
another-example-service		30	IN	A	10.120.0.2
				30	IN	A	10.120.1.2
DNS_ZONE_FILE
}
`, name)
}
*/
