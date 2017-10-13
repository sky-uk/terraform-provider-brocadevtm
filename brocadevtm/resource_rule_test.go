package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"log"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMRuleBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	ruleName := fmt.Sprintf("acctest_brocadevtm_rule-%d", randomInt)
	ruleResourceName := "brocadevtm_rule.acctest"

	fmt.Printf("\n\nRule is %s.\n\n", ruleName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMRuleCheckDestroy(state, ruleName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMRuleNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccBrocadeVTMRuleNoRule(ruleName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccBrocadeVTMRuleCreate(ruleName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMRuleExists(ruleName, ruleResourceName),
					resource.TestCheckResourceAttr(ruleResourceName, "name", ruleName),
					resource.TestCheckResourceAttr(ruleResourceName, "rule", "if( string.ipmaskmatch( request.getremoteip(), \"192.168.11.13\" ) ){\n    connection.discard();\n}\n"),
				),
			},
			{
				Config: testAccBrocadeVTMRuleUpdate(ruleName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMRuleExists(ruleName, ruleResourceName),
					resource.TestCheckResourceAttr(ruleResourceName, "name", ruleName),
					resource.TestCheckResourceAttr(ruleResourceName, "rule", "if( string.ipmaskmatch( request.getremoteip(), \"10.78.12.34\" ) ){\n    connection.discard();\n}\n"),
				),
			},
		},
	})
}

func testAccBrocadeVTMRuleCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})

	// We need to copy the client as we want to specify different headers for rule which will conflict with other resources.
	vtmClient := config["octetClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_rule" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		allRules, err := vtmClient.GetAllResources("rules")
		if err != nil {
			return fmt.Errorf("Error: Brocade vTM error occurred while retrieving list of rules, %v", err)
		}
		for _, childRule := range allRules {
			if childRule["name"] == name {
				return fmt.Errorf("Error: Brocade vTM Rule %s still exists", name)
			}
		}
	}

	return nil
}

func testAccBrocadeVTMRuleExists(ruleName, ruleResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[ruleResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Rule %s wasn't found in resources", ruleName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Rule ID not set for %s in resources", ruleName)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		vtmClient := config["octetClient"].(*api.Client)
		allRules, err := vtmClient.GetAllResources("rules")
		if err != nil {
			return fmt.Errorf("Brocade vTM Rule - error while retrieving a list of all rules: %v", err)
		}
		for _, childRule := range allRules {
			if childRule["name"] == ruleName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM Rule %s not found on remote vTM", ruleName)
	}
}

func testAccBrocadeVTMRuleNoName() string {
	return fmt.Sprintf(`
resource "brocadevtm_rule" "acctest" {
rule = <<RULE
if( string.ipmaskmatch( request.getremoteip(), "192.168.11.13" ) ){
    connection.discard();
}
RULE
}
`)
}

func testAccBrocadeVTMRuleNoRule(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_rule" "acctest" {
name = "%s"
}
`, name)
}

func testAccBrocadeVTMRuleCreate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_rule" "acctest" {
name = "%s"
rule = <<RULE
if( string.ipmaskmatch( request.getremoteip(), "192.168.11.13" ) ){
    connection.discard();
}
RULE
}
`, name)
}

func testAccBrocadeVTMRuleUpdate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_rule" "acctest" {
name = "%s"
rule = <<RULE
if( string.ipmaskmatch( request.getremoteip(), "10.78.12.34" ) ){
    connection.discard();
}
RULE
}
`, name)
}
