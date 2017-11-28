package brocadevtm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func resourceApplianceNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplianceNatCreate,
		Read:   resourceApplianceNatRead,
		Update: resourceApplianceNatUpdate,
		Delete: resourceApplianceNatDelete,

		Schema: map[string]*schema.Schema{
			"many_to_one_all_ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_number": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A unique rule identifier",
						},
						"pool": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Pool of a \"many to one overload\" type NAT rule.",
						},
						"tip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TIP Group of a \"many to one overload\" type NAT rule.",
						},
					},
				},
			},
			"many_to_one_port_locked": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_number": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A unique rule identifier",
						},
						"pool": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Pool of a \"many to one port locked\" type NAT rule.",
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Port number of a \"many to one port locked\" type NAT rule.",
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol allowed for a \"many to one port locked\" type NAT rule.",
							ValidateFunc: validation.StringInSlice([]string{
								"icmp",
								"sctp",
								"tcp",
								"udp",
								"udplite",
							}, false),
						},
						"tip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TIP Group of a \"many to one port locked\" type NAT rule.",
						},
					},
				},
			},
			"one_to_one": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_number": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A unique rule identifier",
						},
						"enable_inbound": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enabling the inbound part of a \"one to one\" type NAT rule",
						},
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP address of a \"one-to-one\" type NAT rule.",
						},
						"tip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TIP Group of a \"one-to-one \" type NAT rule.",
						},
					},
				},
			},
			"port_mapping": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_number": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A unique rule identifier",
						},
						"dport_first": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "First port of the dest. port range",
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"dport_last": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Last port of the dest. port range",
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"virtual_server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target Virtual Server of a \"port mapping\" rule",
						},
					},
				},
			},
		},
	}
}

func resourceApplianceNatCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	natResource := make(map[string]interface{})
	properties := make(map[string]interface{})
	basic := make(map[string]interface{})

	keys := []string{"many_to_one_all_ports", "many_to_one_port_locked", "one_to_one", "port_mapping"}
	for _, key := range keys {
		if v, ok := d.GetOk(key); ok {
			basic[key] = v.(*schema.Set).List()
		}
	}
	properties["basic"] = basic
	natResource["properties"] = properties

	err := client.Set("appliance/nat", "", natResource, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM Appliance/Nat error whilst creating: %s", err)
	}
	d.SetId("appliance_nat")
	return resourceApplianceNatRead(d, m)
}

func resourceApplianceNatRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()

	natResource := make(map[string]map[string]interface{})
	err := client.GetByName("appliance/nat", "", &natResource)
	basic := natResource["properties"]["basic"].(map[string]interface{})
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM Appliance/Nat error whilst retrieving: %s", err)
	}

	resource := resourceApplianceNat()
	for key := range resource.Schema {
		d.Set(key, basic[key])
	}
	return nil
}

func resourceApplianceNatUpdate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	natResource := make(map[string]map[string]map[string]interface{})
	natResource["properties"] = make(map[string]map[string]interface{})
	natResource["properties"]["basic"] = make(map[string]interface{})

	keys := []string{"many_to_one_all_ports", "many_to_one_port_locked", "one_to_one", "port_mapping"}
	for _, key := range keys {
		if d.HasChange(key) {
			natResource["properties"]["basic"][key] = make([]interface{}, 0)
			if v, ok := d.GetOk(key); ok {
				natResource["properties"]["basic"][key] = v.(*schema.Set).List()
			}
		}
	}

	err := client.Set("appliance/nat", "", natResource, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM ApplianceNat error whilst creating: %s", err)
	}
	return resourceApplianceNatRead(d, m)
}

// Actually you can't delete the resource on the Brocade server
// what we can do is delete all the NAT rules
// this does not prevent terraform to anyway delete the resource in the
// state file...
func resourceApplianceNatDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	natResource := make(map[string]interface{})
	properties := make(map[string]interface{})
	basic := make(map[string]interface{})

	keys := []string{"many_to_one_all_ports", "many_to_one_port_locked", "one_to_one", "port_mapping"}
	for _, key := range keys {
		basic[key] = make([]interface{}, 0)
	}
	properties["basic"] = basic
	natResource["properties"] = properties
	err := client.Set("appliance/nat", "", natResource, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM ApplianceNat error whilst deleting all NAT rules: %s", err)
	}
	d.SetId("")
	return nil
}
