package pulsevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-pulse-vtm/api"
	"regexp"
	"testing"
)

func TestAccPulseVTMLocationBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	locationName := fmt.Sprintf("acctest_pulsevtm_location-%d", randomInt)
	locationResourceName := "pulsevtm_location.acctest"

	fmt.Printf("\nLocation Name is %s.\n\n", locationName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccPulseVTMLocationCheckDestroy(state, locationName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccPulseVTMLocationNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPulseVTMLocationNoLocationIDTemplate(locationName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPulseVTMLocationInvalidLatitudeTemplate(locationName),
				ExpectError: regexp.MustCompile(`must be between -90 and 90 degrees inclusive`),
			},
			{
				Config:      testAccPulseVTMLocationInvalidLatitude2Template(locationName),
				ExpectError: regexp.MustCompile(`must be between -90 and 90 degrees inclusive`),
			},
			{
				Config:      testAccPulseVTMLocationInvalidLongitudeTemplate(locationName),
				ExpectError: regexp.MustCompile(`must be between -180 and 180 degrees inclusive`),
			},
			{
				Config:      testAccPulseVTMLocationInvalidLongitude2Template(locationName),
				ExpectError: regexp.MustCompile(`must be between -180 and 180 degrees inclusive`),
			},
			{
				Config:      testAccPulseVTMLocationInvalidTypeTemplate(locationName),
				ExpectError: regexp.MustCompile(`must be one of config or glb`),
			},
			{
				Config: testAccPulseVTMLocationCreateTemplate(locationName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMLocationExists(locationName, locationResourceName),
					resource.TestCheckResourceAttr(locationResourceName, "name", locationName),
					resource.TestCheckResourceAttr(locationResourceName, "location_id", "32001"),
					resource.TestCheckResourceAttr(locationResourceName, "latitude", "-36.353417"),
					resource.TestCheckResourceAttr(locationResourceName, "longitude", "146.687568"),
					resource.TestCheckResourceAttr(locationResourceName, "note", "Acceptance test location"),
					resource.TestCheckResourceAttr(locationResourceName, "type", "config"),
				),
			},
			{
				Config: testAccPulseVTMLocationUpdateTemplate(locationName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMLocationExists(locationName, locationResourceName),
					resource.TestCheckResourceAttr(locationResourceName, "name", locationName),
					resource.TestCheckResourceAttr(locationResourceName, "location_id", "32102"),
					resource.TestCheckResourceAttr(locationResourceName, "latitude", "51.503607"),
					resource.TestCheckResourceAttr(locationResourceName, "longitude", "-0.307904"),
					resource.TestCheckResourceAttr(locationResourceName, "note", "Acceptance test location - updated"),
					resource.TestCheckResourceAttr(locationResourceName, "type", "glb"),
				),
			},
		},
	})
}

func testAccPulseVTMLocationCheckDestroy(state *terraform.State, name string) error {

	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "pulsevtm_location" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}

		locations, err := client.GetAllResources("locations")

		if err != nil {
			return fmt.Errorf("[ERROR] Pulse vTM Location - error occurred whilst retrieving a list of all locations: %+v", err)
		}
		for _, locationChild := range locations {
			if locationChild["name"] == name {
				return fmt.Errorf("[ERROR] Pulse vTM Location %s still exists", name)
			}
		}
	}
	return nil
}

func testAccPulseVTMLocationExists(locationName, locationResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[locationResourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Pulse vTM Location %s wasn't found in resources", locationName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Pulse vTM Location ID not set for %s in resources", locationName)
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		locations, err := client.GetAllResources("locations")

		if err != nil {
			return fmt.Errorf("[ERROR] getting all locations: %+v", err)
		}

		for _, locationChild := range locations {
			if locationChild["name"] == locationName {
				return nil
			}
		}
		return fmt.Errorf("[ERROR] Pulse vTM Location %s not found on remote vTM", locationName)
	}
}

func testAccPulseVTMLocationNoNameTemplate() string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  location_id = 32001
  latitude = -36.353417
  longitude = 146.687568
  note = "Acceptance test location"
  type = "glb"
}
`)
}

func testAccPulseVTMLocationNoLocationIDTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  name = "%s"
  latitude = -36.353417
  longitude = 146.687568
  note = "Acceptance test location"
  type = "glb"
}
`, locationName)
}

func testAccPulseVTMLocationInvalidTypeTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -36.353417
  longitude = 146.687568
  note = "Acceptance test location"
  type = "SOME_INVALID_TYPE"
}
`, locationName)
}

func testAccPulseVTMLocationInvalidLatitudeTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = 180.562456
  longitude = 146.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccPulseVTMLocationInvalidLatitude2Template(locationName string) string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -180.562456
  longitude = 146.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccPulseVTMLocationInvalidLongitudeTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -36.353417
  longitude = 196.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccPulseVTMLocationInvalidLongitude2Template(locationName string) string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -36.353417
  longitude = -196.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccPulseVTMLocationCreateTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -36.353417
  longitude = 146.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccPulseVTMLocationUpdateTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "pulsevtm_location" "acctest" {
  name = "%s"
  location_id = 32102
  latitude = 51.503607
  longitude = -0.307904
  note = "Acceptance test location - updated"
  type = "glb"
}
`, locationName)
}
