package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/virtualserver"
	"regexp"
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
				Config: testAccBrocadeVTMVirtualServerInvalidInt(virtualServerName),
				ExpectError: regexp.MustCompile(`can't be negative`),
			},
			{
				Config: testAccBrocadeVTMVirtualServerInvalidProtocol(virtualServerName),
				ExpectError: regexp.MustCompile(`fsdfdsfsdfsafasdf`),
			},
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

	vtmClient := testAccProvider.Meta().(*brocadevtm.VTMClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_virtual_server" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		api := virtualserver.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return nil
		}
		for _, virtual_server := range api.GetResponse().Children {
			if virtual_server.Name == name {
				return fmt.Errorf("Brocade vTM Virtual Server %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMVirtualServerExists(virtualServerName, virtualServerResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[virtualServerResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Virtual Server %s wasn't found in resources", virtualServerName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Virtual Server ID not set for %s in resources", virtualServerName)
		}

		vtmClient := testAccProvider.Meta().(*brocadevtm.VTMClient)
		api := virtualserver.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		for _, virtualServer := range api.GetResponse().Children {
			if virtualServer.Name == virtualServerName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM Virtual Server %s not found on remote vTM", virtualServerName)
	}
}

func testAccBrocadeVTMVirtualServerInvalidInt(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
enabled = true
listen_on_any = false
listen_traffic_ips = ["test-traffic-ip-group"]
pool = "test-pool"
port = -80
protocol = "http"
request_rules = []
ssl_decrypt = false
connection_keepalive = false
connection_keepalive_timeout = 60
connection_max_client_buffer = 256
connection_max_server_buffer = 256
connection_max_transaction_duration = 100
connection_server_first_banner = "Banner"
connection_timeout = 60
ssl_server_cert_default = "abcdefg"
ssl_server_cert_host_mapping_host = "test-host.example.com"
ssl_server_cert_host_mapping_alt_certificates = ["cert1","cert2"]
ssl_server_cert_host_mapping_certificate = "test-host.example.com"
ssl_support_ssl2 = "use_default"
ssl_support_ssl3 = "use_default"
ssl_support_tls1 = "use_default"
ssl_support_tls1_1 = "use_default"
ssl_support_tls1_2 = "use_default"
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerInvalidProtocol(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
enabled = true
listen_on_any = false
listen_traffic_ips = ["test-traffic-ip-group"]
pool = "test-pool"
port = -80
protocol = "SOME_INVALID_PROTOCOL"
request_rules = []
ssl_decrypt = false
connection_keepalive = false
connection_keepalive_timeout = 60
connection_max_client_buffer = 256
connection_max_server_buffer = 256
connection_max_transaction_duration = 100
connection_server_first_banner = "Banner"
connection_timeout = 60
ssl_server_cert_default = "abcdefg"
ssl_server_cert_host_mapping_host = "test-host.example.com"
ssl_server_cert_host_mapping_alt_certificates = ["cert1","cert2"]
ssl_server_cert_host_mapping_certificate = "test-host.example.com"
ssl_support_ssl2 = "use_default"
ssl_support_ssl3 = "use_default"
ssl_support_tls1 = "use_default"
ssl_support_tls1_1 = "use_default"
ssl_support_tls1_2 = "use_default"
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerCreate(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
enabled = true
listen_on_any = false
listen_traffic_ips = ["test-traffic-ip-group"]
pool = "test-pool"
port = 80
protocol = "http"
request_rules = []
ssl_decrypt = false
connection_keepalive = false
connection_keepalive_timeout = 60
connection_max_client_buffer = 256
connection_max_server_buffer = 256
connection_max_transaction_duration = 100
connection_server_first_banner = "Banner"
connection_timeout = 60
ssl_server_cert_default = "abcdefg"
ssl_server_cert_host_mapping_host = "test-host.example.com"
ssl_server_cert_host_mapping_alt_certificates = ["cert1","cert2"]
ssl_server_cert_host_mapping_certificate = "test-host.example.com"
ssl_support_ssl2 = "use_default"
ssl_support_ssl3 = "use_default"
ssl_support_tls1 = "use_default"
ssl_support_tls1_1 = "use_default"
ssl_support_tls1_2 = "use_default"
}
`, virtualServerName)
}
