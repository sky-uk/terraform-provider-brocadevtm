package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/pool"
	"testing"
)

func TestAccPool_Basic(t *testing.T) {

	randomInt := acctest.RandInt()

	poolName := fmt.Sprintf("acctest_brocadevtm_pool-%d", randomInt)
	poolResourceName := "brocadevtm_pool.acctest"

	fmt.Printf("\n\nPool is %s.\n\n", poolName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBrocadeVTMPoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckVTMServiceConfig(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckBrocadeVTMPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "monitorlist", "[\"ping\"]"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_timeout", "60"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connections_per_node", "10"),
					resource.TestCheckResourceAttr(poolResourceName, "max_queue_size", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "max_reply_time", "60"),
					resource.TestCheckResourceAttr(poolResourceName, "queue_timeout", "60"),
					resource.TestCheckResourceAttr(poolResourceName, "http_keepalive", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "http_keepalive_non_idempotent", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing_priority_enabled", "false"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing_priority_nodes", "8"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp_nagle", "false"),
				),
			},
			resource.TestStep{
				Config: testAccCheckVTMServiceConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					testCheckBrocadeVTMPoolExists(poolResourceName),
					resource.TestCheckResourceAttr(poolResourceName, "name", poolName),
					resource.TestCheckResourceAttr(poolResourceName, "monitorlist", "[\"ping\"]"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_attempts", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "max_idle_connections_pernode", "40"),
					resource.TestCheckResourceAttr(poolResourceName, "max_timed_out_connection_attempts", "40"),
					resource.TestCheckResourceAttr(poolResourceName, "node_close_with_rst", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connection_timeout", "120"),
					resource.TestCheckResourceAttr(poolResourceName, "max_connections_per_node", "20"),
					resource.TestCheckResourceAttr(poolResourceName, "max_queue_size", "40"),
					resource.TestCheckResourceAttr(poolResourceName, "max_reply_time", "120"),
					resource.TestCheckResourceAttr(poolResourceName, "queue_timeout", "120"),
					resource.TestCheckResourceAttr(poolResourceName, "http_keepalive", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "http_keepalive_non_idempotent", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing_priority_enabled", "true"),
					resource.TestCheckResourceAttr(poolResourceName, "load_balancing_priority_nodes", "16"),
					resource.TestCheckResourceAttr(poolResourceName, "tcp_nagle", "true"),
				),
			},
		},
	})
}

func testAccCheckBrocadeVTMPoolDestroy(s *terraform.State) error {
	vtmClient := testAccProvider.Meta().(*brocadevtm.VTMClient)
	var name string
	for _, r := range s.RootModule().Resources {
		if r.Type != "brocadevtm_pool" {
			continue
		}

		if name, ok := r.Primary.Attributes["name"]; ok && name == "" {
			return nil
		}

		api := pool.NewGetSingle(name)
		err := vtmClient.Do(api)

		if err != nil {
			return err
		}
	}
	return nil
}

func testCheckBrocadeVTMPoolExists(name string) resource.TestCheckFunc {
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

		vtmClient := testAccProvider.Meta().(*brocadevtm.VTMClient)

		api := pool.NewGetSingle(rs.Primary.Attributes["name"])
		err := vtmClient.Do(api)

		if err != nil {
			return fmt.Errorf("Received an error retrieving service with name: %s, %s", rs.Primary.Attributes["name"], err)
		}

		return nil
	}
}

func testAccCheckVTMServiceConfig(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "foo" {
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
  tcp_nagle = false
}`, poolName)
}

func testAccCheckVTMServiceConfigUpdated(poolName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "foo" {
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
  node_close_with_rst = true
  max_connection_timeout = 120
  max_connections_per_node = 20
  max_queue_size = 40
  max_reply_time = 120
  queue_timeout = 120
  http_keepalive = true
  http_keepalive_non_idempotent = true
  load_balancing_priority_enabled = true
  load_balancing_priority_nodes = 16
  tcp_nagle = true
}`, poolName)
}
