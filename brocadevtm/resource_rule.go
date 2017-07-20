package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/rule"
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

	vtmClient := m.(*brocadevtm.VTMClient)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Transfer-Encoding"] = "text"
	vtmClient.Headers = headers

	if v, ok := d.GetOk("name"); ok && v != "" {
		vtmRule.Name = v.(string)
	}
	if v, ok := d.GetOk("rule"); ok {
		vtmRule.Script = v.(string)
	}

	createAPI := rule.NewCreate(vtmRule.Name, []byte(fmt.Sprintf(vtmRule.Script)))
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Error while creating rule %s: %+v", vtmRule.Name, err)
	}

	d.SetId(vtmRule.Name)
	return resourceRuleRead(d, m)
}

func resourceRuleRead(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule

	vtmClient := m.(*brocadevtm.VTMClient)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Transfer-Encoding"] = "text"
	vtmClient.Headers = headers

	vtmRule.Name = d.Id()
	readAPI := rule.NewGetRule(vtmRule.Name)
	err := vtmClient.Do(readAPI)
	if readAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error while retrieving rule %s: %v", vtmRule.Name, err)
	}
	vtmRule.Script = readAPI.GetResponse()
	d.SetId(vtmRule.Name)
	d.Set("rule", vtmRule.Script)

	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	hasChanges := false
	vtmRule.Name = d.Id()

	if d.HasChange("rule") {
		if v, ok := d.GetOk("rule"); ok {
			vtmRule.Script = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		vtmClient := m.(*brocadevtm.VTMClient)
		headers := make(map[string]string)
		headers["Content-Type"] = "application/octet-stream"
		headers["Content-Transfer-Encoding"] = "text"
		vtmClient.Headers = headers

		updateAPI := rule.NewUpdate(vtmRule.Name, []byte(fmt.Sprintf(vtmRule.Script)))
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("Error while updating rule %s: %v", vtmRule.Name, err)
		}
		d.SetId(vtmRule.Name)
		d.Set("rule", vtmRule.Script)

	}
	return resourceRuleRead(d, m)
}

func resourceRuleDelete(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	vtmClient := m.(*brocadevtm.VTMClient)

	vtmRule.Name = d.Id()
	deleteAPI := rule.NewDelete(vtmRule.Name)
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("Error while deleting rule %s: %v", vtmRule.Name, err)
	}

	d.SetId("")
	return nil
}
