package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"net/http"
)

func resourceAptimizerProfile() *schema.Resource {
	return &schema.Resource{

		Create: resourceAptimizerProfileCreate,
		Read:   resourceAptimizerProfileRead,
		Update: resourceAptimizerProfileUpdate,
		Delete: resourceAptimizerProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The name of the Web Accelerator Profile",
				ValidateFunc: util.NoZeroValues,
			},
			"background_after": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "If Web Accelerator can finish optimizing the resource within this time limit then serve the optimized content to the client, otherwise complete the optimization in the background and return the original content to the client. 0 = Web accelerator will always wait for the optimization to complete before sending a response to the client",
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"background_on_additional_resources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If a web page contains resources that have not yet been optimized, fetch and optimize those resources in the background and send a partially optimized web page to clients until all resources on that page are ready.",
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "active",
				Description:  "Set the Web Accelerator mode to turn acceleration on, off or in stealth mode",
				ValidateFunc: validateMode,
			},
			"show_info_bar": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: " Show the Web Accelerator information bar on optimized web pages. This requires HTML optimization to be enabled in the acceleration settings.",
			},
		},
	}
}

func validateMode(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"active",
		"idle",
		"stealth":
		return
	}
	errors = append(errors, fmt.Errorf("%s must be one of active, idle or stealth", k))
	return
}

func resourceAptimizerProfileCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	aptimizerProfileBasicConfig := make(map[string]interface{})
	aptimizerProfilePropertiesConfig := make(map[string]interface{})
	aptimizerProfileConfig := make(map[string]interface{})

	client := config["jsonClient"].(*api.Client)

	name := d.Get("name").(string)
	aptimizerProfileBasicConfig["background_after"] = uint(d.Get("background_after").(int))
	aptimizerProfileBasicConfig["background_on_additional_resources"] = d.Get("background_on_additional_resources").(bool)
	aptimizerProfileBasicConfig["mode"] = d.Get("mode").(string)
	aptimizerProfileBasicConfig["show_info_bar"] = d.Get("show_info_bar").(bool)

	aptimizerProfilePropertiesConfig["basic"] = aptimizerProfileBasicConfig
	aptimizerProfileConfig["properties"] = aptimizerProfilePropertiesConfig

	err := client.Set("aptimizer/profiles", name, aptimizerProfileConfig, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst creating Aptimizer Profile %s: %v", name, err)
	}
	d.SetId(name)
	return resourceAptimizerProfileRead(d, m)
}

func resourceAptimizerProfileRead(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	name := d.Id()
	aptimizerProfileConfig := make(map[string]interface{})

	err := client.GetByName("aptimizer/profiles", name, &aptimizerProfileConfig)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst retrieving Aptimizer Profile %s: %v", name, err)
	}

	aptimizerProfilePropertiesConfig := aptimizerProfileConfig["properties"].(map[string]interface{})
	aptimizerProfileBasicConfig := aptimizerProfilePropertiesConfig["basic"].(map[string]interface{})

	d.Set("background_after", aptimizerProfileBasicConfig["background_after"])
	d.Set("background_on_additional_resources", aptimizerProfileBasicConfig["background_on_additional_resources"])
	d.Set("mode", aptimizerProfileBasicConfig["mode"])
	d.Set("show_info_bar", aptimizerProfileBasicConfig["show_info_bar"])

	return nil
}

func resourceAptimizerProfileUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	aptimizerProfileBasicConfig := make(map[string]interface{})
	aptimizerProfilePropertiesConfig := make(map[string]interface{})
	aptimizerProfileConfig := make(map[string]interface{})

	client := config["jsonClient"].(*api.Client)
	if d.HasChange("background_after") {
		aptimizerProfileBasicConfig["background_after"] = uint(d.Get("background_after").(int))
	}
	if d.HasChange("background_on_additional_resources") {
		aptimizerProfileBasicConfig["background_on_additional_resources"] = d.Get("background_on_additional_resources").(bool)
	}
	if d.HasChange("mode") {
		aptimizerProfileBasicConfig["mode"] = d.Get("mode").(string)
	}
	if d.HasChange("show_info_bar") {
		aptimizerProfileBasicConfig["show_info_bar"] = d.Get("show_info_bar").(bool)
	}

	aptimizerProfilePropertiesConfig["basic"] = aptimizerProfileBasicConfig
	aptimizerProfileConfig["properties"] = aptimizerProfilePropertiesConfig

	err := client.Set("aptimizer/profiles", d.Id(), aptimizerProfileConfig, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst updating Aptimizer Profile %s: %v", d.Id(), err)
	}

	return resourceAptimizerProfileRead(d, m)
}

func resourceAptimizerProfileDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("aptimizer/profiles", d, m)
}
