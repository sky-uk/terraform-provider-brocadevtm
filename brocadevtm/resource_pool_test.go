package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api/pool"
	"github.com/sky-uk/go-rest-api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPoolCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolNodeInvalidAlgo(poolName),
				//ExpectError: regexp.MustCompile(`must be one of fastest_response_time, least_connections, perceptive, random, round_robin, weighted_least_connections, weighted_round_robin`),
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
				Config:      testAccPoolInvalidNodeDeleteBehaviour(poolName),
				ExpectError: regexp.MustCompile(`must be one of immediate or drain`),
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
				),
			},
		},
	})
}

func testAccPoolCheckDestroy(s *terraform.State) error {
	vtmClient := testAccProvider.Meta().(*rest.Client)
	var name string
	for _, r := range s.RootModule().Resources {
		if r.Type != "brocadevtm_pool" {
			continue
		}

		if name, ok := r.Primary.Attributes["name"]; ok && name == "" {
			return nil
		}

		api := pool.NewGet(name)
		err := vtmClient.Do(api)

		if err != nil {
			return err
		}
	}
	return nil
}

func testCheckPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pool name is set")
		}

		if name, ok := rs.Primary.Attributes["name"]; ok && name == "" {
			return fmt.Errorf("No pool name is set")
		}

		vtmClient := testAccProvider.Meta().(*rest.Client)

		api := pool.NewGet(rs.Primary.Attributes["name"])
		err := vtmClient.Do(api)

		if err != nil {
			return fmt.Errorf("Received an error retrieving service with name: %s, %s", rs.Primary.Attributes["name"], err)
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
}`, poolName)
}
