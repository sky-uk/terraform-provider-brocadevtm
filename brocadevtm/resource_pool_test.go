package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/pool"
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
				Config:      testAccPoolNodeListHasNoNodeAttribute(poolName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPoolNodeListHasNoPriorityAttribute(poolName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPoolNodeListHasNoWeightAttribute(poolName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPoolNodeListHasNoStateAttribute(poolName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPoolInvalidNode(poolName),
				ExpectError: regexp.MustCompile(`Must be a valid IP and port seperated by a colon. i.e 127.0.0.1:80`),
			},

			{
				Config:      testAccPoolInvalidNodeNoIP(poolName),
				ExpectError: regexp.MustCompile(`Must be a valid IP and port seperated by a colon. i.e 127.0.0.1:80`),
			},
			{
				Config:      testAccPoolInvalidNodeNoPort(poolName),
				ExpectError: regexp.MustCompile(`Must be a valid IP and port seperated by a colon. i.e 127.0.0.1:80`),
			},
			{
				Config: testAccPoolCheckVTMServiceConfig(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "monitorlist.0", "ping"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_timeout", "60"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connections_per_node", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "max_queue_size", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "max_reply_time", "60"),
					resource.TestCheckResourceAttr(poolResourceName, "queue_timeout", "60"),
					resource.TestCheckResourceAttr(poolResourceName, "http_keepalive", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "http_keepalive_non_idempotent", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing_priority_enabled", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing_priority_nodes", "8"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp_nagle", "true"),
				),
			},
			{
				Config: testAccPoolCheckVTMServiceConfigUpdated(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "monitorlist.0", "ping"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "40"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "40"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_timeout", "120"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connections_per_node", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "max_queue_size", "40"),
					resource.TestCheckResourceAttr(poolResourceName, "max_reply_time", "120"),
					resource.TestCheckResourceAttr(poolResourceName, "queue_timeout", "120"),
					resource.TestCheckResourceAttr(poolResourceName, "http_keepalive", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "http_keepalive_non_idempotent", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing_priority_enabled", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing_priority_nodes", "16"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp_nagle", "false"),
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
			return fmt.Errorf("Resource has no name")
		}

		var poolObj pool.Pool
		err := client.GetByName("pools", name, &poolObj)
		if client.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Brocade vTM Pools: error checking resource existance: %s", err)
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
		if name, ok := rs.Primary.Attributes["name"]; ok && name == "" {
			return fmt.Errorf("No pool name is set")
		}

		var poolObj pool.Pool
		err := client.GetByName("pools", name, &poolObj)
		if client.StatusCode != http.StatusOK {
			return fmt.Errorf("Resource %s not found, status code: %d", name, client.StatusCode)
		}
		if err != nil {
			return fmt.Errorf("Error retrieving resource %s: %s", name, err)
		}
		return nil
	}
}

func testAccPoolNodeInvalidAlgo(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    priority=1
    state="active"
    weight=1
  }
  max_connection_attempts = -10
  max_idle_connections_pernode = 20
  max_timed_out_connection_attempts = 20
  node_close_with_rst = false
  max_connection_timeout = 60
  max_connections_per_node = 10
  max_queue_size = 20
  max_reply_time = 60
  queue_timeout = 60
  http_keepalive = false
  http_keepalive_non_idempotent = false
  load_balancing_priority_enabled = false
  load_balancing_priority_nodes = 8
  load_balancing_algorithm = "INVALID_ALGO"
  tcp_nagle = false
}`, poolName)
}

func testAccPoolNodeUnsignedInt(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    priority=1
    state="active"
    weight=1
  }
  max_connection_attempts = -10
  max_idle_connections_pernode = 20
  max_timed_out_connection_attempts = 20
  node_close_with_rst = false
  max_connection_timeout = 60
  max_connections_per_node = 10
  max_queue_size = 20
  max_reply_time = 60
  queue_timeout = 60
  http_keepalive = false
  http_keepalive_non_idempotent = false
  load_balancing_priority_enabled = false
  load_balancing_priority_nodes = 8
  load_balancing_algorithm = "least_connections"
  tcp_nagle = false
}`, poolName)
}

