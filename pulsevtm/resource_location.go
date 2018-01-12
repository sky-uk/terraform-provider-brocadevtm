package pulsevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-pulse-vtm/api"
	"github.com/sky-uk/terraform-provider-pulsevtm/pulsevtm/util"
	"net/http"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocationCreate,
		Read:   resourceLocationRead,
		Update: resourceLocationUpdate,
		Delete: resourceLocationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name of the location",
				ForceNew:    true,
			},
			"location_id": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The location identifier",
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"latitude": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Default:      0.0,
				Description:  "The latitude of the location",
				ValidateFunc: checkLatitudeWithinRange,
			},
			"longitude": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Default:      0.0,
				Description:  "The longitude of the location",
				ValidateFunc: checkLongitudeWithinRange,
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A note regarding the location",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Is the location used by traffic managers or GLBs?",
				ValidateFunc: checkValidLocationType,
			},
		},
	}
}

func checkValidLocationType(v interface{}, k string) (ws []string, errors []error) {
	locationType := v.(string)
	if locationType != "config" && locationType != "glb" {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be one of config or glb", k))
	}
	return
}

func checkLatitudeWithinRange(v interface{}, k string) (ws []string, errors []error) {
	latitude := v.(float64)
	if latitude < -90 || latitude > 90 {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be between -90 and 90 degrees inclusive", k))
	}
	return
}

func checkLongitudeWithinRange(v interface{}, k string) (ws []string, errors []error) {
	longitude := v.(float64)
	if longitude < -180 || longitude > 180 {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be between -180 and 180 degrees inclusive", k))
	}
	return
}
func locationAttribute(name string) string {
	if name == "location_id" {
		return "id"
	}
	if name == "id" {
		return "location_id"
	}
	return name
}

func resourceLocationCreate(d *schema.ResourceData, m interface{}) error {

	var name string
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	locationBasicConfiguration := make(map[string]interface{})
	locationPropertiesConfiguration := make(map[string]interface{})
	locationConfiguration := make(map[string]interface{})

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	if v, ok := d.GetOk("location_id"); ok {
		locationBasicConfiguration["id"] = uint(v.(int))
	}
	if v, ok := d.GetOk("latitude"); ok {
		locationBasicConfiguration["latitude"] = v.(float64)
	}
	if v, ok := d.GetOk("longitude"); ok {
		locationBasicConfiguration["longitude"] = v.(float64)
	}
	if v, ok := d.GetOk("note"); ok && v != "" {
		locationBasicConfiguration["note"] = v.(string)
	}
	if v, ok := d.GetOk("type"); ok && v != "" {
		locationBasicConfiguration["type"] = v.(string)
	}

	locationPropertiesConfiguration["basic"] = locationBasicConfiguration
	locationConfiguration["properties"] = locationPropertiesConfiguration

	err := client.Set("locations", name, locationConfiguration, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM Location error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceLocationRead(d, m)
}

func resourceLocationRead(d *schema.ResourceData, m interface{}) error {

	name := d.Id()
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	locationBasicConfiguration := make(map[string]interface{})
	locationPropertiesConfiguration := make(map[string]interface{})
	locationConfiguration := make(map[string]interface{})

	client.WorkWithConfigurationResources()
	err := client.GetByName("locations", name, &locationConfiguration)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		d.SetId("")
		return fmt.Errorf("[ERROR] PulseVTM location error whilst retrieving %s: %v", name, err)
	}
	locationPropertiesConfiguration = locationConfiguration["properties"].(map[string]interface{})
	locationBasicConfiguration = locationPropertiesConfiguration["basic"].(map[string]interface{})

	for _, attribute := range []string{"id", "latitude", "longitude", "note", "type"} {
		err := d.Set(locationAttribute(attribute), locationBasicConfiguration[attribute])
		if err != nil {
			return fmt.Errorf("[ERROR] PulseVTM location error whilst setting attribute %s: %v ", attribute, err)
		}

	}

	return nil
}

func resourceLocationUpdate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	var name string

	locationBasicConfiguration := make(map[string]interface{})
	locationPropertiesConfiguration := make(map[string]interface{})
	locationConfiguration := make(map[string]interface{})

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if d.HasChange("location_id") {
		locationBasicConfiguration["id"] = uint(d.Get("location_id").(int))
	}
	if d.HasChange("latitude") {
		locationBasicConfiguration["latitude"] = d.Get("latitude").(float64)
	}
	if d.HasChange("longitude") {
		locationBasicConfiguration["longitude"] = d.Get("longitude").(float64)
	}
	if d.HasChange("note") {
		locationBasicConfiguration["note"] = d.Get("note").(string)
	}
	if d.HasChange("type") {
		locationBasicConfiguration["type"] = d.Get("type").(string)
	}

	locationPropertiesConfiguration["basic"] = locationBasicConfiguration
	locationConfiguration["properties"] = locationPropertiesConfiguration

	err := client.Set("locations", name, locationConfiguration, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM locations error whilst updating %s: %v", name, err)
	}

	return resourceLocationRead(d, m)
}

func resourceLocationDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("locations", d, m)
}
