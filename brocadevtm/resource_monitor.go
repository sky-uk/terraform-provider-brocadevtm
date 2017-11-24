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
			"http": {
				Type:        schema.TypeSet,
				Description: "HTTP section",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"body_regex": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status_regex": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"rtsp": {
				Type:        schema.TypeSet,
				Description: "RTSP section",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body_regex": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status_regex": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"script_arguments": {
				Type:        schema.TypeSet,
				Description: "Script arguments to script program",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"script_program": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sip": {
				Type:        schema.TypeSet,
				Description: "SIP section",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body_regex": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status_regex": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"transport": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"tcp",
								"udp",
							}, false),
						},
					},
				},
			},
			"tcp": {
				Type:        schema.TypeSet,
				Description: "TCP section",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"close_string": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"max_response_len": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2048,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"response_regex": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  ".+",
						},
						"write_string": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"udp": {
				Type:        schema.TypeSet,
				Description: "UDP section",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accept_all": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func basicMonitorKeys() []string {
	return []string{
		"back_off",
		"delay",
		"failures",
		"machine",
		"note",
		"scope",
		"timeout",
		"type",
		"use_ssl",
		"verbose",
	}
}

func resourceMonitorSet(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	monitorRequest := make(map[string]interface{})
	monitorProperties := make(map[string]interface{})

	name := d.Get("name").(string)
	// basic section
	util.GetSection(d, "basic", monitorProperties, basicMonitorKeys())

	// all other sections apart from script
	for _, section := range []string{
		"http",
		"rtsp",
		"sip",
		"tcp",
		"udp",
	} {
		if d.HasChange(section) {
			monitorProperties[section] = d.Get(section).(*schema.Set).List()[0]
		}
	}

	// script section
	monitorScriptSection := make(map[string]interface{})
	if d.HasChange("script_arguments") {
		monitorScriptSection["arguments"] = d.Get("script_arguments").(*schema.Set).List()
	}
	if d.HasChange("script_program") {
		monitorScriptSection["program"] = d.Get("script_program").(string)
	}
	monitorProperties["script"] = monitorScriptSection

	monitorRequest["properties"] = monitorProperties
	util.TraverseMapTypes(monitorRequest)
	err := client.Set("monitors", name, monitorRequest, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM Monitor error whilst creating %s: %v", name, err)
	}
	d.SetId(name)
	return resourceMonitorRead(d, m)
}

func resourceMonitorRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	monitorResponse := make(map[string]interface{})
	name := d.Id()

	client.WorkWithConfigurationResources()
	err := client.GetByName("monitors", name, &monitorResponse)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM Monitor error whilst retrieving %s: %v", name, err)
	}
	monitorProperties := monitorResponse["properties"].(map[string]interface{})
	monitorBasic := monitorProperties["basic"].(map[string]interface{})

	// basic section
	for _, key := range basicMonitorKeys() {
		err := d.Set(key, monitorBasic[key])
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM Monitor error whilst setting key %s: %v", key, err)
		}
	}

	// all other sections apart from script
	for _, sectionName := range []string{
		"http",
		"rtsp",
		"sip",
		"tcp",
		"udp",
	} {
		set := make([]map[string]interface{}, 0)
		readSectionMap, err := util.BuildReadMap(monitorProperties[sectionName].(map[string]interface{}))
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM Monitor error whilst building section map to set: %v", err)
		}
		set = append(set, readSectionMap)
		err = d.Set(sectionName, set)
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM Monitor error whilst setting section %s: %v", sectionName, err)
		}
	}

	// script section
	err = d.Set("script_program", monitorProperties["script"].(map[string]interface{})["program"].(string))
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM Monitor error whilst setting script_program: %v", err)
	}
	scriptArgumentSet := make([]map[string]interface{}, 0)
	for _, item := range monitorProperties["script"].(map[string]interface{})["arguments"].([]interface{}) {
		scriptArgument := make(map[string]interface{})
		for key, value := range item.(map[string]interface{}) {
			scriptArgument[key] = value.(string)
		}
		scriptArgumentSet = append(scriptArgumentSet, scriptArgument)
	}
	err = d.Set("script_arguments", scriptArgumentSet)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM Monitor error whilst setting script_arguments: %v", err)
	}

	return nil
}

func resourceMonitorDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("monitors", d, m)
}
