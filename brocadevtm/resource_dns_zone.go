package brocadevtm

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
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

	var name string
	res := make(map[string]interface{})
	prop := make(map[string]interface{})
	basic := make(map[string]interface{})

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	if v, ok := d.GetOk("origin"); ok && v != "" {
		basic["origin"] = v.(string)
	}
	if v, ok := d.GetOk("zone_file"); ok && v != "" {
		basic["zonefile"] = v.(string)
	}

	prop["basic"] = basic
	res["properties"] = prop
	err := client.Set("dns_server/zones", name, res, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS zone error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceDNSZoneRead(d, m)
}

func resourceDNSZoneRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	name := d.Id()

	res := make(map[string]interface{})

	client.WorkWithConfigurationResources()
	err := client.GetByName("dns_server/zones", name, &res)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		d.SetId("")
		return fmt.Errorf("BrocadeVTM DNS zone error whilst reading %s: %v", name, err)
	}
	d.SetId(name)
	props := res["properties"].(map[string]interface{})
	basic := props["basic"].(map[string]interface{})
	d.Set("origin", basic["origin"])
	d.Set("zone_file", basic["zonefile"])

	return nil
}

func resourceDNSZoneUpdate(d *schema.ResourceData, m interface{}) error {

	name := d.Id()
	res := make(map[string]interface{})
	props := make(map[string]interface{})
	basic := make(map[string]interface{})

	for _, key := range []string{"origin", "zone_file"} {
		if d.HasChange(key) {
			basic[key] = d.Get(key).(string)
		}
	}

	props["basic"] = basic
	res["properties"] = props

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	err := client.Set("dns_server/zones", name, &res, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS zone error whilst updating %s: %v", name, err)
	}

	return resourceDNSZoneRead(d, m)
}

func resourceDNSZoneDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("dns_server/zones", d, m)
}
