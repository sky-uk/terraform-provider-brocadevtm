package pulsevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-pulse-vtm/api"
	"net/http"
	"testing"
)

func TestAccPulseVTMBasicSSLTicketKey(t *testing.T) {

	sslTicketKeyName := acctest.RandomWithPrefix("acctest_pulsevtm_ssl_ticket_key")
	resourceName := "pulsevtm_ssl_ticket_key.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccPulseVTMSSLTicketKeyCheckDestroy(state, sslTicketKeyName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPulseSSLTicketKeyCreateTemplate(sslTicketKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMSSLTicketKeyExists(sslTicketKeyName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", sslTicketKeyName),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "aes_256_cbc_hmac_sha256"),
					resource.TestCheckResourceAttr(resourceName, "identifier", "c65361fa7e53e20ae704844795446d4c"),
					resource.TestCheckResourceAttr(resourceName, "key", "197dc85083e0f4390cd4f4a6ee3b866d"),
					resource.TestCheckResourceAttr(resourceName, "validity_end", "1514853000"),
					resource.TestCheckResourceAttr(resourceName, "validity_start", "1515767356"),
				),
			},
			{
				Config: testAccPulseSSLTicketKeyUpdateTemplate(sslTicketKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMSSLTicketKeyExists(sslTicketKeyName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", sslTicketKeyName),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "aes_256_cbc_hmac_sha256"),
					resource.TestCheckResourceAttr(resourceName, "identifier", "c65361fa7e53e20ae704844795446d4d"),
					resource.TestCheckResourceAttr(resourceName, "key", "197dc85083e0f4390cd4f4a6ee3b866e"),
					resource.TestCheckResourceAttr(resourceName, "validity_end", "1514856099"),
					resource.TestCheckResourceAttr(resourceName, "validity_start", "1515787399"),
				),
			},
		},
	})
}

func testAccPulseVTMSSLTicketKeyCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "pulsevtm_ssl_ticket_key" {
			continue
		}
		sslTicketKeyConfiguration := make(map[string]interface{})

		err := client.GetByName("ssl/ticket_keys", rs.Primary.ID, &sslTicketKeyConfiguration)
		if client.StatusCode == http.StatusOK {
			return fmt.Errorf("[ERROR] Pulse vTM Check Destroy Error: SSL Ticket Key %s still exists", name)
		}
		if client.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("[ERROR] Pulse vTM Check Destroy Error: SSL Ticket Key %+v ", err)
	}
	return nil
}

func testAccPulseVTMSSLTicketKeyExists(name, resourceName string) resource.TestCheckFunc {
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
		sslTicketKeyConfiguration := make(map[string]interface{})
		err := client.GetByName("ssl/ticket_keys", name, &sslTicketKeyConfiguration)
		if client.StatusCode != http.StatusOK {
			return fmt.Errorf("[ERROR] Pulse vTM error whilst retrieving VTM SSL Ticket Key: %+v", err)
		}
		return nil
	}
}

func testAccPulseSSLTicketKeyCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_ssl_ticket_key" "acctest" {
  name = "%s"
  algorithm = "aes_256_cbc_hmac_sha256"
  identifier = "c65361fa7e53e20ae704844795446d4c"
  key = "197dc85083e0f4390cd4f4a6ee3b866d"
  validity_end = 1514853000
  validity_start = 1515767356
}
`, name)
}

func testAccPulseSSLTicketKeyUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_ssl_ticket_key" "acctest" {
  name = "%s"
  algorithm = "aes_256_cbc_hmac_sha256"
  identifier = "c65361fa7e53e20ae704844795446d4d"
  key = "197dc85083e0f4390cd4f4a6ee3b866e"
  validity_end = 1514856099
  validity_start = 1515787399
}
`, name)
}
