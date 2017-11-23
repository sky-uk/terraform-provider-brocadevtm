package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"net/http"
)

func resourceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorSet,
		Read:   resourceMonitorRead,
		Update: resourceMonitorSet,
		Delete: resourceMonitorDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"back_off": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Should the monitor slowly increase the delay after it has failed?",
			},
			"delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: util.ValidateUnsignedInteger,
				Description:  "The minimum time between calls to a monitor.",
			},
			"failures": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: util.ValidateUnsignedInteger,
				Description:  "The number of times in a row that a node must fail execution of the monitor before it is classed as unavailable.",
			},
			"machine": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The machine to monitor, where relevant this should be in the form <hostname>:<port>, for \"ping\" monitors the <port> part must not be specified.",
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the montitor.",
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pernode",
				ValidateFunc: validation.StringInSlice([]string{
					"pernode",  // Node: Monitor each node in the pool separately
					"poolwide", // Pool/GLB: Monitor a specified machine
				}, false),
				Description: "A monitor can either monitor each node in the pool separately and disable an individual node if it fails, or it can monitor a specific machine and disable the entire pool if that machine fails. GLB location monitors must monitor a specific machine.",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: util.ValidateUnsignedInteger,
				Description:  "The maximum runtime for an individual instance of the monitor.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ping",
				ValidateFunc: validation.StringInSlice([]string{
					"connect",         // TCP Connect monitor
					"http",            // HTTP monitor
					"ping",            // Ping monitor
					"program",         //  External program monitor
					"rtsp",            // RTSP monitor
					"sip",             // SIP monitor
					"tcp_transaction", // TCP transaction monitor
				}, false),
				Description: "The internal monitor implementation of this monitor.",
			},
			"use_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the monitor should connect using SSL.",
			},
			"verbose": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the monitor should emit verbose logging. This is useful for diagnosing problems",
			},
			"http_host_header": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_body_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_status_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rtsp_body_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rtsp_status_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rtsp_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"script_program": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"script_arguments": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"sip_body_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sip_status_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sip_transport": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tcp",
					"udp",
				}, false),
			},
			"tcp_close_string": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tcp_max_response_len": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"tcp_response_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tcp_write_string": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"udp_accept_all": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func buildScriptArgumentsSection(scriptArguments interface{}) []map[string]string {

	monitorScriptArguments := make([]map[string]string, 0)

	for _, item := range scriptArguments.([]interface{}) {
		scriptArgumentItem := item.(map[string]interface{})
		monitorScriptArgument := make(map[string]string)
		scriptArgumentAttributes := []string{"name", "description", "value"}

		for _, argumentAttribute := range scriptArgumentAttributes {
			if v, ok := scriptArgumentItem[argumentAttribute].(string); ok {
				monitorScriptArgument[argumentAttribute] = v
			}
		}
		monitorScriptArguments = append(monitorScriptArguments, monitorScriptArgument)
	}
	return monitorScriptArguments
}

func getMonitorMapAttributeList(mapName string) []string {

	var attributes []string

	switch mapName {
	case "basic":
		attributes = []string{"back_off", "delay", "failures", "machine", "note", "scope", "timeout", "type", "use_ssl", "verbose"}
	case "http":
		attributes = []string{"authentication", "host_header", "body_regex", "path", "status_regex"}
	case "rtsp":
		attributes = []string{"body_regex", "status_regex", "path"}
	case "script":
		attributes = []string{"program"}
	case "sip":
		attributes = []string{"body_regex", "status_regex", "transport"}
	case "tcp":
		attributes = []string{"close_string", "max_response_len", "response_regex", "write_string"}
	case "udp":
		attributes = []string{"accept_all"}
	default:
		attributes = []string{}
	}
	return attributes
}

func resourceMonitorSet(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	monitorConfiguration := make(map[string]interface{})
	monitorPropertiesConfiguration := make(map[string]interface{})

	name := d.Get("name").(string)

	// Basic section
	monitorBasicConfiguration := make(map[string]interface{})
	monitorBasicConfiguration = util.AddSimpleGetAttributesToMap(d, monitorBasicConfiguration, "", getMonitorMapAttributeList("basic"))
	monitorPropertiesConfiguration["basic"] = monitorBasicConfiguration

	// Script section
	monitorScriptConfiguration := make(map[string]interface{})
	monitorScriptConfiguration = util.AddSimpleGetOkAttributesToMap(d, monitorScriptConfiguration, "script_", getMonitorMapAttributeList("script"))
	if d.HasChange("script_arguments") {
		if v, ok := d.GetOk("script_arguments"); ok {
			monitorScriptConfiguration["arguments"] = buildScriptArgumentsSection(v.(*schema.Set).List())
		}
		monitorPropertiesConfiguration["script"] = monitorScriptConfiguration
	}

	// All other sections
	for _, sectionName := range []string{
		"http",
		"rtsp",
		"sip",
		"tcp",
		"udp",
	} {
		section := make(map[string]interface{})
		section = util.AddSimpleGetOkAttributesToMap(d, section, fmt.Sprintf("%s_", sectionName), getMonitorMapAttributeList(sectionName))
		monitorPropertiesConfiguration[sectionName] = section
	}

	monitorConfiguration["properties"] = monitorPropertiesConfiguration
	err := client.Set("monitors", name, monitorConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Monitor error whilst creating %s: %v", name, err)
	}
	d.SetId(name)
	return resourceMonitorRead(d, m)
}

func resourceMonitorRead(d *schema.ResourceData, m interface{}) error {

	name := d.Id()
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	monitorConfiguration := make(map[string]interface{})

	client.WorkWithConfigurationResources()
	err := client.GetByName("monitors", name, &monitorConfiguration)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM Monitor error whilst retrieving %s: %v", name, err)
	}

	monitorPropertiesConfiguration := monitorConfiguration["properties"].(map[string]interface{})

	// Basic Section
	monitorBasicConfiguration := monitorPropertiesConfiguration["basic"].(map[string]interface{})
	util.SetSimpleAttributesFromMap(d, monitorBasicConfiguration, "", getMonitorMapAttributeList("basic"))

	// Script Section
	monitorScriptConfiguration := monitorPropertiesConfiguration["script"].(map[string]interface{})
	d.Set("script_program", monitorScriptConfiguration["program"])
	d.Set("script_arguments", buildScriptArgumentsSection(monitorScriptConfiguration["arguments"]))

	for _, sectionName := range []string{
		"http",
		"rtsp",
		"sip",
		"tcp",
		"udp",
	} {
		section := monitorPropertiesConfiguration[sectionName].(map[string]interface{})
		util.SetSimpleAttributesFromMap(d, section, fmt.Sprintf("%s_", sectionName), getMonitorMapAttributeList(sectionName))
	}

	return nil
}

func resourceMonitorDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("monitors", d, m)
}
