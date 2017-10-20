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

func resourceSSLServerKeyCreate(d *schema.ResourceData, m interface{}) error {
	err := util.SSLKeyCreate(d, m, "ssl/server_keys")
	if err != nil {
		return err
	}
	return resourceSSLServerKeyRead(d, m)
}

func resourceSSLServerKeyRead(d *schema.ResourceData, m interface{}) error {
	err := util.SSLKeyRead(d, m, "ssl/server_keys")
	if err != nil {
		return err
	}
	return nil
}

func resourceSSLServerKeyUpdate(d *schema.ResourceData, m interface{}) error {
	err := util.SSLKeyUpdate(d, m, "ssl/server_keys")
	if err != nil {
		return err
	}
	return resourceSSLServerKeyRead(d, m)
}

func resourceSSLServerKeyDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("ssl/server_keys", d, m)
}
