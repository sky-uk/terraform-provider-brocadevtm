package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/pool"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"net/http"
	"regexp"
	"testing"
)

func TestAccPool_Basic(t *testing.T) {

	randomInt := acctest.RandInt()
	poolName := fmt.Sprintf("acctest_brocadevtm_pool-%d", randomInt)
	poolResourceName := "brocadevtm_pool.acctest"

	nodePattern := regexp.MustCompile(`nodes_table\.[0-9]+\.node`)
	priorityPattern := regexp.MustCompile(`nodes_table\.[0-9]+\.priority`)
	sourceIPPattern := regexp.MustCompile(`nodes_table\.[0-9]+\.source_ip`)
	statePattern := regexp.MustCompile(`nodes_table\.[0-9]+\.state`)
	weightPattern := regexp.MustCompile(`nodes_table\.[0-9]+\.weight`)
	securityIDPattern := regexp.MustCompile(`auto_scaling\.[0-9]+\.securitygroupids\.[0-9]+`)
	subnetIDPattern := regexp.MustCompile(`auto_scaling\.[0-9]+\.subnetids\.[0-9]+`)
	dnsAutoScaleHostnamePattern := regexp.MustCompile(`dns_autoscale\.[0-9]+\.hostnames\.[0-9]+`)
	sslCNPattern := regexp.MustCompile(`ssl\.[0-9]+\.common_name_match\.[0-9]+`)
	sslEllipticCurvesPattern := regexp.MustCompile(`ssl\.[0-9]+\.elliptic_curves\.[0-9]+`)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPoolCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPoolNodeInvalidAlgo(poolName),
				ExpectError: regexp.MustCompile(`must be one of fastest_response_time, least_connections, perceptive, random, round_robin, weighted_least_connections, weighted_round_robin`),
			},

			{
				Config:      testAccPoolNodeUnsignedInt(poolName),
				ExpectError: regexp.MustCompile(`can't be negative`),
			},
			{
				Config:      testAccPoolNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPoolNoNodes(poolName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPoolInvalidNode(poolName),
				ExpectError: regexp.MustCompile(`must be a valid IP/Hostname and port seperated by a colon. i.e 127.0.0.1:80`),
			},
			{
				Config:      testAccPoolInvalidNodeNoIP(poolName),
				ExpectError: regexp.MustCompile(`must be a valid IP/Hostname and port seperated by a colon. i.e 127.0.0.1:80`),
			},
			{
				Config:      testAccPoolInvalidNodeNoPort(poolName),
				ExpectError: regexp.MustCompile(`must be a valid IP/Hostname and port seperated by a colon. i.e 127.0.0.1:80`),
			},
			{
				Config:      testAccPoolInvalidNodeDeleteBehaviour(poolName),
				ExpectError: regexp.MustCompile(`must be one of immediate or drain`),
			},
			{
				Config:      testAccPoolOneItemList(poolName),
				ExpectError: regexp.MustCompile(`attribute supports 1 item maximum, config has 2 declared`),
			},
			{
				Config:      testAccPoolValidAWSSGPrefix(poolName),
				ExpectError: regexp.MustCompile(`one or more items in the list of strings doesn't match the prefix sg-`),
			},
			{
				Config:      testAccPoolValidAWSSubnetPrefix(poolName),
				ExpectError: regexp.MustCompile(`one or more items in the list of strings doesn't match the prefix subnet-`),
			},
			{
				Config: testAccPoolCreateTemplate(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "nodes_table.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, nodePattern, "192.168.10.10:80"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, priorityPattern, "5"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, statePattern, "draining"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, weightPattern, "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sourceIPPattern, "192.168.120.6"),
					resource.TestCheckResourceAttr(poolResourceName, "bandwidth_class", "example"),
					resource.TestCheckResourceAttr(poolResourceName, "failure_pool", "test-pool"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "100"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "8"),
					resource.TestCheckResourceAttr(poolResourceName, "monitors.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "monitors.0", "Full HTTP"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "node_connection_attempts", "6"),
					resource.TestCheckResourceAttr(poolResourceName, "node_delete_behaviour", "immediate"),
					resource.TestCheckResourceAttr(poolResourceName, "node_drain_to_delete_timeout", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "note", "example test pool"),
					resource.TestCheckResourceAttr(poolResourceName, "passive_monitoring", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "persistence_class", "example"),
					resource.TestCheckResourceAttr(poolResourceName, "transparent", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.addnode_delaytime", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.cloud_credentials", "example"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.cluster", "10.0.0.1"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.data_center", "vCentre server"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.external", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.hysteresis", "100"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.imageid", "image id"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.ips_to_use", "private_ips"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.last_node_idle_time", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.max_nodes", "100"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.min_nodes", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.name", "example"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.port", "8980"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.refactory", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.response_time", "100"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.scale_down_level", "90"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.scale_up_level", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.securitygroupids.#", "3"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, securityIDPattern, "sg-12345"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, securityIDPattern, "sg-23456"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, securityIDPattern, "sg-34567"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.size_id", "sizeID"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.subnetids.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, subnetIDPattern, "subnet-xxxx"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, subnetIDPattern, "subnet-xxxx"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.max_connect_time", "4"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.max_connections_per_node", "100"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.max_queue_size", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.max_reply_time", "12"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.queue_timeout", "14"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.0.enabled", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.0.hostnames.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, dnsAutoScaleHostnamePattern, "example01.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, dnsAutoScaleHostnamePattern, "example02.example.com"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.0.port", "8080"),
					resource.TestCheckResourceAttr(poolResourceName, "ftp.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "ftp.0.support_rfc_2428", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "http.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "http.0.keepalive", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "http.0.keepalive_non_idempotent", "true"),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.#", "1"),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.0.principle", ""),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.0.target", ""),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.0.algorithm", "weighted_least_connections"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.0.priority_enabled", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.0.priority_nodes", "3"),
					resource.TestCheckResourceAttr(poolResourceName, "node.0.close_on_death", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "node.0.retry_fail_time", "30"),
					resource.TestCheckResourceAttr(poolResourceName, "smtp.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "smtp.0.send_starttls", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.client_auth", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.common_name_match.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sslCNPattern, "example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sslCNPattern, "another-example.com"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.elliptic_curves.#", "3"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sslEllipticCurvesPattern, "P384"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sslEllipticCurvesPattern, "P256"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sslEllipticCurvesPattern, "P521"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.enable", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.enhance", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.send_close_alerts", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.server_name", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.signature_algorithms", "ECDSA_SHA224 DSA_SHA256"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_ciphers", "SSL_ECDHE_RSA_WITH_AES_128_CBC_SHA SSL_ECDHE_RSA_WITH_AES_256_GCM_SHA384"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_ssl2", "enabled"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_ssl3", "enabled"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_tls1", "enabled"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_tls1_1", "enabled"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_tls1_2", "enabled"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.strict_verify", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp.0.nagle", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.0.accept_from", "all"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.0.accept_from_mask", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.0.response_timeout", "0"),
				),
			},
			{
				Config: testAccPoolUpdateTemplate(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "nodes_table.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, nodePattern, "192.168.10.10:80"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, priorityPattern, "5"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, statePattern, "draining"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, weightPattern, "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sourceIPPattern, "192.168.120.6"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, nodePattern, "192.168.10.12:80"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, priorityPattern, "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, statePattern, "active"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, weightPattern, "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sourceIPPattern, "192.168.120.6"),
					resource.TestCheckResourceAttr(poolResourceName, "bandwidth_class", "another-example"),
					resource.TestCheckResourceAttr(poolResourceName, "failure_pool", "test-pool2"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "55"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "4"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "5"),
					resource.TestCheckResourceAttr(poolResourceName, "monitors.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "monitors.0", "Full HTTPS"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "node_connection_attempts", "3"),
					resource.TestCheckResourceAttr(poolResourceName, "node_delete_behaviour", "drain"),
					resource.TestCheckResourceAttr(poolResourceName, "node_drain_to_delete_timeout", "4"),
					resource.TestCheckResourceAttr(poolResourceName, "note", "example test pool - updated"),
					resource.TestCheckResourceAttr(poolResourceName, "passive_monitoring", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "persistence_class", "another-example"),
					resource.TestCheckResourceAttr(poolResourceName, "transparent", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.addnode_delaytime", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.cloud_credentials", "another-example"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.cluster", "10.0.2.100"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.data_center", "another vCentre server"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.enabled", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.external", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.hysteresis", "200"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.imageid", "another image id"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.ips_to_use", "publicips"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.last_node_idle_time", "78"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.max_nodes", "200"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.min_nodes", "50"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.name", "anotherExample"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.port", "9980"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.refactory", "56"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.response_time", "89"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.scale_down_level", "75"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.scale_up_level", "15"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.securitygroupids.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, securityIDPattern, "sg-23456"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, securityIDPattern, "sg-34567"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.size_id", "sizeID"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.0.subnetids.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, subnetIDPattern, "subnet-aaaa"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, subnetIDPattern, "subnet-cccc"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.max_connect_time", "5"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.max_connections_per_node", "110"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.max_queue_size", "8"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.max_reply_time", "7"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.0.queue_timeout", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.0.enabled", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.0.hostnames.#", "3"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, dnsAutoScaleHostnamePattern, "example01.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, dnsAutoScaleHostnamePattern, "example02.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, dnsAutoScaleHostnamePattern, "example03.example.com"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.0.port", "8090"),
					resource.TestCheckResourceAttr(poolResourceName, "ftp.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "ftp.0.support_rfc_2428", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "http.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "http.0.keepalive", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "http.0.keepalive_non_idempotent", "false"),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.#", "1"),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.0.principle", ""),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.0.target", ""),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.0.algorithm", "weighted_round_robin"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.0.priority_enabled", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.0.priority_nodes", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "node.0.close_on_death", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "node.0.retry_fail_time", "45"),
					resource.TestCheckResourceAttr(poolResourceName, "smtp.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "smtp.0.send_starttls", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.client_auth", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.common_name_match.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sslCNPattern, "another-example.com"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.elliptic_curves.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sslEllipticCurvesPattern, "P256"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, sslEllipticCurvesPattern, "P521"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.enable", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.enhance", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.send_close_alerts", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.server_name", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.signature_algorithms", "RSA_SHA224 ECDSA_SHA224 DSA_SHA256"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_ciphers", "SSL_ECDHE_RSA_WITH_AES_128_CBC_SHA SSL_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384 SSL_ECDHE_RSA_WITH_AES_256_GCM_SHA384"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_ssl2", "use_default"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_ssl3", "use_default"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_tls1", "use_default"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_tls1_1", "use_default"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.ssl_support_tls1_2", "use_default"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.0.strict_verify", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp.0.nagle", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.#", "1"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.0.accept_from", "dest_ip_only"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.0.accept_from_mask", "192.168.0.1/24"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.0.response_timeout", "5"),
				),
			},
		},
	})
}

