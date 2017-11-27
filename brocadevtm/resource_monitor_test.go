package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"os"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMMonitorBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	monitorName := fmt.Sprintf("acctest_brocadevtm_monitor-%d", randomInt)
	monitorResourceName := "brocadevtm_monitor.acctest"
	fmt.Printf("\n\nMonitor Name is %s.\n\n", monitorName)
	var usedVersion = "3.8"
	if os.Getenv("BROCADEVTM_API_VERSION") != "" {
		usedVersion = os.Getenv("BROCADEVTM_API_VERSION")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMMonitorCheckDestroy(state, monitorName)
		},
		Steps: []resource.TestStep{
			{ // Step 0
				Config:      testAccBrocadeVTMMonitorInvalidName(),
				ExpectError: regexp.MustCompile(`BrocadeVTM Monitor error whilst creating ../virtual_servers/some_random_virtual_server: The path '/api/tm/` + usedVersion + `/config/active/monitors/../virtual_servers/some_random_virtual_server' is invalid`),
			},
			{ // Step 1
				Config: testAccBrocadeVTMMonitorCreateTemplate(monitorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName),
					resource.TestCheckResourceAttr(monitorResourceName, "name", monitorName),
					resource.TestCheckResourceAttr(monitorResourceName, "delay", "6"),
					resource.TestCheckResourceAttr(monitorResourceName, "timeout", "2"),
					resource.TestCheckResourceAttr(monitorResourceName, "failures", "7"),
					resource.TestCheckResourceAttr(monitorResourceName, "verbose", "true"),
					resource.TestCheckResourceAttr(monitorResourceName, "use_ssl", "true"),
					resource.TestCheckResourceAttr(monitorResourceName, "http.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("http", "host_header"), "some_other_header"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("http", "authentication"), "admin:password"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("http", "body_regex"), "^ok"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("http", "path"), "/some/status/page"),
					resource.TestCheckResourceAttr(monitorResourceName, "rtsp.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("rtsp", "status_regex"), "^[234][0-9][0-9]$"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("rtsp", "body_regex"), ""),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("rtsp", "path"), "/"),
					resource.TestCheckResourceAttr(monitorResourceName, "script_program", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "sip.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("sip", "body_regex"), ""),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("sip", "status_regex"), "^[234][0-9][0-9]$"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("sip", "transport"), "udp"),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "close_string"), ""),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "max_response_len"), "4048"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "response_regex"), ".*"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "write_string"), ""),
					resource.TestCheckResourceAttr(monitorResourceName, "udp.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "accept_all"), "false"),
				),
			},
			{ // Step 2
				Config: testAccBrocadeVTMMonitorUpdateTemplate(monitorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName),
					resource.TestCheckResourceAttr(monitorResourceName, "name", monitorName),
					resource.TestCheckResourceAttr(monitorResourceName, "delay", "5"),
					resource.TestCheckResourceAttr(monitorResourceName, "timeout", "5"),
					resource.TestCheckResourceAttr(monitorResourceName, "failures", "9"),
					resource.TestCheckResourceAttr(monitorResourceName, "verbose", "false"),
					resource.TestCheckResourceAttr(monitorResourceName, "use_ssl", "false"),
					resource.TestCheckResourceAttr(monitorResourceName, "http.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("http", "host_header"), "some_header"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("http", "authentication"), "some_authentication"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("http", "body_regex"), "^healthy"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("http", "path"), "/some/other/status/page"),
					resource.TestCheckResourceAttr(monitorResourceName, "rtsp.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("rtsp", "status_regex"), "^[234][0-9][0-9]$"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("rtsp", "body_regex"), "something"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("rtsp", "path"), "/"),
					resource.TestCheckResourceAttr(monitorResourceName, "script_program", "dns.pl"),
					resource.TestCheckResourceAttr(monitorResourceName, "sip.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("sip", "body_regex"), ""),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("sip", "status_regex"), "^[234][0-9][0-9]$"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("sip", "transport"), "udp"),
					resource.TestCheckResourceAttr(monitorResourceName, "tcp.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "close_string"), ""),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "max_response_len"), "2048"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "response_regex"), ".+"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("tcp", "write_string"), ""),
					resource.TestCheckResourceAttr(monitorResourceName, "udp.#", "1"),
					util.AccTestCheckValueInKeyPattern(monitorResourceName, util.AccTestCreateRegexPatternForSetItems("udp", "accept_all"), "false"),
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
			return fmt.Errorf("[ERROR] Brocade vTM Monitor - error occurred whilst retrieving a list of all monitors: %+v", err)
		}
		for _, monitorChild := range monitors {
			if monitorChild["name"] == name {
				return fmt.Errorf("[ERROR] Brocade vTM monitor %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[monitorResourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Brocade vTM Monitor resource %s not found in resources", monitorResourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Brocade vTM Monitor ID not set in resources")
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		monitors, err := client.GetAllResources("monitors")

		if err != nil {
			return fmt.Errorf("[ERROR] getting all monitors: %+v", err)
		}
		for _, monitorChild := range monitors {
			if monitorChild["name"] == monitorName {
				return nil
			}
		}
		return fmt.Errorf("[ERROR] Brocade vTM Monitor %s not found on remote vTM", monitorName)
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
  http = [
    {
      host_header = "some_other_header"
      authentication = "admin:password"
      body_regex = "^ok"
      path = "/some/status/page"
    },
  ]
  rtsp = [
    {
      status_regex = "^[234][0-9][0-9]$"
      body_regex = ""
      path = "/"
    },
  ]
  script_arguments {
    name="test1"
    description="paas test"
    value="dns.pl"
  }
  script_program = ""
  sip = [
    {
      body_regex = ""
      status_regex = "^[234][0-9][0-9]$"
      transport = "udp"
    },
  ]
  tcp = [
    {
      close_string = ""
      max_response_len = "4048"
      response_regex = ".*"
      write_string = ""
    },
  ]
  udp = [
    {
      accept_all = false
    },
  ]
}
`, monitorName)
}

func testAccBrocadeVTMMonitorUpdateTemplate(monitorName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_monitor" "acctest" {
  name = "%s"
  back_off = false
  delay = 5
  failures = 9
  machine = "10.93.59.24:9090"
  note = "a description of this monitor..."
  scope = "poolwide"
  timeout = 5
  type = "tcp_transaction"
  verbose = false
  use_ssl = false
  http = [
    {
      host_header = "some_header"
      authentication = "some_authentication"
      body_regex = "^healthy"
      path = "/some/other/status/page"
    },
  ]
  rtsp = [
    {
      status_regex = "^[234][0-9][0-9]$"
      body_regex = "something"
      path = "/"
    },
  ]
  script_arguments {
    name="test2"
    description="paas test2"
    value="bla.pl"
  }
  script_program = "dns.pl"
  sip = [
    {
      body_regex = ""
      status_regex = "^[234][0-9][0-9]$"
      transport = "udp"
    },
  ]
  tcp = [
    {
      close_string = ""
      max_response_len = "2048"
      response_regex = ".+"
      write_string = ""
    },
  ]
  udp = [
    {
      accept_all = false
    },
  ]
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
