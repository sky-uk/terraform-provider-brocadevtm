package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMTrafficIpGroupBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	trafficIPGroupName := fmt.Sprintf("acctest_brocadevtm_traffic_ip_group-%d", randomInt)
	trafficIPGroupResourceName := "brocadevtm_traffic_ip_group.acctest"
	backendTrafficIpsPattern := regexp.MustCompile(`backend_traffic_ips\.[0-9]+`)
	ipMappingIPPattern := regexp.MustCompile(`ip_mapping\.[0-9]+\.ip`)
	ipMappingTMPattern := regexp.MustCompile(`ip_mapping\.[0-9]+\.traffic_manager`)
	ipAddressesPattern := regexp.MustCompile(`ipaddresses\.[0-9]+`)
	slavesPattern := regexp.MustCompile(`slaves\.[0-9]+`)

	fmt.Printf("\n\nTraffic IP Group is %s.\n\n", trafficIPGroupName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMTrafficIPGroupCheckDestroy(state, trafficIPGroupName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMTrafficIPGroupNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidBackendTrafficIPs(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`backend_traffic_ips: should be a list`),
			},
			{
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidIPAssignmentMode(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`expected ip_assignment_mode to be one of \[alphabetic balanced\], got SOME_INVALID_MODE`),
			},
			{
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidMode(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`expected mode to be one of \[singlehosted ec2elastic ec2vpcelastic ec2vpcprivate multihosted rhi\], got SOME_INVALID_MODE`),
			},
			{
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidIPAddress(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`ipaddresses: should be a list`),
			},
			{
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidMulticastIP(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`must be a valid multicast IP \(224.0.0.0 - 239.255.255.255\)`),
			},
			{
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidUnsignedInt(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`can't be negative`),
			},
			{
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidRHIProtocol(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`expected rhi_protocols to be one of \[ospf bgp\], got INVALID_PROTOCOL`),
			},
			{
				Config: testAccBrocadeVTMTrafficIPGroupCreateTemplate(trafficIPGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMTrafficIPGroupExists(trafficIPGroupName, trafficIPGroupResourceName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "name", trafficIPGroupName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "backend_traffic_ips.#", "1"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, backendTrafficIpsPattern, "192.168.10.49"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "hash_source_port", "true"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ip_assignment_mode", "alphabetic"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ip_mapping.#", "1"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingIPPattern, "192.168.34.56"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingTMPattern, "10.93.59.27"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ipaddresses.#", "1"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipAddressesPattern, "192.168.100.10"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "keeptogether", "true"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "location", "10"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "mode", "singlehosted"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "multicast", "232.123.23.45"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "note", "Acceptance test - create"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_bgp_metric_base", "5"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_bgp_passive_metric_offset", "2"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_ospfv2_metric_base", "7"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_ospfv2_passive_metric_offset", "3"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_protocols", "ospf"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "slaves.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, slavesPattern, "192.168.34.45"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, slavesPattern, "192.168.34.46"),
				),
			},
			{
				Config: testAccBrocadeVTMTrafficIPGroupUpdateTemplate(trafficIPGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMTrafficIPGroupExists(trafficIPGroupName, trafficIPGroupResourceName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "name", trafficIPGroupName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "backend_traffic_ips.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, backendTrafficIpsPattern, "192.168.10.92"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, backendTrafficIpsPattern, "192.168.10.86"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "hash_source_port", "false"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ip_assignment_mode", "balanced"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ip_mapping.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingIPPattern, "192.168.34.56"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingIPPattern, "192.168.34.64"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingTMPattern, "10.93.59.27"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingTMPattern, "10.93.59.27"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ipaddresses.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipAddressesPattern, "192.168.100.11"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipAddressesPattern, "192.168.100.12"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "keeptogether", "false"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "location", "12"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "mode", "multihosted"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "multicast", "232.123.23.143"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "note", "Acceptance test - update"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_bgp_metric_base", "15"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_bgp_passive_metric_offset", "3"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_ospfv2_metric_base", "17"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_ospfv2_passive_metric_offset", "5"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_protocols", "bgp"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "slaves.#", "3"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, slavesPattern, "192.168.34.47"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, slavesPattern, "192.168.34.46"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, slavesPattern, "192.168.34.45"),
				),
			},
			{
				Config: testAccBrocadeVTMTrafficIPGroupUpdate2Template(trafficIPGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMTrafficIPGroupExists(trafficIPGroupName, trafficIPGroupResourceName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "name", trafficIPGroupName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "backend_traffic_ips.#", "1"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, backendTrafficIpsPattern, "192.168.10.28"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "hash_source_port", "true"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ip_assignment_mode", "alphabetic"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ip_mapping.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingIPPattern, "192.168.34.64"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingIPPattern, "192.168.34.56"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingTMPattern, "10.93.59.27"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipMappingTMPattern, "10.93.59.27"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ipaddresses.#", "1"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, ipAddressesPattern, "192.168.100.12"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "keeptogether", "true"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "location", "5"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "mode", "singlehosted"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "multicast", "232.123.23.48"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "note", "Acceptance test - update 2"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_bgp_metric_base", "12"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_bgp_passive_metric_offset", "4"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_ospfv2_metric_base", "14"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_ospfv2_passive_metric_offset", "4"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "rhi_protocols", "ospf"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "slaves.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, slavesPattern, "192.168.34.46"),
					util.AccTestCheckValueInKeyPattern(trafficIPGroupResourceName, slavesPattern, "192.168.34.45"),
				),
			},
		},
	})
}

