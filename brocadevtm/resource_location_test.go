package brocadevtm

/*
import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api/location"
	"github.com/sky-uk/go-rest-api"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMLocationBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	locationName := fmt.Sprintf("acctest_brocadevtm_location-%d", randomInt)
	locationResourceName := "brocadevtm_location.acctest"

	fmt.Printf("\nLocation Name is %s.\n\n", locationName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMLocationCheckDestroy(state, locationName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMLocationNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMLocationNoLocationIDTemplate(locationName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMLocationInvalidLatitudeTemplate(locationName),
				ExpectError: regexp.MustCompile(`must be between -90 and 90 degrees inclusive`),
			},
			{
				Config:      testAccBrocadeVTMLocationInvalidLatitude2Template(locationName),
				ExpectError: regexp.MustCompile(`must be between -90 and 90 degrees inclusive`),
			},
			{
				Config:      testAccBrocadeVTMLocationInvalidLongitudeTemplate(locationName),
				ExpectError: regexp.MustCompile(`must be between -180 and 180 degrees inclusive`),
			},
			{
				Config:      testAccBrocadeVTMLocationInvalidLongitude2Template(locationName),
				ExpectError: regexp.MustCompile(`must be between -180 and 180 degrees inclusive`),
			},
			{
				Config:      testAccBrocadeVTMLocationInvalidTypeTemplate(locationName),
				ExpectError: regexp.MustCompile(`must be one of config or glb`),
			},
			{
				Config: testAccBrocadeVTMLocationCreateTemplate(locationName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMLocationExists(locationName, locationResourceName),
					resource.TestCheckResourceAttr(locationResourceName, "name", locationName),
					resource.TestCheckResourceAttr(locationResourceName, "location_id", "32001"),
					resource.TestCheckResourceAttr(locationResourceName, "latitude", "-36.353417"),
					resource.TestCheckResourceAttr(locationResourceName, "longitude", "146.687568"),
					resource.TestCheckResourceAttr(locationResourceName, "note", "Acceptance test location"),
					resource.TestCheckResourceAttr(locationResourceName, "type", "config"),
				),
			},
			{
				Config: testAccBrocadeVTMLocationUpdateTemplate(locationName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMLocationExists(locationName, locationResourceName),
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

func testAccBrocadeVTMLocationCheckDestroy(state *terraform.State, name string) error {

	vtmClient := testAccProvider.Meta().(*rest.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_location" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		api := location.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Brocade vTM Location - error occurred whilst retrieving a list of all locations")
		}
		for _, location := range api.ResponseObject().(*location.Locations).Children {
			if location.Name == name {
				return fmt.Errorf("Brocade vTM Location %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMLocationExists(locationName, locationResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[locationResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Location %s wasn't found in resources", locationName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Location ID not set for %s in resources", locationName)
		}

		vtmClient := testAccProvider.Meta().(*rest.Client)
		api := location.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		for _, location := range api.ResponseObject().(*location.Locations).Children {
			if location.Name == locationName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM Location %s not found on remote vTM", locationName)
	}
}

func testAccBrocadeVTMLocationNoNameTemplate() string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  location_id = 32001
  latitude = -36.353417
  longitude = 146.687568
  note = "Acceptance test location"
  type = "glb"
}
`)
}

func testAccBrocadeVTMLocationNoLocationIDTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  name = "%s"
  latitude = -36.353417
  longitude = 146.687568
  note = "Acceptance test location"
  type = "glb"
}
`, locationName)
}

func testAccBrocadeVTMLocationInvalidTypeTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -36.353417
  longitude = 146.687568
  note = "Acceptance test location"
  type = "SOME_INVALID_TYPE"
}
`, locationName)
}

func testAccBrocadeVTMLocationInvalidLatitudeTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = 180.562456
  longitude = 146.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccBrocadeVTMLocationInvalidLatitude2Template(locationName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -180.562456
  longitude = 146.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccBrocadeVTMLocationInvalidLongitudeTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -36.353417
  longitude = 196.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccBrocadeVTMLocationInvalidLongitude2Template(locationName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -36.353417
  longitude = -196.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccBrocadeVTMLocationCreateTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  name = "%s"
  location_id = 32001
  latitude = -36.353417
  longitude = 146.687568
  note = "Acceptance test location"
  type = "config"
}
`, locationName)
}

func testAccBrocadeVTMLocationUpdateTemplate(locationName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_location" "acctest" {
  name = "%s"
  location_id = 32102
  latitude = 51.503607
  longitude = -0.307904
  note = "Acceptance test location - updated"
  type = "glb"
}
`, locationName)
}
*/
