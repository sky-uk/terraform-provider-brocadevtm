package brocadevtm

import (
	"fmt"

	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func resourceUserAuthenticator() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAuthenticatorSet,
		Update: resourceUserAuthenticatorSet,
		Read:   resourceUserAuthenticatorRead,
		Delete: resourceUserAuthenticatorDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ldap", "radius", "tacacs_plus"}, false),
			},
			"ldap": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_dn": {
							Type:        schema.TypeString,
							Description: "The base DN (Distinguished Name) under which directory searches will be applied. ",
							Optional:    true,
						},
						"bind_dn": {
							Type:        schema.TypeString,
							Description: "Template to construct the bind DN (Distinguished Name) from the username. ",
							Optional:    true,
						},
						"dn_method": {
							Type:         schema.TypeString,
							Description:  "FQDN of the member pair",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"construct", "none", "search"}, false),
						},
						"fallback_group": {
							Type:        schema.TypeString,
							Description: "If the group attribute is not defined, or returns no results for the user logging in, the group named here will be used.",
							Optional:    true,
						},
						"filter": {
							Type:        schema.TypeString,
							Description: "A filter that can be used to extract a unique user record located under the base DN (Distinguished Name).",
							Optional:    true,
						},
						"group_attribute": {
							Type:        schema.TypeString,
							Description: "The LDAP attribute that gives a user's group. If there are multiple entries for the attribute all will be extracted and they'll be lexicographically sorted, then the first one to match a Permission Group name will be used.",
							Optional:    true,
						},
						"group_field": {
							Type:        schema.TypeString,
							Description: "The sub-field of the group attribute that gives a user's group.",
							Optional:    true,
						},
						"group_filter": {
							Type:        schema.TypeString,
							Description: "If the user record returned by filter does not contain the required group information you may specify an alternative group search filter here.",
							Optional:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "The port to connect to the LDAP server on.",
							Optional:    true,
							Default:     389,
						},
						"search_dn": {
							Type:        schema.TypeString,
							Description: "The bind DN (Distinguished Name) to use when searching the directory for a user's bind DN. ",
							Optional:    true,
						},
						"search_password": {
							Type:        schema.TypeString,
							Description: "If binding to the LDAP server using search_dn requires a password, enter it here.",
							Sensitive:   true,
							Optional:    true,
						},
						"server": {
							Type:        schema.TypeString,
							Description: "The IP or hostname of the LDAP server.",
							Optional:    true,
						},
						"timeout": {
							Type:        schema.TypeInt,
							Description: "Connection timeout in seconds.",
							Optional:    true,
							Default:     30,
						},
					},
				},
			},
			"radius": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fallback_group": {
							Type:        schema.TypeString,
							Description: "If no group is found using the vendor and group identifiers, or the group found is not valid, the group specified here will be used.",
							Optional:    true,
						},
						"group_attribute": {
							Type:        schema.TypeInt,
							Description: "The RADIUS identifier for the attribute that specifies an account's group.",
							Optional:    true,
							Default:     1,
						},
						"group_vendor": {
							Type:        schema.TypeInt,
							Description: "The RADIUS identifier for the vendor of the RADIUS attribute that specifies an account's group.",
							Optional:    true,
							Default:     7146,
						},
						"nas_identifier": {
							Type:        schema.TypeString,
							Description: "This value is sent to the RADIUS server",
							Optional:    true,
						},
						"nas_ip_address": {
							Type:        schema.TypeString,
							Description: "This value is sent to the RADIUS server, if left blank the address of the interface used to connect to the server will be used.",
							Optional:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "The port to connect to the RADIUS server on.",
							Optional:    true,
							Default:     1812,
						},
						"secret": {
							Type:        schema.TypeString,
							Description: "Secret key shared with the RADIUS server.",
							Sensitive:   true,
							Optional:    true,
						},
						"server": {
							Type:        schema.TypeString,
							Description: "The IP or hostname of the RADIUS server.",
							Optional:    true,
						},
						"timeout": {
							Type:        schema.TypeInt,
							Description: "Connection timeout in seconds.",
							Optional:    true,
							Default:     30,
						},
					},
				},
			},
			"tacacs_plus": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_type": {
							Type:         schema.TypeString,
							Description:  "Authentication type to use",
							Optional:     true,
							Default:      "pap",
							ValidateFunc: validation.StringInSlice([]string{"ascii", "pap"}, false),
						},
						"fallback_group": {
							Type:        schema.TypeString,
							Description: "If group_service is not used, or no group value is provided for the user by the TACACS+ server, the group specified here will be used.",
							Optional:    true,
						},
						"group_field": {
							Type:        schema.TypeString,
							Description: "The TACACS+ 'service' field that provides each user's group.",
							Optional:    true,
							Default:     "permission-group",
						},
						"group_service": {
							Type:        schema.TypeString,
							Description: "The TACACS+ 'service' that provides each user's group field.",
							Optional:    true,
							Default:     "zeus",
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "The port to connect to the TACACS+ server on.",
							Optional:    true,
							Default:     49,
						},
						"secret": {
							Type:        schema.TypeString,
							Description: "Secret key shared with the TACACS+ server.",
							Sensitive:   true,
							Optional:    true,
						},
						"server": {
							Type:        schema.TypeString,
							Description: " The IP or hostname of the TACACS+ server.",
							Optional:    true,
						},
						"timeout": {
							Type:        schema.TypeInt,
							Description: "Connection timeout in seconds.",
							Optional:    true,
							Default:     30,
						},
					},
				},
			},
		},
	}
}

func resourceUserAuthenticatorSet(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	res := make(map[string]interface{})
	props := make(map[string]interface{})
	basic := make(map[string]interface{})

	name := d.Get("name").(string)

	util.AddChangedSimpleAttributesToMap(d, basic, "", []string{
		"description",
		"enabled",
		"type",
	})
	props["basic"] = basic

	if d.HasChange("ldap") {
		props["ldap"] = d.Get("ldap").([]interface{})[0]
	}
	if d.HasChange("radius") {
		props["radius"] = d.Get("radius").([]interface{})[0]
	}
	if d.HasChange("tacacs_plus") {
		props["tacacs_plus"] = d.Get("tacacs_plus").([]interface{})[0]
	}
	res["properties"] = props

	err := client.Set("user_authenticators", name, res, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst creating user authenticator %s: %v", name, err)
	}
	d.SetId(name)
	return resourceUserAuthenticatorRead(d, m)
}

func resourceUserAuthenticatorRead(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	res := make(map[string]interface{})

	client.WorkWithConfigurationResources()
	err := client.GetByName("user_authenticators", d.Id(), &res)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst retrieving user authenticator %s: %v", d.Id(), err)
	}

	props := res["properties"].(map[string]interface{})
	basic := props["basic"].(map[string]interface{})

	d.Set("description", basic["description"])
	d.Set("enabled", basic["enabled"])
	d.Set("type", basic["type"])
	d.Set("ldap", props["ldap"])
	d.Set("radius", props["radius"])
	d.Set("tacacs_plus", props["tacacs_plus"])
	return nil
}

func resourceUserAuthenticatorDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("user_authenticators", d, m)
}
