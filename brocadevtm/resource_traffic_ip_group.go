package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/traffic_ip_group"
	"github.com/sky-uk/go-brocade-vtm/api/traffic_ip_group_manager"
	"github.com/sky-uk/go-rest-api"
	"net/http"
	"regexp"
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
				Type:        schema.TypeList,
				Description: "List of IP addresses to raise on the traffic managers",
				Optional:    true,
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
				Type:         schema.TypeString,
				Description:  "The method used to distribute traffic IPs across machines in the cluster",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateTrafficIPGroupMode,
			},
			"multicastip": {
				Type:         schema.TypeString,
				Description:  "Multicast IP address",
				Optional:     true,
				ValidateFunc: validateTrafficIPGroupMulticastIP,
			},
		},
	}
}

func validateTrafficIPGroupMode(v interface{}, k string) (ws []string, errors []error) {
	mode := v.(string)
	modeOptions := regexp.MustCompile(`^(singlehosted|ec2elastic|ec2vpcelastic|ec2vpcprivate|multihosted|rhi)$`)
	if !modeOptions.MatchString(mode) {
		errors = append(errors, fmt.Errorf("%q must be one of singlehosted, ec2elastic, ec2vpcelastic, ec2vpcprivate, multihosted or rhi", k))
	}
	return
}

func validateTrafficIPGroupMulticastIP(v interface{}, k string) (ws []string, errors []error) {
	multicastIP := v.(string)
	validMulticastIPs := regexp.MustCompile(`^2[2-3][0-9]\.[0-9]+\.[0-9]+\.[0-9]+$`)
	if !validMulticastIPs.MatchString(multicastIP) {
		errors = append(errors, fmt.Errorf("%q must be a valid multicast IP (224.0.0.0 - 239.255.255.255)", k))
	}
	return
}

func getTrafficManagers(m interface{}) ([]string, error) {
	vtmClient := m.(*rest.Client)
	getTrafficManagersAPI := trafficIpGroupManager.NewGetAll()
	var trafficManagers []string

	err := vtmClient.Do(getTrafficManagersAPI)
	if err != nil {
		return trafficManagers, fmt.Errorf("error retrieving a list of traffic managers")
	}

	response := getTrafficManagersAPI.ResponseObject().(*trafficIpGroupManager.TrafficManagerChildren)
	for _, trafficManager := range response.Children {
		trafficManagers = append(trafficManagers, trafficManager.Name)
	}
	return trafficManagers, nil
}

func buildIPAddresses(ipAddresses interface{}) []string {
	ipAddressList := make([]string, len(ipAddresses.([]interface{})))
	for idx, ipAddress := range ipAddresses.([]interface{}) {
		ipAddressList[idx] = ipAddress.(string)
	}
	return ipAddressList
}

