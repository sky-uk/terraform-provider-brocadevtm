package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/monitor"
)

func resourceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorCreate,
		Read:   resourceMonitorRead,
		Update: resourceMonitorUpdate,
		Delete: resourceMonitorDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"failures": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"verbose": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"use_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"http_host_header": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_body_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceMonitorCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var createMonitor monitor.Monitor
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("BrocadeVTM Create Error: name argument required")
	}
	if v, ok := d.GetOk("delay"); ok {
		createMonitor.Properties.Basic.Delay = v.(int)
	}
	if v, ok := d.GetOk("timeout"); ok {
		createMonitor.Properties.Basic.Timeout = v.(int)
	}
	if v, ok := d.GetOk("failures"); ok {
		createMonitor.Properties.Basic.Failures = v.(int)
	}
	if v, ok := d.GetOk("verbose"); ok {
		monitorVerbosity := v.(bool)
		createMonitor.Properties.Basic.Verbose = &monitorVerbosity
	}
	if v, ok := d.GetOk("use_ssl"); ok {
		monitorSSL := v.(bool)
		createMonitor.Properties.Basic.UseSSL = &monitorSSL
	}
	if v, ok := d.GetOk("http_host_header"); ok {
		createMonitor.Properties.HTTP.HostHeader = v.(string)
	}
	if v, ok := d.GetOk("http_path"); ok {
		createMonitor.Properties.HTTP.URIPath = v.(string)
	}
	if v, ok := d.GetOk("http_authentication"); ok {
		createMonitor.Properties.HTTP.Authentication = v.(string)
	}
	if v, ok := d.GetOk("http_body_regex"); ok {
		createMonitor.Properties.HTTP.BodyRegex = v.(string)
	}

	createAPI := monitor.NewCreate(name, createMonitor)

	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Create Error: %+v", err)
	}

	if createAPI.StatusCode() != 201 && createAPI.StatusCode() != 200 {
		return fmt.Errorf("BrocadeVTM Create Error: Invalid HTTP response code %+v returned. Response object was %+v", createAPI.StatusCode(), createAPI.ResponseObject())
	}

	d.SetId(name)
	return resourceMonitorRead(d, m)

}

func resourceMonitorRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var readName string

	if v, ok := d.GetOk("name"); ok {
		readName = v.(string)
	} else {
		return fmt.Errorf("BrocadeVTM Read Error: name argument required")
	}

	getAllAPI := monitor.NewGetAll()
	err := vtmClient.Do(getAllAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Read Error: %+v", err)
	}
	getChildMonitor := getAllAPI.GetResponse().FilterByName(readName)
	if getChildMonitor.Name != readName {
		d.SetId("")
		return nil
	}
	getSingleMonitorAPI := monitor.NewGetSingleMonitor(getChildMonitor.Name)
	getMonitorProperties := getSingleMonitorAPI.GetResponse()
	err = vtmClient.Do(getSingleMonitorAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Read Error: %+v", err)
	}

	d.Set("name", getChildMonitor.Name)
	d.Set("delay", getMonitorProperties.Properties.Basic.Delay)
	d.Set("timeout", getMonitorProperties.Properties.Basic.Timeout)
	d.Set("failures", getMonitorProperties.Properties.Basic.Failures)
	d.Set("verbose", getMonitorProperties.Properties.Basic.Verbose)
	d.Set("use_ssl", getMonitorProperties.Properties.Basic.UseSSL)
	d.Set("http_host_header", getMonitorProperties.Properties.HTTP.HostHeader)
	d.Set("http_path", getMonitorProperties.Properties.HTTP.URIPath)
	d.Set("http_authentication", getMonitorProperties.Properties.HTTP.Authentication)
	d.Set("http_body_regex", getMonitorProperties.Properties.HTTP.BodyRegex)
	return nil
}

func resourceMonitorUpdate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var readName string
	var updateMonitor monitor.Monitor
	hasChanges := false

	if v, ok := d.GetOk("name"); ok {
		readName = v.(string)
	} else {
		return fmt.Errorf("BrocadeVTM Update Error: name argument required")
	}
	if d.HasChange("delay") {
		if v, ok := d.GetOk("delay"); ok {
			updateMonitor.Properties.Basic.Delay = v.(int)
		}
		hasChanges = true
	}
	if d.HasChange("timeout") {
		if v, ok := d.GetOk("timeout"); ok {
			updateMonitor.Properties.Basic.Timeout = v.(int)
		}
		hasChanges = true
	}
	if d.HasChange("failures") {
		if v, ok := d.GetOk("failures"); ok {
			updateMonitor.Properties.Basic.Failures = v.(int)
		}
		hasChanges = true
	}
	if d.HasChange("verbose") {
		monitorVerbosity := d.Get("verbose").(bool)
		updateMonitor.Properties.Basic.Verbose = &monitorVerbosity
		hasChanges = true
	}
	if d.HasChange("use_ssl") {
		monitorSSL := d.Get("use_ssl").(bool)
		updateMonitor.Properties.Basic.UseSSL = &monitorSSL
		hasChanges = true
	}
	if d.HasChange("http_host_header") {
		if v, ok := d.GetOk("http_host_header"); ok {
			updateMonitor.Properties.HTTP.HostHeader = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("http_path") {
		if v, ok := d.GetOk("http_path"); ok {
			updateMonitor.Properties.HTTP.URIPath = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("http_authentication") {
		if v, ok := d.GetOk("http_authentication"); ok {
			updateMonitor.Properties.HTTP.Authentication = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("http_body_regex") {
		if v, ok := d.GetOk("http_body_regex"); ok {
			updateMonitor.Properties.HTTP.BodyRegex = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		updateAPI := monitor.NewUpdate(readName, updateMonitor)
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Update Error: %+v", err)
		}

		if updateAPI.StatusCode() != 201 && updateAPI.StatusCode() != 200 {
			return fmt.Errorf("BrocadeVTM Update Error: Invalid HTTP response code %+v returned. Response object was %+v", updateAPI.StatusCode(), updateAPI.ResponseObject())
		}
	}
	return resourceMonitorRead(d, m)
}

func resourceMonitorDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var readName string

	if v, ok := d.GetOk("name"); ok {
		readName = v.(string)
	} else {
		return fmt.Errorf("BrocadeVTM Delete Error: name argument required")
	}

	getAllAPI := monitor.NewGetSingleMonitor(readName)
	err := vtmClient.Do(getAllAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Delete: Error fetching monitor %s", readName)
	}
	if getAllAPI.StatusCode() == 404 {
		d.SetId("")
		return nil
	}

	deleteAPI := monitor.NewDelete(readName)
	err = vtmClient.Do(deleteAPI)
	if err != nil || deleteAPI.StatusCode() != 204 {
		return fmt.Errorf("BrocadeVTM Delete: Error deleting monitor %s. Return code != 204. Error: %+v", readName, err)
	}

	d.SetId("")
	return nil
}
