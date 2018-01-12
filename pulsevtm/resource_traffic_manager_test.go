package pulsevtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-pulse-vtm/api"
	"net/http"
)

func TestAccPulseVTMTrafficManagerBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	trafficManagerName := fmt.Sprintf("acctest_pulsevtm_traffic_manager-%d", randomInt)
	trafficManagerResourceName := "pulsevtm_traffic_manager.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccPulseVTMTrafficManagerCheckDestroy(state, trafficManagerName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccPulseTrafficManagerNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{

				Config: testAccPulseTrafficManagerCreate(trafficManagerName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMTrafficManagerExists(trafficManagerName, trafficManagerResourceName),
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
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.0.name", "networkinterface1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.0.networks.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.0.networks.0", "0.0.0.0/24"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.0.networks.1", "1.1.1.1/24"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.1.name", "networkinterface2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.1.networks.0", "0.0.0.0/20"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.1.networks.1", "1.1.1.1/20"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "updater_ip", "0.0.0.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.gateway_ipv4", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.gateway_ipv6", "2001:db8:0:1234:0:567:8:2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.0.name", "host1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.0.ip_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.1.name", "host2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.1.ip_address", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.name", "if1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.autoneg", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.bond", "bond1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.duplex", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.mtu", "1500"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.speed", "100"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.name", "if2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.autoneg", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.bond", "bond2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.duplex", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.mtu", "1550"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.speed", "10"),
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
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.manageservices", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.managevpcconf", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.0", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.1", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.0", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.1", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.0.name", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.0.gw", "192.168.4.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.0.if", "if1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.0.mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.0", "searchdomain1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.1", "searchdomain2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_enabled", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_password_allowed", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_port", "22"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.timezone", "UTC"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.0", "vlan.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.1", "vlan.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.0.allow_update", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.0.bind_ip", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.0.external_ip", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.0.port", "9080"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.bgp_router_id", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_ip", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_neighbor_addrs.#", "3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_neighbor_addrs.0", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_neighbor_addrs.1", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_neighbor_addrs.2", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.lss_dedicated_ips.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.lss_dedicated_ips.0", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.lss_dedicated_ips.1", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptables.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptables.0.config_enabled", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptrans.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptrans.0.fwmark", "50"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptrans.0.iptables_enabled", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptrans.0.routing_table", "300"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "java.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "java.0.port", "1500"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "remote_licensing.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "remote_licensing.0.email_address", "test1@testemail.dev"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "remote_licensing.0.message", "message1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.0.port", "500"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.0.bind_ips.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.0.bind_ips.0", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.0.bind_ips.1", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.allow.#", "3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.allow.0", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.allow.1", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.allow.2", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.auth_password", "testpassword"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.bind_ip", "127.0.0.5"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.community", "public"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.enabled", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.hash_algorithm", "md5"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.port", "50"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.priv_password", "privpassword"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.security_level", "noauthnopriv"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.username", "username1"),
				),
			},

			{
				Config: testAccPulseTrafficManagerUpdate(trafficManagerName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMTrafficManagerExists(trafficManagerName, trafficManagerResourceName),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "name", trafficManagerName),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "admin_master_xmlip", "1.1.1.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "admin_slave_xmlip", "1.1.1.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.name", "cardOneUpdated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.interfaces.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.interfaces.0", "interface5"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.interfaces.1", "interface6"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.0.label", "eth0:2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.name", "cardTwoUpdated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.interfaces.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.interfaces.0", "interface7"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.interfaces.1", "interface8"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_card.1.label", "eth0:3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.0.sysctl", "sysctrl1Updated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.0.description", "sysctrl1 description Updated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.0.value", "valueOneUpdated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.1.sysctl", "sysctrl2Updated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.1.description", "sysctrl2 description Updated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance_sysctl.1.value", "valueTwoUpdated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "authentication_server_ip", "0.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "location", "locationtestUpdated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "nameip", "2.2.2.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "num_aptimizer_threads", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "num_children", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "number_of_cpus", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_server_port", "15"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.0.name", "networkinterface1Updated"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.0.networks.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.0.networks.0", "1.0.0.0/24"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.0.networks.1", "2.1.1.1/24"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.1.name", "networkinterface2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.1.networks.0", "1.0.0.0/20"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "trafficip.1.networks.1", "2.1.1.1/20"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "updater_ip", "0.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.gateway_ipv4", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.gateway_ipv6", "2001:db8:0:1234:0:567:8:3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.0.name", "host3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.0.ip_address", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.1.name", "host4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.hosts.1.ip_address", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.name", "if3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.autoneg", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.bond", "bond3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.duplex", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.mtu", "3000"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.0.speed", "1000"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.name", "if4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.autoneg", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.bond", "bond4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.duplex", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.mtu", "1800"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.if.1.speed", "100"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.0.name", "eth0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.0.addr", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.0.isexternal", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.0.mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.1.name", "eth1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.1.addr", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.1.isexternal", "true"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ip.1.mask", "255.255.254.0"),
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
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.managereservedports", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.managereturnpath", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.manageservices", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.managevpcconf", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.0", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.name_servers.1", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.0", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ntpservers.1", "127.0.0.5"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.0.name", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.0.gw", "192.168.4.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.0.if", "if2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.routes.0.mask", "255.255.254.0"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.0", "searchdomain3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.search_domains.1", "searchdomain4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_enabled", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_password_allowed", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.ssh_port", "44"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.timezone", "GMT"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.0", "vlan.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "appliance.0.vlans.1", "vlan.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.0.allow_update", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.0.bind_ip", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.0.external_ip", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "cluster_comms.0.port", "9081"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.bgp_router_id", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_ip", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_neighbor_addrs.#", "3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_neighbor_addrs.0", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_neighbor_addrs.1", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.ospfv2_neighbor_addrs.2", "127.0.0.5"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.lss_dedicated_ips.#", "3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.lss_dedicated_ips.0", "127.0.0.1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.lss_dedicated_ips.1", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "fault_tolerance.0.lss_dedicated_ips.2", "127.0.0.2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptables.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptables.0.config_enabled", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptrans.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptrans.0.fwmark", "100"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptrans.0.iptables_enabled", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "iptrans.0.routing_table", "600"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "java.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "java.0.port", "3000"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "remote_licensing.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "remote_licensing.0.email_address", "test2@testemail.dev"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "remote_licensing.0.message", "message2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.0.port", "1000"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.0.bind_ips.#", "2"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.0.bind_ips.0", "127.0.0.3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "rest_api.0.bind_ips.1", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.#", "1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.allow.#", "3"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.allow.0", "127.0.0.4"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.allow.1", "127.0.0.5"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.allow.2", "127.0.0.6"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.auth_password", "testpasswordupdate"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.bind_ip", "127.0.0.6"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.community", "private"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.enabled", "false"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.hash_algorithm", "sha1"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.port", "100"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.priv_password", "privpasswordupdate"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.security_level", "authpriv"),
					resource.TestCheckResourceAttr(trafficManagerResourceName, "snmp.0.username", "usernameupdated"),
				),
			},
		},
	})
}

func testAccPulseVTMTrafficManagerCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "pulsevtm_traffic_manager" {
			continue
		}
		trafficManagerConfiguration := make(map[string]interface{})

		err := client.GetByName("traffic_managers", rs.Primary.ID, &trafficManagerConfiguration)
		if client.StatusCode == http.StatusOK {
			return fmt.Errorf("[ERROR] Pulse vTM Check Destroy Error: Traffic Manager %s still exists", name)
		}
		if client.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("[ERROR] Pulse vTM Check Destroy Error: Traffic Manager %+v ", err)
	}
	return nil
}

func testAccPulseVTMTrafficManagerExists(name, resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("[ERROR] Pulse vTM error whilst retrieving VTM Traffic Managers: %+v", err)
		}
		return nil
	}
}

func testAccPulseTrafficManagerCreate(name string) string {
	return fmt.Sprintf(`
       resource "pulsevtm_traffic_manager" "acctest" {
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
			mode = "static"
          },
          {
            name = "if2"
            autoneg = false
            bond = "bond2"
            duplex = true
            mtu = 1550
            speed = "10"
			mode = "static"
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
	       ssh_enabled = true
	       ssh_password_allowed = true
	       ssh_port = 22
	       timezone = "UTC"
	       vlans = ["vlan.1","vlan.2"]

       }

       cluster_comms = {
		allow_update = true
		bind_ip = "127.0.0.1"
		external_ip = "127.0.0.2"
		port = 9080
       }

	fault_tolerance = {
		bgp_router_id = "127.0.0.1"
		lss_dedicated_ips = ["127.0.0.3", "127.0.0.2"],
		ospfv2_ip = "127.0.0.1"
		ospfv2_neighbor_addrs = ["127.0.0.1","127.0.0.2","127.0.0.3"]
	}

	iptables = {
		config_enabled = true
	}

	iptrans = {
		fwmark = 50
		iptables_enabled = true
		routing_table = 300
	}

	java = {
		port = 1500
	}

	remote_licensing = {
		email_address = "test1@testemail.dev"
		message = "message1"
	}

	rest_api = {
		port = 500
		bind_ips = ["127.0.0.1","127.0.0.2"]
	}

	snmp = {
		allow = ["127.0.0.1","127.0.0.2","127.0.0.3"]
		auth_password = "testpassword"
		bind_ip = "127.0.0.5"
		community = "public"
		enabled = true
		hash_algorithm = "md5"
		port = 50
		priv_password = "privpassword"
		security_level = "noauthnopriv"
		username = "username1"
	}

}
`, name)
}

