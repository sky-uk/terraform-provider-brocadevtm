package brocadevtm

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
)

func resourceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleSet,
		Read:   resourceRuleRead,
		Update: resourceRuleSet,
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

func resourceRuleSet(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["octetClient"].(*api.Client)

	name := d.Get("name").(string)
	rule := d.Get("rule").(string)

	err := client.Set("rules", name, []byte(rule), nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Rule error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceRuleRead(d, m)
}

func resourceRuleRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["octetClient"].(*api.Client)

	ruleText := new([]byte)
	client.WorkWithConfigurationResources()
	err := client.GetByName("rules", d.Id(), ruleText)

	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM Rule error whilst retrieving %s: %v", d.Id(), err)
	}

	d.Set("rule", string(*ruleText))
	return nil
}

func resourceRuleDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("rules", d, m)
}
