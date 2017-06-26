package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/traffic_ip_group"
)

func resourceTrafficIPGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrafficIPGroupCreate,
		Read:   resourceTrafficIPGroupRead,
		Update: resourceTrafficIPGroupUpdate,
		Delete: resourceTrafficIPGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machines": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return nil
}

func resourceTrafficIPGroupCreate(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*brocadevtm.VTMClient)
	var createTrafficIPGroup trafficIpGroups.TrafficIPGroup
	var tipgName string

	if v, ok := d.GetOk("name"); ok {
		tipgName = v.(string)
	} else {
		return fmt.Errorf("Name argument required")
	}

	if v, ok := d.GetOk("machines"); ok {
		machineList := v.([]interface{})
		machinesToAdd := make([]string, len(machineList))
		for i, value := range machineList {
			tagID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in machines")
			}
			machinesToAdd[i] = tagID
		}
		createTrafficIPGroup.Properties.Basic.Machines = machinesToAdd
	}

	createAPI := trafficIpGroups.NewCreate(tipgName, createTrafficIPGroup)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Could not create traffic IP group: %+v", err)
	}
	if createAPI.StatusCode() != 201 && createAPI.StatusCode() != 200 {
		return fmt.Errorf("Invalid HTTP response code %+v returned. Response object was %+v", createAPI.StatusCode(), createAPI.ResponseObject())
	}
	d.SetId(tipgName)

	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupRead(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*brocadevtm.VTMClient)
	var createTrafficIPGroup trafficIpGroups.TrafficIPGroup
	var tipgName string
	if v, ok := d.GetOk("name"); ok {
		tipgName = v.(string)
	} else {
		return fmt.Errorf("Name argument required")
	}
	if v, ok := d.GetOk("machines"); ok {
		machineList := v.([]interface{})
		machinesToAdd := make([]string, len(machineList))
		for i, value := range machineList {
			tagID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in machines")
			}
			machinesToAdd[i] = tagID
		}
		createTrafficIPGroup.Properties.Basic.Machines = machinesToAdd
	}

	getSingleAPI := trafficIpGroups.NewGetSingle(tipgName)
	err := vtmClient.Do(getSingleAPI)
	if err != nil {
		return fmt.Errorf("Error reading Traffic IP Group:", err)
	}
	d.Set("name", tipgName)
	d.Set("machines", createTrafficIPGroup.Properties.Basic.Machines)

	return nil
}

func resourceTrafficIPGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTrafficIPGroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
