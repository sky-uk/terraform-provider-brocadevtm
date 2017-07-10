package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccBrocadeVTMVirtualServerBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	virtualServerName := fmt.Sprintf("acctest_brocadevtm_virtual_server-%d", randomInt)
	virtualServerResourceName := "brocadevtm_virtual_server.acctest"

	fmt.Printf("\n\nVirtual Server is %s.\n\n", virtualServerName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMVirtualServerCheckDestroy(state, virtualServerName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBrocadeVTMVirtualServerCreate(virtualServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMVirtualServerExists(virtualServerName, virtualServerResourceName),
					resource.TestCheckResourceAttr(virtualServerResourceName, "name", virtualServerName),
				),
			},
		},
	})
}

func testAccBrocadeVTMVirtualServerCheckDestroy(state *terraform.State, name string) error {
	return nil
}

func testAccBrocadeVTMVirtualServerExists(virtualServerName, virtualServerResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		return nil
	}
}

func testAccBrocadeVTMVirtualServerCreate(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
}
`, virtualServerName)
}
