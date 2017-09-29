package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/dns_zone_file"
	"github.com/sky-uk/go-rest-api"
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
			},
			"dns_zone_file": {
				Type:        schema.TypeString,
				Description: "DNS zone file",
				Required:    true,
			},
		},
	}
}

func resourceDNSZoneFileCreate(d *schema.ResourceData, m interface{}) error {

	var dnsZoneFileObject dnsZoneFile.DNSZoneFile
	headers := make(map[string]string)

	// We need to copy the client as we want to specify different headers for DNS Zone File which will conflict with other resources.
	client := m.(*rest.Client)
	vtmClient := *client
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Transfer-Encoding"] = "text"
	vtmClient.Headers = headers

	if v, ok := d.GetOk("name"); ok && v != "" {
		dnsZoneFileObject.Name = v.(string)
	}
	if v, ok := d.GetOk("dns_zone_file"); ok {
		dnsZoneFileObject.FileName = v.(string)
	}

	createAPI := dnsZoneFile.NewCreate(dnsZoneFileObject.Name, []byte(fmt.Sprintf(dnsZoneFileObject.FileName)))
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM DNS Zone File error whilst creating %s: %v", dnsZoneFileObject.Name, err)
	}

	d.SetId(dnsZoneFileObject.Name)

	return resourceDNSZoneFileRead(d, m)
}

func resourceDNSZoneFileRead(d *schema.ResourceData, m interface{}) error {

	var dnsZoneFileObject dnsZoneFile.DNSZoneFile
	headers := make(map[string]string)

	// We need to copy the client as we want to specify different headers for DNS Zone File which will conflict with other resources.
	client := m.(*rest.Client)
	vtmClient := *client
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Transfer-Encoding"] = "text"
	vtmClient.Headers = headers

	dnsZoneFileObject.Name = d.Id()
	readAPI := dnsZoneFile.NewGet(dnsZoneFileObject.Name)
	err := vtmClient.Do(readAPI)
	if err != nil {
		if readAPI.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM DNS Zone File error whilst retrieving %s: %v", dnsZoneFileObject.Name, err)
	}

	response := readAPI.ResponseObject().(*[]byte)
	dnsZoneFileObject.FileName = string(*response)

	d.SetId(dnsZoneFileObject.Name)
	d.Set("dns_zone_file", dnsZoneFileObject.FileName)

	return nil
}

func resourceDNSZoneFileUpdate(d *schema.ResourceData, m interface{}) error {

	var dnsZoneFileObject dnsZoneFile.DNSZoneFile
	headers := make(map[string]string)
	hasChanges := false
	dnsZoneFileObject.Name = d.Id()

	if d.HasChange("dns_zone_file") {
		if v, ok := d.GetOk("dns_zone_file"); ok {
			dnsZoneFileObject.FileName = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		// We need to copy the client as we want to specify different headers for DNS Zone File which will conflict with other resources.
		client := m.(*rest.Client)
		vtmClient := *client
		headers["Content-Type"] = "application/octet-stream"
		headers["Content-Transfer-Encoding"] = "text"
		vtmClient.Headers = headers

		updateAPI := dnsZoneFile.NewUpdate(dnsZoneFileObject.Name, []byte(fmt.Sprintf(dnsZoneFileObject.FileName)))
		err := vtmClient.Do(updateAPI)

		if err != nil {
			return fmt.Errorf("BrocadeVTM DNS Zone File error whilst updating %s: %vv", dnsZoneFileObject.Name, err)
		}
		d.SetId(dnsZoneFileObject.Name)
		d.Set("dns_zone_file", dnsZoneFileObject.FileName)
	}

	return resourceDNSZoneFileRead(d, m)
}

func resourceDNSZoneFileDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)

	dnsZoneFileObjectName := d.Id()
	deleteAPI := dnsZoneFile.NewDelete(dnsZoneFileObjectName)
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("BrocadeVTM DNS Zone File error whilst deleting %s: %v", dnsZoneFileObjectName, err)
	}

	d.SetId("")
	return nil
}
