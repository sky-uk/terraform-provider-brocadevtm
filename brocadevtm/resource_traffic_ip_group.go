package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/traffic_ip_group"
	"net/http"
)

func resourceTrafficIPGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrafficIPGroupCreate,
		Read:   resourceTrafficIPGroupRead,
		Update: resourceTrafficIPGroupUpdate,
		Delete: resourceTrafficIPGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Traffic IP group name",
				Required:    true,
				ForceNew:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the traffic IP group should be enabled",
				Optional:    true,
				Computed:    true,
			},
			"hashsourceport": {
				Type:        schema.TypeBool,
				Description: "Whether or not the source port should be taken into account when deciding which traffic manager should handle a request.",
				Optional:    true,
				Computed:    true,
			},

			"ipaddresses": {
				// Check API doco re updates.
				Type:        schema.TypeList,
				Description: "List of IP addresses to raise on the traffic managers - typically this is one IP address",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"trafficmanagers": {
				Type:        schema.TypeList,
				Description: "List of traffic managers on which to raise this traffic IP - automatically retrieved from vTM",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"mode": {
				Type:        schema.TypeString,
				Description: "The method used to distribute traffic IPs across machines in the cluster - multihosted when using multicast",
				Optional:    true,
				Computed:    true,
			},
			"multicastip": {
				// Check API doco re updates.
				Type:        schema.TypeString,
				Description: "Multicast IP address",
				Required:    true,
			},
		},
	}
}

func getTrafficManagers(m interface{}) ([]string, error) {
	vtmClient := m.(*brocadevtm.VTMClient)
	getTrafficManagersAPI := trafficIpGroups.NewGetTrafficManagerList()
	var trafficManagers []string

	err := vtmClient.Do(getTrafficManagersAPI)
	if err != nil {
		return trafficManagers, fmt.Errorf("Error retrieving a list of traffic managers")
	}

	response := getTrafficManagersAPI.GetResponse()
	for _, trafficManager := range response.Children {
		trafficManagers = append(trafficManagers, trafficManager.Name)
	}
	return trafficManagers, nil
}

func resourceTrafficIPGroupCreate(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*brocadevtm.VTMClient)
	var createTrafficIPGroup trafficIpGroups.TrafficIPGroup
	var tipgName string

	// Retrieve the list of Brocade vTM traffic managers and assign it to Machines
	trafficManagers, err := getTrafficManagers(m)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	createTrafficIPGroup.Properties.Basic.Machines = trafficManagers

	if v, ok := d.GetOk("name"); ok && v != "" {
		tipgName = v.(string)
	}
	if v, _ := d.GetOk("enabled"); v != nil {
		enableTrafficIPGroup := v.(bool)
		createTrafficIPGroup.Properties.Basic.Enabled = &enableTrafficIPGroup
	}
	if v, ok := d.GetOk("hashsourceport"); ok {
		hashSourcePort := v.(bool)
		createTrafficIPGroup.Properties.Basic.HashSourcePort = &hashSourcePort
	}
	if v, ok := d.GetOk("ipaddresses"); ok && v != "" {
		ipAddresses := make([]string, len(v.([]interface{})))
		for idx, ipAddress := range v.([]interface{}) {
			ipAddresses[idx] = ipAddress.(string)
		}
		createTrafficIPGroup.Properties.Basic.IPAddresses = ipAddresses
	} else {
		return fmt.Errorf("ipaddresses argument required")
	}
	if v, ok := d.GetOk("mode"); ok && v != "" {
		createTrafficIPGroup.Properties.Basic.Mode = v.(string)
	}
	if v, ok := d.GetOk("multicastip"); ok && v != "" {
		createTrafficIPGroup.Properties.Basic.Multicast = v.(string)
	} else {
		return fmt.Errorf("multicastip argument required")
	}

	createAPI := trafficIpGroups.NewCreate(tipgName, createTrafficIPGroup)
	err = vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Error creating traffic IP group %s: %+v", tipgName, err)
	}
	if createAPI.StatusCode() != http.StatusCreated && createAPI.StatusCode() != http.StatusOK {
		return fmt.Errorf("Invalid HTTP response code %d returned. Response object was %+v", createAPI.StatusCode(), createAPI.ResponseObject())
	}

	d.SetId(tipgName)
	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupRead(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*brocadevtm.VTMClient)
	var readTrafficIPGroup trafficIpGroups.TrafficIPGroup
	var tipgName string

	if v, ok := d.GetOk("name"); ok {
		tipgName = v.(string)
	} else {
		return fmt.Errorf("Name argument required")
	}

	getSingleAPI := trafficIpGroups.NewGetSingle(tipgName)
	err := vtmClient.Do(getSingleAPI)
	if err != nil {
		return fmt.Errorf("Error reading Traffic IP Group: %+v", err)
	}
	if getSingleAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	readTrafficIPGroup = *getSingleAPI.GetResponse()
	d.Set("name", tipgName)
	d.Set("enabled", *readTrafficIPGroup.Properties.Basic.Enabled)
	d.Set("hashsourceport", *readTrafficIPGroup.Properties.Basic.HashSourcePort)
	d.Set("ipaddresses", readTrafficIPGroup.Properties.Basic.IPAddresses)
	d.Set("trafficmanagers", readTrafficIPGroup.Properties.Basic.Machines)
	d.Set("mode", readTrafficIPGroup.Properties.Basic.Mode)
	d.Set("multicastip", readTrafficIPGroup.Properties.Basic.Multicast)

	return nil
}

