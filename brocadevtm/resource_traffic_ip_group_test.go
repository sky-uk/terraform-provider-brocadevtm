package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/traffic_ip_group"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMTrafficIpGroupBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	trafficIPGroupName := fmt.Sprintf("acctest_brocadevtm_traffic_ip_group-%d", randomInt)
	trafficIPGroupResourceName := "brocadevtm_traffic_ip_group.acctest"

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
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidMode(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`must be one of singlehosted, ec2elastic, ec2vpcelastic, ec2vpcprivate, multihosted or rhi`),
			},
			{
				Config:      testAccBrocadeVTMTrafficIPGroupInvalidIPAddress(trafficIPGroupName),
				ExpectError: regexp.MustCompile(`ipaddresses: should be a list`),
			},
			{
				Config: testAccBrocadeVTMTrafficIPGroupCreateTemplate(trafficIPGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMTrafficIPGroupExists(trafficIPGroupName, trafficIPGroupResourceName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "name", trafficIPGroupName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "hashsourceport", "false"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ipaddresses.0", "192.168.100.10"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "mode", "singlehosted"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "multicastip", "192.168.100.11"),
					// TODO: can we get the name of the test traffic manager rather than hard coding?
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "trafficmanagers.0", "h1ist01-v00.paas.d50.ovp.bskyb.com"),
				),
			},
			{
				Config: testAccBrocadeVTMTrafficIPGroupUpdateTemplate(trafficIPGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMTrafficIPGroupExists(trafficIPGroupName, trafficIPGroupResourceName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "name", trafficIPGroupName),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "hashsourceport", "true"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "ipaddresses.0", "192.168.100.11"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "mode", "multihosted"),
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "multicastip", "192.168.100.12"),
					// TODO: can we get the name of the test traffic manager rather than hard coding?
					resource.TestCheckResourceAttr(trafficIPGroupResourceName, "trafficmanagers.0", "h1ist01-v00.paas.d50.ovp.bskyb.com"),
				),
			},
		},
	})
}

func testAccBrocadeVTMTrafficIPGroupCheckDestroy(state *terraform.State, name string) error {

	vtmClient := testAccProvider.Meta().(*brocadevtm.VTMClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_traffic_ip_group" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}

		api := trafficIpGroups.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return nil
		}
		if api.GetResponse().FilterByName(name).Name == name {
			return fmt.Errorf("Brocade vTM traffic IP group %s still exists", name)
		}
	}
	return nil
}

func testAccBrocadeVTMTrafficIPGroupExists(trafficIPGroupName, trafficIPGroupResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[trafficIPGroupResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Traffic IP Group resource %s not found in resources", trafficIPGroupName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Traffic IP Group ID not set in resources")
		}

		vtmClient := testAccProvider.Meta().(*brocadevtm.VTMClient)
		getAllAPI := trafficIpGroups.NewGetAll()

		err := vtmClient.Do(getAllAPI)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		foundTrafficIPGroup := getAllAPI.GetResponse().FilterByName(trafficIPGroupName)
		if foundTrafficIPGroup.Name != trafficIPGroupName {
			return fmt.Errorf("Brocade vTM Traffic IP Group %s not found on remote vTM", trafficIPGroupName)
		}
		return nil
	}
}

func testAccBrocadeVTMTrafficIPGroupNoName() string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
enabled = true
hashsourceport = false
ipaddresses = ["192.168.100.10"]
mode = "singlehosted"
multicastip = "192.168.100.11"
}
`)
}

func testAccBrocadeVTMTrafficIPGroupInvalidMode(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
name = "%s"
enabled = false
hashsourceport = false
ipaddresses = ["192.168.100.10"]
mode = "SOME_INVALID_MODE"
multicastip = "192.168.100.11"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupInvalidIPAddress(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
name = "%s"
enabled = false
hashsourceport = false
ipaddresses = "192.168.100.10"
mode = "multihosted"
multicastip = "192.168.100.11"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupCreateTemplate(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
name = "%s"
enabled = false
hashsourceport = false
ipaddresses = ["192.168.100.10"]
mode = "singlehosted"
multicastip = "192.168.100.11"
}
`, trafficIPGroupName)
}

func testAccBrocadeVTMTrafficIPGroupUpdateTemplate(trafficIPGroupName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_traffic_ip_group" "acctest" {
name = "%s"
enabled = true
hashsourceport = true
ipaddresses = ["192.168.100.11"]
mode = "multihosted"
multicastip = "192.168.100.12"
}
`, trafficIPGroupName)
}
