package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/monitor"
	"github.com/sky-uk/go-rest-api"
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
			"delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateMonitorUnsignedInteger,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateMonitorUnsignedInteger,
			},
			"failures": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateMonitorUnsignedInteger,
			},
			"verbose": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"use_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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

func validateMonitorUnsignedInteger(v interface{}, k string) (ws []string, errors []error) {
	ttl := v.(int)
	if ttl < 0 {
		errors = append(errors, fmt.Errorf("%q can't be negative", k))
	}
	return
}

func resourceMonitorCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var createMonitor monitor.Monitor
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("delay"); ok {
		delay := v.(int)
		createMonitor.Properties.Basic.Delay = uint(delay)
	}
	if v, ok := d.GetOk("timeout"); ok {
		timeout := v.(int)
		createMonitor.Properties.Basic.Timeout = uint(timeout)
	}
	if v, ok := d.GetOk("failures"); ok {
		failures := v.(int)
		createMonitor.Properties.Basic.Failures = uint(failures)
	}
	if v, ok := d.GetOk("verbose"); ok {
		monitorVerbosity := v.(bool)
		createMonitor.Properties.Basic.Verbose = &monitorVerbosity
	}
	if v, ok := d.GetOk("use_ssl"); ok {
		monitorSSL := v.(bool)
		createMonitor.Properties.Basic.UseSSL = &monitorSSL
	}
	if v, ok := d.GetOk("http_host_header"); ok {
		createMonitor.Properties.HTTP.HostHeader = v.(string)
	}
	if v, ok := d.GetOk("http_path"); ok {
		createMonitor.Properties.HTTP.URIPath = v.(string)
	}
	if v, ok := d.GetOk("http_authentication"); ok {
		createMonitor.Properties.HTTP.Authentication = v.(string)
	}
	if v, ok := d.GetOk("http_body_regex"); ok {
		createMonitor.Properties.HTTP.BodyRegex = v.(string)
	}
	if v, ok := d.GetOk("http_status_regex"); ok {
		createMonitor.Properties.HTTP.StatusRegex = v.(string)
	}
	if v, ok := d.GetOk("rtsp_body_regex"); ok {
		createMonitor.Properties.RTSP.BodyRegex = v.(string)
	}
	if v, ok := d.GetOk("rtsp_status_regex"); ok {
		createMonitor.Properties.RTSP.StatusRegex = v.(string)
	}
	if v, ok := d.GetOk("rtsp_path"); ok {
		createMonitor.Properties.RTSP.URIPath = v.(string)
	}
	if v, ok := d.GetOk("script_program"); ok {
		createMonitor.Properties.SCRIPT.Program = v.(string)
	}
	if v, ok := d.GetOk("script_arguments"); ok {
		if arguments, ok := v.(*schema.Set); ok {
			argumentsList := []monitor.ArgumentIssue{}
			for _, value := range arguments.List() {
				argumentsObject := value.(map[string]interface{})
				newArguments := monitor.ArgumentIssue{}
				if nameValue, ok := argumentsObject["name"].(string); ok {
					newArguments.Name = nameValue
				}
				if descriptionValue, ok := argumentsObject["description"].(string); ok {
					newArguments.Description = descriptionValue
				}
				if valueValue, ok := argumentsObject["value"].(string); ok {
					newArguments.Value = valueValue
				}
				argumentsList = append(argumentsList, newArguments)

			}
			createMonitor.Properties.SCRIPT.Arguments = argumentsList
		}
	}

	if v, ok := d.GetOk("sip_body_regex"); ok {
		createMonitor.Properties.SIP.BodyRegex = v.(string)
	}
	if v, ok := d.GetOk("sip_status_regex"); ok {
		createMonitor.Properties.SIP.StatusRegex = v.(string)
	}
	if v, ok := d.GetOk("sip_transport"); ok {
		createMonitor.Properties.SIP.Transport = v.(string)
	}
	if v, ok := d.GetOk("tcp_close_string"); ok {
		createMonitor.Properties.TCP.CloseString = v.(string)
	}
	if v, ok := d.GetOk("tcp_max_response_len"); ok {
		createMonitor.Properties.TCP.MaxResponseLen = uint(v.(int))
	}
	if v, ok := d.GetOk("tcp_response_regex"); ok {
		createMonitor.Properties.TCP.ResponseRegex = v.(string)
	}
	if v, ok := d.GetOk("tcp_write_string"); ok {
		createMonitor.Properties.TCP.WriteString = v.(string)
	}
	if v, ok := d.GetOk("udp_accept_all"); ok {
		monitorAcceptAll := v.(bool)
		createMonitor.Properties.UDP.AcceptAll = &monitorAcceptAll
	}
	createAPI := monitor.NewCreate(name, createMonitor)

	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Monitor error whilst creating %s: %v", name, err)
	}

	d.SetId(name)
	return resourceMonitorRead(d, m)

}

func resourceMonitorRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var readName string

	if v, ok := d.GetOk("name"); ok {
		readName = v.(string)
	}

	getSingleMonitorAPI := monitor.NewGet(readName)
	err := vtmClient.Do(getSingleMonitorAPI)
	if err != nil {
		if getSingleMonitorAPI.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM Monitor error whilst retrieving %s: %v", readName, err)
	}

	getMonitorProperties := getSingleMonitorAPI.ResponseObject().(*monitor.Monitor)
	d.Set("name", readName)
	d.Set("delay", getMonitorProperties.Properties.Basic.Delay)
	d.Set("timeout", getMonitorProperties.Properties.Basic.Timeout)
	d.Set("failures", getMonitorProperties.Properties.Basic.Failures)
	d.Set("verbose", getMonitorProperties.Properties.Basic.Verbose)
	d.Set("use_ssl", getMonitorProperties.Properties.Basic.UseSSL)
	d.Set("http_host_header", getMonitorProperties.Properties.HTTP.HostHeader)
	d.Set("http_path", getMonitorProperties.Properties.HTTP.URIPath)
	d.Set("http_authentication", getMonitorProperties.Properties.HTTP.Authentication)
	d.Set("http_body_regex", getMonitorProperties.Properties.HTTP.BodyRegex)
	d.Set("http_status_regex", getMonitorProperties.Properties.HTTP.StatusRegex)
	d.Set("rtsp_body_regex", getMonitorProperties.Properties.RTSP.BodyRegex)
	d.Set("rtsp_status_regex", getMonitorProperties.Properties.RTSP.StatusRegex)
	d.Set("rtsp_path", getMonitorProperties.Properties.RTSP.URIPath)
	d.Set("script_program", getMonitorProperties.Properties.SCRIPT.Program)
	d.Set("script_arguments", getMonitorProperties.Properties.SCRIPT.Arguments)
	d.Set("sip_body_regex", getMonitorProperties.Properties.SIP.BodyRegex)
	d.Set("sip_status_regex", getMonitorProperties.Properties.SIP.StatusRegex)
	d.Set("sip_transport", getMonitorProperties.Properties.SIP.Transport)
	d.Set("tcp_close_string", getMonitorProperties.Properties.TCP.CloseString)
	d.Set("tcp_max_response_len", getMonitorProperties.Properties.TCP.MaxResponseLen)
	d.Set("tcp_response_regex", getMonitorProperties.Properties.TCP.ResponseRegex)
	d.Set("tcp_write_string", getMonitorProperties.Properties.TCP.WriteString)
	d.Set("udp_accept_all", getMonitorProperties.Properties.UDP.AcceptAll)
	return nil
}

