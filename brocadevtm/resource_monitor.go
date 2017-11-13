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
		Create: resourceMonitorCreate,
		Read:   resourceMonitorRead,
		Update: resourceMonitorUpdate,
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

func resourceMonitorCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	monitorConfiguration := make(map[string]interface{})
	monitorPropertiesConfiguration := make(map[string]interface{})

	name := d.Get("name").(string)

	// Basic section
	monitorBasicConfiguration := make(map[string]interface{})
	monitorBasicConfiguration = util.AddSimpleGetAttributesToMap(d, monitorBasicConfiguration, "", []string{"back_off", "delay", "failures", "machine", "note", "scope", "timeout", "type", "use_ssl", "verbose"})
	monitorPropertiesConfiguration["basic"] = monitorBasicConfiguration

	// HTTP Section
	monitorHTTPConfiguration := make(map[string]interface{})
	monitorHTTPConfiguration = util.AddSimpleGetOkAttributesToMap(d, monitorHTTPConfiguration, "http_", []string{"authentication", "host_header", "body_regex", "path", "status_regex"})
	monitorPropertiesConfiguration["http"] = monitorHTTPConfiguration

	// RTSP section
	monitorRTSPConfiguration := make(map[string]interface{})
	monitorRTSPConfiguration = util.AddSimpleGetOkAttributesToMap(d, monitorRTSPConfiguration, "rtsp_", []string{"body_regex", "status_regex", "path"})
	monitorPropertiesConfiguration["rtsp"] = monitorRTSPConfiguration

	// Script section
	monitorScriptConfiguration := make(map[string]interface{})
	monitorScriptConfiguration = util.AddSimpleGetOkAttributesToMap(d, monitorScriptConfiguration, "script_", []string{"program"})
	if v, ok := d.GetOk("script_arguments"); ok {
		monitorScriptConfiguration["arguments"] = buildScriptArgumentsSection(v.(*schema.Set).List())
	}
	monitorPropertiesConfiguration["script"] = monitorScriptConfiguration

	// SIP Section
	monitorSIPConfiguration := make(map[string]interface{})
	monitorSIPConfiguration = util.AddSimpleGetOkAttributesToMap(d, monitorSIPConfiguration, "sip_", []string{"body_regex", "status_regex", "transport"})
	monitorPropertiesConfiguration["sip"] = monitorSIPConfiguration

	// TCP Section
	monitorTCPConfiguration := make(map[string]interface{})
	monitorTCPConfiguration = util.AddSimpleGetOkAttributesToMap(d, monitorTCPConfiguration, "tcp_", []string{"close_string", "max_response_len", "response_regex", "write_string"})
	monitorPropertiesConfiguration["tcp"] = monitorTCPConfiguration

	// UDP Section
	monitorUDPConfiguration := make(map[string]interface{})
	monitorUDPConfiguration = util.AddSimpleGetOkAttributesToMap(d, monitorUDPConfiguration, "udp_", []string{"accept_all"})
	monitorPropertiesConfiguration["udp"] = monitorUDPConfiguration

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
	util.SetSimpleAttributesFromMap(d, monitorBasicConfiguration, "", []string{"back_off", "delay", "failures", "machine", "note", "scope", "timeout", "type", "use_ssl", "verbose"})

	// HTTP Section
	monitorHTTPConfiguration := monitorPropertiesConfiguration["http"].(map[string]interface{})
	util.SetSimpleAttributesFromMap(d, monitorHTTPConfiguration, "http_", []string{"host_header", "path", "authentication", "body_regex", "status_regex"})

	// RTSP Section
	monitorRTSPConfiguration := monitorPropertiesConfiguration["rtsp"].(map[string]interface{})
	util.SetSimpleAttributesFromMap(d, monitorRTSPConfiguration, "rtsp_", []string{"body_regex", "status_regex", "path"})

	// Script Section
	monitorScriptConfiguration := monitorPropertiesConfiguration["script"].(map[string]interface{})
	d.Set("script_program", monitorScriptConfiguration["program"])
	d.Set("script_arguments", buildScriptArgumentsSection(monitorScriptConfiguration["arguments"]))

	// SIP Section
	monitorSIPConfiguration := monitorPropertiesConfiguration["sip"].(map[string]interface{})
	util.SetSimpleAttributesFromMap(d, monitorSIPConfiguration, "sip_", []string{"body_regex", "status_regex", "transport"})

	// TCP Section
	monitorTCPConfiguration := monitorPropertiesConfiguration["tcp"].(map[string]interface{})
	util.SetSimpleAttributesFromMap(d, monitorTCPConfiguration, "tcp_", []string{"close_string", "max_response_len", "response_regex", "write_string"})

	// UDP Section
	monitorUDPConfiguration := monitorPropertiesConfiguration["udp"].(map[string]interface{})
	d.Set("udp_accept_all", monitorUDPConfiguration["accept_all"])

	return nil
}

