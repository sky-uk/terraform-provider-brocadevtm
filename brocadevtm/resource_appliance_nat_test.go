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

func compileRegex(key, attr string) *regexp.Regexp {
	return regexp.MustCompile(key + `\.[0-9]+\.` + attr)
}

func TestAccBrocadeVTMApplianceNat(t *testing.T) {

	randomInt := acctest.RandInt()
	name := fmt.Sprintf("acctest_brocadevtm_appliance_nat-%d", randomInt)
	resourceName := fmt.Sprintf("brocadevtm_appliance_nat.acctest")
	fmt.Printf("\n\nAppliance Nat is %s.\n\n", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMApplianceNatCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBrocadeVTMApplianceNatEmptyResource(),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
				),
			},
			{
				Config: manyToOneCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "rule_number"), "4765348"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "tip"), name),
				),
			},
			{
				Config: manyToOneUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "rule_number"), "58673"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "tip"), name),
				),
			},
			{
				Config: manyToOnePortLockedCreateTemplate(name),
				//Destroy:                   false,
				//PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "rule_number"), "20001"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "port"), "2728"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "tip"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "protocol"), "tcp"),
				),
			},
			{
				Config: manyToOnePortLockedUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "rule_number"), "20001"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "port"), "2730"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "tip"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "protocol"), "tcp"),
				),
			},
			{
				Config: oneToOneCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "rule_number"), "20001"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "port"), "2728"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "tip"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "protocol"), "tcp"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("one_to_one", "rule_number"), "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("one_to_one", "rule_number"), "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("one_to_one", "ip"), "192.168.10.10"),
				),
			},
			{
				Config: oneToOneUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "rule_number"), "20002"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "port"), "2730"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "tip"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "protocol"), "tcp"),
				),
			},
			{
				Config: portMappingCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("port_mapping", "dport_first"), "10"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("port_mapping", "dport_last"), "20"),
				),
			},
			{
				Config: portMappingUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("port_mapping", "dport_first"), "10"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("port_mapping", "dport_last"), "30"),
				),
			},
		},
	})
}

func testAccBrocadeVTMApplianceNatCheckDestroy(state *terraform.State, name string) error {

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_appliance_nat" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()

		resources, err := client.GetAllResources("appliance/nat")
		if err != nil {
			return fmt.Errorf("[ERROR] Brocade vTM error whilst retrieving appliance nat: %+v", err)
		}
		for _, resource := range resources {
			if resource["name"] == name {
				return fmt.Errorf("[ERROR] Brocade vTM Appliance nat %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMApplianceNatExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Brocade vTM Appliance Nat %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Brocade vTM Appliance Nat ID not set for %s in resources", name)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()
		nat := make(map[string]interface{})
		err := client.GetByName("appliance/nat", "", &nat)
		if err != nil {
			return fmt.Errorf("[ERROR] Brocade vTM error whilst retriving appliance nat: %v", err)
		}
		return nil
	}
}

func testAccBrocadeVTMApplianceNatEmptyResource() string {
	return fmt.Sprintf(`
resource "brocadevtm_appliance_nat" "acctest" {
  many_to_one_all_ports = []
  many_to_one_port_locked = []
  one_to_one = []
  port_mapping = []
}
`)
}

func manyToOneCreateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = [{
	  rule_number = 4765348
	  pool = "%s"
	  tip = "%s"
  }]
  many_to_one_port_locked = []
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func manyToOneUpdateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = [{
	  rule_number = 58673
	  pool = "%s"
	  tip = "%s"
  }]
  many_to_one_port_locked = []
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func manyToOnePortLockedCreateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2728
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func manyToOnePortLockedUpdateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2730
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func oneToOneCreateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2728
	  protocol = "tcp"
  }]
  one_to_one = [{
	  rule_number = 1
	  enable_inbound = false
	  ip = "192.168.10.10"
	  tip = "%s"
	},{
	  rule_number = 2
	  enable_inbound = false
	  ip = "192.168.10.11"
	  tip = "%s"
	}
  ]
  port_mapping = []
}
`, name, name, name, name, name, name)
}

func oneToOneUpdateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20002
	  pool = "%s"
	  tip = "%s"
	  port = 2730
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func portMappingCreateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_virtual_server" "acctest" {

	name = "%s"
	bandwidth_class = "testUpdate"
	completion_rules = ["completionRule2","completionRule3"]
	connect_timeout = 100
	enabled = false
	glb_services = ["testservice3","testservice4"]
	listen_on_any = false
	listen_on_hosts = ["host3","host4"]
	listen_on_traffic_ips = ["ip1"]
	note = "update acceptance test"
	pool = "test-pool"
	port = 100
	protection_class = "testProtectionClassUpdate"
	protocol = "ftp"
	request_rules = ["ruleThree"]
	response_rules = ["ruleFour"]
	slm_class = "testClassUpdate"
	ssl_decrypt = false
	transparent = false

}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = [ "brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest", "brocadevtm_virtual_server.acctest" ]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2728
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = [{
	  rule_number = 1
	  dport_first = 10
	  dport_last = 20
	  virtual_server = "%s"
  }]
}
`, name, name, name, name, name, name)
}

func portMappingUpdateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_virtual_server" "acctest" {

	name = "%s"
	bandwidth_class = "testUpdate"
	completion_rules = ["completionRule2","completionRule3"]
	connect_timeout = 100
	enabled = false
	glb_services = ["testservice3","testservice4"]
	listen_on_any = false
	listen_on_hosts = ["host3","host4"]
	listen_on_traffic_ips = ["ip1"]
	note = "update acceptance test"
	pool = "test-pool"
	port = 100
	protection_class = "testProtectionClassUpdate"
	protocol = "ftp"
	request_rules = ["ruleThree"]
	response_rules = ["ruleFour"]
	slm_class = "testClassUpdate"
	ssl_decrypt = false
	transparent = false

}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = [ "brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest", "brocadevtm_virtual_server.acctest" ]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2728
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = [{
	  rule_number = 1
	  dport_first = 10
	  dport_last = 30
	  virtual_server = "%s"
  }]
}
`, name, name, name, name, name, name)
}
