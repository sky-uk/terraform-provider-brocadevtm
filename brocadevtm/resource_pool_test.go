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
				ExpectError: regexp.MustCompile(`expected node_delete_behaviour to be one of \[drain immediate\]`),
			},
			{
				Config:      testAccPoolOneItemList(poolName),
				ExpectError: regexp.MustCompile(`attribute supports 1 item maximum, config has 2 declared`),
			},
			{
				Config: testAccPoolCreateTemplate(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "nodes_table.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "node"), "192.168.10.10:80"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "priority"), "5"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "state"), "draining"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "weight"), "2"),
					resource.TestCheckResourceAttr(poolResourceName, "bandwidth_class", "example"),
					resource.TestCheckResourceAttr(poolResourceName, "failure_pool", "test-pool"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "100"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "8"),
					resource.TestCheckResourceAttr(poolResourceName, "monitors.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("monitors"), "Full HTTP"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "node_connection_attempts", "6"),
					resource.TestCheckResourceAttr(poolResourceName, "node_delete_behaviour", "immediate"),
					resource.TestCheckResourceAttr(poolResourceName, "node_drain_to_delete_timeout", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "note", "example test pool"),
					resource.TestCheckResourceAttr(poolResourceName, "passive_monitoring", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "persistence_class", "example"),
					resource.TestCheckResourceAttr(poolResourceName, "transparent", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "auto_scaling.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "addnode_delaytime"), "10"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "cloud_credentials"), "example"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "cluster"), "10.0.0.1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "data_center"), "vCentre server"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "enabled"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "external"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "hysteresis"), "100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "imageid"), "image id"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "ips_to_use"), "private_ips"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "last_node_idle_time"), "10"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "max_nodes"), "100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "min_nodes"), "20"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "name"), "example"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "port"), "8980"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "refractory"), "10"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "response_time"), "100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "scale_down_level"), "90"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "scale_up_level"), "20"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "securitygroupids.#"), "3"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "securitygroupids"), "sg-12345"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "securitygroupids"), "sg-23456"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "securitygroupids"), "sg-34567"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "size_id"), "sizeID"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "subnetids.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "subnetids"), "subnet-xxxx"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "subnetids"), "subnet-xxxx"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_connect_time"), "4"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_connections_per_node"), "100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_queue_size"), "10"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_reply_time"), "12"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "queue_timeout"), "14"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "enabled"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "hostnames.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("dns_autoscale", "hostnames"), "example01.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("dns_autoscale", "hostnames"), "example02.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "port"), "8080"),
					resource.TestCheckResourceAttr(poolResourceName, "ftp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ftp", "support_rfc_2428"), "true"),
					resource.TestCheckResourceAttr(poolResourceName, "http.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("http", "keepalive"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("http", "keepalive_non_idempotent"), "true"),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.#", "1"),
					//util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("kerberos_protocol_transition", "principle"), ""),
					//util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("kerberos_protocol_transition", "target"), ""),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "algorithm"), "weighted_least_connections"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "priority_enabled"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "priority_nodes"), "3"),
					resource.TestCheckResourceAttr(poolResourceName, "node.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("node", "close_on_death"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("node", "retry_fail_time"), "30"),
					resource.TestCheckResourceAttr(poolResourceName, "smtp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("smtp", "send_starttls"), "true"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "client_auth"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "common_name_match.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "common_name_match"), "example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "common_name_match"), "another-example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "elliptic_curves.#"), "3"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "elliptic_curves"), "P384"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "elliptic_curves"), "P256"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "elliptic_curves"), "P521"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "enable"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "enhance"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "send_close_alerts"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_name"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "signature_algorithms"), "ECDSA_SHA224 DSA_SHA256"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_ciphers"), "SSL_ECDHE_RSA_WITH_AES_128_CBC_SHA SSL_ECDHE_RSA_WITH_AES_256_GCM_SHA384"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_ssl2"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_ssl3"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1_1"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1_2"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "strict_verify"), "true"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "nagle"), "true"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "accept_from"), "all"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "accept_from_mask"), "10.0.0.0/8"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "response_timeout"), "0"),
				),
			},
			{
				Config: testAccPoolUpdateTemplate(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "nodes_table.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "node"), "192.168.10.10:80"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "priority"), "5"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "state"), "draining"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "weight"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "node"), "192.168.10.12:80"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "priority"), "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "state"), "active"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "weight"), "2"),
					resource.TestCheckResourceAttr(poolResourceName, "bandwidth_class", "another-example"),
					resource.TestCheckResourceAttr(poolResourceName, "failure_pool", "test-pool2"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "55"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "4"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "5"),
					resource.TestCheckResourceAttr(poolResourceName, "monitors.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("monitors"), "Full HTTPS"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "node_connection_attempts", "3"),
					resource.TestCheckResourceAttr(poolResourceName, "node_delete_behaviour", "drain"),
					resource.TestCheckResourceAttr(poolResourceName, "node_drain_to_delete_timeout", "4"),
					resource.TestCheckResourceAttr(poolResourceName, "note", "example test pool - updated"),
					resource.TestCheckResourceAttr(poolResourceName, "passive_monitoring", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "persistence_class", "another-example"),
					resource.TestCheckResourceAttr(poolResourceName, "transparent", "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "addnode_delaytime"), "20"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "cloud_credentials"), "another-example"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "cluster"), "10.0.2.100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "data_center"), "another vCentre server"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "enabled"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "external"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "hysteresis"), "200"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "imageid"), "another image id"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "ips_to_use"), "publicips"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "last_node_idle_time"), "78"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "max_nodes"), "200"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "min_nodes"), "50"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "name"), "anotherExample"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "port"), "9980"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "refractory"), "56"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "response_time"), "89"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "scale_down_level"), "75"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "scale_up_level"), "15"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "securitygroupids.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "securitygroupids"), "sg-23456"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "securitygroupids"), "sg-34567"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "size_id"), "sizeID"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "subnetids.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "subnetids"), "subnet-aaaa"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "subnetids"), "subnet-cccc"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_connect_time"), "5"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_connections_per_node"), "110"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_queue_size"), "8"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_reply_time"), "7"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "queue_timeout"), "20"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "enabled"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "hostnames.#"), "3"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("dns_autoscale", "hostnames"), "example01.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("dns_autoscale", "hostnames"), "example02.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("dns_autoscale", "hostnames"), "example03.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "port"), "8090"),
					resource.TestCheckResourceAttr(poolResourceName, "ftp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ftp", "support_rfc_2428"), "false"),
					resource.TestCheckResourceAttr(poolResourceName, "http.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("http", "keepalive"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("http", "keepalive_non_idempotent"), "false"),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.#", "1"),
					//util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("kerberos_protocol_transition", "principle"), ""),
					//util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("kerberos_protocol_transition", "target"), ""),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "algorithm"), "weighted_round_robin"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "priority_enabled"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "priority_nodes"), "1"),
					resource.TestCheckResourceAttr(poolResourceName, "node.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("node", "close_on_death"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("node", "retry_fail_time"), "45"),
					resource.TestCheckResourceAttr(poolResourceName, "smtp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("smtp", "send_starttls"), "false"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "client_auth"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "common_name_match.#"), "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "common_name_match"), "another-example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "elliptic_curves.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "elliptic_curves"), "P256"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "elliptic_curves"), "P521"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "enable"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "enhance"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "send_close_alerts"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_name"), "false"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "signature_algorithms"), "RSA_SHA224 ECDSA_SHA224 DSA_SHA256"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_ciphers"), "SSL_ECDHE_RSA_WITH_AES_128_CBC_SHA SSL_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384 SSL_ECDHE_RSA_WITH_AES_256_GCM_SHA384"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_ssl2"), "use_default"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_ssl3"), "use_default"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1"), "use_default"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1_1"), "use_default"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1_2"), "use_default"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "strict_verify"), "false"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "nagle"), "false"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "accept_from"), "dest_ip_only"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "accept_from_mask"), "192.168.0.1/24"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "response_timeout"), "5"),
				),
			},
			{
				Config: testAccPoolCreateTemplateNodesList(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "nodes_table.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "node"), "192.168.10.10:80"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "priority"), "5"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "state"), "draining"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("nodes_table", "weight"), "2"),
					resource.TestCheckResourceAttr(poolResourceName, "bandwidth_class", "example"),
					resource.TestCheckResourceAttr(poolResourceName, "failure_pool", "test-pool"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "100"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "8"),
					resource.TestCheckResourceAttr(poolResourceName, "monitors.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("monitors"), "Full HTTP"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "node_connection_attempts", "6"),
					resource.TestCheckResourceAttr(poolResourceName, "node_delete_behaviour", "immediate"),
					resource.TestCheckResourceAttr(poolResourceName, "node_drain_to_delete_timeout", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "note", "example test pool"),
					resource.TestCheckResourceAttr(poolResourceName, "passive_monitoring", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "persistence_class", "example"),
					resource.TestCheckResourceAttr(poolResourceName, "transparent", "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "addnode_delaytime"), "10"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "cloud_credentials"), "example"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "cluster"), "10.0.0.1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "data_center"), "vCentre server"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "enabled"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "external"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "hysteresis"), "100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "imageid"), "image id"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "ips_to_use"), "private_ips"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "last_node_idle_time"), "10"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "max_nodes"), "100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "min_nodes"), "20"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "name"), "example"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "port"), "8980"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "refractory"), "10"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "response_time"), "100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "scale_down_level"), "90"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "scale_up_level"), "20"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "securitygroupids.#"), "3"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "securitygroupids"), "sg-12345"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "securitygroupids"), "sg-23456"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "securitygroupids"), "sg-34567"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "size_id"), "sizeID"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("auto_scaling", "subnetids.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "subnetids"), "subnet-xxxx"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("auto_scaling", "subnetids"), "subnet-xxxx"),
					resource.TestCheckResourceAttr(poolResourceName, "pool_connection.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_connect_time"), "4"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_connections_per_node"), "100"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_queue_size"), "10"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "max_reply_time"), "12"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("pool_connection", "queue_timeout"), "14"),
					resource.TestCheckResourceAttr(poolResourceName, "dns_autoscale.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "enabled"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "hostnames.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("dns_autoscale", "hostnames"), "example01.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("dns_autoscale", "hostnames"), "example02.example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("dns_autoscale", "port"), "8080"),
					resource.TestCheckResourceAttr(poolResourceName, "ftp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ftp", "support_rfc_2428"), "true"),
					resource.TestCheckResourceAttr(poolResourceName, "http.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("http", "keepalive"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("http", "keepalive_non_idempotent"), "true"),
					//resource.TestCheckResourceAttr(poolResourceName, "kerberos_protocol_transition.#", "1"),
					//util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("kerberos_protocol_transition", "principle"), ""),
					//util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSet("kerberos_protocol_transition", "target"), ""),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "algorithm"), "weighted_least_connections"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "priority_enabled"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("load_balancing", "priority_nodes"), "3"),
					resource.TestCheckResourceAttr(poolResourceName, "node.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("node", "close_on_death"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("node", "retry_fail_time"), "30"),
					resource.TestCheckResourceAttr(poolResourceName, "smtp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("smtp", "send_starttls"), "true"),
					resource.TestCheckResourceAttr(poolResourceName, "ssl.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "client_auth"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "common_name_match.#"), "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "common_name_match"), "example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "common_name_match"), "another-example.com"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "elliptic_curves.#"), "3"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "elliptic_curves"), "P384"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "elliptic_curves"), "P256"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForNestedSets("ssl", "elliptic_curves"), "P521"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "enable"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "enhance"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "send_close_alerts"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_name"), "true"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "signature_algorithms"), "ECDSA_SHA224 DSA_SHA256"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_ciphers"), "SSL_ECDHE_RSA_WITH_AES_128_CBC_SHA SSL_ECDHE_RSA_WITH_AES_256_GCM_SHA384"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_ssl2"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_ssl3"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1_1"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ssl_support_tls1_2"), "enabled"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "strict_verify"), "true"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "nagle"), "true"),
					resource.TestCheckResourceAttr(poolResourceName, "udp.#", "1"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "accept_from"), "all"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "accept_from_mask"), "10.0.0.0/8"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "response_timeout"), "0"),
				),
			},
			{
				Config: testAccPoolUpdateTemplateNodesList(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "nodes_list.#", "2"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, regexp.MustCompile("nodes_list."), "192.168.10.11:80"),
					util.AccTestCheckValueInKeyPattern(poolResourceName, regexp.MustCompile("nodes_list."), "192.168.10.12:80"),
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
      refractory = 10
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
    },
    {
      node = "192.168.10.12:80"
      priority = 1
      state = "active"
      weight = 2
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
      refractory = 56
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

func testAccPoolCreateTemplateNodesList(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_list = [ "192.168.10.11:80" ]
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
      refractory = 10
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

func testAccPoolUpdateTemplateNodesList(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_list = [ "192.168.10.11:80","192.168.10.12:80" ]
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
      refractory = 10
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