func resourceMonitorUpdate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var readName string
	var updateMonitor monitor.Monitor
	hasChanges := false

	if v, ok := d.GetOk("name"); ok {
		readName = v.(string)
	}
	if d.HasChange("delay") {
		if v, ok := d.GetOk("delay"); ok {
			delay := v.(int)
			updateMonitor.Properties.Basic.Delay = uint(delay)
		}
		hasChanges = true
	}
	if d.HasChange("timeout") {
		if v, ok := d.GetOk("timeout"); ok {
			timeout := v.(int)
			updateMonitor.Properties.Basic.Timeout = uint(timeout)
		}
		hasChanges = true
	}
	if d.HasChange("failures") {
		if v, ok := d.GetOk("failures"); ok {
			failures := v.(int)
			updateMonitor.Properties.Basic.Failures = uint(failures)
		}
		hasChanges = true
	}
	if d.HasChange("verbose") {
		monitorVerbosity := d.Get("verbose").(bool)
		updateMonitor.Properties.Basic.Verbose = &monitorVerbosity
		hasChanges = true
	}
	if d.HasChange("use_ssl") {
		monitorSSL := d.Get("use_ssl").(bool)
		updateMonitor.Properties.Basic.UseSSL = &monitorSSL
		hasChanges = true
	}
	if d.HasChange("http_host_header") {
		if v, ok := d.GetOk("http_host_header"); ok {
			updateMonitor.Properties.HTTP.HostHeader = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("http_path") {
		if v, ok := d.GetOk("http_path"); ok {
			updateMonitor.Properties.HTTP.URIPath = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("http_authentication") {
		if v, ok := d.GetOk("http_authentication"); ok {
			updateMonitor.Properties.HTTP.Authentication = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("http_body_regex") {
		if v, ok := d.GetOk("http_body_regex"); ok {
			updateMonitor.Properties.HTTP.BodyRegex = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("http_status_regex") {
		if v, ok := d.GetOk("http_status_regex"); ok {
			updateMonitor.Properties.HTTP.StatusRegex = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("rtsp_status_regex") {
		if v, ok := d.GetOk("rtsp_status_regex"); ok {
			updateMonitor.Properties.RTSP.StatusRegex = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("rtsp_body_regex") {
		if v, ok := d.GetOk("rtsp_body_regex"); ok {
			updateMonitor.Properties.RTSP.BodyRegex = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("rtsp_path") {
		if v, ok := d.GetOk("rtsp_path"); ok {
			updateMonitor.Properties.RTSP.URIPath = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("script_arguments") {
		if v, ok := d.GetOk("script_arguments"); ok {
			if arguments, ok := v.(*schema.Set); ok {
				argumentsList := []monitor.ArgumentIssue{}
				for _, value := range arguments.List() {
					argumentsObject := value.(map[string]interface{})
					newArguments := monitor.ArgumentIssue{}
					if nameValue, ok := argumentsObject["name"].(string); ok {
						newArguments.Name = nameValue
					}
					if descriptionValue, ok := argumentsObject["description"].(string); ok {
						newArguments.Description = descriptionValue
					}
					if valueValue, ok := argumentsObject["value"].(string); ok {
						newArguments.Value = valueValue
					}
					argumentsList = append(argumentsList, newArguments)

				}
				updateMonitor.Properties.SCRIPT.Arguments = argumentsList
			}
			hasChanges = true
		}
	}
	if d.HasChange("script_program") {
		if v, ok := d.GetOk("script_program"); ok {
			updateMonitor.Properties.SCRIPT.Program = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("sip_body_regex") {
		if v, ok := d.GetOk("sip_body_regex"); ok {
			updateMonitor.Properties.SIP.BodyRegex = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("sip_status_regex") {
		if v, ok := d.GetOk("sip_status_regex"); ok {
			updateMonitor.Properties.SIP.StatusRegex = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("sip_transport") {
		if v, ok := d.GetOk("sip_transport"); ok {
			updateMonitor.Properties.SIP.StatusRegex = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("tcp_close_string") {
		if v, ok := d.GetOk("tcp_close_string"); ok {
			updateMonitor.Properties.TCP.CloseString = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("tcp_max_response_len") {
		if v, ok := d.GetOk("tcp_max_response_len"); ok {
			updateMonitor.Properties.TCP.MaxResponseLen = uint(v.(int))
		}
		hasChanges = true
	}
	if d.HasChange("tcp_response_regex") {
		if v, ok := d.GetOk("tcp_response_regex"); ok {
			updateMonitor.Properties.TCP.ResponseRegex = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("tcp_write_string") {
		if v, ok := d.GetOk("tcp_write_string"); ok {
			updateMonitor.Properties.TCP.WriteString = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("udp_accept_all") {
		if v, ok := d.GetOk("udp_accept_all"); ok {
			monitorAcceptAll := v.(bool)
			updateMonitor.Properties.UDP.AcceptAll = &monitorAcceptAll
		}
		hasChanges = true
	}
	if hasChanges {
		updateAPI := monitor.NewUpdate(readName, updateMonitor)
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Monitor error whilst updating %s: %v", readName, err)
		}
	}
	return resourceMonitorRead(d, m)
}

func resourceMonitorDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var readName string

	if v, ok := d.GetOk("name"); ok {
		readName = v.(string)
	}

	deleteAPI := monitor.NewDelete(readName)
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("BrocadeVTM Monitor error whilst deleting %s: %v", readName, err)
	}

	d.SetId("")
	return nil
}
