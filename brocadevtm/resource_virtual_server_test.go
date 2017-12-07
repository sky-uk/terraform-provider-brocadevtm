package brocadevtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func TestAccBrocadeVTMVirtualServerBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	virtualServerName := fmt.Sprintf("acctest_brocadevtm_virtual_server-%d", randomInt)
	resourceName := "brocadevtm_virtual_server.acctest"

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
				Config: testAccBrocadeVTMVirtualServerCreate(virtualServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMVirtualServerExists(virtualServerName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", virtualServerName),
					resource.TestCheckResourceAttr(resourceName, "add_cluster_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "add_x_forwarded_for", "true"),
					resource.TestCheckResourceAttr(resourceName, "add_x_forwarded_proto", "true"),
					resource.TestCheckResourceAttr(resourceName, "autodetect_upgrade_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "close_with_rst", "true"),
					resource.TestCheckResourceAttr(resourceName, "completionrules.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("completionrules"), "completionRule1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("completionrules"), "completionRule2"),
					resource.TestCheckResourceAttr(resourceName, "connect_timeout", "50"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ftp_force_server_secure", "true"),
					resource.TestCheckResourceAttr(resourceName, "glb_services.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("glb_services"), "testservice"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("glb_services"), "testservice2"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_any", "true"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_hosts.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("listen_on_hosts"), "host1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("listen_on_hosts"), "host2"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_traffic_ips.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("listen_on_traffic_ips"), "ip1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("listen_on_traffic_ips"), "ip2"),
					resource.TestCheckResourceAttr(resourceName, "note", "create acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "pool", "test-pool"),
					resource.TestCheckResourceAttr(resourceName, "port", "50"),
					resource.TestCheckResourceAttr(resourceName, "protection_class", "testProtectionClass"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "dns"),
					resource.TestCheckResourceAttr(resourceName, "request_rules.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("request_rules"), "ruleOne"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("request_rules"), "ruleTwo"),
					resource.TestCheckResourceAttr(resourceName, "response_rules.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("response_rules"), "ruleOne"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("response_rules"), "ruleTwo"),
					resource.TestCheckResourceAttr(resourceName, "slm_class", "testClass"),
					resource.TestCheckResourceAttr(resourceName, "so_nagle", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl_client_cert_headers", "all"),
					resource.TestCheckResourceAttr(resourceName, "ssl_decrypt", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl_honor_fallback_scsv", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "transparent", "true"),
					resource.TestCheckResourceAttr(resourceName, "connection_errors.0.error_file", "testErrorFile"),
					resource.TestCheckResourceAttr(resourceName, "smtp.0.expect_starttls", "true"),
					resource.TestCheckResourceAttr(resourceName, "tcp.0.proxy_close", "true"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_class", "test"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.keepalive", "true"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.keepalive_timeout", "500"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_client_buffer", "5000"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_server_buffer", "7000"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_transaction_duration", "5"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.server_first_banner", "testbanner"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.timeout", "4000"),
					resource.TestCheckResourceAttr(resourceName, "cookie.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.domain", "no_rewrite"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.new_domain", "testdomain"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.path_regex", "testregex"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.path_replace", "testreplace"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.secure", "no_modify"),
					resource.TestCheckResourceAttr(resourceName, "dns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.edns_client_subnet", "true"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.edns_udpsize", "3000"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.max_udpsize", "4000"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.rrset_order", "cyclic"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.verbose", "true"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.zones.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("dns", "zones"), "testzone1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("dns", "zones"), "testzone2"),
					resource.TestCheckResourceAttr(resourceName, "ftp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.data_source_port", "10"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.force_client_secure", "true"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.port_range_high", "50"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.port_range_low", "5"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.ssl_data", "true"),
					resource.TestCheckResourceAttr(resourceName, "gzip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.compress_level", "5"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.etag_rewrite", "weaken"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.include_mime.#", "3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("gzip", "include_mime"), "mimetype1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("gzip", "include_mime"), "mimetype2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("gzip", "include_mime"), "mimetype3"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.max_size", "4000"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.min_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.no_size", "true"),
					resource.TestCheckResourceAttr(resourceName, "http.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "http.0.chunk_overhead_forwarding", "eager"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_regex", "testregex"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_replace", "testlocationreplace"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_rewrite", "never"),
					resource.TestCheckResourceAttr(resourceName, "http.0.mime_default", "text/html"),
					resource.TestCheckResourceAttr(resourceName, "http.0.mime_detect", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.connect_timeout", "50"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.data_frame_size", "200"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.header_table_size", "4096"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_blacklist.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("http2", "headers_index_blacklist"), "header1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("http2", "headers_index_blacklist"), "header2"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_default", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_whitelist.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("http2", "headers_index_whitelist"), "header3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("http2", "headers_index_whitelist"), "header4"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.idle_timeout_no_streams", "60"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.idle_timeout_open_streams", "120"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_concurrent_streams", "10"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_frame_size", "20000"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_header_padding", "10"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.merge_cookie_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.stream_window_size", "200"),
					resource.TestCheckResourceAttr(resourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log.0.client_connection_failures", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.format", "logfile/format"),
					resource.TestCheckResourceAttr(resourceName, "log.0.save_all", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.server_connection_failures", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.session_persistence_verbose", "true"),
					resource.TestCheckResourceAttr(resourceName, "log.0.ssl_failures", "true"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.0.save_all", "true"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.0.trace_io", "true"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_port_range_high", "50"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_port_range_low", "20"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_timeout", "35"),
					resource.TestCheckResourceAttr(resourceName, "sip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.dangerous_requests", "forbid"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.follow_route", "true"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.max_connection_mem", "50"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.mode", "full_gateway"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.rewrite_uri", "true"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_port_range_high", "60"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_port_range_low", "40"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.timeout_messages", "true"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.transaction_timeout", "20"),
					resource.TestCheckResourceAttr(resourceName, "ssl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.add_http_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.client_cert_cas.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "client_cert_cas"), "cas1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "client_cert_cas"), "cas2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.elliptic_curves.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "elliptic_curves"), "P256"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "elliptic_curves"), "P384"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.#", "3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "issued_certs_never_expire"), "cas1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "issued_certs_never_expire"), "cas2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "issued_certs_never_expire"), "cas3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ocsp_issuers"), "issuerName"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "ocsp_issuers"), "issuerName2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_max_response_age", "50"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_stapling", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_time_tolerance", "50"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_timeout", "20"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.prefer_sslv3", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.request_client_cert", "request"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.send_close_alerts", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_alt_certificates.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_alt_certificates"), "testssl001"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_default", "testssl002"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_host_mapping.#", "3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "fakehost1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "fakehost2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "fakehost3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert4"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert5"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert6"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert7"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert8"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert9"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.signature_algorithms", "RSA_SHA256"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_ciphers", "SSL_RSA_WITH_AES_128_CBC_SHA256"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_ssl2", "use_default"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_ssl3", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1_1", "use_default"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1_2", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.trust_magic", "true"),
					resource.TestCheckResourceAttr(resourceName, "syslog.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.format", "syslog/format"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.ip_end_point", "127.0.0.1:515"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.msg_len_limit", "500"),
					resource.TestCheckResourceAttr(resourceName, "udp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.end_point_persistence", "true"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.port_smp", "true"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.response_datagrams_expected", "-1"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.timeout", "25"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.control_out", "testcontrolout"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.error_page_time", "20"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.max_time", "50"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.refresh_time", "4"),
				),
			},
			{
				Config: testAccBrocadeVTMVirtualServerUpdate(virtualServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMVirtualServerExists(virtualServerName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", virtualServerName),
					resource.TestCheckResourceAttr(resourceName, "add_cluster_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "add_x_forwarded_for", "false"),
					resource.TestCheckResourceAttr(resourceName, "add_x_forwarded_proto", "false"),
					resource.TestCheckResourceAttr(resourceName, "autodetect_upgrade_headers", "false"),
					resource.TestCheckResourceAttr(resourceName, "close_with_rst", "false"),
					resource.TestCheckResourceAttr(resourceName, "completionrules.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("completionrules"), "completionRule2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("completionrules"), "completionRule3"),
					resource.TestCheckResourceAttr(resourceName, "connect_timeout", "100"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ftp_force_server_secure", "false"),
					resource.TestCheckResourceAttr(resourceName, "glb_services.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("glb_services"), "testservice3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("glb_services"), "testservice4"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_any", "false"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_hosts.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("listen_on_hosts"), "host3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("listen_on_hosts"), "host4"),
					resource.TestCheckResourceAttr(resourceName, "listen_on_traffic_ips.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("listen_on_traffic_ips"), "ip1"),
					resource.TestCheckResourceAttr(resourceName, "note", "update acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "pool", "test-pool"),
					resource.TestCheckResourceAttr(resourceName, "port", "100"),
					resource.TestCheckResourceAttr(resourceName, "protection_class", "testProtectionClassUpdate"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "ftp"),
					resource.TestCheckResourceAttr(resourceName, "request_rules.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("request_rules"), "ruleThree"),
					resource.TestCheckResourceAttr(resourceName, "response_rules.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSet("response_rules"), "ruleFour"),
					resource.TestCheckResourceAttr(resourceName, "slm_class", "testClassUpdate"),
					resource.TestCheckResourceAttr(resourceName, "so_nagle", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl_client_cert_headers", "simple"),
					resource.TestCheckResourceAttr(resourceName, "ssl_decrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl_honor_fallback_scsv", "use_default"),
					resource.TestCheckResourceAttr(resourceName, "transparent", "false"),
					resource.TestCheckResourceAttr(resourceName, "connection_errors.0.error_file", "testErrorFileUpdate"),
					resource.TestCheckResourceAttr(resourceName, "smtp.0.expect_starttls", "false"),
					resource.TestCheckResourceAttr(resourceName, "tcp.0.proxy_close", "false"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aptimizer.0.profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_class", "testUpdate"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.keepalive", "false"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.keepalive_timeout", "250"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_client_buffer", "5050"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_server_buffer", "7070"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.max_transaction_duration", "10"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.server_first_banner", "testbannerupdate"),
					resource.TestCheckResourceAttr(resourceName, "vs_connection.0.timeout", "4050"),
					resource.TestCheckResourceAttr(resourceName, "cookie.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.domain", "set_to_named"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.new_domain", "testdomainupdate"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.path_regex", "testregexupdate"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.path_replace", "testreplaceupdate"),
					resource.TestCheckResourceAttr(resourceName, "cookie.0.secure", "unset_secure"),
					resource.TestCheckResourceAttr(resourceName, "dns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.edns_client_subnet", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.edns_udpsize", "3050"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.max_udpsize", "4050"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.rrset_order", "fixed"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.verbose", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns.0.zones.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("dns", "zones"), "testzone2"),
					resource.TestCheckResourceAttr(resourceName, "ftp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.data_source_port", "15"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.force_client_secure", "false"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.port_range_high", "55"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.port_range_low", "7"),
					resource.TestCheckResourceAttr(resourceName, "ftp.0.ssl_data", "false"),
					resource.TestCheckResourceAttr(resourceName, "gzip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.compress_level", "8"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.etag_rewrite", "wrap"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.include_mime.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("gzip", "include_mime"), "mimetype3"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.max_size", "4050"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.min_size", "50"),
					resource.TestCheckResourceAttr(resourceName, "gzip.0.no_size", "false"),
					resource.TestCheckResourceAttr(resourceName, "http.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "http.0.chunk_overhead_forwarding", "lazy"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_regex", "testregexupdate"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_replace", "testlocationreplaceupdate"),
					resource.TestCheckResourceAttr(resourceName, "http.0.location_rewrite", "if_host_matches"),
					resource.TestCheckResourceAttr(resourceName, "http.0.mime_default", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "http.0.mime_detect", "false"),
					resource.TestCheckResourceAttr(resourceName, "http2.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.connect_timeout", "75"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.data_frame_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.header_table_size", "4098"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_blacklist.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("http2", "headers_index_blacklist"), "header3"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("http2", "headers_index_blacklist"), "header4"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_default", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.headers_index_whitelist.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("http2", "headers_index_whitelist"), "header1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("http2", "headers_index_whitelist"), "header2"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.idle_timeout_no_streams", "80"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.idle_timeout_open_streams", "150"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_concurrent_streams", "15"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_frame_size", "20050"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.max_header_padding", "13"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.merge_cookie_headers", "false"),
					resource.TestCheckResourceAttr(resourceName, "http2.0.stream_window_size", "201"),
					resource.TestCheckResourceAttr(resourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log.0.client_connection_failures", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.format", "logfile/updateformat"),
					resource.TestCheckResourceAttr(resourceName, "log.0.save_all", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.server_connection_failures", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.session_persistence_verbose", "false"),
					resource.TestCheckResourceAttr(resourceName, "log.0.ssl_failures", "false"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "recent_connections.0.save_all", "false"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "request_tracing.0.trace_io", "false"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_port_range_high", "75"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_port_range_low", "23"),
					resource.TestCheckResourceAttr(resourceName, "rtsp.0.streaming_timeout", "37"),
					resource.TestCheckResourceAttr(resourceName, "sip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.dangerous_requests", "node"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.follow_route", "false"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.max_connection_mem", "60"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.mode", "sip_gateway"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.rewrite_uri", "false"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_port_range_high", "73"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_port_range_low", "45"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.streaming_timeout", "19"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.timeout_messages", "false"),
					resource.TestCheckResourceAttr(resourceName, "sip.0.transaction_timeout", "23"),
					resource.TestCheckResourceAttr(resourceName, "ssl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.add_http_headers", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.client_cert_cas.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "client_cert_cas"), "cas2"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.elliptic_curves.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "elliptic_curves"), "P384"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.issued_certs_never_expire.#", "2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "issued_certs_never_expire"), "cas2"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "issued_certs_never_expire"), "cas3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_issuers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_max_response_age", "55"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_stapling", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_time_tolerance", "55"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ocsp_timeout", "25"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.prefer_sslv3", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.request_client_cert", "require"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.send_close_alerts", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_alt_certificates.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_alt_certificates"), "testssl002"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_default", "testssl001"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.server_cert_host_mapping.#", "1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "fakehost7"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert6"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert5"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert1"),
					util.AccTestCheckValueInKeyPattern(resourceName, util.AccTestCreateRegexPatternForSetItems("ssl", "server_cert_host_mapping"), "altcert3"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.signature_algorithms", "ECDSA_SHA256"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_ciphers", "SSL_RSA_WITH_RC4_128_SHA"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_ssl2", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_ssl3", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1_1", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.ssl_support_tls1_2", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ssl.0.trust_magic", "false"),
					resource.TestCheckResourceAttr(resourceName, "syslog.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.format", "syslog/formatupdate"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.ip_end_point", "127.0.0.1:700"),
					resource.TestCheckResourceAttr(resourceName, "syslog.0.msg_len_limit", "505"),
					resource.TestCheckResourceAttr(resourceName, "udp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.end_point_persistence", "false"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.port_smp", "false"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.response_datagrams_expected", "50"),
					resource.TestCheckResourceAttr(resourceName, "udp.0.timeout", "33"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.control_out", "testcontroloutupdate"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.error_page_time", "25"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.max_time", "75"),
					resource.TestCheckResourceAttr(resourceName, "web_cache.0.refresh_time", "9"),
				),
			},
		},
	})
}

func testAccBrocadeVTMVirtualServerCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "infoblox_virtual_server" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		vs := make(map[string]interface{})
		client.WorkWithConfigurationResources()
		err := client.GetByName("virtual_servers", name, &vs)
		if err != nil {
			return fmt.Errorf("[ERROR] Brocade vTM error occurred while retrieving Virtual Server: %s", err)
		}
		if client.StatusCode == http.StatusOK {
			return fmt.Errorf("[ERROR] Brocade vTM Virtual Server %s still exists", name)
		}
	}
	return nil
}

func testAccBrocadeVTMVirtualServerExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Brocade vTM Virtual Server %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Brocade vTM Virtual Server ID not set for %s in resources", name)
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)

		vs := make(map[string]interface{})
		client.WorkWithConfigurationResources()
		err := client.GetByName("virtual_servers", name, &vs)
		if err != nil {
			return fmt.Errorf("[ERROR] Brocade vTM Virtual Server - error while retrieving virtual server %s: %s", name, err)
		}
		if client.StatusCode == http.StatusOK {
			return nil
		}
		return fmt.Errorf("[ERROR] Brocade vTM Virtual Server %s not found on remote vTM", name)
	}
}

func testAccBrocadeVTMVirtualServerNoName() string {
	return `
resource "brocadevtm_virtual_server" "acctest" {
pool = "test-pool"
port = 80
}
`
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
	listen_on_hosts = ["host2","host1"]
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
	connection_errors = {
		error_file = "testErrorFile"
	}
	smtp = {
		expect_starttls = true
	}
	tcp = {
		proxy_close = true
	}
	aptimizer = {
		enabled = true
		profile = [{
			name = "profile1"
			urls = ["url2","url1"]
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
		include_mime = ["mimetype1","mimetype2","mimetype3"]
		max_size = 4000
		min_size = 5
		no_size = true
	}

	http = {
		chunk_overhead_forwarding = "eager"
		location_regex = "testregex"
		location_replace = "testlocationreplace"
		location_rewrite = "never"
		mime_default = "text/html"
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

	    	server_cert_host_mapping = [
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
		{
		  host = "fakehost3"
		  certificate = "altcert2"
		  alt_certificates = ["altcert7","altcert8","altcert9"]
		}
		]

		signature_algorithms = "RSA_SHA256"
		ssl_ciphers = "SSL_RSA_WITH_AES_128_CBC_SHA256"
		ssl_support_ssl2 = "use_default"
		ssl_support_ssl3 = "disabled"
		ssl_support_tls1 = "enabled"
		ssl_support_tls1_1 = "use_default"
		ssl_support_tls1_2 = "disabled"
		trust_magic = true
	  }

	  syslog = {
	    enabled = true
	    format = "syslog/format"
	    ip_end_point = "127.0.0.1:515"
	    msg_len_limit = 500
	  }

	  udp = {
	    end_point_persistence = true
	    port_smp = true
	    response_datagrams_expected = -1
	    timeout = 25
          }

	web_cache = {
	    control_out = "testcontrolout"
	    enabled = true
	    error_page_time = 20
	    max_time = 50
	    refresh_time = 4
 	 }
}
`, virtualServerName)
}

func testAccBrocadeVTMVirtualServerUpdate(virtualServerName string) string {
	return fmt.Sprintf(`

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
`, virtualServerName)
}
