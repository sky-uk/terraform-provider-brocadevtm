package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/user_authenticators"
	"github.com/sky-uk/go-rest-api"
	"net/http"
	"strings"
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
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"radius", "tacacs_plus"},
				MaxItems:      1,
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
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"ldap", "tacacs_plus"},
				MaxItems:      1,
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
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"ldap", "radius"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_type": {
							Type:         schema.TypeString,
							Description:  "Authentication type to use",
							Optional:     true,
							Default:      "pap",
							ValidateFunc: validateTACASPlusAuthenticationType,
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
		"tacus+plus":
		return
	}
	errors = append(errors, fmt.Errorf("Access level must be one of ldap, radius or tacas_plus"))
	return
}

// validateTACASPlusAuthenticationType : Validates that the authentication type entered is supported
func validateTACASPlusAuthenticationType(v interface{}, k string) (ws []string, errors []error) {
	switch strings.ToLower(v.(string)) {
	case
		"ascii",
		"pap":
		return
	}
	errors = append(errors, fmt.Errorf("Access level must be one of ldap, radius or tacas_plus"))
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
	var userAuthenticator userauthenticators.UserAuthenticator
	headers := make(map[string]string)

	client := m.(*rest.Client)
	vtmClient := *client
	headers["Content-Type"] = "application/json"
	vtmClient.Headers = headers
	vtmClient.Debug = true

	userAuthenticatorName := d.Get("name").(string)

	if v, ok := d.GetOk("description"); ok {
		userAuthenticator.Properties.Basic.Description = v.(string)
	}
	if v, ok := d.GetOk("enabled"); ok {
		userAuthenticator.Properties.Basic.Enabled = v.(bool)
	}
	if v, ok := d.GetOk("type"); ok {
		userAuthenticator.Properties.Basic.Type = v.(string)
		switch v.(string) {
		case "ldap":
			userAuthenticator.LDAP = assignLDAPValues(d)
		case "radius":
			userAuthenticator.Radius = assignRadiusValues(d)
		case "tacas_plus":
			userAuthenticator.TACACSPlus = assignTACASPlusValues(d)
		}

	}
	createAPI := userauthenticators.NewPut(userAuthenticatorName, userAuthenticator)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst creating user authenticator %s: %v", userAuthenticatorName, string(createAPI.RawResponse()))
	}
	d.SetId(userAuthenticatorName)
	return resourceUserAuthenticatorRead(d, m)
}

func resourceUserAuthenticatorUpdate(d *schema.ResourceData, m interface{}) error {
	var updatedUserAuthenticator userauthenticators.UserAuthenticator
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

	authenticationType := d.Get("type").(string)

	if d.HasChange("type") {
		updatedUserAuthenticator.Basic.Type = authenticationType
	}

	switch authenticationType {
	case "ldap":
		if d.HasChange("ldap") {
			updatedUserAuthenticator.Properties.LDAP = assignLDAPValues(d)
		}
	case "radius":
		if d.HasChange("radius") {
			updatedUserAuthenticator.Properties.Radius = assignRadiusValues(d)
		}
	case "tacas_plus":
		if d.HasChange("tacas_plus") {
			updatedUserAuthenticator.Properties.TACACSPlus = assignTACASPlusValues(d)
		}
	}

	if hasChanges {
		headers := make(map[string]string)
		client := m.(*rest.Client)
		vtmClient := *client
		headers["Content-Type"] = "application/json"
		vtmClient.Headers = headers

		updateAPI := userauthenticators.NewPut(d.Id(), updatedUserAuthenticator)
		err := vtmClient.Do(updateAPI)

		if err != nil {
			return fmt.Errorf("BrocadeVTM error whilst updating user authenticator %s: %vv", d.Id(), err)
		}
		return resourceUserAuthenticatorRead(d, m)
	}

	return nil
}

func resourceUserAuthenticatorRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*rest.Client)
	vtmClient := *client
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	vtmClient.Headers = headers
	readAPI := userauthenticators.NewGet(d.Id())
	err := vtmClient.Do(readAPI)
	if err != nil {
		if readAPI.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM error whilst retrieving user authenticator %s: %v", d.Id(), err)
	}

	returnedUserAuthenticator := readAPI.ResponseObject().(*userauthenticators.UserAuthenticator)

	d.Set("description", returnedUserAuthenticator.Properties.Basic.Description)
	d.Set("enabled", returnedUserAuthenticator.Properties.Basic.Enabled)
	d.Set("type", returnedUserAuthenticator.Properties.Basic.Type)
	//d.Set("ldap", returnedUserAuthenticator.Properties.LDAP)
	//d.Set("radius", returnedUserAuthenticator.Properties.Radius)
	//d.Set("tacas_plus", returnedUserAuthenticator.Properties.TACACSPlus)
	return nil
}

