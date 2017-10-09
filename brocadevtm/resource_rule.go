package brocadevtm

/*
import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/rule"
	"github.com/sky-uk/go-rest-api"
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
	headers := make(map[string]string)

	// We need to copy the client as we want to specify different headers for rule which will conflict with other resources.
	client := m.(*rest.Client)
	vtmClient := *client
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
		return fmt.Errorf("BrocadeVTM Rule error whilst creating %s: %v", vtmRule.Name, err)
	}

	d.SetId(vtmRule.Name)

	return resourceRuleRead(d, m)
}

func resourceRuleRead(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	headers := make(map[string]string)

	// We need to copy the client as we want to specify different headers for rule which will conflict with other resources.
	client := m.(*rest.Client)
	vtmClient := *client
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Transfer-Encoding"] = "text"
	vtmClient.Headers = headers

	vtmRule.Name = d.Id()
	readAPI := rule.NewGet(vtmRule.Name)
	err := vtmClient.Do(readAPI)
	if err != nil {
		if readAPI.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM Rule error whilst retrieving %s: %v", vtmRule.Name, err)
	}

	response := readAPI.ResponseObject().(*[]byte)
	vtmRule.Script = string(*response)

	d.SetId(vtmRule.Name)
	d.Set("rule", vtmRule.Script)

	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	headers := make(map[string]string)
	hasChanges := false
	vtmRule.Name = d.Id()

	if d.HasChange("rule") {
		if v, ok := d.GetOk("rule"); ok {
			vtmRule.Script = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		// We need to copy the client as we want to specify different headers for rule which will conflict with other resources.
		client := m.(*rest.Client)
		vtmClient := *client
		headers["Content-Type"] = "application/octet-stream"
		headers["Content-Transfer-Encoding"] = "text"
		vtmClient.Headers = headers

		updateAPI := rule.NewUpdate(vtmRule.Name, []byte(fmt.Sprintf(vtmRule.Script)))
		err := vtmClient.Do(updateAPI)

		if err != nil {
			return fmt.Errorf("BrocadeVTM Rule error whilst updating %s: %vv", vtmRule.Name, err)
		}
		d.SetId(vtmRule.Name)
		d.Set("rule", vtmRule.Script)
	}

	return resourceRuleRead(d, m)
}

func resourceRuleDelete(d *schema.ResourceData, m interface{}) error {

	var vtmRule rule.TrafficScriptRule
	vtmClient := m.(*rest.Client)

	vtmRule.Name = d.Id()
	deleteAPI := rule.NewDelete(vtmRule.Name)
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("BrocadeVTM Rule error whilst deleting %s: %v", vtmRule.Name, err)
	}

	d.SetId("")
	return nil
}
*/
