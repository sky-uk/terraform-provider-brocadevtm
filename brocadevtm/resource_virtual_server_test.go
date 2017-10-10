package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-rest-api"
	"net/http"
	"regexp"
	"testing"
	"github.com/sky-uk/go-brocade-vtm/api/virtualserver"
)

func TestAccBrocadeVTMVirtualServerBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	virtualServerName := fmt.Sprintf("acctest_brocadevtm_virtual_server-%d", randomInt)
	resourceName := "brocadevtm_virtual_server.acctest"

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
				ExpectError: regexp.MustCompile(`Port has to be between 1 and 65535`),
			},
			{
				Config:      testAccBrocadeVTMVirtualServerInvalidProtocol(virtualServerName),
				ExpectError: regexp.MustCompile(`must be one of client_first, dns, dns_tcp, ftp, http, https, imaps, imapv2, imapv3, imapv4, ldap, ldaps, pop3, pop3s, rtsp, server_first, siptcp, sipudp, smtp, ssl, stream, telnet, udp or udpstreaming`),
			},
			{
				Config: testAccBrocadeVTMVirtualServerCreate(virtualServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMVirtualServerExists(virtualServerName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", virtualServerName),
				),
			},
			/*
				{
					Config: testAccBrocadeVTMVirtualServerUpdate(virtualServerName),
					Check: resource.ComposeTestCheckFunc(
						testAccBrocadeVTMVirtualServerExists(virtualServerName, resourceName),
						resource.TestCheckResourceAttr(resourceName, "name", virtualServerName),

					),
				},
			*/
		},
	})
}

func testAccBrocadeVTMVirtualServerCheckDestroy(state *terraform.State, name string) error {
	vtmClient := testAccProvider.Meta().(*rest.Client)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "infoblox_virtual_server" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		api := virtualserver.NewGet(name)
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Error: Brocade vTM error occurred while retrieving Virtual Server: %v", err)
		}
		if api.StatusCode() == http.StatusOK {
			return fmt.Errorf("Error: Brocade vTM Virtual Server %s still exists", name)
		}
	}
	return nil
}

func testAccBrocadeVTMVirtualServerExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Virtual Server %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Virtual Server ID not set for %s in resources", name)
		}
		vtmClient := testAccProvider.Meta().(*rest.Client)
		api := virtualserver.NewGet(name)
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Brocade vTM Virtual Server - error while retrieving User Group: %v", err)
		}
		if api.StatusCode() == http.StatusOK {
			return nil
		}
		return fmt.Errorf("Brocade vTM Virtual Server %s not found on remote vTM", name)
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

