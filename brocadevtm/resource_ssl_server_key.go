package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/ssl_server_key"
)

func resourceSSLServerKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSLServerKeyCreate,
		Read:   resourceSSLServerKeyRead,
		Update: resourceSSLServerKeyUpdate,
		Delete: resourceSSLServerKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"private": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"public": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"request": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSSLServerKeyCreate(d *schema.ResourceData, meta interface{}) error {
	vtmClient := meta.(*brocadevtm.VTMClient)
	var name string
	var payloadObject sslServerKey.SSLServerKey

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("BrocadeVTM Create Error: name argument required")
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

	createSSLServerKey := sslServerKey.NewCreate(name, &payloadObject)
	err := vtmClient.Do(createSSLServerKey)
	if err != nil && createSSLServerKey.StatusCode() != 201 {
		d.SetId("")
		return fmt.Errorf("BrocadeVTM Create Error: %+v", err)
	}

	d.SetId(name)
	return resourceSSLServerKeyRead(d, meta)
}

func resourceSSLServerKeyRead(d *schema.ResourceData, meta interface{}) error {
	vtmClient := meta.(*brocadevtm.VTMClient)
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("BrocadeVTM Read Error: name argument required")
	}
	getSSLServerKey := sslServerKey.NewGet(name)
	err := vtmClient.Do(getSSLServerKey)
	if err != nil && getSSLServerKey.StatusCode() != 200 {
		d.SetId("")
		return fmt.Errorf("BrocadeVTM Read Error: %+v", err)
	}

	sslServerKey := getSSLServerKey.GetResponse()
	d.Set("note", sslServerKey.Properties.Basic.Note)
	// TODO: API doesn't return the private key back, so we ignore it,
	// otherwise plan is always changing it.
	// d.Set("private", sslServerKey.Properties.Basic.Private)
	d.Set("public", sslServerKey.Properties.Basic.Public)
	d.Set("request", sslServerKey.Properties.Basic.Request)

	return nil
}

func resourceSSLServerKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	vtmClient := meta.(*brocadevtm.VTMClient)
	var name string
	var payloadObject sslServerKey.SSLServerKey
	hasChanges := false

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("BrocadeVTM Update Error: name argument required")
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
		updateSSLServerKey := sslServerKey.NewUpdate(name, payloadObject)
		err := vtmClient.Do(updateSSLServerKey)
		if err != nil {
			d.SetId("")
			return fmt.Errorf("BrocadeVTM Update Error: %+v", err)
		}
	}
	return resourceSSLServerKeyRead(d, meta)
}

func resourceSSLServerKeyDelete(d *schema.ResourceData, meta interface{}) error {
	vtmClient := meta.(*brocadevtm.VTMClient)

	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("BrocadeVTM Delete Error: name argument required")
	}

	getSSLServerKey := sslServerKey.NewGet(name)
	err := vtmClient.Do(getSSLServerKey)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Delete: SSL Server key already doesn't exist, %s", name)
	}
	if getSSLServerKey.StatusCode() == 404 {
		d.SetId("")
		return nil
	}

	deleteAPI := sslServerKey.NewDelete(name)
	err = vtmClient.Do(deleteAPI)
	if err != nil || deleteAPI.StatusCode() != 204 {
		return fmt.Errorf("BrocadeVTM Delete: Error deleting SSLServerKey %s. Return code != 204. Error: %+v", name, err)
	}

	d.SetId("")
	return nil
}
