package brocadevtm

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/user_authenticator"
	"net/http"
)

func resourceUserAuthenticator() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAuthenticatorCreate,
		Update: resourceUserAuthenticatorUpdate,
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
				ValidateFunc: validateAuthenticationType,
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
							ValidateFunc: validateDistinguishedNameMethod,
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
							ValidateFunc: validateTACACSPlusAuthenticationType,
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

// validateAuthenticationType : Validates that the authentication type entered is supported
func validateAuthenticationType(v interface{}, k string) (ws []string, errors []error) {
	switch strings.ToLower(v.(string)) {
	case
		"ldap",
		"radius",
		"tacacs_plus":
		return
	}
	errors = append(errors, fmt.Errorf("Access level must be one of ldap, radius or tacacs_plus"))
	return
}

// validateTACACSPlusAuthenticationType : Validates that the authentication type entered is supported
func validateTACACSPlusAuthenticationType(v interface{}, k string) (ws []string, errors []error) {
	switch strings.ToLower(v.(string)) {
	case
		"ascii",
		"pap":
		return
	}
	errors = append(errors, fmt.Errorf("Access level must be one of ascii or pap"))
	return
}

// validateDistinguishedNameMethod : Validates that the Distinguished Name method entered is supported
func validateDistinguishedNameMethod(v interface{}, k string) (ws []string, errors []error) {
	switch strings.ToLower(v.(string)) {
	case
		"construct",
		"none",
		"search":
		return
	}
	errors = append(errors, fmt.Errorf("Access level must be one of construct, none or search"))
	return
}

func resourceUserAuthenticatorCreate(d *schema.ResourceData, m interface{}) error {
	var userAuthenticator userAuthenticator.UserAuthenticator
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	userAuthenticatorName := d.Get("name").(string)

	if v, ok := d.GetOk("description"); ok {
		userAuthenticator.Properties.Basic.Description = v.(string)
	}
	if v, ok := d.GetOk("enabled"); ok {
		userAuthenticator.Properties.Basic.Enabled = v.(bool)
	}
	if v, ok := d.GetOk("type"); ok {
		userAuthenticator.Properties.Basic.Type = v.(string)
	}

	if v, ok := d.GetOk("ldap"); ok {
		ldapList := []map[string]interface{}{}
		for _, ldap := range v.([]interface{}) {
			ldapList = append(ldapList, ldap.(map[string]interface{}))
		}
		userAuthenticator.LDAP = assignLDAPValues(ldapList)
	}

	if v, ok := d.GetOk("radius"); ok {
		radiusList := []map[string]interface{}{}
		for _, radius := range v.([]interface{}) {
			radiusList = append(radiusList, radius.(map[string]interface{}))
		}
		userAuthenticator.Radius = assignRadiusValues(radiusList)
	}

	if v, ok := d.GetOk("tacacs_plus"); ok {
		tacacsPlusList := []map[string]interface{}{}
		for _, tacacsPlus := range v.([]interface{}) {
			tacacsPlusList = append(tacacsPlusList, tacacsPlus.(map[string]interface{}))
		}
		userAuthenticator.TACACSPlus = assignTACACSPlusValues(tacacsPlusList)
	}

	err := client.Set("user_authenticators", userAuthenticatorName, &userAuthenticator, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst creating user authenticator %s: %v", userAuthenticatorName, err)
	}
	d.SetId(userAuthenticatorName)
	return resourceUserAuthenticatorRead(d, m)
}

func resourceUserAuthenticatorUpdate(d *schema.ResourceData, m interface{}) error {
	var updatedUserAuthenticator userAuthenticator.UserAuthenticator
	hasChanges := false

	if d.HasChange("description") {
		updatedUserAuthenticator.Properties.Basic.Description = d.Get("description").(string)
		hasChanges = true
	}

	oldEnabled, newEnabled := d.GetChange("enabled")
	if oldEnabled.(bool) != newEnabled.(bool) {
		updatedUserAuthenticator.Properties.Basic.Enabled = newEnabled.(bool)
		hasChanges = true
	} else {
		updatedUserAuthenticator.Properties.Basic.Enabled = oldEnabled.(bool)
	}

	if d.HasChange("type") {
		updatedUserAuthenticator.Properties.Basic.Type = d.Get("type").(string)
		hasChanges = true
	}

	if d.HasChange("ldap") {
		ldaps := []map[string]interface{}{}
		for _, ldap := range d.Get("ldap").([]interface{}) {
			ldaps = append(ldaps, ldap.(map[string]interface{}))
		}
		updatedUserAuthenticator.LDAP = assignLDAPValues(ldaps)
		hasChanges = true
	}

	if d.HasChange("radius") {
		radiusList := []map[string]interface{}{}
		for _, radius := range d.Get("radius").([]interface{}) {
			radiusList = append(radiusList, radius.(map[string]interface{}))
		}
		updatedUserAuthenticator.Radius = assignRadiusValues(radiusList)
		hasChanges = true
	}

	if d.HasChange("tacacs_plus") {
		tacacsPlusList := []map[string]interface{}{}
		for _, tacacsPlus := range d.Get("tacacs_plus").([]interface{}) {
			tacacsPlusList = append(tacacsPlusList, tacacsPlus.(map[string]interface{}))
		}
		updatedUserAuthenticator.Properties.TACACSPlus = assignTACACSPlusValues(tacacsPlusList)
		hasChanges = true
	}

	if hasChanges {
		config := m.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)

		err := client.Set("user_authenticators", d.Id(), &updatedUserAuthenticator, nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM error whilst updating user authenticator %s: %v", d.Id(), err)
		}
		return resourceUserAuthenticatorRead(d, m)
	}

	return nil
}

