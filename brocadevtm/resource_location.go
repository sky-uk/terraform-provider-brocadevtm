package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/location"
	"github.com/sky-uk/go-brocade-vtm/api/monitor"
	"github.com/sky-uk/go-rest-api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
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
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The location identifier",
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"latitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Default:     0.0,
				Description: "The latitude of the location",
			},
			"longitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Default:     0.0,
				Description: "The longitude of the location",
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
				Description:  "Is the location used by traffic managers or for GLBs?",
				ValidateFunc: checkValidLocationType,
			},
		},
	}
}

func checkValidLocationType(v interface{}, k string) (ws []string, errors []error) {
	locationType := v.(string)
	if locationType != "config" && locationType != "glb" {
		errors = append(errors, fmt.Errorf("%q must be one of config or glb", k))
	}
	return
}

func resourceLocationCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var createLocation location.Location
	var name string

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	if v, ok := d.GetOk("id"); ok && v != "" {
		locationID := v.(int)
		createLocation.Properties.Basic.ID = uint(locationID)
	}
	if v, ok := d.GetOk("latitude"); ok && v != "" {
		createLocation.Properties.Basic.Latitude = v.(float32)
	}
	if v, ok := d.GetOk("longitude"); ok && v != "" {
		createLocation.Properties.Basic.Longitude = v.(float32)
	}
	if v, ok := d.GetOk("note"); ok && v != "" {
		createLocation.Properties.Basic.Note = v.(string)
	}
	if v, ok := d.GetOk("type"); ok && v != "" {
		createLocation.Properties.Basic.Type = v.(string)
	}

	createAPI := location.NewCreate(name, createLocation)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Location error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceLocationRead(d, m)
}

func resourceLocationRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	locationName := d.Id()

	getLocationAPI := location.NewGet(locationName)
	err := vtmClient.Do(getLocationAPI)
	if getLocationAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM Location error whilst retrieving %s: %v", locationName, err)
	}

	getLocationProperties := getLocationAPI.ResponseObject().(*location.Location)
	d.Set("name", locationName)
	d.Set("id", getLocationProperties.Properties.Basic.ID)
	d.Set("latitude", getLocationProperties.Properties.Basic.Latitude)
	d.Set("longitude", getLocationProperties.Properties.Basic.Longitude)
	d.Set("note", getLocationProperties.Properties.Basic.Note)
	d.Set("type", getLocationProperties.Properties.Basic.Type)

	return nil
}

func resourceLocationUpdate(d *schema.ResourceData, m interface{}) error {

	return resourceLocationRead(d, m)
}

func resourceLocationDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	locationName := d.Id()

	deleteAPI := location.NewDelete(locationName)
	err := vtmClient.Do(deleteAPI)
	if deleteAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM Location error whilst deleting %s: %v", locationName, err)
	}

	d.SetId("")
	return nil
}