func resourceUserAuthenticatorDelete(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*rest.Client)
	deleteAPI := userauthenticators.NewDelete(d.Id())
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("BrocadeVTM error whilst deleting user authenticator %s: %v", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func assignLDAPValues(d *schema.ResourceData) (ldapStruct userauthenticators.LDAP) {
	if v, ok := d.GetOk("ldap"); ok {
		vL := v.(*schema.Set).List()[0]
		ldapBlock := vL.(map[string]interface{})

		if v, ok := ldapBlock["base_dn"].(string); ok && v != "" {
			ldapStruct.BaseDN = v
		}
		if v, ok := ldapBlock["bind_dn"].(string); ok && v != "" {
			ldapStruct.BindDN = v
		}
		if v, ok := ldapBlock["dn_method"].(string); ok && v != "" {
			ldapStruct.DNMethod = strings.ToLower(v)
		}
		if v, ok := ldapBlock["fallback_group"].(string); ok && v != "" {
			ldapStruct.FallbackGroup = v
		}
		if v, ok := ldapBlock["filter"].(string); ok && v != "" {
			ldapStruct.Filter = v
		}
		if v, ok := ldapBlock["group_attribute"].(string); ok && v != "" {
			ldapStruct.GroupAttribute = v
		}
		if v, ok := ldapBlock["group_field"].(string); ok && v != "" {
			ldapStruct.GroupField = v
		}
		if v, ok := ldapBlock["group_filter"].(string); ok && v != "" {
			ldapStruct.GroupFilter = v
		}
		if v, ok := ldapBlock["port"].(int); ok {
			ldapStruct.Port = uint(v)
		}
		if v, ok := ldapBlock["search_dn"].(string); ok && v != "" {
			ldapStruct.SearchDN = v
		}
		if v, ok := ldapBlock["search_password"].(string); ok && v != "" {
			ldapStruct.SearchPassword = v
		}
		if v, ok := ldapBlock["server"].(string); ok && v != "" {
			ldapStruct.Server = v
		}
		if v, ok := ldapBlock["timeout"].(int); ok {
			ldapStruct.Timeout = uint(v)
		}
	}
	return
}

func assignRadiusValues(d *schema.ResourceData) (radiusStruct userauthenticators.Radius) {
	if v, ok := d.GetOk("radius"); ok {
		vL := v.(*schema.Set).List()[0]
		radiusBlock := vL.(map[string]interface{})

		if v, ok := radiusBlock["fallback_group"].(string); ok {
			radiusStruct.FallbackGroup = v
		}
		if v, ok := radiusBlock["group_attribute"].(int); ok {
			radiusStruct.GroupAttribute = uint(v)
		}
		if v, ok := radiusBlock["group_vendor"].(int); ok {
			radiusStruct.GroupVendor = uint(v)
		}
		if v, ok := radiusBlock["nas_identifier"].(string); ok {
			radiusStruct.NasIdentifier = v
		}
		if v, ok := radiusBlock["nas_ip_address"].(string); ok {
			radiusStruct.NasIPAddress = v
		}
		if v, ok := radiusBlock["port"].(int); ok {
			radiusStruct.Port = uint(v)
		}
		if v, ok := radiusBlock["secret"].(string); ok {
			radiusStruct.Secret = v
		}
		if v, ok := radiusBlock["server"].(string); ok {
			radiusStruct.Server = v
		}
		if v, ok := radiusBlock["timeout"].(int); ok {
			radiusStruct.Timeout = uint(v)
		}
	}

	return
}

func assignTACASPlusValues(d *schema.ResourceData) (tacasPlusStruct userauthenticators.TACACSPlus) {
	if v, ok := d.GetOk("tacas_plus"); ok {
		vL := v.(*schema.Set).List()[0]
		tacasPlusBlock := vL.(map[string]interface{})

		if v, ok := tacasPlusBlock["auth_type"].(string); ok {
			tacasPlusStruct.AuthType = v
		}
		if v, ok := tacasPlusBlock["fallback_group"].(string); ok {
			tacasPlusStruct.FallbackGroup = v
		}
		if v, ok := tacasPlusBlock["group_field"].(string); ok {
			tacasPlusStruct.GroupField = v
		}
		if v, ok := tacasPlusBlock["group_service"].(string); ok {
			tacasPlusStruct.GroupService = v
		}
		if v, ok := tacasPlusBlock["port"].(int); ok {
			tacasPlusStruct.Port = uint(v)
		}
		if v, ok := tacasPlusBlock["secret"].(string); ok {
			tacasPlusStruct.Secret = v
		}
		if v, ok := tacasPlusBlock["server"].(string); ok {
			tacasPlusStruct.Server = v
		}
		if v, ok := tacasPlusBlock["timeout"].(int); ok {
			tacasPlusStruct.Port = uint(v)
		}
	}
	return
}
