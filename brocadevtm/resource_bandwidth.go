package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
	"regexp"
)

func resourceBandwidth() *schema.Resource {
	return &schema.Resource{
		Create: resourceBandwidthCreate,
		Read:   resourceBandwidthRead,
		Update: resourceBandwidthUpdate,
		Delete: resourceBandwidthDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the bandwidth class",
			},
			"maximum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10000,
				Description: "Maximum bandwidth to allocate to connections that are associated with this bandwidth class",
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A note to assign to this bandwidth class",
			},
			"sharing": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "cluster",
				Description:  "Scope of the bandwidth class",
				ValidateFunc: validateBandwidthSharing,
			},
		},
	}
}

func validateBandwidthSharing(v interface{}, k string) (ws []string, errors []error) {
	sharing := v.(string)
	sharingOptions := regexp.MustCompile(`^(cluster|connection|machine)$`)
	if !sharingOptions.MatchString(sharing) {
		errors = append(errors, fmt.Errorf("%q must be one of cluster, connection, machine", k))
	}
	return
}

func resourceBandwidthCreate(d *schema.ResourceData, m interface{}) error {

	var name string
	config := m.(map[string]interface{})
	bandwidthBasicConfiguration := make(map[string]interface{})
	bandwidthPropertiesConfiguration := make(map[string]interface{})
	bandwidthConfiguration := make(map[string]interface{})

	client := config["jsonClient"].(*api.Client)

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	if v, ok := d.GetOk("maximum"); ok {
		bandwidthBasicConfiguration["maximum"] = uint(v.(int))
	}
	if v, ok := d.GetOk("note"); ok {
		bandwidthBasicConfiguration["note"] = v.(string)
	}
	if v, ok := d.GetOk("sharing"); ok {
		bandwidthBasicConfiguration["sharing"] = v.(string)
	}
	bandwidthPropertiesConfiguration["basic"] = bandwidthBasicConfiguration
	bandwidthConfiguration["properties"] = bandwidthPropertiesConfiguration

	err := client.Set("bandwidth", name, &bandwidthConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Bandwidth error whilst creating %s: %v", name, err)
	}
	d.SetId(name)
	return resourceBandwidthRead(d, m)
}

func resourceBandwidthRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	name := d.Id()
	bandwidthConfiguration := make(map[string]interface{})

	err := client.GetByName("bandwidth", name, &bandwidthConfiguration)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM Bandwidth error whilst retrieving %s: %v", name, err)
	}

	bandwidthPropertiesConfiguration := bandwidthConfiguration["properties"].(map[string]interface{})
	bandwidthBasicConfiguration := bandwidthPropertiesConfiguration["basic"].(map[string]interface{})

	d.Set("maximum", bandwidthBasicConfiguration["maximum"])
	d.Set("note", bandwidthBasicConfiguration["note"])
	d.Set("sharing", bandwidthBasicConfiguration["sharing"])

	return nil
}

func resourceBandwidthUpdate(d *schema.ResourceData, m interface{}) error {

	name := d.Id()
	hasChanges := false
	bandwidthBasicConfiguration := make(map[string]interface{})
	bandwidthPropertiesConfiguration := make(map[string]interface{})
	bandwidthConfiguration := make(map[string]interface{})

	if d.HasChange("maximum") {
		if v, ok := d.GetOk("maximum"); ok {
			bandwidthBasicConfiguration["maximum"] = uint(v.(int))
		}
		hasChanges = true
	}
	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok {
			bandwidthBasicConfiguration["note"] = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("sharing") {
		if v, ok := d.GetOk("sharing"); ok {
			bandwidthBasicConfiguration["sharing"] = v.(string)
		}
		hasChanges = true
	}
	bandwidthPropertiesConfiguration["basic"] = bandwidthBasicConfiguration
	bandwidthConfiguration["properties"] = bandwidthPropertiesConfiguration

	if hasChanges {
		config := m.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		err := client.Set("bandwidth", name, &bandwidthConfiguration, nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Bandwidth error whilst creating %s: %v", name, err)
		}
	}
	d.SetId(name)
	return resourceBandwidthRead(d, m)
}

func resourceBandwidthDelete(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	name := d.Id()

	err := client.Delete("bandwidth", name)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM Bandwidth error whilst deleting %s: %v", name, err)
	}
	d.SetId("")
	return nil
}
