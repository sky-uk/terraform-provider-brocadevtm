package brocadevtm

import (
	"fmt"
        "log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

func resourceDNSZoneFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSZoneFileCreate,
		Read:   resourceDNSZoneFileRead,
		Update: resourceDNSZoneFileUpdate,
		Delete: resourceDNSZoneFileDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the DNS zone",
				Required:    true,
				ForceNew:    true,
			},
			"dns_zone_config": {
				Type:        schema.TypeString,
				Description: "DNS zone configuration section",
				Required:    true,
			},
		},
	}
}

func resourceDNSZoneFileCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["octetClient"].(*api.Client)

	var name, dnsZoneConfig string

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	if v, ok := d.GetOk("dns_zone_config"); ok {
		dnsZoneConfig = v.(string)
	}

	err := client.Set("dns_server/zone_files", name, []byte(dnsZoneConfig), nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS Zone File error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceDNSZoneFileRead(d, m)
}

func resourceDNSZoneFileRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	name := d.Id()
	zoneConfig := new([]byte)
	client.WorkWithConfigurationResources()
	err := client.GetByName("dns_server/zone_files", name, zoneConfig)
	if client.StatusCode == http.StatusNoContent {
		d.SetId("")
		log.Printf("BrocadeVTM DNS zone file %s not found", name)
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS zone file error whilst reading %s: %v", name, err)
	}
	d.Set("dns_zone_config", string(*zoneConfig))
	return nil
}

func resourceDNSZoneFileUpdate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["octetClient"].(*api.Client)

	hasChanges := false
	name := d.Id()
	var zoneConfig string

	if d.HasChange("dns_zone_config") {
		if v, ok := d.GetOk("dns_zone_config"); ok {
			zoneConfig = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		err := client.Set("dns_server/zone_files", name, []byte(zoneConfig), nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM DNS Zone File error whilst updating %s: %v", name, err)
		}
		d.Set("dns_zone_config", zoneConfig)
	}
	return resourceDNSZoneFileRead(d, m)
}

func resourceDNSZoneFileDelete(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	name := d.Id()
	err := client.Delete("dns_server/zone_files", name)
	if client.StatusCode == http.StatusNoContent || client.StatusCode == http.StatusNotFound {
		d.SetId("")
                return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS zone file error whilst deleting %s: %v", name, err)
	}
	return nil
}
