package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"testing"
	"sort"
)

func TestAccBrocadeVTMResourceGlobalSettings(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMGlobalSettingsCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config:                    testAccBrocadeGlobalSettingsCreate(),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMGlobalSettingsExists(),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.monitor_memory_size", "4096"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.so_rbuff_size", "0"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.client_first_opt", "false"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.cluster_identifier", ""),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.accepting_delay", "50"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.afm_enabled", "false"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.chunk_size", "16384"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.tip_class_limit", "10000"),
				),
			},
			{
				Config:                    testAccBrocadeGlobalSettingsUpdate(),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMGlobalSettingsExists(),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.monitor_memory_size", "4096"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.so_rbuff_size", "0"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.client_first_opt", "false"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.cluster_identifier", ""),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.accepting_delay", "100"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.afm_enabled", "false"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.chunk_size", "16384"),
					resource.TestCheckResourceAttr("brocadevtm_global_settings.global_settings", "basic.0.tip_class_limit", "10000"),
				),
			},
		},
	})
}

func testAccBrocadeVTMGlobalSettingsCheckDestroy(state *terraform.State) error {

	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_global_settings" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		gs := make(map[string]interface{})

		apiVersions := client.VersionsSupported
		sort.Sort(sort.Reverse(sort.StringSlice(apiVersions)))
		latestVersion := apiVersions[0]
		err := client.GetByURL("/api/tm/"+latestVersion+"/config/active/global_settings", &gs)
		if err != nil {
			return nil
		}
	}
	return fmt.Errorf("Brocade vTM, global settings still found")
}

func testAccBrocadeVTMGlobalSettingsExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources["brocadevtm_global_settings.global_settings"]
		if !ok {
			return fmt.Errorf("\nBrocade vTM global settings missing")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM ID not set")
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		apiVersions := client.VersionsSupported
		sort.Sort(sort.Reverse(sort.StringSlice(apiVersions)))
		latestVersion := apiVersions[0]
		gs := make(map[string]interface{})
		err := client.GetByURL("/api/tm/"+latestVersion+"/config/active/global_settings", &gs)
		if err != nil {
			return fmt.Errorf("Error getting global settings: %+v", err)
		}
		return nil
	}
}

func testAccBrocadeGlobalSettingsCreate() string {
	return `resource "brocadevtm_global_settings" "global_settings" {
   basic = {
    monitor_memory_size = 4096
    so_rbuff_size = 0
    client_first_opt = false
    cluster_identifier = ""
    accepting_delay = 50
    afm_enabled = false
    chunk_size = 16384
    tip_class_limit = 10000
   }
}`
}

func testAccBrocadeGlobalSettingsUpdate() string {
	return `resource "brocadevtm_global_settings" "global_settings" {
   basic = {
    monitor_memory_size = 4096
    so_rbuff_size = 0
    client_first_opt = false
    cluster_identifier = ""
    accepting_delay = 100
    afm_enabled = false
    chunk_size = 16384
    tip_class_limit = 10000
   }
}`
}
