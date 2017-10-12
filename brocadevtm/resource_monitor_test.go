package brocadevtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
)

func TestAccBrocadeVTMMonitorBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	monitorName := fmt.Sprintf("acctest_brocadevtm_monitor-%d", randomInt)
	monitorResourceName := "brocadevtm_monitor.acctest"

	fmt.Printf("\n\nMonitor Name is %s.\n\n", monitorName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMMonitorCheckDestroy(state, monitorName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMMonitorInvalidName(),
				ExpectError: regexp.MustCompile(`BrocadeVTM Monitor error whilst creating ../virtual_servers/some_random_virtual_server: Response status code: 400`),
			},
			{
				Config: testAccBrocadeVTMMonitorCreateTemplate(monitorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName),
					resource.TestCheckResourceAttr(monitorResourceName, "name", monitorName),
					resource.TestCheckResourceAttr(monitorResourceName, "delay", "6"),
					resource.TestCheckResourceAttr(monitorResourceName, "timeout", "2"),
					resource.TestCheckResourceAttr(monitorResourceName, "failures", "7"),
					resource.TestCheckResourceAttr(monitorResourceName, "verbose", "true"),
					resource.TestCheckResourceAttr(monitorResourceName, "use_ssl", "true"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_host_header", "some_other_header"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_authentication", "admin:password"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_body_regex", "^ok"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_path", "/some/status/page"),
					resource.TestCheckResourceAttr(monitorResourceName, "rtsp_status_regex", "^[234][0-9][0-9]$"),
					resource.TestCheckResourceAttr(monitorResourceName, "rtsp_body_regex", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "rtsp_path", "/"),
					resource.TestCheckResourceAttr(monitorResourceName, "script_program", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "sip_body_regex", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "sip_status_regex", "^[234][0-9][0-9]$"),
					resource.TestCheckResourceAttr(monitorResourceName, "sip_transport", "udp"),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp_close_string", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp_max_response_len", "4048"),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp_response_regex", ".*"),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp_write_string", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "udp_accept_all", "false"),
				),
			},
			{
				Config: testAccBrocadeVTMMonitorUpdateTemplate(monitorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName),
					resource.TestCheckResourceAttr(monitorResourceName, "name", monitorName),
					resource.TestCheckResourceAttr(monitorResourceName, "delay", "5"),
					resource.TestCheckResourceAttr(monitorResourceName, "timeout", "5"),
					resource.TestCheckResourceAttr(monitorResourceName, "failures", "9"),
					resource.TestCheckResourceAttr(monitorResourceName, "verbose", "false"),
					resource.TestCheckResourceAttr(monitorResourceName, "use_ssl", "false"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_host_header", "some_header"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_authentication", "some_authentication"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_body_regex", "^healthy"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_path", "/some/other/status/page"),
					resource.TestCheckResourceAttr(monitorResourceName, "rtsp_status_regex", "^[234][0-9][0-9]$"),
					resource.TestCheckResourceAttr(monitorResourceName, "rtsp_body_regex", "something"),
					resource.TestCheckResourceAttr(monitorResourceName, "rtsp_path", "/"),
					resource.TestCheckResourceAttr(monitorResourceName, "script_program", "dns.pl"),
					resource.TestCheckResourceAttr(monitorResourceName, "sip_body_regex", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "sip_status_regex", "^[234][0-9][0-9]$"),
					resource.TestCheckResourceAttr(monitorResourceName, "sip_transport", "udp"),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp_close_string", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp_max_response_len", "2048"),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp_response_regex", ".+"),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp_write_string", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "udp_accept_all", "false"),
				),
			},
		},
	})
}

func testAccBrocadeVTMMonitorCheckDestroy(state *terraform.State, name string) error {

	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_monitor" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}

		monitors, err := client.GetAllResources("monitors")

		if err != nil {
			return fmt.Errorf("Brocade vTM Monitor - error occurred whilst retrieving a list of all monitors: %+v", err)
		}
		for _, monitorChild := range monitors {
			if monitorChild["name"] == name {
				return fmt.Errorf("Brocade vTM monitor %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[monitorResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Monitor resource %s not found in resources", monitorResourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Monitor ID not set in resources")
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		monitors, err := client.GetAllResources("monitors")

		if err != nil {
			//return fmt.Errorf("Error: %+v", err)
			return fmt.Errorf("Error getting all monitors: %+v", err)
		}
		for _, monitorChild := range monitors {
			if monitorChild["name"] == monitorName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM Monitor %s not found on remote vTM", monitorName)
	}
}

func testAccBrocadeVTMMonitorCreateTemplate(monitorName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_monitor" "acctest" {
  name = "%s"
  delay = 6
  timeout = 2
  failures = 7
  verbose = true
  use_ssl = true
  http_host_header = "some_other_header"
  http_authentication = "admin:password"
  http_body_regex = "^ok"
  http_path = "/some/status/page"
  rtsp_status_regex = "^[234][0-9][0-9]$"
  rtsp_body_regex = ""
  rtsp_path = "/"
  script_arguments {
    name="test1"
    description="paas test"
    value="dns.pl"
  }
  script_program = ""
  sip_body_regex = ""
  sip_status_regex = "^[234][0-9][0-9]$"
  sip_transport = "udp"
  tcp_close_string = ""
  tcp_max_response_len = "4048"
  tcp_response_regex = ".*"
  tcp_write_string = ""
  udp_accept_all = false
}
`, monitorName)
}

func testAccBrocadeVTMMonitorUpdateTemplate(monitorName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_monitor" "acctest" {
  name = "%s"
  delay = 5
  timeout = 5
  failures = 9
  verbose = false
  use_ssl = false
  http_host_header = "some_header"
  http_authentication = "some_authentication"
  http_body_regex = "^healthy"
  http_path = "/some/other/status/page"
  rtsp_status_regex = "^[234][0-9][0-9]$"
  rtsp_body_regex = "something"
  rtsp_path = "/"
  script_arguments {
    name="test2"
    description="paas test2"
    value="bla.pl"
  }
  script_program = "dns.pl"
  sip_body_regex = ""
  sip_status_regex = "^[234][0-9][0-9]$"
  sip_transport = "udp"
  tcp_close_string = ""
  tcp_max_response_len = "2048"
  tcp_response_regex = ".+"
  tcp_write_string = ""
  udp_accept_all = false
}
`, monitorName)
}

func testAccBrocadeVTMMonitorInvalidName() string {
	return fmt.Sprintf(`
resource "brocadevtm_monitor" "acctest" {
  name = "../virtual_servers/some_random_virtual_server"
  delay = 5
}
`)
}
