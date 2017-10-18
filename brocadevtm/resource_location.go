package brocadevtm

/*
import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/location"
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
		errors = append(errors, fmt.Errorf("%q must be one of config or glb", k))
	}
	return
}

func checkLatitudeWithinRange(v interface{}, k string) (ws []string, errors []error) {
	latitude := v.(float64)
	if latitude < -90 || latitude > 90 {
		errors = append(errors, fmt.Errorf("%q must be between -90 and 90 degrees inclusive", k))
	}
	return
}

func checkLongitudeWithinRange(v interface{}, k string) (ws []string, errors []error) {
	longitude := v.(float64)
	if longitude < -180 || longitude > 180 {
		errors = append(errors, fmt.Errorf("%q must be between -180 and 180 degrees inclusive", k))
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
	if v, ok := d.GetOk("location_id"); ok {
		locationID := v.(int)
		createLocation.Properties.Basic.ID = uint(locationID)
	}
	if v, ok := d.GetOk("latitude"); ok {
		createLocation.Properties.Basic.Latitude = v.(float64)
	}
	if v, ok := d.GetOk("longitude"); ok {
		createLocation.Properties.Basic.Longitude = v.(float64)
	}
	if v, ok := d.GetOk("note"); ok && v != "" {
		createLocation.Properties.Basic.Note = v.(string)
	}
	if v, ok := d.GetOk("type"); ok && v != "" {
		createLocation.Properties.Basic.Type = v.(string)
	}

	createLocationAPI := location.NewCreate(name, createLocation)
	err := vtmClient.Do(createLocationAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Location error whilst creating %s: %v", name, createLocationAPI.ErrorObject())
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
		return fmt.Errorf("BrocadeVTM Location error whilst retrieving %s: %v", locationName, getLocationAPI.ErrorObject())
	}

	getLocationProperties := getLocationAPI.ResponseObject().(*location.Location)
	d.Set("name", locationName)
	d.Set("location_id", getLocationProperties.Properties.Basic.ID)
	d.Set("latitude", getLocationProperties.Properties.Basic.Latitude)
	d.Set("longitude", getLocationProperties.Properties.Basic.Longitude)
	d.Set("note", getLocationProperties.Properties.Basic.Note)
	d.Set("type", getLocationProperties.Properties.Basic.Type)

	return nil
}

func resourceLocationUpdate(d *schema.ResourceData, m interface{}) error {

	hasChanges := false
	var updateLocation location.Location
	name := d.Id()

	if d.HasChange("location_id") {
		if v, ok := d.GetOk("location_id"); ok {
			locationID := v.(int)
			updateLocation.Properties.Basic.ID = uint(locationID)
		}
		hasChanges = true
	}
	if d.HasChange("latitude") {
		if v, ok := d.GetOk("latitude"); ok {
			updateLocation.Properties.Basic.Latitude = v.(float64)
		}
		hasChanges = true
	}
	if d.HasChange("longitude") {
		if v, ok := d.GetOk("longitude"); ok {
			updateLocation.Properties.Basic.Longitude = v.(float64)
		}
		hasChanges = true
	}
	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok && v != "" {
			updateLocation.Properties.Basic.Note = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok && v != "" {
			updateLocation.Properties.Basic.Type = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		vtmClient := m.(*rest.Client)
		updateLocationAPI := location.NewUpdate(name, updateLocation)
		err := vtmClient.Do(updateLocationAPI)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Location error whilst updating %s: %v", name, updateLocationAPI.ErrorObject())
		}
	}

	d.SetId(name)
	return resourceLocationRead(d, m)
}

func resourceLocationDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	locationName := d.Id()

	deleteLocationAPI := location.NewDelete(locationName)
	err := vtmClient.Do(deleteLocationAPI)
	if deleteLocationAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM Location error whilst deleting %s: %v", locationName, deleteLocationAPI.ErrorObject())
	}

	d.SetId("")
	return nil
}
*/
