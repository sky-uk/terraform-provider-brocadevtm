package brocadevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

func resourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualServerCreate,
		Read:   resourceVirtualServerRead,
		Update: resourceVirtualServerUpdate,
		Delete: resourceVirtualServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Description: "Name of the virtual server",
				Required: true,
			},
			"listen_traffic_ips": {
				Type:        schema.TypeList,
				Description: "List of traffic IPs to listen on",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"pool": {
				Type: schema.TypeString,
				Description: "Name of the pool to use with the virtual server",
				Required: true,
			},
			"port":{
				Type: schema.TypeInt,
				Description: "Port the virtual server should listen on",
				Required: true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"protocol": {
				Type: schema.TypeString,
				Description: "Protocol to use with the virtual server",
				Required: true,
			},
			"request_rules": {
				Type: schema.TypeList,
				Description: "A list of request rules",
				Optional: true,
			},
			"ssl_decrypt": {
				Type: schema.TypeBool,
				Description: "Whether to enable or disable SSL",
				Optional: true,
			},
			"connection_keepalive": {
				Type: schema.TypeBool,
				Description: "Whether to enable keepalive",
				Optional: true,
			},
			"connection_keepalive_timeout": {
				Type: schema.TypeInt,
				Description: "Keepalive timeout",
				Optional: true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
		},
	}
}

func validateVirtualServerUnsignedInteger(v interface{}, k string) (ws []string, errors []error) {
	ttl := v.(int)
	if ttl < 0 {
		errors = append(errors, fmt.Errorf("%q can't be negative", k))
	}
	return
}

func resourceVirtualServerCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVirtualServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVirtualServerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVirtualServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