func testAccPoolCheckDestroy(s *terraform.State) error {

	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, r := range s.RootModule().Resources {
		if r.Type != "brocadevtm_pool" {
			continue
		}

		var name string
		var ok bool
		if name, ok = r.Primary.Attributes["name"]; !ok {
			return nil
		}

		var poolObj pool.Pool
		err := client.GetByName("pools", name, &poolObj)
		if client.StatusCode == http.StatusNotFound {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testCheckPoolExists(resName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)

		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pool name is set")
		}

		var name string
		if name, ok = rs.Primary.Attributes["name"]; ok && name == "" {
			return fmt.Errorf("No pool name is set")
		}

		var poolObj pool.Pool
		err := client.GetByName("pools", name, &poolObj)
		if err != nil {
			return fmt.Errorf("Received an error retrieving service with name: %s, %s", name, err)
		}
		return nil
	}
}

func testAccPoolNodeInvalidAlgo(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
}`, poolName)
}

func testAccPoolNodeUnsignedInt(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
  max_idle_connections_pernode = -1
}`, poolName)
}

func testAccPoolInvalidNodeNoPort(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
}`, poolName)
}

func testAccPoolInvalidNodeNoIP(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
}`, poolName)
}

func testAccPoolInvalidNode(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "5453535345353"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
}`, poolName)
}

func testAccPoolNoName() string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
}`)
}

