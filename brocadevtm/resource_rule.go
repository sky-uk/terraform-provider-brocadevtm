package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/rule"
	"log"
	"net/http"
	"strings"
)

func resourceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleCreate,
		Read:   resourceRuleRead,
		Update: resourceRuleUpdate,
		Delete: resourceRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the rule",
				Required:    true,
			},
			"rule": {
				Type:        schema.TypeList,
				Description: "The traffic script",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildRule(ruleLines interface{}) string {
	var trafficRule string
	for idx, ruleLine := range ruleLines.([]interface{}) {
		if idx > 0 {
			trafficRule = trafficRule + "\n" + ruleLine.(string)
		} else {
			trafficRule = ruleLine.(string)
		}
	}
	return trafficRule
}

func buildRuleArray(trafficRule string) []string {
	var ruleLines []string
	trafficRuleLines := strings.Split(trafficRule, "\n")
	for _, trafficRuleLine := range trafficRuleLines {
		log.Printf(fmt.Sprintf("[DEBUG] Rule line is %s", trafficRuleLine))
		ruleLines = append(ruleLines, trafficRuleLine)
	}
	return ruleLines
}

func resourceRuleCreate(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	var scriptAsBytes []byte

	vtmClient := m.(*brocadevtm.VTMClient)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Transfer-Encoding"] = "text"
	vtmClient.Headers = headers

	if v, ok := d.GetOk("name"); ok && v != "" {
		vtmRule.Name = v.(string)
	}
	if v, ok := d.GetOk("rule"); ok {
		vtmRule.Script = buildRule(v)
		scriptAsBytes = []byte(fmt.Sprintf(vtmRule.Script))
	}

	createAPI := rule.NewCreate(vtmRule.Name, scriptAsBytes)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Error while creating rule %s: %+v", vtmRule.Name, err)
	}

	d.SetId(vtmRule.Name)
	return resourceRuleRead(d, m)
}

func resourceRuleRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Transfer-Encoding"] = "text"
	vtmClient.Headers = headers

	ruleName := d.Id()
	readAPI := rule.NewGetRule(ruleName)
	err := vtmClient.Do(readAPI)
	if readAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error while retrieving rule %s: %v", ruleName, err)
	}
	ruleArray := buildRuleArray(readAPI.GetResponse())
	d.SetId(ruleName)
	d.Set("rule", ruleArray)

	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	var scriptAsBytes []byte
	hasChanges := false

	vtmRule.Name = d.Id()

	if d.HasChange("rule") {
		if v, ok := d.GetOk("rule"); ok {
			vtmRule.Script = buildRule(v)
			scriptAsBytes = []byte(fmt.Sprintf(vtmRule.Script))
		}
		hasChanges = true
	}

	if hasChanges {
		vtmClient := m.(*brocadevtm.VTMClient)
		headers := make(map[string]string)
		headers["Content-Type"] = "application/octet-stream"
		headers["Content-Transfer-Encoding"] = "text"
		vtmClient.Headers = headers

		updateAPI := rule.NewUpdate(vtmRule.Name, scriptAsBytes)
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("Error while updating rule %s: %v", vtmRule.Name, err)
		}
		d.SetId(vtmRule.Name)
		d.Set("rule", buildRuleArray(vtmRule.Script))

	}
	return resourceRuleRead(d, m)
}

func resourceRuleDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)

	ruleName := d.Id()
	deleteAPI := rule.NewDelete(ruleName)
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("Error while deleting rule %s: %v", ruleName, err)
	}

	d.SetId("")
	return nil
}
