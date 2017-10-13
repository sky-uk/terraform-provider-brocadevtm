package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/ssl_server_key"
	"net/http"
)

func resourceSSLServerKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSLServerKeyCreate,
		Read:   resourceSSLServerKeyRead,
		Update: resourceSSLServerKeyUpdate,
		Delete: resourceSSLServerKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"note": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"private": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"public": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"request": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSSLServerKeyCreate(d *schema.ResourceData, meta interface{}) error {
	var name string
	var payloadObject sslServerKey.SSLServerKey

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("note"); ok {
		payloadObject.Properties.Basic.Note = v.(string)
	}
	if v, ok := d.GetOk("private"); ok {
		payloadObject.Properties.Basic.Private = v.(string)
	}
	if v, ok := d.GetOk("public"); ok {
		payloadObject.Properties.Basic.Public = v.(string)
	}
	if v, ok := d.GetOk("request"); ok {
		payloadObject.Properties.Basic.Request = v.(string)
	}

	config := meta.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	err := client.Set("ssl/server_keys", name, &payloadObject, nil)

	if err != nil {
		return fmt.Errorf("BrocadeVTM SSL Server Key error whilst updating %s: %v", name, err)
	}
	d.SetId(name)
	return resourceSSLServerKeyRead(d, meta)
}

func resourceSSLServerKeyRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	name := d.Id()
	var sslServerKey sslServerKey.SSLServerKey
	err := client.GetByName("ssl/server_keys", name, &sslServerKey)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM SSL Server Key error whilst retrieving %s: %v", name, err)
	}

	d.Set("note", sslServerKey.Properties.Basic.Note)
	// TODO: API doesn't return the private key back, so we ignore it,
	// otherwise plan is always changing it.
	// d.Set("private", sslServerKey.Properties.Basic.Private)
	d.Set("public", sslServerKey.Properties.Basic.Public)
	d.Set("request", sslServerKey.Properties.Basic.Request)

	return nil
}

func resourceSSLServerKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	var name string
	var payloadObject sslServerKey.SSLServerKey
	hasChanges := false

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok {
			payloadObject.Properties.Basic.Note = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("private") {
		if v, ok := d.GetOk("private"); ok {
			payloadObject.Properties.Basic.Private = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("public") {
		if v, ok := d.GetOk("public"); ok {
			payloadObject.Properties.Basic.Public = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("request") {
		if v, ok := d.GetOk("request"); ok {
			payloadObject.Properties.Basic.Request = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {
		config := meta.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		err := client.Set("ssl/server_keys", name, &payloadObject, nil)

		if err != nil {
			return fmt.Errorf("BrocadeVTM SSL Server Key error whilst updating %s: %v", name, err)
		}
		d.SetId(name)
	}
	return resourceSSLServerKeyRead(d, meta)
}

func resourceSSLServerKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	name := d.Id()

	err := client.Delete("ssl/server_keys", name)
	if client.StatusCode == http.StatusNoContent || client.StatusCode == http.StatusNotFound {
		return nil
	}
	return fmt.Errorf("BrocadeVTM SSL Server Key error whilst deleting %s: %v", name, err)
}