func resourceUserAuthenticatorRead(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	var userGroupAuthenticator userAuthenticator.UserAuthenticator

	err := client.GetByName("user_authenticators", d.Id(), &userGroupAuthenticator)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst retrieving user authenticator %s: %v", d.Id(), err)
	}

	d.Set("description", userGroupAuthenticator.Properties.Basic.Description)
	d.Set("enabled", userGroupAuthenticator.Properties.Basic.Enabled)
	d.Set("type", userGroupAuthenticator.Properties.Basic.Type)
	ldapList := []userAuthenticator.LDAP{userGroupAuthenticator.Properties.LDAP}
	d.Set("ldap", ldapList)
	radiusList := []userAuthenticator.Radius{userGroupAuthenticator.Properties.Radius}
	d.Set("radius", radiusList)
	tacacsPlusList := []userAuthenticator.TACACSPlus{userGroupAuthenticator.Properties.TACACSPlus}
	d.Set("tacacs_plus", tacacsPlusList)
	return nil
}

func resourceUserAuthenticatorDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("user_authenticators", d, m)
}

func assignLDAPValues(ldapList []map[string]interface{}) (ldapStruct userAuthenticator.LDAP) {

	if v, ok := ldapList[0]["base_dn"].(string); ok && v != "" {
		ldapStruct.BaseDN = v
	}
	if v, ok := ldapList[0]["bind_dn"].(string); ok && v != "" {
		ldapStruct.BindDN = v
	}
	if v, ok := ldapList[0]["dn_method"].(string); ok && v != "" {
		ldapStruct.DNMethod = strings.ToLower(v)
	}
	if v, ok := ldapList[0]["fallback_group"].(string); ok && v != "" {
		ldapStruct.FallbackGroup = v
	}
	if v, ok := ldapList[0]["filter"].(string); ok && v != "" {
		ldapStruct.Filter = v
	}
	if v, ok := ldapList[0]["group_attribute"].(string); ok && v != "" {
		ldapStruct.GroupAttribute = v
	}
	if v, ok := ldapList[0]["group_field"].(string); ok && v != "" {
		ldapStruct.GroupField = v
	}
	if v, ok := ldapList[0]["group_filter"].(string); ok && v != "" {
		ldapStruct.GroupFilter = v
	}
	if v, ok := ldapList[0]["port"].(int); ok {
		ldapStruct.Port = uint(v)
	}
	if v, ok := ldapList[0]["search_dn"].(string); ok && v != "" {
		ldapStruct.SearchDN = v
	}
	if v, ok := ldapList[0]["search_password"].(string); ok && v != "" {
		ldapStruct.SearchPassword = v
	}
	if v, ok := ldapList[0]["server"].(string); ok && v != "" {
		ldapStruct.Server = v
	}
	if v, ok := ldapList[0]["timeout"].(int); ok {
		ldapStruct.Timeout = uint(v)
	}

	return
}

func assignRadiusValues(radiusList []map[string]interface{}) (radiusStruct userAuthenticator.Radius) {

	if v, ok := radiusList[0]["fallback_group"].(string); ok {
		radiusStruct.FallbackGroup = v
	}
	if v, ok := radiusList[0]["group_attribute"].(int); ok {
		radiusStruct.GroupAttribute = uint(v)
	}
	if v, ok := radiusList[0]["group_vendor"].(int); ok {
		radiusStruct.GroupVendor = uint(v)
	}
	if v, ok := radiusList[0]["nas_identifier"].(string); ok {
		radiusStruct.NasIdentifier = v
	}
	if v, ok := radiusList[0]["nas_ip_address"].(string); ok {
		radiusStruct.NasIPAddress = v
	}
	if v, ok := radiusList[0]["port"].(int); ok {
		radiusStruct.Port = uint(v)
	}
	if v, ok := radiusList[0]["secret"].(string); ok {
		radiusStruct.Secret = v
	}
	if v, ok := radiusList[0]["server"].(string); ok {
		radiusStruct.Server = v
	}
	if v, ok := radiusList[0]["timeout"].(int); ok {
		radiusStruct.Timeout = uint(v)
	}

	return
}

func assignTACACSPlusValues(tacacsPlusList []map[string]interface{}) (tacacsPlusStruct userAuthenticator.TACACSPlus) {

	if v, ok := tacacsPlusList[0]["auth_type"].(string); ok {
		tacacsPlusStruct.AuthType = v
	}
	if v, ok := tacacsPlusList[0]["fallback_group"].(string); ok {
		tacacsPlusStruct.FallbackGroup = v
	}
	if v, ok := tacacsPlusList[0]["group_field"].(string); ok {
		tacacsPlusStruct.GroupField = v
	}
	if v, ok := tacacsPlusList[0]["group_service"].(string); ok {
		tacacsPlusStruct.GroupService = v
	}
	if v, ok := tacacsPlusList[0]["port"].(int); ok {
		tacacsPlusStruct.Port = uint(v)
	}
	if v, ok := tacacsPlusList[0]["secret"].(string); ok {
		tacacsPlusStruct.Secret = v
	}
	if v, ok := tacacsPlusList[0]["server"].(string); ok {
		tacacsPlusStruct.Server = v
	}
	if v, ok := tacacsPlusList[0]["timeout"].(int); ok {
		tacacsPlusStruct.Port = uint(v)
	}

	return
}
