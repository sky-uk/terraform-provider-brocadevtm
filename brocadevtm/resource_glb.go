package brocadevtm

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/glb"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func resourceGLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceGLBCreate,
		Read:   resourceGLBRead,
		Update: resourceGLBUpdate,
		Delete: resourceGLBDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name of the GLB",
				ForceNew:    true,
			},
			"algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "hybrid",
				Description: "GLB Algorithm",
				ValidateFunc: validation.StringInSlice([]string{
					"chained",
					"geo",
					"hybrid",
					"load",
					"round_robin",
					"weighted_random"}, false),
			},
			"all_monitors_needed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether all assigned monitors in a location need to be working",
			},
			"autorecovery": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the last location to fail will be availble once it recovers",
			},
			"chained_auto_failback": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether automatic failback is enabled",
			},
			"disable_on_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Locations which recover from a failure will be disabled",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the GLB service is enabled or not",
			},
			"return_ips_on_fail": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to return all IPs or none during a failure of all locations",
			},
			"geo_effect": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				Description:  "How important the client's location is when deciding which location to use",
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "The TTL for the DNS records handled by the GLB service",
			},
			"chained_location_order": {
				Type:        schema.TypeList,
				Description: "Locations the GLB service operates in and the order in which locations fail",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rules": {
				Type:        schema.TypeList,
				Description: "A list of response rules to be applied to the GLB service",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"domains": {
				Type:        schema.TypeSet,
				Description: "A list of FQDN which should be used with this GLB service",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"last_resort_response": {
				Type:        schema.TypeSet,
				Description: "The response to send when all locations fail",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"location_draining": {
				Type:        schema.TypeSet,
				Description: "List of locations which are draining. No requests will be sent to these locations",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"location_settings": {
				Type:        schema.TypeSet,
				Description: "Table which contains location specific settings",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:        schema.TypeString,
							Description: "Location which the settings apply to",
							Optional:    true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Description:  "Weight to be given to this location when using the weighted random algorithm",
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							Description: "IP addresses in the location",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"monitors": {
							Type:        schema.TypeList,
							Description: "Monitors used in the location",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"dnssec_keys": { // TODO : should be "dnssec_keys"
				Type:        schema.TypeSet,
				Description: "Maps keys to domains",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Description: "Domain related to associated keys",
							Optional:    true,
						},
						"ssl_keys": {
							Type:        schema.TypeList,
							Description: "Keys for the associated domain",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"logging_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether or not to log connections to this GLB service",
				Optional:    true,
			},
			"log_file_name": {
				Type:        schema.TypeString,
				Description: "File to log to",
				Optional:    true,
				Computed:    true,
			},
			"log_format": {
				Type:        schema.TypeString,
				Description: "Format to us in log file",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func buildLocationSettings(locationSettingsSet *schema.Set) []glb.LocationSetting {

	locationSettingObjects := make([]glb.LocationSetting, 0)

	for _, locationSettingItem := range locationSettingsSet.List() {

		locationSetting := locationSettingItem.(map[string]interface{})
		locationSettingObject := glb.LocationSetting{}
		if location, ok := locationSetting["location"].(string); ok {
			locationSettingObject.Location = location
		}
		if weight, ok := locationSetting["weight"].(int); ok {
			locationSettingObject.Weight = uint(weight)
		}
		if ipAddresses, ok := locationSetting["ip_addresses"]; ok {
			locationSettingObject.IPS = util.BuildStringArrayFromInterface(ipAddresses)
		}
		if monitors, ok := locationSetting["monitors"]; ok {
			locationSettingObject.Monitors = util.BuildStringArrayFromInterface(monitors)
		}
		locationSettingObjects = append(locationSettingObjects, locationSettingObject)
	}
	return locationSettingObjects
}

func buildDNSSecKeys(dnsSecKeysSet *schema.Set) []glb.DNSSecKey {

	dnsSecKeyObjects := make([]glb.DNSSecKey, 0)

	for _, dnsSecItem := range dnsSecKeysSet.List() {

		dnsSec := dnsSecItem.(map[string]interface{})
		dnsSecObject := glb.DNSSecKey{}
		if domain, ok := dnsSec["domain"].(string); ok {
			dnsSecObject.Domain = domain
		}
		if sslKeys, ok := dnsSec["ssl_keys"]; ok {
			dnsSecObject.SSLKeys = util.BuildStringArrayFromInterface(sslKeys)
		}
		dnsSecKeyObjects = append(dnsSecKeyObjects, dnsSecObject)
	}
	return dnsSecKeyObjects
}

func resourceGLBCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	res := make(map[string]interface{})
	props := make(map[string]interface{})
	basic := make(map[string]interface{})

	name := d.Get("name").(string)

	basic["algorithm"] = d.Get("algorithm").(string)

	basic["all_monitors_needed"] = d.Get("all_monitors_needed").(bool)
	basic["autorecovery"] = d.Get("autorecovery").(bool)
	basic["chained_auto_failback"] = d.Get("chained_auto_failback").(bool)
	basic["disable_on_failure"] = d.Get("disable_on_failure").(bool)
	basic["enabled"] = d.Get("enabled").(bool)
	basic["return_ips_on_fail"] = d.Get("return_ips_on_fail").(bool)

	basic["geo_effect"] = d.Get("geo_effect").(int)
	basic["ttl"] = d.Get("ttl").(int)

	if v, ok := d.GetOk("chained_location_order"); ok {
		basic["chained_location_order"] = util.BuildStringArrayFromInterface(v)
	}
	if v, ok := d.GetOk("rules"); ok {
		basic["rules"] = util.BuildStringArrayFromInterface(v)
	}
	if v, ok := d.GetOk("domains"); ok {
		basic["domains"] = util.BuildStringListFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("last_resort_response"); ok {
		basic["last_resort_response"] = util.BuildStringListFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("location_draining"); ok {
		basic["location_draining"] = util.BuildStringListFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("location_settings"); ok {
		basic["location_settings"] = buildLocationSettings(v.(*schema.Set))
	}
	if v, ok := d.GetOk("dnssec_keys"); ok {
		basic["dnssec_keys"] = buildDNSSecKeys(v.(*schema.Set))
	}

	log := make(map[string]interface{})
	log["enabled"] = d.Get("logging_enabled").(bool)

	if v, ok := d.GetOk("log_file_name"); ok && v != "" {
		log["filename"] = v.(string)
	}
	if v, ok := d.GetOk("log_format"); ok && v != "" {
		log["format"] = v.(string)
	}

	props["basic"] = basic
	props["log"] = log
	res["properties"] = props

	err := client.Set("glb_services", name, res, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM GLB error whilst creating %s: %v", name, err)
	}
	d.SetId(name)
	return resourceGLBRead(d, m)
}

func resourceGLBRead(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()

	res := make(map[string]interface{})
	err := client.GetByName("glb_services", d.Id(), &res)

	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM GLB error whilst retrieving %s: %v", d.Id(), err)
	}

	// lists
	//chained_location_order
	//dnssec_keys
	//domains
	//last_resort_response
	//location_draining
	//location_settings
	//rules

	// scalars...
	for _, attr := range []string{
		"algorithm",
		"all_monitors_needed",
		"autorecovery",
		"chained_auto_failback",
		"disable_on_failure",
		"enabled",
		"geo_effect",
		"peer_health_timeout",
		"return_ips_on_fail",
		"ttl",
	} {
		d.Set(attr, res[attr])
	}

	//... and the log section...

	d.Set("algorithm", glbObject.Properties.Basic.Algorithm)
	d.Set("all_monitors_needed", glbObject.Properties.Basic.AllMonitorsNeeded)
	d.Set("autorecovery", glbObject.Properties.Basic.AutoRecovery)
	d.Set("chained_auto_failback", glbObject.Properties.Basic.ChainedAutoFailback)
	d.Set("disable_on_failure", glbObject.Properties.Basic.DisableOnFailure)
	d.Set("enabled", glbObject.Properties.Basic.Enabled)
	d.Set("return_ips_on_fail", glbObject.Properties.Basic.ReturnIPSOnFail)
	d.Set("ttl", glbObject.Properties.Basic.TTL)
	d.Set("geo_effect", glbObject.Properties.Basic.GeoEffect)
	d.Set("chained_location_order", glbObject.Properties.Basic.ChainedLocationOrder)
	d.Set("rules", glbObject.Properties.Basic.Rules)
	d.Set("domains", glbObject.Properties.Basic.Domains)
	d.Set("last_resort_response", glbObject.Properties.Basic.LastResortResponse)
	d.Set("location_draining", glbObject.Properties.Basic.LocationDraining)
	d.Set("location_settings", glbObject.Properties.Basic.LocationSettings)
	d.Set("dnssec_keys", glbObject.Properties.Basic.DNSSecKeys)
	d.Set("logging_enabled", glbObject.Properties.Log.Enabled)
	d.Set("log_file_name", glbObject.Properties.Log.Filename)
	d.Set("log_format", glbObject.Properties.Log.Format)

	return nil
}

func resourceGLBUpdate(d *schema.ResourceData, m interface{}) error {

	hasChanges := false
	name := d.Id()
	var updateGLB glb.GLB

	if d.HasChange("algorithm") {
		if v, ok := d.GetOk("algorithm"); ok && v != "" {
			updateGLB.Properties.Basic.Algorithm = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("all_monitors_needed") {
		hasChanges = true
	}
	updateGLB.Properties.Basic.AllMonitorsNeeded = d.Get("all_monitors_needed").(bool)

	if d.HasChange("autorecovery") {
		hasChanges = true
	}
	updateGLB.Properties.Basic.AutoRecovery = d.Get("autorecovery").(bool)

	if d.HasChange("chained_auto_failback") {
		hasChanges = true
	}
	updateGLB.Properties.Basic.ChainedAutoFailback = d.Get("chained_auto_failback").(bool)

	if d.HasChange("disable_on_failure") {
		hasChanges = true
	}
	updateGLB.Properties.Basic.DisableOnFailure = d.Get("disable_on_failure").(bool)

	if d.HasChange("enabled") {
		hasChanges = true
	}
	updateGLB.Properties.Basic.Enabled = d.Get("enabled").(bool)

	if d.HasChange("return_ips_on_fail") {
		hasChanges = true
	}
	updateGLB.Properties.Basic.ReturnIPSOnFail = d.Get("return_ips_on_fail").(bool)

	if d.HasChange("geo_effect") {
		hasChanges = true
	}
	geoEffect := d.Get("geo_effect").(int)
	updateGLB.Properties.Basic.GeoEffect = uint(geoEffect)

	if d.HasChange("ttl") {
		hasChanges = true
	}
	updateGLB.Properties.Basic.TTL = d.Get("ttl").(int)

	if d.HasChange("chained_location_order") {
		if v, ok := d.GetOk("chained_location_order"); ok {
			updateGLB.Properties.Basic.ChainedLocationOrder = util.BuildStringArrayFromInterface(v)
		}
		hasChanges = true
	}
	if d.HasChange("rules") {
		if v, ok := d.GetOk("rules"); ok {
			updateGLB.Properties.Basic.Rules = util.BuildStringArrayFromInterface(v)
		}
		hasChanges = true
	}
	if d.HasChange("domains") {
		if v, ok := d.GetOk("domains"); ok {
			updateGLB.Properties.Basic.Domains = util.BuildStringListFromSet(v.(*schema.Set))
		}
		hasChanges = true
	}
	if d.HasChange("last_resort_response") {
		if v, ok := d.GetOk("last_resort_response"); ok {
			updateGLB.Properties.Basic.LastResortResponse = util.BuildStringListFromSet(v.(*schema.Set))
		}
		hasChanges = true
	}
	if d.HasChange("location_draining") {
		if v, ok := d.GetOk("location_draining"); ok {
			updateGLB.Properties.Basic.LocationDraining = util.BuildStringListFromSet(v.(*schema.Set))
		}
		hasChanges = true
	}
	if d.HasChange("location_settings") {
		if v, ok := d.GetOk("location_settings"); ok {
			updateGLB.Properties.Basic.LocationSettings = buildLocationSettings(v.(*schema.Set))
		}
	}
	if d.HasChange("dnssec_keys") {
		if v, ok := d.GetOk("dnssec_keys"); ok {
			updateGLB.Properties.Basic.DNSSecKeys = buildDNSSecKeys(v.(*schema.Set))
		}
	}
	if d.HasChange("logging_enabled") {
		hasChanges = true
	}
	updateGLB.Properties.Log.Enabled = d.Get("logging_enabled").(bool)

	if d.HasChange("log_file_name") {
		if v, ok := d.GetOk("log_file_name"); ok && v != "" {
			updateGLB.Properties.Log.Filename = v.(string)
		}
	}
	if d.HasChange("log_format") {
		if v, ok := d.GetOk("log_format"); ok && v != "" {
			updateGLB.Properties.Log.Format = v.(string)
		}
	}

	if hasChanges {
		config := m.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		err := client.Set("glb_services", name, updateGLB, nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM GLB error whilst updating %s: %v", name, err)
		}
	}
	d.SetId(name)
	return resourceGLBRead(d, m)
}

func resourceGLBDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("glb_services", d, m)
}
