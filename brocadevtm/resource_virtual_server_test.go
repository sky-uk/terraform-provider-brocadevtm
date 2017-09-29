package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api/virtualserver"
	"github.com/sky-uk/go-rest-api"
	"regexp"
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
				Config:      testAccBrocadeVTMVirtualServerNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerNegativeInt(virtualServerName),
				ExpectError: regexp.MustCompile(`can't be negative`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerInvalidProtocol(virtualServerName),
				ExpectError: regexp.MustCompile(`must be one of client_first, dns, dns_tcp, ftp, http, https, imaps, imapv2, imapv3, imapv4, ldap, ldaps, pop3, pop3s, rtsp, server_first, siptcp, sipudp, smtp, ssl, stream, telnet, udp or udpstreaming`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerInvalidSSLSupportOption(virtualServerName),
				ExpectError: regexp.MustCompile(`must be one of use_default, disabled or enabled`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerInvalidNonce(virtualServerName),
				ExpectError: regexp.MustCompile(`must be one of off, on or strict`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerInvalidOCSPRequired(virtualServerName),
				ExpectError: regexp.MustCompile(`must be one of none, optional, strict`),
			},
			{
				Config: testAccBrocadeVTMVirtualServerCreate(virtualServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMVirtualServerExists(virtualServerName, virtualServerResourceName),
					resource.TestCheckResourceAttr(virtualServerResourceName, "name", virtualServerName),
					resource.TestCheckResourceAttr(virtualServerResourceName, "pool", "test-pool"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "port", "80"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "listen_on_any", "false"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "listen_traffic_ips.0", "test-traffic-ip-group"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "request_rules.0", "ruleOne"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_decrypt", "false"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_keepalive", "false"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_keepalive_timeout", "60"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_max_client_buffer", "1024"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_max_server_buffer", "8192"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_max_transaction_duration", "300"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_server_first_banner", "ACCESS DENIED"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_timeout", "30"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_default", "testsslkey2"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_ssl2", "disabled"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_ssl3", "disabled"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_tls1", "use_default"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_tls1_1", "enabled"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_tls1_2", "disabled"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.#", "2"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.0.ssl_server_alt_certs.#", "0"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.0.ssl_server_cert_host", "s1sta07-v00.devops.prd.ovp.bskyb.com"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.0.ssl_server_cert", "testsslkey1"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.1.ssl_server_alt_certs.#", "2"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.1.ssl_server_cert_host", "h1pipeline.devops.int.ovp.bskyb.com"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.1.ssl_server_cert", "testssl001"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.1.ssl_server_alt_certs.0", "testssl002"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.1.ssl_server_alt_certs.1", "testssl003"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_enable", "true"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.#", "1"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.issuer", "_DEFAULT_"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.aia", "true"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.nonce", "off"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.required", "optional"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.responder_cert", "SOME_RESPONDER_CERT"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.signer", "ANOTHER_SIGNER"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.url", "http://www.sky.uk/"),
				),
			},
			{
				Config: testAccBrocadeVTMVirtualServerUpdate(virtualServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMVirtualServerExists(virtualServerName, virtualServerResourceName),
					resource.TestCheckResourceAttr(virtualServerResourceName, "name", virtualServerName),
					resource.TestCheckResourceAttr(virtualServerResourceName, "pool", "test-pool2"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "port", "443"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "listen_on_any", "true"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "listen_traffic_ips.0", "another-test-traffic-ip-group"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "protocol", "https"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "request_rules.0", "ruleTwo"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_decrypt", "true"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_keepalive", "true"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_keepalive_timeout", "120"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_max_client_buffer", "2048"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_max_server_buffer", "4096"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_max_transaction_duration", "600"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_server_first_banner", "ACCESS ALLOWED"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "connection_timeout", "45"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_default", "testsslkey1"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_ssl2", "enabled"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_ssl3", "enabled"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_tls1", "enabled"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_tls1_1", "use_default"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_support_tls1_2", "enabled"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.#", "2"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.0.ssl_server_alt_certs.#", "1"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.0.ssl_server_cert_host", "h1pipeline.devops.int.ovp.bskyb.com"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.0.ssl_server_cert", "testssl002"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.0.ssl_server_alt_certs.0", "testssl001"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.1.ssl_server_alt_certs.#", "0"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.1.ssl_server_cert_host", "s1sta07-v00.devops.prd.ovp.bskyb.com"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ssl_server_cert_host_mapping.1.ssl_server_cert", "testsslkey1"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_enable", "false"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.#", "1"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.issuer", "_DEFAULT_"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.aia", "false"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.nonce", "on"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.required", "none"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.responder_cert", "ANOTHER_RESPONDER_CERT"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.signer", "SOME_SIGNER"),
					resource.TestCheckResourceAttr(virtualServerResourceName, "ocsp_issuers.0.url", "http://www.sky.com/"),
				),
			},
		},
	})
}

func testAccBrocadeVTMVirtualServerCheckDestroy(state *terraform.State, name string) error {

	vtmClient := testAccProvider.Meta().(*rest.Client)

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
			return fmt.Errorf("Brocade vTM Virtual Server - error occurred while retrieving a list of all virtual servers")
		}
		for _, virtualServer := range api.ResponseObject().(*virtualserver.VirtualServersList).Children {
			if virtualServer.Name == name {
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

		vtmClient := testAccProvider.Meta().(*rest.Client)
		api := virtualserver.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		for _, virtualServer := range api.ResponseObject().(*virtualserver.VirtualServersList).Children {
			if virtualServer.Name == virtualServerName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM Virtual Server %s not found on remote vTM", virtualServerName)
	}
}

func testAccBrocadeVTMVirtualServerNoName() string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
pool = "test-pool"
port = 80
}
`)
}

func testAccBrocadeVTMVirtualServerNegativeInt(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = -80
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerInvalidProtocol(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
protocol = "SOME_INVALID_PROTOCOL"
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerInvalidSSLSupportOption(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
ssl_support_ssl2 = "SOME_INVALID_SSL_OPTION"
ssl_support_ssl3 = "SOME_INVALID_SSL_OPTION"
ssl_support_tls1 = "SOME_INVALID_SSL_OPTION"
ssl_support_tls1_1 = "SOME_INVALID_SSL_OPTION"
ssl_support_tls1_2 = "SOME_INVALID_SSL_OPTION"
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerInvalidNonce(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
ocsp_enable = true
ocsp_enable = true
ocsp_issuers = [
  {
     issuer = "_DEFAULT_"
     aia = true
     nonce = "SOME_INVALID_OPTION"
     required = "optional"
     responder_cert = ""
     signer = ""
     url = ""
  },
]
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerInvalidOCSPRequired(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
ocsp_enable = true
ocsp_issuers = [
  {
     issuer = "_DEFAULT_"
     aia = true
     nonce = "off"
     required = "SOME_INVALID_OPTION"
     responder_cert = ""
     signer = ""
     url = ""
  },
]
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerCreate(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_pool" "acctest" {
  name = "test-pool"
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
  load_balancing_algorithm = "round_robin"
  tcp_nagle = false
}
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool"
port = 80
enabled = false
listen_on_any = false
listen_traffic_ips = ["test-traffic-ip-group"]
protocol = "http"
request_rules = ["ruleOne"]
ssl_decrypt = false
connection_keepalive = false
connection_keepalive_timeout = 60
connection_max_client_buffer = 1024
connection_max_server_buffer = 8192
connection_max_transaction_duration = 300
connection_server_first_banner = "ACCESS DENIED"
connection_timeout = 30
ssl_server_cert_default = "testsslkey2"
ssl_support_ssl2 = "disabled"
ssl_support_ssl3 = "disabled"
ssl_support_tls1 = "use_default"
ssl_support_tls1_1 = "enabled"
ssl_support_tls1_2 = "disabled"
ssl_server_cert_host_mapping = [
 {
  ssl_server_cert_host = "s1sta07-v00.devops.prd.ovp.bskyb.com"
  ssl_server_cert = "testsslkey1"
},
{
  ssl_server_cert_host = "h1pipeline.devops.int.ovp.bskyb.com"
  ssl_server_alt_certs = ["testssl002","testssl003"]
  ssl_server_cert = "testssl001"
}
]
ocsp_enable = true
ocsp_issuers = [
  {
     issuer = "_DEFAULT_"
     aia = true
     nonce = "off"
     required = "optional"
     responder_cert = "SOME_RESPONDER_CERT"
     signer = "ANOTHER_SIGNER"
     url = "http://www.sky.uk/"
  },
]
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerUpdate(virtualServerName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "test-pool2"
port = 443
enabled = true
listen_on_any = true
listen_traffic_ips = ["another-test-traffic-ip-group"]
protocol = "https"
request_rules = ["ruleTwo"]
ssl_decrypt = true
connection_keepalive = true
connection_keepalive_timeout = 120
connection_max_client_buffer = 2048
connection_max_server_buffer = 4096
connection_max_transaction_duration = 600
connection_server_first_banner = "ACCESS ALLOWED"
connection_timeout = 45
ssl_server_cert_default = "testsslkey1"
ssl_support_ssl2 = "enabled"
ssl_support_ssl3 = "enabled"
ssl_support_tls1 = "enabled"
ssl_support_tls1_1 = "use_default"
ssl_support_tls1_2 = "enabled"
ssl_server_cert_host_mapping = [
  {
    ssl_server_cert_host = "h1pipeline.devops.int.ovp.bskyb.com"
    ssl_server_alt_certs = ["testssl001"]
    ssl_server_cert = "testssl002"
  },
  {
    ssl_server_cert_host = "s1sta07-v00.devops.prd.ovp.bskyb.com"
    ssl_server_cert = "testsslkey1"
  },
]
ocsp_enable = false
ocsp_issuers = [
  {
     issuer = "_DEFAULT_"
     aia = false
     nonce = "on"
     required = "none"
     responder_cert = "ANOTHER_RESPONDER_CERT"
     signer = "SOME_SIGNER"
     url = "http://www.sky.com/"
  },
]
}
`, virtualServerName)
}
