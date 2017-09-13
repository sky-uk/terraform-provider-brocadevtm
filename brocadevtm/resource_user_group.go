package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/user_groups"
	"github.com/sky-uk/go-rest-api"
	"net/http"
	"strings"
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

func ValidateAccessLevel(v interface{}, k string) (ws []string, errors []error) {
	switch strings.ToLower(v.(string)) {
	case
		"none",
		"ro",
		"full":
		return
	}
	errors = append(errors, fmt.Errorf("Access level must be one of NONE, RO or FULL", k))
	return
}

func buildPermissionsObject(permissions *schema.Set) []usergroups.Permission {
	permissionValues := []usergroups.Permission{}
	for _, permission := range permissions.List() {
		permissionObject := permission.(map[string]interface{})
		newPermission := usergroups.Permission{}
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
	var userGroup usergroups.UserGroup
	headers := make(map[string]string)

	client := m.(*rest.Client)
	vtmClient := *client
	headers["Content-Type"] = "application/json"
	vtmClient.Headers = headers
	vtmClient.Debug = true

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

	createAPI := usergroups.NewPut(userGroupName, userGroup)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM User Group error whilst creating %s: %v", userGroupName, err)
	}

	d.SetId(userGroupName)

	return resourceUserGroupRead(d, m)
}

func resourceUserGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*rest.Client)
	vtmClient := *client
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	vtmClient.Headers = headers
	readAPI := usergroups.NewGet(d.Id())
	err := vtmClient.Do(readAPI)
	if err != nil {
		if readAPI.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM User Group error whilst retrieving %s: %v", d.Id(), err)
	}

	returnedUserGroup := readAPI.ResponseObject().(*usergroups.UserGroup)

	d.Set("description", returnedUserGroup.Properties.Basic.Description)
	d.Set("password_expire_time", returnedUserGroup.Properties.Basic.PasswordExpireTime)
	d.Set("timeout", returnedUserGroup.Properties.Basic.Timeout)
	d.Set("permissions", returnedUserGroup.Properties.Basic.Permissions)

	return nil
}

func resourceUserGroupUpdate(d *schema.ResourceData, m interface{}) error {
	var updatedUserGroup usergroups.UserGroup
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
		headers := make(map[string]string)
		client := m.(*rest.Client)
		vtmClient := *client
		headers["Content-Type"] = "application/json"
		vtmClient.Headers = headers

		updateAPI := usergroups.NewPut(d.Id(), updatedUserGroup)
		err := vtmClient.Do(updateAPI)

		if err != nil {
			return fmt.Errorf("BrocadeVTM User Group error whilst updating %s: %vv", d.Id(), err)
		}
		return resourceUserGroupRead(d, m)
	}
	return nil
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*rest.Client)
	deleteAPI := usergroups.NewDelete(d.Id())
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("BrocadeVTM User Group error whilst deleting %s: %v", d.Id(), err)
	}
	d.SetId("")
	return nil
}
