package pulsevtm

import (
	"fmt"

	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-pulse-vtm/api"
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntAtLeast(0),
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
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"none",
								"ro",
								"full",
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	res := make(map[string]interface{})
	props := make(map[string]interface{})
	basic := make(map[string]interface{})

	name := d.Get("name").(string)

	if v, ok := d.GetOk("description"); ok && v != "" {
		basic["description"] = v.(string)
	}
	if v, ok := d.GetOk("password_expire_time"); ok {
		basic["password_expire_time"] = uint(v.(int))
	}
	if v, ok := d.GetOk("timeout"); ok {
		basic["timeout"] = uint(v.(int))
	}
	if v, ok := d.GetOk("permissions"); ok {
		basic["permissions"] = v.(*schema.Set).List()
	}

	props["basic"] = basic
	res["properties"] = props

	err := client.Set("user_groups", name, &res, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM User Group error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceUserGroupRead(d, m)
}

func resourceUserGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	res := make(map[string]interface{})

	client.WorkWithConfigurationResources()
	err := client.GetByName("user_groups", d.Id(), &res)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		d.SetId("")
		return fmt.Errorf("[ERROR] PulseVTM User Group error whilst retrieving %s: %v", d.Id(), err)
	}

	props := res["properties"].(map[string]interface{})
	basic := props["basic"].(map[string]interface{})

	for _, attribute := range []string{"description", "password_expire_time", "timeout", "permissions"} {
		err = d.Set(attribute, basic[attribute])
		if err != nil {
			return fmt.Errorf("[ERROR] PulseVTM User Group error whilst setting attributes %s: %v", attribute, err)
		}
	}

	return nil
}

func resourceUserGroupUpdate(d *schema.ResourceData, m interface{}) error {
	res := make(map[string]interface{})
	props := make(map[string]interface{})
	basic := make(map[string]interface{})

	if d.HasChange("description") {
		basic["description"] = d.Get("description").(string)
	}
	if d.HasChange("password_expire_time") {
		basic["password_expire_time"] = uint(d.Get("password_expire_time").(int))
	}
	if d.HasChange("timeout") {
		basic["timeout"] = uint(d.Get("timeout").(int))
	}
	if d.HasChange("permissions") {
		basic["permissions"] = d.Get("permissions").(*schema.Set).List()
	}

	props["basic"] = basic
	res["properties"] = props

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	err := client.Set("user_groups", d.Id(), res, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM User Group error whilst updating %s: %v", d.Id(), err)
	}
	return resourceUserGroupRead(d, m)
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("user_groups", d, m)
}
