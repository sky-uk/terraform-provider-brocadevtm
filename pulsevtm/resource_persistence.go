package pulsevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-pulse-vtm/api"
	"net/http"
	"regexp"
)

func resourcePersistence() *schema.Resource {
	return &schema.Resource{
		Create: resourcePersistenceCreate,
		Read:   resourcePersistenceRead,
		Update: resourcePersistenceUpdate,
		Delete: resourcePersistenceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the bandwidth class",
			},
			"cookie": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cookie to use for tracking session persistence",
			},
			"delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not a session is deleted when a session fails",
			},
			"failure_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "new_node",
				Description:  "Action the pool takes if the session data is invalid or the node can't be contacted",
				ValidateFunc: validateFailureMode,
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Note regarding the session persistence class",
			},
			"subnet_prefix_length_v4": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ensure all requests from this IPv4 subnet, specified as a prefix length, are sent to the same node. If set to 0, requests from different IPv4 addresses will be load-balanced individually.",
			},
			"subnet_prefix_length_v6": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ensure all requests from this IPv6 subnet, specified as a prefix length, are sent to the same node. If set to 0, requests from different IPv6 addresses will be load-balanced individually.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ip",
				Description:  "Type of session persistence to use",
				ValidateFunc: validateType,
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Redirect URL to send clients if session persistence is configured to redirect when a node dies",
			},
		},
	}
}

func validateFailureMode(v interface{}, k string) (ws []string, errors []error) {
	failureMode := v.(string)
	failureModeOptions := regexp.MustCompile(`^(close|new_node|url)$`)
	if !failureModeOptions.MatchString(failureMode) {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be one of close, new_node or url", k))
	}
	return
}

func validateType(v interface{}, k string) (ws []string, errors []error) {
	persistenceType := v.(string)
	persistenceTypeOptions := regexp.MustCompile(`^(asp|cookie|ip|j2ee|named|ssl|transparent|universal|x_zeus)$`)
	if !persistenceTypeOptions.MatchString(persistenceType) {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be one of asp, cookie, ip, j2ee, named, ssl, transparent, universal or x_zeus", k))
	}
	return
}

func resourcePersistenceCreate(d *schema.ResourceData, m interface{}) error {

	var name string
	config := m.(map[string]interface{})
	persistenceBasicConfiguration := make(map[string]interface{})
	persistencePropertiesConfiguration := make(map[string]interface{})
	persistenceConfiguration := make(map[string]interface{})

	client := config["jsonClient"].(*api.Client)

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	if v, ok := d.GetOk("cookie"); ok {
		persistenceBasicConfiguration["cookie"] = v.(string)
	}
	persistenceBasicConfiguration["delete"] = d.Get("delete").(bool)

	if v, ok := d.GetOk("failure_mode"); ok && v != "" {
		persistenceBasicConfiguration["failure_mode"] = v.(string)
	}
	if v, ok := d.GetOk("note"); ok {
		persistenceBasicConfiguration["note"] = v.(string)
	}

	if v, ok := d.GetOk("subnet_prefix_length_v4"); ok {
		persistenceBasicConfiguration["subnet_prefix_length_v4"] = v.(int)
	}

	if v, ok := d.GetOk("subnet_prefix_length_v6"); ok {
		persistenceBasicConfiguration["subnet_prefix_length_v6"] = v.(int)
	}

	if v, ok := d.GetOk("type"); ok && v != "" {
		persistenceBasicConfiguration["type"] = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		persistenceBasicConfiguration["url"] = v.(string)
	}
	persistencePropertiesConfiguration["basic"] = persistenceBasicConfiguration
	persistenceConfiguration["properties"] = persistencePropertiesConfiguration

	err := client.Set("persistence", name, &persistenceConfiguration, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM Persistence error whilst creating %s: %v", name, err)
	}
	d.SetId(name)
	return resourcePersistenceRead(d, m)
}

func resourcePersistenceRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	name := d.Id()
	persistenceConfiguration := make(map[string]interface{})

	err := client.GetByName("persistence", name, &persistenceConfiguration)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM Persistence error whilst retrieving %s: %v", name, err)
	}
	persistencePropertiesConfiguration := persistenceConfiguration["properties"].(map[string]interface{})
	persistenceBasicConfiguration := persistencePropertiesConfiguration["basic"].(map[string]interface{})

	for _, attribute := range []string{"cookie", "delete", "failure_mode", "note", "type", "url"} {
		err := d.Set(attribute, persistenceBasicConfiguration[attribute])
		if err != nil {
			return fmt.Errorf("[ERROR] PulseVTM Persistence error whilst setting attributes %s: %v", attribute, err)
		}
	}
	return nil
}

func resourcePersistenceUpdate(d *schema.ResourceData, m interface{}) error {

	name := d.Id()
	persistenceBasicConfiguration := make(map[string]interface{})
	persistencePropertiesConfiguration := make(map[string]interface{})
	persistenceConfiguration := make(map[string]interface{})

	if d.HasChange("cookie") {
		if v, ok := d.GetOk("cookie"); ok {
			persistenceBasicConfiguration["cookie"] = v.(string)
		}
	}
	if d.HasChange("delete") {
		persistenceBasicConfiguration["delete"] = d.Get("delete").(bool)
	}
	if d.HasChange("failure_mode") {
		if v, ok := d.GetOk("failure_mode"); ok && v != "" {
			persistenceBasicConfiguration["failure_mode"] = v.(string)
		}
	}
	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok {
			persistenceBasicConfiguration["note"] = v.(string)
		}
	}

	if d.HasChange("subnet_prefix_length_v4") {
		if v, ok := d.GetOk("subnet_prefix_length_v4"); ok {
			persistenceBasicConfiguration["subnet_prefix_length_v4"] = v.(int)
		}
	}

	if d.HasChange("subnet_prefix_length_v6") {
		if v, ok := d.GetOk("subnet_prefix_length_v6"); ok {
			persistenceBasicConfiguration["subnet_prefix_length_v6"] = v.(int)
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok && v != "" {
			persistenceBasicConfiguration["type"] = v.(string)
		}
	}
	if d.HasChange("url") {
		if v, ok := d.GetOk("url"); ok {
			persistenceBasicConfiguration["url"] = v.(string)
		}
	}
	persistencePropertiesConfiguration["basic"] = persistenceBasicConfiguration
	persistenceConfiguration["properties"] = persistencePropertiesConfiguration
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	err := client.Set("persistence", name, &persistenceConfiguration, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM Persistence error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourcePersistenceRead(d, m)
}

func resourcePersistenceDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("persistence", d, m)
}
