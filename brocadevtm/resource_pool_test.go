package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/pool"
	"testing"
)

func TestAccPool_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBrocadeVTMPoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckVTMServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckBrocadeVTMPoolExists("brocadevtm_pool.foo"),
					resource.TestCheckResourceAttr(
						"brocadevtm_pool.foo", "name", "pool_foo"),
				),
			},
			resource.TestStep{
				Config: testAccCheckVTMServiceConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					testCheckBrocadeVTMPoolExists("brocadevtm_pool.foo"),
					resource.TestCheckResourceAttr(
						"brocadevtm_pool.foo", "name", "pool_bar"),
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

const testAccCheckVTMServiceConfig = `
resource "brocadevtm_pool" "foo" {
  name = "pool_foo"
  monitorlist = ["ping"]
  node {
    node="127.0.0.1:80"
    priority=1
    state="active"
    weight=1
  }
  max_connection_attempts = 5
}`

const testAccCheckVTMServiceConfigUpdated = `
resource "brocadevtm_pool" "foo" {
  name = "pool_bar"
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
  max_connection_attempts = 5
}`