func testAccBrocadeVTMTrafficIPGroupCheckDestroy(state *terraform.State, name string) error {

	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_traffic_ip_group" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}

		trafficIPGroups, err := client.GetAllResources("traffic_ip_groups")
		if err != nil {
			return fmt.Errorf("[ERROR] Brocade vTM traffic IP group error retrieving the list of traffic IP groups")
		}
		for _, trafficIPGroupItem := range trafficIPGroups {
			if trafficIPGroupItem["name"] == name {
				return fmt.Errorf("[ERROR] Brocade vTM traffic IP group %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMTrafficIPGroupExists(trafficIPGroupName, trafficIPGroupResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[trafficIPGroupResourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Brocade vTM Traffic IP Group resource %s not found in resources", trafficIPGroupName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Brocade vTM Traffic IP Group ID not set in resources")
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		trafficIPGroups, err := client.GetAllResources("traffic_ip_groups")

		if err != nil {
			return fmt.Errorf("[ERROR] %+v", err)
		}
		for _, trafficIPGroupItem := range trafficIPGroups {
			if trafficIPGroupItem["name"] == trafficIPGroupName {
				return nil
			}
		}
		return fmt.Errorf("[ERROR] Brocade vTM Traffic IP Group %s not found on remote vTM", trafficIPGroupName)
	}
}

func testAccBrocadeVTMTrafficIPGroupNoName() string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
}
`)
}

func testAccBrocadeVTMTrafficIPGroupInvalidBackendTrafficIPs(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  backend_traffic_ips = "192.168.10.49"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupInvalidIPAssignmentMode(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  ip_assignment_mode = "SOME_INVALID_MODE"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupInvalidMode(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  mode = "SOME_INVALID_MODE"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupInvalidIPAddress(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  ipaddresses = "192.168.100.10"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupInvalidMulticastIP(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  multicast = "192.168.100.11"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupInvalidUnsignedInt(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  rhi_bgp_metric_base = -1
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupInvalidRHIProtocol(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  rhi_protocols = "INVALID_PROTOCOL"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupCreateTemplate(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  backend_traffic_ips=["192.168.10.49"]
  enabled = true
  hash_source_port = true
  ip_assignment_mode = "alphabetic"
  ip_mapping = [
    {
      ip = "192.168.34.56"
      traffic_manager = "10.93.59.27"
    },
  ]
  ipaddresses = ["192.168.100.10"]
  keeptogether = true
  location = 10
  mode = "singlehosted"
  multicast = "232.123.23.45"
  note = "Acceptance test - create"
  rhi_bgp_metric_base = 5
  rhi_bgp_passive_metric_offset = 2
  rhi_ospfv2_metric_base = 7
  rhi_ospfv2_passive_metric_offset = 3
  rhi_protocols = "ospf"
  slaves = [ "192.168.34.45", "192.168.34.46" ]
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupUpdateTemplate(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  backend_traffic_ips=["192.168.10.92","192.168.10.86"]
  enabled = false
  hash_source_port = false
  ip_assignment_mode = "balanced"
  ip_mapping = [
    {
      ip = "192.168.34.56"
      traffic_manager = "10.93.59.27"
    },
    {
      ip = "192.168.34.64"
      traffic_manager = "10.93.59.27"
    },
  ]
  ipaddresses = ["192.168.100.11", "192.168.100.12"]
  keeptogether = false
  location = 12
  mode = "multihosted"
  multicast = "232.123.23.143"
  note = "Acceptance test - update"
  rhi_bgp_metric_base = 15
  rhi_bgp_passive_metric_offset = 3
  rhi_ospfv2_metric_base = 17
  rhi_ospfv2_passive_metric_offset = 5
  rhi_protocols = "bgp"
  slaves = [ "192.168.34.47", "192.168.34.46", "192.168.34.45" ]
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupUpdate2Template(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  backend_traffic_ips=["192.168.10.28"]
  enabled = true
  hash_source_port = true
  ip_assignment_mode = "alphabetic"
  ip_mapping = [
    {
      ip = "192.168.34.64"
      traffic_manager = "10.93.59.27"
    },
    {
      ip = "192.168.34.56"

      traffic_manager = "10.93.59.27"
    },
  ]
  ipaddresses = ["192.168.100.12"]
  keeptogether = true
  location = 5
  machines = [ "192.168.10.11", "10.93.59.27" ]
  mode = "singlehosted"
  multicast = "232.123.23.48"
  note = "Acceptance test - update 2"
  rhi_bgp_metric_base = 12
  rhi_bgp_passive_metric_offset = 4
  rhi_ospfv2_metric_base = 14
  rhi_ospfv2_passive_metric_offset = 4
  rhi_protocols = "ospf"
  slaves = [ "192.168.34.46", "192.168.34.45" ]
}
`, trafficIPGroupName)
}
