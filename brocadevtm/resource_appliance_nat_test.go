package brocadevtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func compileRegex(key, attr string) *regexp.Regexp {
	return regexp.MustCompile(key + `\.[0-9]+\.` + attr)
}

func TestAccBrocadeVTMApplianceNat(t *testing.T) {

	randomInt := acctest.RandInt()
	name := fmt.Sprintf("acctest_brocadevtm_appliance_nat-%d", randomInt)
	resourceName := fmt.Sprintf("brocadevtm_appliance_nat.acctest")
	fmt.Printf("\n\nAppliance Nat is %s.\n\n", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMApplianceNatCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBrocadeVTMApplianceNatEmptyResource(),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
				),
			},
			{
				Config: manyToOneCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "rule_number"), "4765348"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "tip"), name),
				),
			},
			{
				Config: manyToOneUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "rule_number"), "58673"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_all_ports", "tip"), name),
				),
			},
			{
				Config: manyToOnePortLockedCreateTemplate(name),
				//Destroy:                   false,
				//PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "rule_number"), "20001"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "port"), "2728"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "tip"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "protocol"), "tcp"),
				),
			},
			{
				Config: manyToOnePortLockedUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "rule_number"), "20001"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "port"), "2730"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "tip"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "protocol"), "tcp"),
				),
			},
			{
				Config: oneToOneCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "rule_number"), "20001"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "port"), "2728"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "tip"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "protocol"), "tcp"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("one_to_one", "rule_number"), "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("one_to_one", "rule_number"), "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("one_to_one", "ip"), "192.168.10.10"),
				),
			},
			{
				Config: oneToOneUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "rule_number"), "20002"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "port"), "2730"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "pool"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "tip"), name),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("many_to_one_port_locked", "protocol"), "tcp"),
				),
			},
			{
				Config: portMappingCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("port_mapping", "dport_first"), "10"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("port_mapping", "dport_last"), "20"),
				),
			},
			{
				Config: portMappingUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMApplianceNatExists(name, resourceName),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("port_mapping", "dport_first"), "10"),
					util.AccTestCheckValueInKeyPattern(resourceName, compileRegex("port_mapping", "dport_last"), "30"),
				),
			},
		},
	})
}

func testAccBrocadeVTMApplianceNatCheckDestroy(state *terraform.State, name string) error {

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_appliance_nat" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()

		resources, err := client.GetAllResources("appliance/nat")
		if err != nil {
			return fmt.Errorf("[ERROR] Brocade vTM error whilst retrieving appliance nat: %+v", err)
		}
		for _, resource := range resources {
			if resource["name"] == name {
				return fmt.Errorf("[ERROR] Brocade vTM Appliance nat %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMApplianceNatExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Brocade vTM Appliance Nat %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Brocade vTM Appliance Nat ID not set for %s in resources", name)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()
		nat := make(map[string]interface{})
		err := client.GetByName("appliance/nat", "", &nat)
		if err != nil {
			return fmt.Errorf("[ERROR] Brocade vTM error whilst retriving appliance nat: %v", err)
		}
		return nil
	}
}

func testAccBrocadeVTMApplianceNatEmptyResource() string {
	return fmt.Sprintf(`
resource "brocadevtm_appliance_nat" "acctest" {
  many_to_one_all_ports = []
  many_to_one_port_locked = []
  one_to_one = []
  port_mapping = []
}
`)
}

func manyToOneCreateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = [{
	  rule_number = 4765348
	  pool = "%s"
	  tip = "%s"
  }]
  many_to_one_port_locked = []
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func manyToOneUpdateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = [{
	  rule_number = 58673
	  pool = "%s"
	  tip = "%s"
  }]
  many_to_one_port_locked = []
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func manyToOnePortLockedCreateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2728
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func manyToOnePortLockedUpdateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2730
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func oneToOneCreateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2728
	  protocol = "tcp"
  }]
  one_to_one = [{
	  rule_number = 1
	  enable_inbound = false
	  ip = "192.168.10.10"
	  tip = "%s"
	},{
	  rule_number = 2
	  enable_inbound = false
	  ip = "192.168.10.11"
	  tip = "%s"
	}
  ]
  port_mapping = []
}
`, name, name, name, name, name, name)
}

func oneToOneUpdateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = ["brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest"]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20002
	  pool = "%s"
	  tip = "%s"
	  port = 2730
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = []
}
`, name, name, name, name)
}

func portMappingCreateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_virtual_server" "acctest" {

	name = "%s"
	add_cluster_ip = false
	add_x_forwarded_for = false
	add_x_forwarded_proto = false
	autodetect_upgrade_headers = false
	bandwidth_class = "testUpdate"
	close_with_rst = false
	completionrules = ["completionRule2","completionRule3"]
	connect_timeout = 100
	enabled = false
	ftp_force_server_secure = false
	glb_services = ["testservice3","testservice4"]
	listen_on_any = false
	listen_on_hosts = ["host3","host4"]
	listen_on_traffic_ips = ["ip1"]
	note = "update acceptance test"
	pool = "test-pool"
	port = 100
	protection_class = "testProtectionClassUpdate"
	protocol = "ftp"
	request_rules = ["ruleThree"]
	response_rules = ["ruleFour"]
	slm_class = "testClassUpdate"
	so_nagle = false
	ssl_client_cert_headers = "simple"
	ssl_decrypt = false
	ssl_honor_fallback_scsv = "use_default"
	transparent = false
	connection_errors = {
		error_file = "testErrorFileUpdate"
	}
	smtp = {
		expect_starttls = false
	}
	tcp = {
		proxy_close = false
	}
	aptimizer = {
		enabled = true
		profile = [{
			name = "profile1"
			urls = ["url4","url3"]
		}
		]
	}

	vs_connection = {
		keepalive = false
		keepalive_timeout = 250
		max_client_buffer = 5050
		max_server_buffer = 7070
		max_transaction_duration = 10
		server_first_banner = "testbannerupdate"
		timeout = 4050
	}

	cookie = {
		domain = "set_to_named"
		new_domain = "testdomainupdate"
		path_regex = "testregexupdate"
		path_replace = "testreplaceupdate"
		secure = "unset_secure"
	}

	dns = {
		edns_client_subnet = false
		edns_udpsize = 3050
		max_udpsize  = 4050
		rrset_order = "fixed"
		verbose = false
		zones = ["testzone2"]
	}

	ftp = {
		data_source_port = 15
		force_client_secure = false
		port_range_high = 55
		port_range_low = 7
		ssl_data = false
	}

	gzip = {
		compress_level = 8
		enabled = false
		etag_rewrite = "wrap"
		include_mime = ["mimetype3"]
		max_size = 4050
		min_size = 50
		no_size = false
	}

	http = {
		chunk_overhead_forwarding = "lazy"
		location_regex = "testregexupdate"
		location_replace = "testlocationreplaceupdate"
		location_rewrite = "if_host_matches"
		mime_default = "application/json"
		mime_detect = false
	}

	http2 = {
		connect_timeout = 75
		data_frame_size = 100
		enabled = false
		header_table_size = 4098
		headers_index_blacklist = ["header3","header4"]
		headers_index_default = true
		headers_index_whitelist = ["header1","header2"]
		idle_timeout_no_streams = 80
		idle_timeout_open_streams = 150
		max_concurrent_streams = 15
		max_frame_size = 20050
		max_header_padding = 13
		merge_cookie_headers = false
		stream_window_size = 201
	}

	log = {
		client_connection_failures = false
		enabled = false
		format = "logfile/updateformat"
		save_all = false
		server_connection_failures = false
		session_persistence_verbose = false
		ssl_failures = false
	}

	recent_connections = {
		enabled = false
		save_all = false
	}

	request_tracing = {
		enabled = false
		trace_io = false
	}

	rtsp = {
		streaming_port_range_high = 75
		streaming_port_range_low = 23
		streaming_timeout = 37
	}

	sip = {
		dangerous_requests = "node"
		follow_route = false
		max_connection_mem = 60
		mode = "sip_gateway"
		rewrite_uri = false
		streaming_port_range_high = 73
		streaming_port_range_low = 45
		streaming_timeout = 19
		timeout_messages = false
		transaction_timeout = 23
	}

	ssl = {
		add_http_headers = false
		client_cert_cas = ["cas2"]
		elliptic_curves = ["P384"]
		issued_certs_never_expire = ["cas2","cas3"]
		ocsp_enable = false

		ocsp_issuers = [
		{
			issuer = "issuerNameUpdated"
			aia = true
			nonce = "off"
			required = "optional"
			responder_cert = "respondercert"
			signer = "fakesigner"
			url = "fake.url"
		},
		]
		ocsp_max_response_age = 55
		ocsp_stapling = false
	    ocsp_time_tolerance = 55
	    ocsp_timeout = 25
	    prefer_sslv3 = false
	    request_client_cert = "require"
	    send_close_alerts = false
	    server_cert_alt_certificates = ["testssl002"]
	    server_cert_default = "testssl001"
	    server_cert_host_mapping = [
		{
		  host = "fakehost7"
		  certificate = "altcert6"
		  alt_certificates = ["altcert5","altcert1","altcert3"]

		},
		]

		signature_algorithms = "ECDSA_SHA256"
		ssl_ciphers = "SSL_RSA_WITH_RC4_128_SHA"
		ssl_support_ssl2 = "disabled"
		ssl_support_ssl3 = "disabled"
		ssl_support_tls1 = "disabled"
		ssl_support_tls1_1 = "disabled"
		ssl_support_tls1_2 = "disabled"
		trust_magic = false
	  }

	  syslog = {
	    enabled = false
	    format = "syslog/formatupdate"
	    ip_end_point = "127.0.0.1:700"
	    msg_len_limit = 505
	  }

	  udp = {
	    end_point_persistence = false
	    port_smp = false
	    response_datagrams_expected = 50
	    timeout = 33
          }

	web_cache = {
	    control_out = "testcontroloutupdate"
	    enabled = false
	    error_page_time = 25
	    max_time = 75
	    refresh_time = 9
 	 }
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = [ "brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest", "brocadevtm_virtual_server.acctest" ]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2728
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = [{
	  rule_number = 1
	  dport_first = 10
	  dport_last = 20
	  virtual_server = "%s"
  }]
}
`, name, name, name, name, name, name)
}

func portMappingUpdateTemplate(name string) string {
	return fmt.Sprintf(`

