package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/dns_zone"
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
				Required:    true,
				Description: "The origin",
			},
			"zone_file": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the DNS zone file to use",
			},
		},
	}
}

func resourceDNSZoneCreate(d *schema.ResourceData, m interface{}) error {

	var dnsZoneName string
	var dnsZoneObject dnsZone.DNSZone
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	if v, ok := d.GetOk("name"); ok && v != "" {
		dnsZoneName = v.(string)
	}
	if v, ok := d.GetOk("origin"); ok && v != "" {
		dnsZoneObject.Properties.Basic.Origin = v.(string)
	}
	if v, ok := d.GetOk("zone_file"); ok && v != "" {
		dnsZoneObject.Properties.Basic.ZoneFile = v.(string)
	}

	err := client.Set("dns_server/zones", dnsZoneName, dnsZoneObject, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS zone error whilst creating %s: %v", dnsZoneName, err)
	}

	d.SetId(dnsZoneName)
	return resourceDNSZoneRead(d, m)
}

func resourceDNSZoneRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	dnsZoneName := d.Id()
	var dnsZoneObject dnsZone.DNSZone

	client.WorkWithConfigurationResources()
	err := client.GetByName("dns_server/zones", dnsZoneName, &dnsZoneObject)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		d.SetId("")
		return fmt.Errorf("BrocadeVTM DNS zone error whilst reading %s: %v", dnsZoneName, err)
	}
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
		config := m.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		err := client.Set("dns_server/zones", dnsZoneName, dnsZoneObject, nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM DNS zone error whilst updating %s: %v", dnsZoneName, err)
		}
	}

	return resourceDNSZoneRead(d, m)
}

func resourceDNSZoneDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("dns_server/zones", d, m)
}