func testAccPoolNoNodes(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
}`, poolName)
}

func testAccPoolInvalidNodeDeleteBehaviour(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
  node_delete_behaviour = "INVALID_BEHAVIOUR"
}`, poolName)
}

func testAccPoolOneItemList(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
  auto_scaling = [
    {
      addnode_delaytime = 10
    },
    {
      addnode_delaytime = 10
    },
  ]
}`, poolName)
}

func testAccPoolValidAWSSGPrefix(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest"{
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
  auto_scaling = [
    {
      securitygroupids = [ "INVALID_PREFIX-1234567" ]
    },
  ]
}`, poolName)
}

func testAccPoolValidAWSSubnetPrefix(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest"{
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
  auto_scaling = [
    {
      subnetids = [ "INVALID_SUBNET-12345" ]
    },
  ]
}`, poolName)
}

func testAccPoolCreateTemplate(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
  bandwidth_class = "example"
  failure_pool = "test-pool"
  max_connection_attempts = 100
  max_idle_connections_pernode = 10
  max_timed_out_connection_attempts = 8
  monitors = [ "Full HTTP" ]
  node_close_with_rst = true
  node_connection_attempts = 6
  node_delete_behaviour = "immediate"
  node_drain_to_delete_timeout = 10
  note = "example test pool"
  passive_monitoring = true
  persistence_class = "example"
  transparent = true

  auto_scaling = [
    {
      addnode_delaytime = 10
      cloud_credentials = "example"
      cluster = "10.0.0.1"
      data_center = "vCentre server"
      enabled = true
      external = true
      hysteresis = 100
      imageid = "image id"
      ips_to_use = "private_ips"
      last_node_idle_time = 10
      max_nodes = 100
      min_nodes = 20
      name = "example"
      port = 8980
      refactory = 10
      response_time = 100
      scale_down_level = 90
      scale_up_level = 20
      securitygroupids = [ "sg-12345", "sg-23456", "sg-34567" ]
      size_id = "sizeID"
      subnetids = [ "subnet-xxxx", "subnet-yyyyy" ]
    },
  ]
  pool_connection = [
    {
      max_connect_time = 4
      max_connections_per_node = 100
      max_queue_size = 10
      max_reply_time = 12
      queue_timeout = 14
    },
  ]
  dns_autoscale = [
    {
      enabled = true
      hostnames = [ "example01.example.com", "example02.example.com" ]
      port = 8080
    },
  ]
  ftp = [
    {
      support_rfc_2428 = true
    },
  ]
  http = [
    {
      keepalive = true
      keepalive_non_idempotent = true
    },
  ]
  /*
  kerberos_protocol_transition = [
    {
      principle = ""
      target = ""
    },
  ]
  */
  load_balancing = [
    {
      algorithm = "weighted_least_connections"
      priority_enabled = true
      priority_nodes = 3
    },
  ]
  node = [
    {
      close_on_death = true
      retry_fail_time = 30
    },
  ]
  smtp = [
    {
      send_starttls = true
    },
  ]
  ssl = [
    {
       client_auth = true
       common_name_match = [ "example.com", "another-example.com" ]
       elliptic_curves = [ "P384", "P256", "P521" ]
       enable = true
       enhance = true
       send_close_alerts = true
       server_name = true
       signature_algorithms = "ECDSA_SHA224 DSA_SHA256"
       ssl_ciphers = "SSL_ECDHE_RSA_WITH_AES_128_CBC_SHA SSL_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
       ssl_support_ssl2 = "enabled"
       ssl_support_ssl3 = "enabled"
       ssl_support_tls1 = "enabled"
       ssl_support_tls1_1 = "enabled"
       ssl_support_tls1_2 = "enabled"
       strict_verify = true
    },
  ]
  tcp = [
    {
      nagle = true
    },
  ]
  udp = [
    {
      accept_from = "all"
      accept_from_mask = "10.0.0.0/8"
      response_timeout = 0
    },
  ]
}`, poolName)
}

