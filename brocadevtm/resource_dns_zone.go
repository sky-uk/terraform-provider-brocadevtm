package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/dns_zone"
	"github.com/sky-uk/go-rest-api"
	"log"
	"net/http"
)

func resourceDNSZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSZoneCreate,
		Read:   resourceDNSZoneRead,
		Update: resourceDNSZoneUpdate,
		Delete: resourceDNSZoneDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the DNS zone",
			},
			"origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain origin for the zone",
			},
			"zone_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the DNS zone file to use",
			},
		},
	}
}

func resourceDNSZoneCreate(d *schema.ResourceData, m interface{}) error {

	var dnsZoneName string
	var dnsZoneObject dnsZone.DNSZone
	vtmClient := m.(*rest.Client)

	if v, ok := d.GetOk("name"); ok && v != "" {
		dnsZoneName = v.(string)
	}
	if v, ok := d.GetOk("origin"); ok && v != "" {
		dnsZoneObject.Properties.Basic.Origin = v.(string)
	}
	if v, ok := d.GetOk("zone_file"); ok && v != "" {
		dnsZoneObject.Properties.Basic.ZoneFile = v.(string)
	}

	createAPI := dnsZone.NewCreate(dnsZoneName, dnsZoneObject)

	log.Printf(fmt.Sprintf("[DEBUG] Object is %+v", dnsZoneObject))

	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS zone error whilst creating %s: %v", dnsZoneName, err)
	}

	d.SetId(dnsZoneName)
	return resourceDNSZoneRead(d, m)
}

func resourceDNSZoneRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	dnsZoneName := d.Id()
	var dnsZoneObject dnsZone.DNSZone

	getAPI := dnsZone.NewGet(dnsZoneName)
	err := vtmClient.Do(getAPI)
	if getAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS zone error whilst reading %s: %v", dnsZoneName, err)
	}
	dnsZoneObject = *getAPI.ResponseObject().(*dnsZone.DNSZone)
	d.SetId(dnsZoneName)
	d.Set("origin", dnsZoneObject.Properties.Basic.Origin)
	d.Set("zone_file", dnsZoneObject.Properties.Basic.ZoneFile)

	return nil
}

func resourceDNSZoneUpdate(d *schema.ResourceData, m interface{}) error {

	hasChanges := false
	dnsZoneName := d.Id()
	var dnsZoneObject dnsZone.DNSZone

	if d.HasChange("origin") {
		if v, ok := d.GetOk("origin"); ok && v != "" {
			dnsZoneObject.Properties.Basic.Origin = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("zone_file") {
		if v, ok := d.GetOk("zone_file"); ok && v != "" {
			dnsZoneObject.Properties.Basic.ZoneFile = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		vtmClient := m.(*rest.Client)
		updateAPI := dnsZone.NewUpdate(dnsZoneName, dnsZoneObject)
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("BrocadeVTM DNS zone error whilst updating %s: %v", dnsZoneName, err))
		}
	}

	return resourceDNSZoneRead(d, m)
}

func resourceDNSZoneDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	dnsZoneName := d.Id()

	deleteAPI := dnsZone.NewDelete(dnsZoneName)
	err := vtmClient.Do(deleteAPI)
	if deleteAPI.StatusCode() != http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("BrocadeVTM DNS zone error whilst deleting %s: %v", dnsZoneName, err))
	}

	d.SetId("")

	return nil
}
