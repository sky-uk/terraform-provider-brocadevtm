package pulsevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-pulse-vtm/api"
	"os"
	"testing"
)

func TestAccPulseVTMResourceGlobalSettings(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccPulseVTMGlobalSettingsCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config:                    testAccPulseGlobalSettingsCreate(),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMGlobalSettingsExists(),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.monitor_memory_size", "4096"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.so_rbuff_size", "0"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.client_first_opt", "false"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.cluster_identifier", ""),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.accepting_delay", "50"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.afm_enabled", "false"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.chunk_size", "16384"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.tip_class_limit", "10000"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.data_plane_acceleration_cores", "two"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.data_plane_acceleration_mode", "true"),
				),
			},
			{
				Config:                    testAccPulseGlobalSettingsUpdate(),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMGlobalSettingsExists(),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.monitor_memory_size", "4096"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.so_rbuff_size", "0"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.client_first_opt", "false"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.cluster_identifier", ""),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.accepting_delay", "100"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.afm_enabled", "false"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.chunk_size", "16384"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.tip_class_limit", "10000"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.data_plane_acceleration_cores", "four"),
					resource.TestCheckResourceAttr("pulsevtm_global_settings.global_settings", "basic.0.data_plane_acceleration_mode", "false"),
				),
			},
		},
	})
}

func testAccPulseVTMGlobalSettingsCheckDestroy(state *terraform.State) error {

	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "pulsevtm_global_settings" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		gs := make(map[string]interface{})

		var usedVersion = "5.1"
		if os.Getenv("PULSEVTM_API_VERSION") != "" {
			usedVersion = os.Getenv("PULSEVTM_API_VERSION")
		}
		err := client.GetByURL("/api/tm/"+usedVersion+"/config/active/global_settings", &gs)
		if err != nil {
			return nil
		}
	}
	return fmt.Errorf("[ERROR] Pulse vTM, global settings still found")
}

func testAccPulseVTMGlobalSettingsExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources["pulsevtm_global_settings.global_settings"]
		if !ok {
			return fmt.Errorf("\n[ERROR] Pulse vTM global settings missing")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Pulse vTM ID not set")
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		var usedVersion = "5.1"
		if os.Getenv("PULSEVTM_API_VERSION") != "" {
			usedVersion = os.Getenv("PULSEVTM_API_VERSION")
		}
		gs := make(map[string]interface{})
		err := client.GetByURL("/api/tm/"+usedVersion+"/config/active/global_settings", &gs)
		if err != nil {
			return fmt.Errorf("[ERROR] getting global settings: %+v", err)
		}
		return nil
	}
}

func testAccPulseGlobalSettingsCreate() string {
	return `resource "pulsevtm_global_settings" "global_settings" {
   basic = {
    monitor_memory_size = 4096
    so_rbuff_size = 0
    client_first_opt = false
    cluster_identifier = ""
    accepting_delay = 50
    afm_enabled = false
    chunk_size = 16384
    tip_class_limit = 10000
    data_plane_acceleration_cores = "two"
    data_plane_acceleration_mode = true
   }
}`
}

func testAccPulseGlobalSettingsUpdate() string {
	return `resource "pulsevtm_global_settings" "global_settings" {
   basic = {
    monitor_memory_size = 4096
    so_rbuff_size = 0
    client_first_opt = false
    cluster_identifier = ""
    accepting_delay = 100
    afm_enabled = false
    chunk_size = 16384
    tip_class_limit = 10000
    data_plane_acceleration_cores = "four"
    data_plane_acceleration_mode = false
   }
}`
}