func testAccPoolUpdateTemplate(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      priority = 5
      state = "draining"
      weight = 2
      source_ip = "192.168.120.6"
    },
    {
      node = "192.168.10.12:80"
      priority = 1
      state = "active"
      weight = 2
      source_ip = "192.168.120.6"
    },
  ]
  bandwidth_class = "another-example"
  failure_pool = "test-pool2"
  max_connection_attempts = 55
  max_idle_connections_pernode = 4
  max_timed_out_connection_attempts = 5
  monitors = [ "Full HTTPS" ]
  node_close_with_rst = false
  node_connection_attempts = 3
  node_delete_behaviour = "drain"
  node_drain_to_delete_timeout = 4
  note = "example test pool - updated"
  passive_monitoring = false
  persistence_class = "another-example"
  transparent = false

  auto_scaling = [
    {
      addnode_delaytime = 20
      cloud_credentials = "another-example"
      cluster = "10.0.2.100"
      data_center = "another vCentre server"
      enabled = false
      external = false
      hysteresis = 200
      imageid = "another image id"
      ips_to_use = "publicips"
      last_node_idle_time = 78
      max_nodes = 200
      min_nodes = 50
      name = "anotherExample"
      port = 9980
      refactory = 56
      response_time = 89
      scale_down_level = 75
      scale_up_level = 15
      securitygroupids = [ "sg-23456", "sg-34567" ]
      size_id = "sizeID"
      subnetids = [ "subnet-aaaa", "subnet-cccc" ]
    },
  ]
  pool_connection = [
    {
      max_connect_time = 5
      max_connections_per_node = 110
      max_queue_size = 8
      max_reply_time = 7
      queue_timeout = 20
    },
  ]
  dns_autoscale = [
    {
      enabled = false
      hostnames = [ "example01.example.com", "example02.example.com", "example03.example.com" ]
      port = 8090
    },
  ]
  ftp = [
    {
      support_rfc_2428 = false
    },
  ]
  http = [
    {
      keepalive = false
      keepalive_non_idempotent = false
    },
  ]
  /*
  kerberos_protocol_transition = [
    {
      principle = ""
      target = ""
    },
  ]
  */
  load_balancing = [
    {
      algorithm = "weighted_round_robin"
      priority_enabled = false
      priority_nodes = 1
    },
  ]
  node = [
    {
      close_on_death = false
      retry_fail_time = 45
    },
  ]
  smtp = [
    {
      send_starttls = false
    },
  ]
  ssl = [
    {
       client_auth = false
       common_name_match = [ "another-example.com" ]
       elliptic_curves = [ "P256", "P521" ]
       enable = false
       enhance = false
       send_close_alerts = false
       server_name = false
       signature_algorithms = "RSA_SHA224 ECDSA_SHA224 DSA_SHA256"
       ssl_ciphers = "SSL_ECDHE_RSA_WITH_AES_128_CBC_SHA SSL_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384 SSL_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
       ssl_support_ssl2 = "use_default"
       ssl_support_ssl3 = "use_default"
       ssl_support_tls1 = "use_default"
       ssl_support_tls1_1 = "use_default"
       ssl_support_tls1_2 = "use_default"
       strict_verify = false
    },
  ]
  tcp = [
    {
      nagle = false
    },
  ]
  udp = [
    {
      accept_from = "dest_ip_only"
      accept_from_mask = "192.168.0.1/24"
      response_timeout = 5
    },
  ]
}`, poolName)
}
