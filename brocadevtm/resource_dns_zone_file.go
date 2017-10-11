package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"log"
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

	var name, dns_zone_config string

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	if v, ok := d.GetOk("dns_zone_config"); ok {
		dns_zone_config = v.(string)
	}

	err := client.Set("dns_server/zone_files", name, []byte(dns_zone_config), nil)
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
	var zone_config []byte
	client.WorkWithConfigurationResources()
	err := client.GetByName("dns_server/zone_files", name, &zone_config)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("BrocadeVTM DNS zone file error whilst reading %s: %v", name, err)
	}
	log.Println("Going to set dns_zone_config to: ", string(zone_config))
	d.Set("dns_zone_config", string(zone_config))
	return nil
}

func resourceDNSZoneFileUpdate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["octetClient"].(*api.Client)

	hasChanges := false
	name := d.Id()
	var zone_config string

	if d.HasChange("dns_zone_config") {
		if v, ok := d.GetOk("dns_zone_config"); ok {
			zone_config = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		err := client.Set("dns_server/zone_files", name, []byte(zone_config), nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM DNS Zone File error whilst updating %s: %v", name, err)
		}
		d.Set("dns_zone_config", zone_config)
	}
	return resourceDNSZoneFileRead(d, m)
}

func resourceDNSZoneFileDelete(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	name := d.Id()
	err := client.Delete("dns_server/zone_files", name)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("BrocadeVTM DNS zone file error whilst deleting %s: %v", name, err)
	}
	d.SetId("")
	return nil
}
