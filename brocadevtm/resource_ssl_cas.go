package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

func resourceSSLCasFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSLCasCreate,
		Read:   resourceSSLCasRead,
		Update: resourceSSLCasUpdate,
		Delete: resourceSSLCasDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the SSL Trusted Certificate",
				Required:    true,
				ForceNew:    true,
			},
			"ssl_cas_config": {
				Type:        schema.TypeString,
				Description: "The CA including CRL",
				Required:    true,
			},
		},
	}
}

func resourceSSLCasCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["octetClient"].(*api.Client)

	var name, sslCasConfig string

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	if v, ok := d.GetOk("ssl_cas_config"); ok {
		sslCasConfig = v.(string)
	}

	err := client.Set("ssl/cas", name, []byte(sslCasConfig), nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM SSL cas config file error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceSSLCasRead(d, m)
}

func resourceSSLCasRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["octetClient"].(*api.Client)

	name := d.Id()
	sslCasConfig := new([]byte)
	client.WorkWithConfigurationResources()

	err := client.GetByName("ssl/cas", name, sslCasConfig)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM SSL cas config file error whilst reading %s: %v", name, err)
	}
	d.Set("ssl_cas_config", string(*sslCasConfig))

	return nil
}

func resourceSSLCasUpdate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["octetClient"].(*api.Client)

	name := d.Id()
	var sslCasConfig string

	if d.HasChange("ssl_cas_config") {
		if v, ok := d.GetOk("ssl_cas_config"); ok {
			sslCasConfig = v.(string)
		}
	}

	err := client.Set("ssl/cas", name, []byte(sslCasConfig), nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM SSL cas config file error whilst updating %s: %v", name, err)
	}
	d.Set("ssl_cas_config", sslCasConfig)
	return resourceSSLCasRead(d, m)
}

func resourceSSLCasDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("ssl/cas", d, m)
}
