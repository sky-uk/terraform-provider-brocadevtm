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
	"net/http"
)

func TestAccBrocadeVTMTrafficManagerBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	trafficManagerName := fmt.Sprintf("acctest_brocadevtm_traffic_manager-%d", randomInt)
	trafficManagerResourceName := "brocadevtm_traffic_manager.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMTrafficManagerCheckDestroy(state, trafficManagerName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeTrafficManagerNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{

				Config: testAccBrocadeTrafficManagerCreate(trafficManagerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMTrafficManagerExists(trafficManagerName, trafficManagerResourceName),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "name", trafficManagerName),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "admin_master_xmlip", "1.1.1.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "admin_slave_xmlip", "1.1.1.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.name", "cardOne"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.interfaces.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.interfaces.0", "interface1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.interfaces.1", "interface2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.label", "eth0:0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.name", "cardTwo"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.interfaces.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.interfaces.0", "interface3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.interfaces.1", "interface4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.label", "eth0:1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.0.sysctl", "sysctrl1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.0.description", "sysctrl1 description"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.0.value", "valueOne"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.1.sysctl", "sysctrl2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.1.description", "sysctrl2 description"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.1.value", "valueTwo"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "authentication_server_ip", "0.0.0.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "location", "locationtest"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "nameip", "1.1.1.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "num_aptimizer_threads", "0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "num_children", "0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "number_of_cpus", "0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_server_port", "10"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "name"), "networkinterface1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.#"), "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.0"), "0.0.0.0/24"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.1"), "1.1.1.1/24"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "name"), "networkinterface2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.0"), "0.0.0.0/20"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.1"), "1.1.1.1/20"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "updater_ip", "0.0.0.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.gateway_ipv4", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.gateway_ipv6", "2001:db8:0:1234:0:567:8:2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.hosts", "name"), "host1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.hosts", "ip_address"), "127.0.0.1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.hosts", "name"), "host2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.hosts", "ip_address"), "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "name"), "if1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "autoneg"), "true"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "bond"), "bond1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "duplex"), "false"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "mtu"), "1500"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "speed"), "100"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "name"), "if2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "autoneg"), "false"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "bond"), "bond2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "duplex"), "true"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "mtu"), "1550"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "speed"), "10"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.0.name", "eth0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.0.addr", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.0.isexternal", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.0.mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.1.name", "eth1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.1.addr", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.1.isexternal", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.1.mask", "255.255.254.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_access", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_addr", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_gateway", "192.168.4.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_ipsrc", "dhcp"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipv4_forwarding", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipv6_forwarding", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.licence_agreed", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.manageazureroutes", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.manageec2conf", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.manageiptrans", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.managereturnpath", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.managevpcconf", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.0", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.1", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.0", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.1", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.#", "1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.routes", "name"), "127.0.0.1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.routes", "gw"), "192.168.4.3"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.routes", "if"), "if1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.routes", "mask"), "255.255.255.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.0", "searchdomain1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.1", "searchdomain2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_client_id", "id1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_client_key", "key1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_enabled", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_ips", "ip1,ip2,ip3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_load_balance", "priority"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_log_level", "debug"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_mode", "local"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_portal_url", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_proxy_host", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_proxy_port", "444"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_enabled", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_password_allowed", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_port", "22"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.timezone", "UTC"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.0", "vlan.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.1", "vlan.2"),
				),
			},

			{
				Config: testAccBrocadeTrafficManagerUpdate(trafficManagerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMTrafficManagerExists(trafficManagerName, trafficManagerResourceName),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "name", trafficManagerName),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "admin_master_xmlip", "1.1.1.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "admin_slave_xmlip", "1.1.1.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "name"), "cardOneUpdated"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "interfaces.#"), "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "interfaces.0"), "interface5"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "interfaces.1"), "interface6"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "label"), "eth0:2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "name"), "cardTwoUpdated"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "interfaces.0"), "interface7"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "interfaces.1"), "interface8"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_card", "label"), "eth0:3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_sysctl", "sysctl"), "sysctrl1Updated"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_sysctl", "description"), "sysctrl1 description Updated"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_sysctl", "value"), "valueOneUpdated"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_sysctl", "sysctl"), "sysctrl2Updated"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_sysctl", "description"), "sysctrl2 description Updated"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance_sysctl", "value"), "valueTwoUpdated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "authentication_server_ip", "0.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "location", "locationtestUpdated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "nameip", "2.2.2.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "num_aptimizer_threads", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "num_children", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "number_of_cpus", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_server_port", "15"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "name"), "networkinterface1Updated"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.#"), "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.0"), "1.0.0.0/24"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.1"), "2.1.1.1/24"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "name"), "networkinterface2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.0"), "1.0.0.0/20"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("trafficip", "networks.1"), "2.1.1.1/20"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "updater_ip", "0.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.gateway_ipv4", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.gateway_ipv6", "2001:db8:0:1234:0:567:8:3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.hosts", "name"), "host3"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.hosts", "ip_address"), "127.0.0.3"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.hosts", "name"), "host4"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.hosts", "ip_address"), "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "name"), "if3"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "autoneg"), "false"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "bond"), "bond3"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "duplex"), "true"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "mtu"), "3000"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "speed"), "1000"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "name"), "if4"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "autoneg"), "true"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "bond"), "bond4"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "duplex"), "false"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "mtu"), "1800"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.if", "speed"), "100"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.#", "2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.ip", "name"), "eth0"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.ip", "addr"), "127.0.0.1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.ip", "isexternal"), "false"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.ip", "mask"), "255.255.255.0"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.ip", "name"), "eth1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.ip", "addr"), "127.0.0.2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.ip", "isexternal"), "false"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.ip", "mask"), "255.255.254.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_access", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_addr", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_gateway", "192.168.4.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_ipsrc", "static"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipmi_lan_mask", "255.255.254.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipv4_forwarding", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ipv6_forwarding", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.licence_agreed", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.manageazureroutes", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.manageec2conf", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.manageiptrans", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.managereturnpath", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.managevpcconf", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.0", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.1", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.0", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.1", "127.0.0.5"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.#", "1"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.routes", "name"), "127.0.0.2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.routes", "gw"), "192.168.4.4"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.routes", "if"), "if2"),
					util.AccTestCheckValueInKeyPattern(trafficManagerResourceName, compileRegex("appliance.0.routes", "mask"), "255.255.254.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.0", "searchdomain3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.1", "searchdomain4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_client_id", "id2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_client_key", "key2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_enabled", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_ips", "ip1,ip2,ip3,ip4,ip5"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_load_balance", "round_robin"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_log_level", "info"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_mode", "portal"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_portal_url", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_proxy_host", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.shim_proxy_port", "445"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_enabled", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_password_allowed", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_port", "44"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.timezone", "GMT"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.0", "vlan.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.1", "vlan.4"),
				),
			},
		},
	})
}

func testAccBrocadeVTMTrafficManagerCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_traffic_manager" {
			continue
		}
		trafficManagerConfiguration := make(map[string]interface{})

		err := client.GetByName("traffic_managers", rs.Primary.ID, &trafficManagerConfiguration)
		if client.StatusCode == http.StatusOK {
			return fmt.Errorf("[ERROR] Brocade vTM Check Destroy Error: Traffic Manager %s still exists", name)
		}
		if client.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("[ERROR] Brocade vTM Check Destroy Error: Traffic Manager %+v ", err)
	}
	return nil
}

func testAccBrocadeVTMTrafficManagerExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("[ERROR] No ID is set")
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()
		trafficManagerConfiguration := make(map[string]interface{})
		err := client.GetByName("traffic_managers", name, &trafficManagerConfiguration)
		if client.StatusCode != http.StatusOK {
			return fmt.Errorf("[ERROR] Brocade vTM error whilst retrieving VTM Traffic Managers: %+v", err)
		}
		return nil
	}
}

func testAccBrocadeTrafficManagerCreate(name string) string {
	return fmt.Sprintf(`
       resource "brocadevtm_traffic_manager" "acctest" {
       name = "%s"
       admin_master_xmlip = "1.1.1.1"
       admin_slave_xmlip = "1.1.1.2"
       appliance_card= [
       		{
       			name = "cardOne"
       			interfaces = ["interface1","interface2"]
       			label = "eth0:0"
       		},
       		{
       			name = "cardTwo"
       			interfaces = ["interface3","interface4"]
       			label = "eth0:1"
       		}
       ]
       appliance_sysctl= [
       		{
       			sysctl = "sysctrl1"
       			description = "sysctrl1 description"
       			value = "valueOne"
       		},
       		{
       			sysctl = "sysctrl2"
       			description = "sysctrl2 description"
       			value = "valueTwo"
       		}
       ]

       authentication_server_ip = "0.0.0.0"
       location = "locationtest"
       nameip = "1.1.1.1"
       num_aptimizer_threads = 0
       num_children = 0
       number_of_cpus = 0
       rest_server_port = 10

        trafficip= [
          {
            name = "networkinterface1"
            networks = ["0.0.0.0/24","1.1.1.1/24"]
          },
          {
            name = "networkinterface2"
            networks = ["0.0.0.0/20","1.1.1.1/20"]
          }
       ]

       updater_ip = "0.0.0.0"

       appliance = {
           gateway_ipv4 = "127.0.0.1"
	   gateway_ipv6 = "2001:db8:0:1234:0:567:8:2"

	   hosts = [
	   	{
	   		name = "host1"
	   		ip_address = "127.0.0.1"
	   	},
	   	{
	   		name = "host2"
	   		ip_address = "127.0.0.2"

	   	}
	   ]
	   if = [
          {
            name = "if1"
            autoneg = true
            bond = "bond1"
            duplex = false
            mtu = 1500
            speed = "100"
          },
          {
            name = "if2"
            autoneg = false
            bond = "bond2"
            duplex = true
            mtu = 1550
            speed = "10"
          }

       ]
       ip = [
        {
         name = "eth0"
         addr = "127.0.0.1"
         isexternal = false
         mask = "255.255.255.0"
        },
        {
         name = "eth1"
         addr = "127.0.0.2"
         isexternal = true
         mask = "255.255.254.0"
        }
       ]

        ipmi_lan_access = true
       ipmi_lan_addr = "127.0.0.1"
       ipmi_lan_gateway = "192.168.4.3"
       ipmi_lan_ipsrc = "dhcp"
       ipmi_lan_mask = "255.255.255.0"
       ipv4_forwarding = true
       ipv6_forwarding = true
       licence_agreed = true
       manageazureroutes = true
       manageec2conf = true
       manageiptrans = true
       managereturnpath = true
       managevpcconf = true
       name_servers = ["127.0.0.1","127.0.0.2"]
       ntpservers = ["127.0.0.3","127.0.0.4"]

       routes = [
          {
            name = "127.0.0.1"
            gw = "192.168.4.3"
            if = "if1"
            mask = "255.255.255.0"
          }
       ]

       search_domains = ["searchdomain1","searchdomain2"]
       shim_client_id = "id1"
       shim_client_key = "key1"
       shim_enabled = true
       shim_ips = "ip1,ip2,ip3"
       shim_load_balance = "priority"
       shim_log_level = "debug"
       shim_mode = "local"
       shim_portal_url = "127.0.0.1"
       shim_proxy_host = "127.0.0.2"
       shim_proxy_port = 444
       ssh_enabled = true
       ssh_password_allowed = true
       ssh_port = 22
       timezone = "UTC"
       vlans = ["vlan.1","vlan.2"]

       }
}
`, name)
}

