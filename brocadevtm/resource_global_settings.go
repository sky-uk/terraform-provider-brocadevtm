package brocadevtm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func resourceGlobalSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlobalSettingsCreate,
		Read:   resourceGlobalSettingsRead,
		Update: resourceGlobalSettingsUpdate,
		Delete: resourceGlobalSettingsDelete,

		Schema: map[string]*schema.Schema{
			"basic": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepting_delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      50,
							ValidateFunc: util.ValidateUnsignedInteger,
							Description:  "How often, in milliseconds, each traffic manager child process (that isn't listening for new connections) checks to see whether it should start listening for new connections.",
						},
						"afm_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Is the application firewall enabled.",
						},
						"chunk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      16384,
							ValidateFunc: util.ValidateUnsignedInteger,
							Description:  "The default chunk size for reading/writing requests",
						},
						"client_first_opt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not your traffic manager should make use of TCP optimisations to defer the processing of new client-first connections until the client has sent some data.",
						},
						"cluster_identifier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster identifier. Generally supplied by Services Director.",
						},
						"license_servers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of license servers for FLA licensing. A license server should be specified as a <ip/host>:<port> pair.",
						},
						"max_fds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1048576,
							ValidateFunc: util.ValidateUnsignedInteger,
							Description:  "The maximum number of file descriptors that your traffic manager will allocate.",
						},
						"monitor_memory_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: util.ValidateUnsignedInteger,
							Default:      4096,
							Description:  "The maximum number of each of nodes, pools or locations that can be monitored. The memory used to store information about nodes, pools and locations is allocated at start-up, so the traffic manager must be restarted after changing this setting.",
						},
						"rate_class_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: util.ValidateUnsignedInteger,
							Default:      25000,
							Description:  "The maximum number of Rate classes that can be created. Approximately 100 bytes will be pre-allocated per Rate class.",
						},
						"shared_pool_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "10MB",
							Description: "The size of the shared memory pool used for shared storage across worker processes (e.g. bandwidth shared data).This is specified as either a percentage of system RAM, 5% for example, or an absolute size such as 10MB.",
						},
						"slm_class_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1024,
							ValidateFunc: util.ValidateUnsignedInteger,
							Description:  "The maximum number of SLM classes that can be created. Approximately 100 bytes will be pre-allocated per SLM class.",
						},
						"so_rbuff_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: util.ValidateUnsignedInteger,
							Default:      0,
							Description:  "The size of the operating system's read buffer. A value of 0 (zero) means to use the OS default; in normal circumstances this is what should be used.",
						},
						"so_wbuff_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: util.ValidateUnsignedInteger,
							Default:      0,
							Description:  "The size of the operating system's write buffer. A value of 0 (zero) means to use the OS default; in normal circumstances this is what should be used.",
						},
						"socket_optimizations": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "auto",
							ValidateFunc: validateSocketOptimizations,
							Description:  "Whether or not the traffic manager should use potential network socket optimisations. If set to auto, a decision will be made based on the host platform.",
						},
						"tip_class_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: util.ValidateUnsignedInteger,
							Default:      10000,
							Description:  "The maximum number of Traffic IP Groups that can be created.",
						},
					},
				},
			},
		},
	}
}

func validateSocketOptimizations(v interface{}, k string) (ws []string, errors []error) {
	so := v.(string)
	if so != "auto" && so != "no" && so != "yes" {
		errors = append(errors, fmt.Errorf("socket_optimizations value not valid (must be either \"auto\" or \"no\" or \"yes\""))
	}
	return
}

func resourceGlobalSettingsCreate(d *schema.ResourceData, m interface{}) error {

	// This resource can't actually be created
	globalSettings := make(map[string]interface{})
	properties := make(map[string]interface{})
	basic := d.Get("basic").([]interface{})
	properties["basic"] = basic[0]
	globalSettings["properties"] = properties

	d.SetId("global_settings")
	return resourceGlobalSettingsRead(d, m)
}

func resourceGlobalSettingsRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	globalSettings := make(map[string]interface{})
	client.WorkWithConfigurationResources()
	err := client.GetByName("global_settings", "", &globalSettings)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("BrocadeVTM error whilst reading %s: %v", "", err)
	}
	properties := globalSettings["properties"].(map[string]interface{})
	log.Println("Properties:\n", properties)
	d.SetId("global_settings")
	d.Set("basic", properties["basic"])

	return nil
}

func resourceGlobalSettingsUpdate(d *schema.ResourceData, m interface{}) error {

	hasChanges := false
	globalSettings := make(map[string]interface{})
	properties := make(map[string]interface{})

	if d.HasChange("basic") {
		basic := d.Get("basic").([]interface{})
		properties["basic"] = basic[0]
		hasChanges = true
	}
	if hasChanges {
		config := m.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		globalSettings["properties"] = properties
		err := client.Set("global_settings", "", globalSettings, nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM DNS error whilst updating %s: %v", "", err)
		}
	}
	return resourceGlobalSettingsRead(d, m)
}

func resourceGlobalSettingsDelete(d *schema.ResourceData, m interface{}) error {
	// this resource can't actually be deleted
	d.SetId("")
	return nil
}
