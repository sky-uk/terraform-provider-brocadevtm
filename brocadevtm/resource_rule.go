package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/rule"
	"net/http"
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
				Type:        schema.TypeString,
				Description: "The rule in traffic script language as a here document",
				Required:    true,
			},
		},
	}
}

func resourceRuleCreate(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	config := m.(map[string]interface{})

	client := config["octetClient"].(*api.Client)

	if v, ok := d.GetOk("name"); ok && v != "" {
		vtmRule.Name = v.(string)
	}
	if v, ok := d.GetOk("rule"); ok {
		vtmRule.Script = v.(string)
	}

	err := client.Set("rules", vtmRule.Name, []byte(vtmRule.Script), nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Rule error whilst creating %s: %v", vtmRule.Name, err)
	}

	d.SetId(vtmRule.Name)

	return resourceRuleRead(d, m)
}

func resourceRuleRead(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	config := m.(map[string]interface{})

	client := config["octetClient"].(*api.Client)
	vtmRule.Name = d.Id()
	client.WorkWithConfigurationResources()
	ruleText := new([]byte)
	err := client.GetByName("rules", vtmRule.Name, ruleText)

	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
	}

	if err != nil {
		return fmt.Errorf("BrocadeVTM Rule error whilst retrieving %s: %v", vtmRule.Name, err)
	}

	d.SetId(vtmRule.Name)
	d.Set("rule", string(*ruleText))
	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	config := m.(map[string]interface{})
	hasChanges := false
	vtmRule.Name = d.Id()

	if d.HasChange("rule") {
		if v, ok := d.GetOk("rule"); ok {
			vtmRule.Script = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {

		client := config["octetClient"].(*api.Client)
		err := client.Set("rules", vtmRule.Name, []byte(vtmRule.Script), nil)

		if err != nil {
			return fmt.Errorf("BrocadeVTM Rule error whilst updating %s: %vv", vtmRule.Name, err)
		}
		d.SetId(vtmRule.Name)
		d.Set("rule", vtmRule.Script)
	}

	return resourceRuleRead(d, m)
}

func resourceRuleDelete(d *schema.ResourceData, m interface{}) error {
	err := DeleteResource("rules", d, m)
	if err != nil {
		return err
	}
	return nil
}