func resourceTrafficIPGroupUpdate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var trafficIPGroupName string
	var updateTrafficIPGroup trafficIpGroups.TrafficIPGroup
	hasChanges := false

	//TODO ipaddresses, trafficmanagers, mode

	if v, ok := d.GetOk("name"); ok && v != "" {
		trafficIPGroupName = v.(string)
	} else {
		return fmt.Errorf("Error Traffic IP Group Update requires name attribute")
	}
	if d.HasChange("enabled") {
		enableTrafficIPGroup := d.Get("enabled").(bool)
		updateTrafficIPGroup.Properties.Basic.Enabled = &enableTrafficIPGroup
		hasChanges = true
	}
	if d.HasChange("hashsourceport") {
		hashSourcePort := d.Get("hashsourceport").(bool)
		updateTrafficIPGroup.Properties.Basic.HashSourcePort = &hashSourcePort
		hasChanges = true
	}
	if d.HasChange("multicastip") {
		if v, ok := d.GetOk("multicastip"); ok {
			updateTrafficIPGroup.Properties.Basic.Multicast = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		updateAPI := trafficIpGroups.NewUpdate(trafficIPGroupName, updateTrafficIPGroup)
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("Error Updating Traffic IP Group %s", trafficIPGroupName)
		}
		if updateAPI.StatusCode() != http.StatusOK {
			return fmt.Errorf("Error updating Traffic IP Group %s invalid http return code", trafficIPGroupName)
		}

		updateResponse := updateAPI.GetResponse()
		d.SetId(trafficIPGroupName)
		d.Set("enabled", *updateResponse.Properties.Basic.Enabled)
		d.Set("hashsourceport", *updateResponse.Properties.Basic.HashSourcePort)
		d.Set("multicastip", updateResponse.Properties.Basic.Multicast)
	}
	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("Name argument required")
	}

	getTrafficIPGroup := trafficIpGroups.NewGetSingle(name)
	err := vtmClient.Do(getTrafficIPGroup)
	if err != nil {
		return fmt.Errorf("Error fetching traffic IP group %s", name)
	}
	if getTrafficIPGroup.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	deleteAPI := trafficIpGroups.NewDelete(name)
	err = vtmClient.Do(deleteAPI)
	if err != nil {
		return fmt.Errorf("Error deleting traffic IP group %s", name)
	}
	if deleteAPI.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("Error deleting traffic IP group %s - status code != 204", name)
	}

	d.SetId("")
	return nil
}
