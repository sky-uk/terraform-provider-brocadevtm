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

func TestAccPulseVTMRuleBasic(t *testing.T) {

	randomInt := acctest.RandInt()

	ruleName := fmt.Sprintf("acctest_pulsevtm_rule-%d", randomInt)
	ruleResourceName := "pulsevtm_rule.acctest"

	fmt.Printf("\n\nRule is %s.\n\n", ruleName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccPulseVTMRuleCheckDestroy(state, ruleName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccPulseVTMRuleNoName(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccPulseVTMRuleNoRule(ruleName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccPulseVTMRuleCreate(ruleName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMRuleExists(ruleName, ruleResourceName),
					resource.TestCheckResourceAttr(ruleResourceName, "name", ruleName),
					resource.TestCheckResourceAttr(ruleResourceName, "rule", "if( string.ipmaskmatch( request.getremoteip(), \"192.168.11.13\" ) ){\n    connection.discard();\n}\n"),
				),
			},
			{
				Config: testAccPulseVTMRuleUpdate(ruleName),
				Check: resource.ComposeTestCheckFunc(
					testAccPulseVTMRuleExists(ruleName, ruleResourceName),
					resource.TestCheckResourceAttr(ruleResourceName, "name", ruleName),
					resource.TestCheckResourceAttr(ruleResourceName, "rule", "if( string.ipmaskmatch( request.getremoteip(), \"10.78.12.34\" ) ){\n    connection.discard();\n}\n"),
				),
			},
		},
	})
}

func testAccPulseVTMRuleCheckDestroy(state *terraform.State, name string) error {
	config := testAccProvider.Meta().(map[string]interface{})

	// We need to copy the client as we want to specify different headers for rule which will conflict with other resources.
	vtmClient := config["octetClient"].(*api.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "pulsevtm_rule" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		allRules, err := vtmClient.GetAllResources("rules")
		if err != nil {
			return fmt.Errorf("[ERROR] Error: Pulse vTM error occurred while retrieving list of rules, %v", err)
		}
		for _, childRule := range allRules {
			if childRule["name"] == name {
				return fmt.Errorf("[ERROR] Error: Pulse vTM Rule %s still exists", name)
			}
		}
	}

	return nil
}

func testAccPulseVTMRuleExists(ruleName, ruleResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[ruleResourceName]
		if !ok {
			return fmt.Errorf("\n[ERROR] Pulse vTM Rule %s wasn't found in resources", ruleName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\n[ERROR] Pulse vTM Rule ID not set for %s in resources", ruleName)
		}
		config := testAccProvider.Meta().(map[string]interface{})
		vtmClient := config["octetClient"].(*api.Client)
		allRules, err := vtmClient.GetAllResources("rules")
		if err != nil {
			return fmt.Errorf("[ERROR] Pulse vTM Rule - error while retrieving a list of all rules: %v", err)
		}
		for _, childRule := range allRules {
			if childRule["name"] == ruleName {
				return nil
			}
		}
		return fmt.Errorf("[ERROR] Pulse vTM Rule %s not found on remote vTM", ruleName)
	}
}

func testAccPulseVTMRuleNoName() string {
	return fmt.Sprintf(`
resource "pulsevtm_rule" "acctest" {
rule = <<RULE
if( string.ipmaskmatch( request.getremoteip(), "192.168.11.13" ) ){
    connection.discard();
}
RULE
}
`)
}

func testAccPulseVTMRuleNoRule(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_rule" "acctest" {
name = "%s"
}
`, name)
}

func testAccPulseVTMRuleCreate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_rule" "acctest" {
name = "%s"
rule = <<RULE
if( string.ipmaskmatch( request.getremoteip(), "192.168.11.13" ) ){
    connection.discard();
}
RULE
}
`, name)
}

func testAccPulseVTMRuleUpdate(name string) string {
	return fmt.Sprintf(`
resource "pulsevtm_rule" "acctest" {
name = "%s"
rule = <<RULE
if( string.ipmaskmatch( request.getremoteip(), "10.78.12.34" ) ){
    connection.discard();
}
RULE
}
`, name)
}
