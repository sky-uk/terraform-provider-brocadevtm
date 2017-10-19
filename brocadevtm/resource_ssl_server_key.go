package brocadevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func resourceSSLServerKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSLServerKeyCreate,
		Read:   resourceSSLServerKeyRead,
		Update: resourceSSLServerKeyUpdate,
		Delete: resourceSSLServerKeyDelete,

		Schema: util.SchemaSSLKey(),
	}
}

func resourceSSLServerKeyCreate(d *schema.ResourceData, meta interface{}) error {
	err := util.SSLKeyCreate(d, meta, "ssl/server_keys")
	if err != nil {
		return err
	}
	return resourceSSLServerKeyRead(d, meta)
}

func resourceSSLServerKeyRead(d *schema.ResourceData, meta interface{}) error {
	err := util.SSLKeyRead(d, meta, "ssl/server_keys")
	if err != nil {
		return err
	}
	return nil
}

func resourceSSLServerKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	err := util.SSLKeyUpdate(d, meta, "ssl/server_keys")
	if err != nil {
		return err
	}
	return resourceSSLServerKeyRead(d, meta)
}

func resourceSSLServerKeyDelete(d *schema.ResourceData, meta interface{}) error {
	err := util.SSLKeyDelete(d, meta, "ssl/server_keys")
	if err != nil {
		return err
	}
	return nil
}
