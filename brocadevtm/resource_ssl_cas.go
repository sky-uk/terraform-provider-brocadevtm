package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"log"
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
		return fmt.Errorf("BrocadeVTM SSL Cas File error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceSSLCasRead(d, m)
}

func resourceSSLCasRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	name := d.Id()
	sslCasConfig := new([]byte)
	client.WorkWithConfigurationResources()

	err := client.GetByName("ssl/cas", name, sslCasConfig)
	if client.StatusCode == http.StatusNoContent {
		d.SetId("")
		log.Printf("BrocadeVTM SSL Cas config file %s not found", name)
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM SSL Cas config file error whilst reading %s: %v", name, err)
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
		return fmt.Errorf("BrocadeVTM SSL Cas Config File error whilst updating %s: %v", name, err)
	}
	d.Set("ssl_cas_config", sslCasConfig)

	return resourceSSLCasRead(d, m)
}

func resourceSSLCasDelete(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	name := d.Id()
	err := client.Delete("ssl/cas", name)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM SSL Cas Config file error whilst deleting %s: %v", name, err)
	}
	return nil
}