func testAccBrocadeVTMVirtualServerCreate(virtualServerName string) string {
	return fmt.Sprintf(`



resource "brocadevtm_virtual_server" "acctest" {

	name = "%s"
	add_cluster_ip = true
	add_x_forwarded_for = true
	add_x_forwarded_proto = true
	autodetect_upgrade_headers = true
	bandwidth_class = "test"
	close_with_rst = true
	completionrules = ["completionRule1","completionRule2"]
	connect_timeout = 50
	enabled = true
	ftp_force_server_secure = true
	glb_services = ["testservice","testservice2"]
	listen_on_any = true
	listen_on_hosts = ["host1","host2"]
	listen_on_traffic_ips = ["ip1","ip2"]
	note = "create acceptance test"
	pool = "test-pool"
	port = 50
	protection_class = "testProtectionClass"
	protocol = "dns"
	request_rules = ["ruleOne","ruleTwo"]
	response_rules = ["ruleOne","ruleTwo"]
	slm_class = "testClass"
	so_nagle = true
	ssl_client_cert_headers = "all"
	ssl_decrypt = true
	ssl_honor_fallback_scsv = "enabled"
	transparent = true
	error_file = "testErrorFile"
	expect_starttls = true
	proxy_close = true

	aptimizer = {
		enabled = true
		profile = [{
			name = "profile1"
			urls = ["url1","url2"]
		},
		{
			name = "profile2"
			urls = ["url3","url4","url5"]
		},
		{
			name = "profile3"
			urls = ["url6","url7","url8","url9"]
		}
		]
	}

	vs_connection = {
		keepalive = true
		keepalive_timeout = 500
		max_client_buffer = 5000
		max_server_buffer = 7000
		max_transaction_duration = 5
		server_first_banner = "testbanner"
		timeout = 4000
	}

	cookie = {
		domain = "no_rewrite"
		new_domain = "testdomain"
		path_regex = "testregex"
		path_replace = "testreplace"
		secure = "no_modify"
	}

	dns = {
		edns_client_subnet = true
		edns_udpsize = 3000
		max_udpsize  = 4000
		rrset_order = "cyclic"
		verbose = true
		zones = ["testzone1","testzone2"]
	}

	ftp = {
		data_source_port = 10
		force_client_secure = true
		port_range_high = 50
		port_range_low = 5
		ssl_data = true
	}

	gzip = {
		compress_level = 5
		enabled = true
		etag_rewrite = "weaken"
		include_mime = ["mimetyp1","mimetype2","mimetype3"]
		max_size = 4000
		min_size = 5
		no_size = true
	}

	http = {
		chunk_overhead_forwarding = "eager"
		location_regex = "testregex"
		location_replace = "testlocationreplace"
		location_rewrite = "never"
		mime_default = "text/html text/plain"
		mime_detect = true
	}

	http2 = {
		connect_timeout = 50
		data_frame_size = 200
		enabled = true
		header_table_size = 4096
		headers_index_blacklist = ["header1","header2"]
		headers_index_default = true
		headers_index_whitelist = ["header3","header4"]
		idle_timeout_no_streams = 60
		idle_timeout_open_streams = 120
		max_concurrent_streams = 10
		max_frame_size = 20000
		max_header_padding = 10
		merge_cookie_headers = true
		stream_window_size = 200
	}

	log = {
		client_connection_failures = true
		enabled = true
		filename = "file/name"
		format = "logfile/format"
		save_all = true
		server_connection_failures = true
		session_persistence_verbose = true
		ssl_failures = true
	}

	recent_connections = {
		enabled = true
		save_all = true
	}

	request_tracing = {
		enabled = true
		trace_io = true
	}

	rtsp = {
		streaming_port_range_high = 50
		streaming_port_range_low = 20
		streaming_timeout = 35
	}

	sip = {
		dangerous_requests = "forbid"
		follow_route = true
		max_connection_mem = 50
		mode = "full_gateway"
		rewrite_uri = true
		streaming_port_range_high = 60
		streaming_port_range_low = 40
		streaming_timeout = 15
		timeout_messages = true
		transaction_timeout = 20
	}

	ssl = {
		add_http_headers = true
		client_cert_cas = ["cas1","cas2"]
		elliptic_curves = ["P256","P384"]
		issued_certs_never_expire = ["cas1","cas2","cas3"]
		ocsp_enable = true

		ocsp_issuers = [
		{
			issuer = "issuerName"
			aia = true
			nonce = "strict"
			required = "optional"
			responder_cert = "respondercert"
			signer = "fakesigner"
			url = "fake.url"
		},
		{
			issuer = "issuerName2"
			aia = true
			nonce = "strict"
			required = "optional"
			responder_cert = "respondercert2"
			signer = "fakesigner2"
			url = "fake2.url"
		},
		]
		ocsp_max_response_age = 50
		ocsp_stapling = true
	    	ocsp_time_tolerance = 50
	    	ocsp_timeout = 20
	    	prefer_sslv3 = true
	    	request_client_cert = "request"
	    	send_close_alerts = true
	    	server_cert_alt_certificates = ["testssl001"]
	    	server_cert_default = "testssl002"

	    	ssl_server_cert_host_mapping = [
		{
		  host = "fakehost1"
		  certificate = "altcert4"
		  alt_certificates = ["altcert1","altcert2","altcert3"]

		},
		{
		  host = "fakehost2"
		  certificate = "altcert1"
		  alt_certificates = ["altcert4","altcert5","altcert6"]
		},
		{ host = "fakehost3"
		  certificate = "altcert2"
		  alt_certificates = ["altcert7","altcert8","altcert9"]
		}
		]

		signature_algorithms = "RSA_SHA256"
		ssl_ciphers = " SSL_RSA_WITH_AES_128_CBC_SHA256"
		ssl_support_ssl2 = "use_default"
		ssl_support_ssl3 = "disabled"
		ssl_support_tls1 = "enabled"
		ssl_support_tls1_1 = "use_default"
		ssl_support_tls1_2 = "disabled"
		trust_magic = true
	  }



}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerUpdate(virtualServerName string) string {
	return fmt.Sprintf(`

resource "brocadevtm_virtual_server" "acctest" {
name = "%s"
pool = "pool-demo"
port = 443
}
`, virtualServerName)
}
