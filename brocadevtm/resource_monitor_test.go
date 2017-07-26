package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-rest-api"
	"github.com/sky-uk/go-brocade-vtm/api/monitor"
	"regexp"
	"testing"
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
				ExpectError: regexp.MustCompile(`Response object was \{\"error_id\":\"http.invalid_path\",\"error_text\":\"The path \'/api/tm/3.8/config/active/monitors/../virtual_servers/some_random_virtual_server\' is invalid`),
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
				),
			},
		},
	})
}

func testAccBrocadeVTMMonitorCheckDestroy(state *terraform.State, name string) error {

	vtmClient := testAccProvider.Meta().(*rest.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_monitor" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}

		api := monitor.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return nil
		}
		for _, monitorChild := range api.ResponseObject().(*monitor.MonitorsList).Children {
			if monitorChild.Name == name {
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

		vtmClient := testAccProvider.Meta().(*rest.Client)
		getAllAPI := monitor.NewGetAll()

		err := vtmClient.Do(getAllAPI)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		for _, monitorChild := range getAllAPI.ResponseObject().(*monitor.MonitorsList).Children {
			if monitorChild.Name == monitorName {
				return fmt.Errorf("Brocade vTM Monitor %s not found on remote vTM", monitorName)
			}
		}
		return nil
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