func testAccPulseTrafficManagerUpdate(name string) string {
	return fmt.Sprintf(`
       resource "pulsevtm_traffic_manager" "acctest" {
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
			mode = "static"
          },
          {
            name = "if4"
            autoneg = true
            bond = "bond4"
            duplex = false
            mtu = 1800
            speed = "100"
			mode = "static"
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
	   managereservedports = false
       managereturnpath = false
	   manageservices = false
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
       ssh_enabled = false
       ssh_password_allowed = false
       ssh_port = 44
       timezone = "GMT"
       vlans = ["vlan.3","vlan.4"]
       }

        cluster_comms = {
		allow_update = false
		bind_ip = "127.0.0.3"
		external_ip = "127.0.0.4"
		port = 9081
       }

	fault_tolerance = {
		bgp_router_id = "127.0.0.2"
		lss_dedicated_ips = ["127.0.0.1", "127.0.0.3", "127.0.0.2"],
		ospfv2_ip = "127.0.0.2"
		ospfv2_neighbor_addrs = ["127.0.0.3","127.0.0.4","127.0.0.5"]
	}

	iptables = {
		config_enabled = false
	}

	iptrans = {
		fwmark = 100
		iptables_enabled = false
		routing_table = 600
	}

	java = {
		port = 3000
	}

	remote_licensing = {
		email_address = "test2@testemail.dev"
		message = "message2"
	}

	rest_api = {
		port = 1000
		bind_ips = ["127.0.0.3","127.0.0.4"]
	}

	snmp = {
		allow = ["127.0.0.4","127.0.0.5","127.0.0.6"]
		auth_password = "testpasswordupdate"
		bind_ip = "127.0.0.6"
		community = "private"
		enabled = false
		hash_algorithm = "sha1"
		port = 100
		priv_password = "privpasswordupdate"
		security_level = "authpriv"
		username = "usernameupdated"
	}
}
`, name)
}

func testAccPulseTrafficManagerNoName() string {
	return (`
       resource "pulsevtm_traffic_manager" "acctest" {

	}
`)
}
