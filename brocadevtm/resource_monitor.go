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
				Default:     true,
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
				Computed: true,
			},
		},
	}
}

func buildScriptArgumentsSection(scriptArguments interface{}) []map[string]string {

	monitorScriptArguments := make([]map[string]string, 0)

	for _, item := range scriptArguments.([]interface{}) {
		scriptArgumentItem := item.(map[string]interface{})
		monitorScriptArgument := make(map[string]string)
		argumentOptions := []string{"name", "description", "value"}

		for _, option := range argumentOptions {
			if v, ok := scriptArgumentItem[option].(string); ok {
				monitorScriptArgument[option] = v
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
	monitorBasicConfiguration["back_off"] = d.Get("back_off").(bool)
	monitorBasicConfiguration["delay"] = d.Get("delay").(int)
	monitorBasicConfiguration["failures"] = d.Get("failures").(int)
	if v, ok := d.GetOk("machine"); ok {
		monitorBasicConfiguration["machine"] = v.(string)
	}
	if v, ok := d.GetOk("note"); ok {
		monitorBasicConfiguration["note"] = v.(string)
	}
	monitorBasicConfiguration["scope"] = d.Get("scope").(string)
	monitorBasicConfiguration["timeout"] = d.Get("timeout").(int)
	monitorBasicConfiguration["type"] = d.Get("type").(string)
	monitorBasicConfiguration["use_ssl"] = d.Get("use_ssl").(bool)
	monitorBasicConfiguration["verbose"] = d.Get("verbose").(bool)
	monitorPropertiesConfiguration["basic"] = monitorBasicConfiguration

	// HTTP Section
	monitorHTTPConfiguration := make(map[string]interface{})
	if v, ok := d.GetOk("http_host_header"); ok {
		monitorHTTPConfiguration["host_header"] = v.(string)
	}
	if v, ok := d.GetOk("http_path"); ok {
		monitorHTTPConfiguration["path"] = v.(string)
	}
	if v, ok := d.GetOk("http_authentication"); ok {
		monitorHTTPConfiguration["authentication"] = v.(string)
	}
	if v, ok := d.GetOk("http_body_regex"); ok {
		monitorHTTPConfiguration["body_regex"] = v.(string)
	}
	if v, ok := d.GetOk("http_status_regex"); ok {
		monitorHTTPConfiguration["status_regex"] = v.(string)
	}
	monitorPropertiesConfiguration["http"] = monitorHTTPConfiguration

	// RTSP section
	monitorRTSPConfiguration := make(map[string]interface{})
	if v, ok := d.GetOk("rtsp_body_regex"); ok {
		monitorRTSPConfiguration["body_regex"] = v.(string)
	}
	if v, ok := d.GetOk("rtsp_status_regex"); ok {
		monitorRTSPConfiguration["status_regex"] = v.(string)
	}
	if v, ok := d.GetOk("rtsp_path"); ok {
		monitorRTSPConfiguration["path"] = v.(string)
	}
	monitorPropertiesConfiguration["rtsp"] = monitorRTSPConfiguration

	// Script section
	monitorScriptConfiguration := make(map[string]interface{})
	if v, ok := d.GetOk("script_program"); ok {
		monitorScriptConfiguration["program"] = v.(string)
	}
	if v, ok := d.GetOk("script_arguments"); ok {
		monitorScriptConfiguration["arguments"] = buildScriptArgumentsSection(v.(*schema.Set).List())
	}
	monitorPropertiesConfiguration["script"] = monitorScriptConfiguration

	// SIP Section
	monitorSIPConfiguration := make(map[string]interface{})
	if v, ok := d.GetOk("sip_body_regex"); ok {
		monitorSIPConfiguration["body_regex"] = v.(string)
	}
	if v, ok := d.GetOk("sip_status_regex"); ok {
		monitorSIPConfiguration["status_regex"] = v.(string)
	}
	if v, ok := d.GetOk("sip_transport"); ok {
		monitorSIPConfiguration["transport"] = v.(string)
	}
	monitorPropertiesConfiguration["sip"] = monitorSIPConfiguration

	// TCP Section
	monitorTCPConfiguration := make(map[string]interface{})
	if v, ok := d.GetOk("tcp_close_string"); ok {
		monitorTCPConfiguration["close_string"] = v.(string)
	}
	if v, ok := d.GetOk("tcp_max_response_len"); ok {
		monitorTCPConfiguration["max_response_len"] = v.(int)
	}
	if v, ok := d.GetOk("tcp_response_regex"); ok {
		monitorTCPConfiguration["response_regex"] = v.(string)
	}
	if v, ok := d.GetOk("tcp_write_string"); ok {
		monitorTCPConfiguration["write_string"] = v.(string)
	}
	monitorPropertiesConfiguration["tcp"] = monitorTCPConfiguration

	// UDP Section
	monitorUDPConfiguration := make(map[string]interface{})
	if v, ok := d.GetOk("udp_accept_all"); ok {
		monitorUDPConfiguration["accept_all"] = v.(bool)
	}
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
	d.Set("back_off", monitorBasicConfiguration["back_off"])
	d.Set("delay", monitorBasicConfiguration["delay"])
	d.Set("failures", monitorBasicConfiguration["failures"])
	d.Set("machine", monitorBasicConfiguration["machine"])
	d.Set("note", monitorBasicConfiguration["note"])
	d.Set("scope", monitorBasicConfiguration["scope"])
	d.Set("timeout", monitorBasicConfiguration["timeout"])
	d.Set("type", monitorBasicConfiguration["type"])
	d.Set("use_ssl", monitorBasicConfiguration["use_ssl"])
	d.Set("verbose", monitorBasicConfiguration["verbose"])

	// HTTP Section
	monitorHTTPConfiguration := monitorPropertiesConfiguration["http"].(map[string]interface{})
	d.Set("http_host_header", monitorHTTPConfiguration["host_header"])
	d.Set("http_path", monitorHTTPConfiguration["path"])
	d.Set("http_authentication", monitorHTTPConfiguration["authentication"])
	d.Set("http_body_regex", monitorHTTPConfiguration["body_regex"])
	d.Set("http_status_regex", monitorHTTPConfiguration["status_regex"])

	// RTSP Section
	monitorRTSPConfiguration := monitorPropertiesConfiguration["rtsp"].(map[string]interface{})
	d.Set("rtsp_body_regex", monitorRTSPConfiguration["body_regex"])
	d.Set("rtsp_status_regex", monitorRTSPConfiguration["status_regex"])
	d.Set("rtsp_path", monitorRTSPConfiguration["path"])

	// Script Section
	monitorScriptConfiguration := monitorPropertiesConfiguration["script"].(map[string]interface{})
	d.Set("script_program", monitorScriptConfiguration["program"])
	d.Set("script_arguments", buildScriptArgumentsSection(monitorScriptConfiguration["arguments"]))

	// SIP Section
	monitorSIPConfiguration := monitorPropertiesConfiguration["sip"].(map[string]interface{})
	d.Set("sip_body_regex", monitorSIPConfiguration["body_regex"])
	d.Set("sip_status_regex", monitorSIPConfiguration["status_regex"])
	d.Set("sip_transport", monitorSIPConfiguration["transport"])

	// TCP Section
	monitorTCPConfiguration := monitorPropertiesConfiguration["tcp"].(map[string]interface{})
	d.Set("tcp_close_string", monitorTCPConfiguration["close_string"])
	d.Set("tcp_max_response_len", monitorTCPConfiguration["max_response_len"])
	d.Set("tcp_response_regex", monitorTCPConfiguration["response_regex"])
	d.Set("tcp_write_string", monitorTCPConfiguration["write_string"])

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
	if d.HasChange("back_off") {
		monitorBasicConfiguration["back_off"] = d.Get("back_off").(bool)
	}
	if d.HasChange("delay") {
		monitorBasicConfiguration["delay"] = d.Get("delay").(int)
	}
	if d.HasChange("failures") {
		monitorBasicConfiguration["failures"] = d.Get("failures").(int)
	}
	if d.HasChange("machine") {
		monitorBasicConfiguration["machine"] = d.Get("machine").(string)
	}
	if d.HasChange("note") {
		monitorBasicConfiguration["note"] = d.Get("note").(string)
	}
	if d.HasChange("scope") {
		monitorBasicConfiguration["scope"] = d.Get("scope").(string)
	}
	if d.HasChange("timeout") {
		monitorBasicConfiguration["timeout"] = d.Get("timeout").(int)
	}
	if d.HasChange("type") {
		monitorBasicConfiguration["type"] = d.Get("type").(string)
	}
	if d.HasChange("verbose") {
		monitorBasicConfiguration["verbose"] = d.Get("verbose").(bool)
	}
	if d.HasChange("use_ssl") {
		monitorBasicConfiguration["use_ssl"] = d.Get("use_ssl").(bool)
	}
	monitorPropertiesConfiguration["basic"] = monitorBasicConfiguration

	// HTTP Section
	monitorHTTPConfiguration := make(map[string]interface{})
	if d.HasChange("http_host_header") {
		monitorHTTPConfiguration["host_header"] = d.Get("http_host_header").(string)
	}
	if d.HasChange("http_path") {
		monitorHTTPConfiguration["path"] = d.Get("http_path").(string)
	}
	if d.HasChange("http_authentication") {
		monitorHTTPConfiguration["authentication"] = d.Get("http_authentication").(string)
	}
	if d.HasChange("http_body_regex") {
		monitorHTTPConfiguration["body_regex"] = d.Get("http_body_regex").(string)
	}
	if d.HasChange("http_status_regex") {
		monitorHTTPConfiguration["status_regex"] = d.Get("http_status_regex").(string)
	}
	monitorPropertiesConfiguration["http"] = monitorHTTPConfiguration

	// RTSP Section
	monitorRTSPConfiguration := make(map[string]interface{})
	if d.HasChange("rtsp_status_regex") {
		monitorRTSPConfiguration["status_regex"] = d.Get("rtsp_status_regex").(string)
	}
	if d.HasChange("rtsp_body_regex") {
		monitorRTSPConfiguration["body_regex"] = d.Get("rtsp_body_regex").(string)
	}
	if d.HasChange("rtsp_path") {
		monitorRTSPConfiguration["path"] = d.Get("rtsp_path").(string)
	}
	monitorPropertiesConfiguration["rtsp"] = monitorRTSPConfiguration

	// Script Section
	monitorScriptConfiguration := make(map[string]interface{})
	if d.HasChange("script_arguments") {
		monitorScriptConfiguration["arguments"] = buildScriptArgumentsSection(d.Get("script_arguments").(*schema.Set).List())
	}
	if d.HasChange("script_program") {
		monitorScriptConfiguration["program"] = d.Get("script_program").(string)
	}
	monitorPropertiesConfiguration["script"] = monitorScriptConfiguration

	// SIP Section
	monitorSIPConfiguration := make(map[string]interface{})
	if d.HasChange("sip_body_regex") {
		monitorSIPConfiguration["body_regex"] = d.Get("sip_body_regex").(string)
	}
	if d.HasChange("sip_status_regex") {
		monitorSIPConfiguration["status_regex"] = d.Get("sip_status_regex").(string)
	}
	if d.HasChange("sip_transport") {
		monitorSIPConfiguration["transport"] = d.Get("sip_transport").(string)
	}
	monitorPropertiesConfiguration["sip"] = monitorSIPConfiguration

	// TCP Section
	monitorTCPConfiguration := make(map[string]interface{})
	if d.HasChange("tcp_close_string") {
		monitorTCPConfiguration["close_string"] = d.Get("tcp_close_string").(string)
	}
	if d.HasChange("tcp_max_response_len") {
		monitorTCPConfiguration["max_response_len"] = d.Get("tcp_max_response_len").(int)
	}
	if d.HasChange("tcp_response_regex") {
		monitorTCPConfiguration["response_regex"] = d.Get("tcp_response_regex").(string)
	}
	if d.HasChange("tcp_write_string") {
		monitorTCPConfiguration["write_string"] = d.Get("tcp_write_string").(string)
	}
	monitorPropertiesConfiguration["tcp"] = monitorTCPConfiguration

	// UDP Section
	monitorUDPConfiguration := make(map[string]interface{})
	if d.HasChange("udp_accept_all") {
		monitorUDPConfiguration["accept_all"] = d.Get("udp_accept_all").(bool)
	}
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