resource "brocadevtm_traffic_ip_group" "acctest" {
  name = "%s"
  enabled = true
  ipaddresses = ["192.168.100.10"]
}

resource "brocadevtm_pool" "acctest" {
  name = "%s"
  nodes_table = [
    {
      node = "192.168.10.10:80"
      state = "active"
    },
  ]
}

resource "brocadevtm_virtual_server" "acctest" {

	name = "%s"
	add_cluster_ip = false
	add_x_forwarded_for = false
	add_x_forwarded_proto = false
	autodetect_upgrade_headers = false
	bandwidth_class = "testUpdate"
	close_with_rst = false
	completionrules = ["completionRule2","completionRule3"]
	connect_timeout = 100
	enabled = false
	ftp_force_server_secure = false
	glb_services = ["testservice3","testservice4"]
	listen_on_any = false
	listen_on_hosts = ["host3","host4"]
	listen_on_traffic_ips = ["ip1"]
	note = "update acceptance test"
	pool = "test-pool"
	port = 100
	protection_class = "testProtectionClassUpdate"
	protocol = "ftp"
	request_rules = ["ruleThree"]
	response_rules = ["ruleFour"]
	slm_class = "testClassUpdate"
	so_nagle = false
	ssl_client_cert_headers = "simple"
	ssl_decrypt = false
	ssl_honor_fallback_scsv = "use_default"
	transparent = false
	connection_errors = {
		error_file = "testErrorFileUpdate"
	}
	smtp = {
		expect_starttls = false
	}
	tcp = {
		proxy_close = false
	}
	aptimizer = {
		enabled = true
		profile = [{
			name = "profile1"
			urls = ["url4","url3"]
		}
		]
	}

	vs_connection = {
		keepalive = false
		keepalive_timeout = 250
		max_client_buffer = 5050
		max_server_buffer = 7070
		max_transaction_duration = 10
		server_first_banner = "testbannerupdate"
		timeout = 4050
	}

	cookie = {
		domain = "set_to_named"
		new_domain = "testdomainupdate"
		path_regex = "testregexupdate"
		path_replace = "testreplaceupdate"
		secure = "unset_secure"
	}

	dns = {
		edns_client_subnet = false
		edns_udpsize = 3050
		max_udpsize  = 4050
		rrset_order = "fixed"
		verbose = false
		zones = ["testzone2"]
	}

	ftp = {
		data_source_port = 15
		force_client_secure = false
		port_range_high = 55
		port_range_low = 7
		ssl_data = false
	}

	gzip = {
		compress_level = 8
		enabled = false
		etag_rewrite = "wrap"
		include_mime = ["mimetype3"]
		max_size = 4050
		min_size = 50
		no_size = false
	}

	http = {
		chunk_overhead_forwarding = "lazy"
		location_regex = "testregexupdate"
		location_replace = "testlocationreplaceupdate"
		location_rewrite = "if_host_matches"
		mime_default = "application/json"
		mime_detect = false
	}

	http2 = {
		connect_timeout = 75
		data_frame_size = 100
		enabled = false
		header_table_size = 4098
		headers_index_blacklist = ["header3","header4"]
		headers_index_default = true
		headers_index_whitelist = ["header1","header2"]
		idle_timeout_no_streams = 80
		idle_timeout_open_streams = 150
		max_concurrent_streams = 15
		max_frame_size = 20050
		max_header_padding = 13
		merge_cookie_headers = false
		stream_window_size = 201
	}

	log = {
		client_connection_failures = false
		enabled = false
		format = "logfile/updateformat"
		save_all = false
		server_connection_failures = false
		session_persistence_verbose = false
		ssl_failures = false
	}

	recent_connections = {
		enabled = false
		save_all = false
	}

	request_tracing = {
		enabled = false
		trace_io = false
	}

	rtsp = {
		streaming_port_range_high = 75
		streaming_port_range_low = 23
		streaming_timeout = 37
	}

	sip = {
		dangerous_requests = "node"
		follow_route = false
		max_connection_mem = 60
		mode = "sip_gateway"
		rewrite_uri = false
		streaming_port_range_high = 73
		streaming_port_range_low = 45
		streaming_timeout = 19
		timeout_messages = false
		transaction_timeout = 23
	}

	ssl = {
		add_http_headers = false
		client_cert_cas = ["cas2"]
		elliptic_curves = ["P384"]
		issued_certs_never_expire = ["cas2","cas3"]
		ocsp_enable = false

		ocsp_issuers = [
		{
			issuer = "issuerNameUpdated"
			aia = true
			nonce = "off"
			required = "optional"
			responder_cert = "respondercert"
			signer = "fakesigner"
			url = "fake.url"
		},
		]
		ocsp_max_response_age = 55
		ocsp_stapling = false
	    ocsp_time_tolerance = 55
	    ocsp_timeout = 25
	    prefer_sslv3 = false
	    request_client_cert = "require"
	    send_close_alerts = false
	    server_cert_alt_certificates = ["testssl002"]
	    server_cert_default = "testssl001"
	    server_cert_host_mapping = [
		{
		  host = "fakehost7"
		  certificate = "altcert6"
		  alt_certificates = ["altcert5","altcert1","altcert3"]

		},
		]

		signature_algorithms = "ECDSA_SHA256"
		ssl_ciphers = "SSL_RSA_WITH_RC4_128_SHA"
		ssl_support_ssl2 = "disabled"
		ssl_support_ssl3 = "disabled"
		ssl_support_tls1 = "disabled"
		ssl_support_tls1_1 = "disabled"
		ssl_support_tls1_2 = "disabled"
		trust_magic = false
	  }

	  syslog = {
	    enabled = false
	    format = "syslog/formatupdate"
	    ip_end_point = "127.0.0.1:700"
	    msg_len_limit = 505
	  }

	  udp = {
	    end_point_persistence = false
	    port_smp = false
	    response_datagrams_expected = 50
	    timeout = 33
          }

	web_cache = {
	    control_out = "testcontroloutupdate"
	    enabled = false
	    error_page_time = 25
	    max_time = 75
	    refresh_time = 9
 	 }
}

resource "brocadevtm_appliance_nat" "acctest" {
  depends_on = [ "brocadevtm_traffic_ip_group.acctest", "brocadevtm_pool.acctest", "brocadevtm_virtual_server.acctest" ]
  many_to_one_all_ports = []
  many_to_one_port_locked = [{
	  rule_number = 20001
	  pool = "%s"
	  tip = "%s"
	  port = 2728
	  protocol = "tcp"
  }]
  one_to_one = []
  port_mapping = [{
	  rule_number = 1
	  dport_first = 10
	  dport_last = 30
	  virtual_server = "%s"
  }]
}
`, name, name, name, name, name, name)
}