func resourceTrafficIPGroupCreate(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*rest.Client)
	var createTrafficIPGroup trafficIpGroups.TrafficIPGroup
	var tipgName string

	// Retrieve the list of Brocade vTM traffic managers and assign it to Machines
	trafficManagers, err := getTrafficManagers(m)
	if err != nil {
		return fmt.Errorf("Traffic IP Group create %v", err)
	}
	createTrafficIPGroup.Properties.Basic.Machines = trafficManagers

	if v, ok := d.GetOk("name"); ok && v != "" {
		tipgName = v.(string)
	} else {
		return fmt.Errorf("Traffic IP Group create requires name argument")
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
		createTrafficIPGroup.Properties.Basic.IPAddresses = buildIPAddresses(v)
	}
	if v, ok := d.GetOk("mode"); ok && v != "" {
		createTrafficIPGroup.Properties.Basic.Mode = v.(string)
	}
	if v, ok := d.GetOk("multicastip"); ok && v != "" {
		createTrafficIPGroup.Properties.Basic.Multicast = v.(string)
	}

	createAPI := trafficIpGroups.NewCreate(tipgName, createTrafficIPGroup)
	err = vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Traffic IP Group create error when creating traffic IP group %s: %+v", tipgName, err)
	}
	if createAPI.StatusCode() != http.StatusCreated && createAPI.StatusCode() != http.StatusOK {
		return fmt.Errorf("Invalid HTTP response code %d returned when creating traffic IP group %s. Response object was %+v", createAPI.StatusCode(), tipgName, createAPI.ResponseObject())
	}

	d.SetId(tipgName)
	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupRead(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*rest.Client)
	var readTrafficIPGroup trafficIpGroups.TrafficIPGroup
	var tipgName string

	if v, ok := d.GetOk("name"); ok {
		tipgName = v.(string)
	} else {
		return fmt.Errorf("Traffic IP Group read error: name argument required")
	}

	getSingleAPI := trafficIpGroups.NewGet(tipgName)
	err := vtmClient.Do(getSingleAPI)
	if err != nil {
		return fmt.Errorf("Traffic IP Group read error while reading Traffic IP Group %s: %+v", tipgName, err)
	}
	if getSingleAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	readTrafficIPGroup = *getSingleAPI.ResponseObject().(*trafficIpGroups.TrafficIPGroup)
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

	vtmClient := m.(*rest.Client)
	var trafficIPGroupName string
	var updateTrafficIPGroup trafficIpGroups.TrafficIPGroup
	hasChanges := false

	if v, ok := d.GetOk("name"); ok && v != "" {
		trafficIPGroupName = v.(string)
	} else {
		return fmt.Errorf("Traffic IP Group update error: name attribute required")
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
	if d.HasChange("ipaddresses") {
		if v, ok := d.GetOk("ipaddresses"); ok && v != "" {
			updateTrafficIPGroup.Properties.Basic.IPAddresses = buildIPAddresses(v)
		}
		hasChanges = true
	}
	if d.HasChange("mode") {
		if v, ok := d.GetOk("mode"); ok && v != "" {
			updateTrafficIPGroup.Properties.Basic.Mode = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("multicastip") {
		if v, ok := d.GetOk("multicastip"); ok && v != "" {
			updateTrafficIPGroup.Properties.Basic.Multicast = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		// On all updates refresh the list of traffic managers
		trafficManagers, err := getTrafficManagers(m)
		if err != nil {
			return fmt.Errorf("Traffic IP Group update error while updating %s: %v", trafficIPGroupName, err)
		}
		updateTrafficIPGroup.Properties.Basic.Machines = trafficManagers

		updateAPI := trafficIpGroups.NewUpdate(trafficIPGroupName, updateTrafficIPGroup)
		err = vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("Traffic IP Group update error while updating %s: %v", trafficIPGroupName, err)
		}
		if updateAPI.StatusCode() != http.StatusOK {
			return fmt.Errorf("Traffic IP Group update error while updating %s: received invalid http return code %d", trafficIPGroupName, updateAPI.StatusCode())
		}

		updateResponse := updateAPI.ResponseObject().(*trafficIpGroups.TrafficIPGroup)
		d.SetId(trafficIPGroupName)
		d.Set("enabled", *updateResponse.Properties.Basic.Enabled)
		d.Set("hashsourceport", *updateResponse.Properties.Basic.HashSourcePort)
		d.Set("ipaddresses", updateResponse.Properties.Basic.IPAddresses)
		d.Set("trafficmanagers", updateResponse.Properties.Basic.Machines)
		d.Set("mode", updateResponse.Properties.Basic.Mode)
		d.Set("multicastip", updateResponse.Properties.Basic.Multicast)
	}
	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("Traffic IP Group delete error: name argument required")
	}

	getTrafficIPGroup := trafficIpGroups.NewGet(name)
	err := vtmClient.Do(getTrafficIPGroup)
	if err != nil {
		return fmt.Errorf("Traffic IP Group delete error while fetching traffic IP group %s: %v", name, err)
	}
	if getTrafficIPGroup.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	deleteAPI := trafficIpGroups.NewDelete(name)
	err = vtmClient.Do(deleteAPI)
	if err != nil {
		return fmt.Errorf("Traffic IP Group delete error while deleting traffic IP group %s: %v", name, err)
	}
	if deleteAPI.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("Traffic IP Group delete error: received invalid http return code %d while deleting %s", deleteAPI.StatusCode(), name)
	}

	d.SetId("")
	return nil
}