func testAccBrocadeTrafficManagerUpdate(name string) string {
	return fmt.Sprintf(`
       resource "brocadevtm_traffic_manager" "acctest" {
       name = "%s"
       admin_master_xmlip = "1.1.1.3"
       admin_slave_xmlip = "1.1.1.4"
       appliance_card= [
       		{
       			name = "cardOneUpdated"
       			interfaces = ["interface5","interface6"]
       			label = "eth0:2"
       		},
       		{
       			name = "cardTwoUpdated"
       			interfaces = ["interface7","interface8"]
       			label = "eth0:3"
       		}
       ]
       appliance_sysctl= [
       		{
       			sysctl = "sysctrl1Updated"
       			description = "sysctrl1 description Updated"
       			value = "valueOneUpdated"
       		},
       		{
       			sysctl = "sysctrl2Updated"
       			description = "sysctrl2 description Updated"
       			value = "valueTwoUpdated"
       		}
       ]

       authentication_server_ip = "0.0.0.1"
       location = "locationtestUpdated"
       nameip = "2.2.2.2"
       num_aptimizer_threads = 1
       num_children = 1
       number_of_cpus = 1
       rest_server_port = 15

        trafficip= [
          {
            name = "networkinterface1Updated"
            networks = ["1.0.0.0/24","2.1.1.1/24"]
          },
          {
            name = "networkinterface2"
            networks = ["1.0.0.0/20","2.1.1.1/20"]
          }
       ]

       updater_ip = "0.0.0.1"
	appliance = {
           gateway_ipv4 = "127.0.0.2"
	   gateway_ipv6 = "2001:db8:0:1234:0:567:8:3"

	   hosts = [
	   	{
	   		name = "host3"
	   		ip_address = "127.0.0.3"
	   	},
	   	{
	   		name = "host4"
	   		ip_address = "127.0.0.4"

	   	}
	   ]
	   if = [
          {
            name = "if3"
            autoneg = false
            bond = "bond3"
            duplex = true
            mtu = 3000
            speed = "1000"
          },
          {
            name = "if4"
            autoneg = true
            bond = "bond4"
            duplex = false
            mtu = 1800
            speed = "100"
          }

       ]
       ip = [
        {
         name = "eth0"
         addr = "127.0.0.1"
         isexternal = false
         mask = "255.255.255.0"
        },
        {
         name = "eth1"
         addr = "127.0.0.2"
         isexternal = true
         mask = "255.255.254.0"
        }
       ]

       ipmi_lan_access = false
       ipmi_lan_addr = "127.0.0.2"
       ipmi_lan_gateway = "192.168.4.4"
       ipmi_lan_ipsrc = "static"
       ipmi_lan_mask = "255.255.254.0"
       ipv4_forwarding = false
       ipv6_forwarding = false
       licence_agreed = false
       manageazureroutes = false
       manageec2conf = false
       manageiptrans = false
       managereturnpath = false
       managevpcconf = false
       name_servers = ["127.0.0.3","127.0.0.4"]
       ntpservers = ["127.0.0.4","127.0.0.5"]

       routes = [
          {
            name = "127.0.0.2"
            gw = "192.168.4.4"
            if = "if2"
            mask = "255.255.254.0"
          }
       ]

       search_domains = ["searchdomain3","searchdomain4"]
       shim_client_id = "id2"
       shim_client_key = "key2"
       shim_enabled = false
       shim_ips = "ip1,ip2,ip3,ip4,ip5"
       shim_load_balance = "round_robin"
       shim_log_level = "info"
       shim_mode = "portal"
       shim_portal_url = "127.0.0.2"
       shim_proxy_host = "127.0.0.3"
       shim_proxy_port = 445
       ssh_enabled = false
       ssh_password_allowed = false
       ssh_port = 44
       timezone = "GMT"
       vlans = ["vlan.3","vlan.4"]
       }
}
`, name)
}

func testAccBrocadeTrafficManagerNoName() string {
	return (`
       resource "brocadevtm_traffic_manager" "acctest" {

	}
`)
}