func testAccPoolInvalidNodeNoPort(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1"
    priority=1
    state="active"
    weight=1
  }
  max_connection_attempts = 10
  max_idle_connections_pernode = 20
  max_timed_out_connection_attempts = 20
  node_close_with_rst = false
  max_connection_timeout = 60
  max_connections_per_node = 10
  max_queue_size = 20
  max_reply_time = 60
  queue_timeout = 60
  http_keepalive = false
  http_keepalive_non_idempotent = false
  load_balancing_priority_enabled = false
  load_balancing_priority_nodes = 8
  load_balancing_algorithm = "least_connections"
  tcp_nagle = false
}`, poolName)
}

func testAccPoolInvalidNodeNoIP(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="8080"
    priority=1
    state="active"
    weight=1
  }
  max_connection_attempts = 10
  max_idle_connections_pernode = 20
  max_timed_out_connection_attempts = 20
  node_close_with_rst = false
  max_connection_timeout = 60
  max_connections_per_node = 10
  max_queue_size = 20
  max_reply_time = 60
  queue_timeout = 60
  http_keepalive = false
  http_keepalive_non_idempotent = false
  load_balancing_priority_enabled = false
  load_balancing_priority_nodes = 8
  load_balancing_algorithm = "least_connections"
  tcp_nagle = false
}`, poolName)
}

func testAccPoolInvalidNode(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="325345234534"
    priority=1
    state="active"
    weight=1
  }
}`, poolName)
}

func testAccPoolNoName() string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    priority=1
    state="active"
    weight=1
  }
  max_connection_attempts = 10
  max_idle_connections_pernode = 20
  max_timed_out_connection_attempts = 20
  node_close_with_rst = false
  max_connection_timeout = 60
  max_connections_per_node = 10
  max_queue_size = 20
  max_reply_time = 60
  queue_timeout = 60
  http_keepalive = false
  http_keepalive_non_idempotent = false
  load_balancing_priority_enabled = false
  load_balancing_priority_nodes = 8
  load_balancing_algorithm = "least_connections"
  tcp_nagle = false
}`)
}

func testAccPoolNoNodes(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  max_connection_attempts = 10
  max_idle_connections_pernode = 20
  max_timed_out_connection_attempts = 20
  node_close_with_rst = false
  max_connection_timeout = 60
  max_connections_per_node = 10
  max_queue_size = 20
  max_reply_time = 60
  queue_timeout = 60
  http_keepalive = false
  http_keepalive_non_idempotent = false
  load_balancing_priority_enabled = false
  load_balancing_priority_nodes = 8
  load_balancing_algorithm = "least_connections"
  tcp_nagle = false
}`, poolName)
}

func testAccPoolNodeListHasNoNodeAttribute(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    priority=1
    state="active"
    weight=1
  }
}`, poolName)
}

func testAccPoolNodeListHasNoPriorityAttribute(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    state="active"
    weight=1
  }
}`, poolName)
}

func testAccPoolNodeListHasNoStateAttribute(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    priority=1
    weight=1
  }
}`, poolName)
}

func testAccPoolNodeListHasNoWeightAttribute(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    priority=1
    state="active"
  }
}`, poolName)
}

func testAccPoolCheckVTMServiceConfig(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    priority=1
    state="active"
    weight=1
  }
  max_connection_attempts = 10
  max_idle_connections_pernode = 20
  max_timed_out_connection_attempts = 20
  node_close_with_rst = true
  max_connection_timeout = 60
  max_connections_per_node = 10
  max_queue_size = 20
  max_reply_time = 60
  queue_timeout = 60
  http_keepalive = true
  http_keepalive_non_idempotent = true
  load_balancing_priority_enabled = true
  load_balancing_priority_nodes = 8
  load_balancing_algorithm = "least_connections"
  tcp_nagle = true
}`, poolName)
}

func testAccPoolCheckVTMServiceConfigUpdated(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "%s"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    priority=1
    state="active"
    weight=1
  }
  node {
    node="127.0.0.2:80"
    priority=1
    state="active"
    weight=1
  }
  max_connection_attempts = 20
  max_idle_connections_pernode = 40
  max_timed_out_connection_attempts = 40
  node_close_with_rst = false
  max_connection_timeout = 120
  max_connections_per_node = 20
  max_queue_size = 40
  max_reply_time = 120
  queue_timeout = 120
  http_keepalive = false
  http_keepalive_non_idempotent = false
  load_balancing_priority_enabled = false
  load_balancing_priority_nodes = 16
  load_balancing_algorithm = "least_connections"
  tcp_nagle = false
}`, poolName)
}
