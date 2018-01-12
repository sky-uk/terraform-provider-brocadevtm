package pulsevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"pulsevtm": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PULSEVTM_USERNAME"); v == "" {
		t.Fatal("PULSEVTM_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("PULSEVTM_PASSWORD"); v == "" {
		t.Fatal("PULSEVTM_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("PULSEVTM_SERVER"); v == "" {
		t.Fatal("PULSEVTM_SERVER must be set for acceptance tests")
	}
}