func resourceMonitorUpdate(d *schema.ResourceData, m interface{}) error {

	name := d.Id()
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	monitorConfiguration := make(map[string]interface{})
	monitorPropertiesConfiguration := make(map[string]interface{})

	// Basic Section
	monitorBasicConfiguration := make(map[string]interface{})
	monitorBasicConfiguration = util.AddChangedSimpleAttributesToMap(d, monitorBasicConfiguration, "", []string{"back_off", "delay", "failures", "machine", "note", "scope", "timeout", "type", "verbose", "use_ssl"})
	monitorPropertiesConfiguration["basic"] = monitorBasicConfiguration

	// HTTP Section
	monitorHTTPConfiguration := make(map[string]interface{})
	monitorHTTPConfiguration = util.AddChangedSimpleAttributesToMap(d, monitorHTTPConfiguration, "http_", []string{"host_header", "path", "authentication", "body_regex", "status_regex"})
	monitorPropertiesConfiguration["http"] = monitorHTTPConfiguration

	// RTSP Section
	monitorRTSPConfiguration := make(map[string]interface{})
	monitorRTSPConfiguration = util.AddChangedSimpleAttributesToMap(d, monitorRTSPConfiguration, "rtsp_", []string{"status_regex", "body_regex", "path"})
	monitorPropertiesConfiguration["rtsp"] = monitorRTSPConfiguration

	// Script Section
	monitorScriptConfiguration := make(map[string]interface{})
	monitorScriptConfiguration = util.AddChangedSimpleAttributesToMap(d, monitorScriptConfiguration, "script_", []string{"program"})
	if d.HasChange("script_arguments") {
		monitorScriptConfiguration["arguments"] = buildScriptArgumentsSection(d.Get("script_arguments").(*schema.Set).List())
	}
	monitorPropertiesConfiguration["script"] = monitorScriptConfiguration

	// SIP Section
	monitorSIPConfiguration := make(map[string]interface{})
	monitorSIPConfiguration = util.AddChangedSimpleAttributesToMap(d, monitorSIPConfiguration, "sip_", []string{"body_regex", "status_regex", "transport"})
	monitorPropertiesConfiguration["sip"] = monitorSIPConfiguration

	// TCP Section
	monitorTCPConfiguration := make(map[string]interface{})
	monitorTCPConfiguration = util.AddChangedSimpleAttributesToMap(d, monitorTCPConfiguration, "tcp_", []string{"close_string", "max_response_len", "response_regex", "write_string"})
	monitorPropertiesConfiguration["tcp"] = monitorTCPConfiguration

	// UDP Section
	monitorUDPConfiguration := make(map[string]interface{})
	monitorUDPConfiguration = util.AddChangedSimpleAttributesToMap(d, monitorUDPConfiguration, "udp_", []string{"accept_all"})
	monitorPropertiesConfiguration["udp"] = monitorUDPConfiguration

	monitorConfiguration["properties"] = monitorPropertiesConfiguration
	err := client.Set("monitors", name, monitorConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Monitor error whilst updating %s: %s", name, err)
	}
	return resourceMonitorRead(d, m)
}

func resourceMonitorDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("monitors", d, m)
}
