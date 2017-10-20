package brocadevtm

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/user_group"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Update: resourceUserGroupUpdate,
		Delete: resourceUserGroupDelete,

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

			"password_expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_level": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateAccessLevel,
						},
					},
				},
			},
		},
	}
}

// ValidateAccessLevel : Validates that the access level entered is correct
func ValidateAccessLevel(v interface{}, k string) (ws []string, errors []error) {
	switch strings.ToLower(v.(string)) {
	case
		"none",
		"ro",
		"full":
		return
	}
	errors = append(errors, fmt.Errorf("Access level must be one of NONE, RO or FULL"))
	return
}

func buildPermissionsObject(permissions *schema.Set) []userGroup.Permission {
	permissionValues := []userGroup.Permission{}
	for _, permission := range permissions.List() {
		permissionObject := permission.(map[string]interface{})
		newPermission := userGroup.Permission{}
		if conigurationElement, ok := permissionObject["name"].(string); ok {
			newPermission.Name = conigurationElement
		}
		if accessLevel, ok := permissionObject["access_level"].(string); ok {
			newPermission.AccessLevel = strings.ToLower(accessLevel)
		}
		permissionValues = append(permissionValues, newPermission)
	}
	return permissionValues
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {
	var userGroup userGroup.UserGroup

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	userGroupName := d.Get("name").(string)

	if v, ok := d.GetOk("description"); ok && v != "" {
		userGroup.Properties.Basic.Description = v.(string)
	}
	if v, ok := d.GetOk("password_expire_time"); ok {
		userGroup.Properties.Basic.PasswordExpireTime = uint(v.(int))
	}
	if v, ok := d.GetOk("timeout"); ok {
		userGroup.Properties.Basic.Timeout = uint(v.(int))
	}

	if v, ok := d.GetOk("permissions"); ok {
		if permissions, ok := v.(*schema.Set); ok {
			userGroup.Properties.Basic.Permissions = buildPermissionsObject(permissions)
		}
	}

	err := client.Set("user_groups", userGroupName, &userGroup, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM User Group error whilst creating %s: %v", userGroupName, err)
	}

	d.SetId(userGroupName)

	return resourceUserGroupRead(d, m)
}

func resourceUserGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	var userGroupObject userGroup.UserGroup

	err := client.GetByName("user_groups", d.Id(), &userGroupObject)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("BrocadeVTM User Group error whilst retrieving %s: %v", d.Id(), err)
	}

	d.Set("description", userGroupObject.Properties.Basic.Description)
	d.Set("password_expire_time", userGroupObject.Properties.Basic.PasswordExpireTime)
	d.Set("timeout", userGroupObject.Properties.Basic.Timeout)
	d.Set("permissions", userGroupObject.Properties.Basic.Permissions)

	return nil
}

func resourceUserGroupUpdate(d *schema.ResourceData, m interface{}) error {
	var updatedUserGroup userGroup.UserGroup
	hasChanges := false

	if d.HasChange("description") {
		updatedUserGroup.Properties.Basic.Description = d.Get("description").(string)
		hasChanges = true
	}
	if d.HasChange("password_expire_time") {
		updatedUserGroup.Properties.Basic.PasswordExpireTime = uint(d.Get("password_expire_time").(int))
		hasChanges = true
	}
	if d.HasChange("timeout") {
		updatedUserGroup.Properties.Basic.Timeout = uint(d.Get("timeout").(int))
		hasChanges = true
	}
	if d.HasChange("permissions") {
		updatedUserGroup.Properties.Basic.Permissions = buildPermissionsObject(d.Get("permissions").(*schema.Set))
		hasChanges = true
	}
	if hasChanges {
		config := m.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		err := client.Set("user_groups", d.Id(), updatedUserGroup, nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM User Group error whilst updating %s: %v", d.Id(), err)
		}
	}
	return resourceUserGroupRead(d, m)
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("user_groups", d, m)
}
