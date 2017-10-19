package brocadevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func resourceSSLClientKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSLClientKeyCreate,
		Read:   resourceSSLClientKeyRead,
		Update: resourceSSLClientKeyUpdate,
		Delete: resourceSSLClientKeyDelete,

		Schema: util.SchemaSSLKey(),
	}
}

func resourceSSLClientKeyCreate(d *schema.ResourceData, meta interface{}) error {
	err := util.SSLKeyCreate(d, meta, "ssl/client_keys")
	if err != nil {
		return err
	}
	return resourceSSLServerKeyRead(d, meta)
}

func resourceSSLClientKeyRead(d *schema.ResourceData, meta interface{}) error {
	err := util.SSLKeyRead(d, meta, "ssl/client_keys")
	if err != nil {
		return err
	}
	return nil
}

func resourceSSLClientKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	err := util.SSLKeyUpdate(d, meta, "ssl/client_keys")
	if err != nil {
		return err
	}
	return nil
}

func resourceSSLClientKeyDelete(d *schema.ResourceData, meta interface{}) error {
	err := util.SSLKeyDelete(d, meta, "ssl/client_keys")
	if err != nil {
		return err
	}
	return nil
}
