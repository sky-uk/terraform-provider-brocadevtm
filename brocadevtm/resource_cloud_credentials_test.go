package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMCloudCredentialsBasic(t *testing.T) {

	cloudCredentialsName := acctest.RandomWithPrefix("acctest_brocadevtm_cloud_credentials")
	resourceName := "brocadevtm_cloud_credentials.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMCloudCredentialsCheckDestroy(state, cloudCredentialsName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMCloudCredentialsNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccBrocadeCloudCredentialCreateTemplate(cloudCredentialsName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMCloudCredentialsExists(cloudCredentialsName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", cloudCredentialsName),
					resource.TestCheckResourceAttr(resourceName, "api_server", "testServer"),
					resource.TestCheckResourceAttr(resourceName, "cloud_api_timeout", "50"),
					resource.TestCheckResourceAttr(resourceName, "cred1", "testCred1"),
					resource.TestCheckResourceAttr(resourceName, "cred2", "testCred2"),
					resource.TestCheckResourceAttr(resourceName, "cred3", "testCred3"),
					resource.TestCheckResourceAttr(resourceName, "script", "testscript"),
					resource.TestCheckResourceAttr(resourceName, "update_interval", "50"),
				),
			},
			{
				Config: testAccBrocadeCloudCredentialUpdateTemplate(cloudCredentialsName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMCloudCredentialsExists(cloudCredentialsName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", cloudCredentialsName),
					resource.TestCheckResourceAttr(resourceName, "api_server", "testServerUpdated"),
					resource.TestCheckResourceAttr(resourceName, "cloud_api_timeout", "100"),
					resource.TestCheckResourceAttr(resourceName, "cred1", "testCred1Updated"),
					resource.TestCheckResourceAttr(resourceName, "cred2", "testCred2Updated"),
					resource.TestCheckResourceAttr(resourceName, "cred3", "testCred3Updated"),
					resource.TestCheckResourceAttr(resourceName, "script", "testscript2"),
					resource.TestCheckResourceAttr(resourceName, "update_interval", "100"),
				),
			},
		},
	})
}

func testAccBrocadeVTMCloudCredentialsCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_cloud_credentials" {
			continue
		}
		cloudCredentialsConfiguration := make(map[string]interface{})

		err := client.GetByName("cloud_api_credentials", rs.Primary.ID, &cloudCredentialsConfiguration)
		if client.StatusCode == http.StatusOK {
			return fmt.Errorf("[ERROR] Brocade vTM Check Destroy Error: Cloud Credential %s still exists", name)
		}
		if client.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("[ERROR] Brocade vTM Check Destroy Error: Cloud Credential %+v ", err)
	}
	return nil
}

func testAccBrocadeVTMCloudCredentialsExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("[ERROR] No ID is set")
		}

		config := testAccProvider.Meta().(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		client.WorkWithConfigurationResources()
		cloudCredentialsConfiguration := make(map[string]interface{})
		err := client.GetByName("cloud_api_credentials", name, &cloudCredentialsConfiguration)
		if client.StatusCode != http.StatusOK {
			return fmt.Errorf("[ERROR] Brocade vTM error whilst retrieving VTM Cloud Credentials: %+v", err)
		}
		return nil
	}
}

func testAccBrocadeVTMCloudCredentialsNoNameTemplate() string {
	return `
resource "brocadevtm_cloud_credentials" "acctest" {
  api_server = "testServer"
  cloud_api_timeout = 50
  cred1 = "testCred1"
  cred2 = "testCred2"
  cred3 = "testCred3"
  script = "fakeScript"
  update_interval = 50
}
`
}

func testAccBrocadeCloudCredentialCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_cloud_credentials" "acctest" {
  name = "%s"
  api_server = "testServer"
  cloud_api_timeout = 50
  cred1 = "testCred1"
  cred2 = "testCred2"
  cred3 = "testCred3"
  script = "testscript"
  update_interval = 50
}
`, name)
}

func testAccBrocadeCloudCredentialUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_cloud_credentials" "acctest" {
  name = "%s"
  api_server = "testServerUpdated"
  cloud_api_timeout = 100
  cred1 = "testCred1Updated"
  cred2 = "testCred2Updated"
  cred3 = "testCred3Updated"
  script = "testscript2"
  update_interval = 100
}
`, name)
}
