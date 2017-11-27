package brocadevtm

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
)

func resourceGLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceGLBSet,
		Read:   resourceGLBRead,
		Update: resourceGLBSet,
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
				Type:        schema.TypeSet,
				Description: "Locations the GLB service operates in and the order in which locations fail",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rules": {
				Type:        schema.TypeSet,
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
						"ips": {
							Type:        schema.TypeSet,
							Description: "IP addresses in the location",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"monitors": {
							Type:        schema.TypeSet,
							Description: "Monitors used in the location",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"dnssec_keys": {
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
						"ssl_key": {
							Type:        schema.TypeSet,
							Description: "Keys for the associated domain",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"log": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether or not to log connections to this GLB service",
							Optional:    true,
							Default:     false,
						},
						"filename": {
							Type:        schema.TypeString,
							Description: "The filename the verbose query information should be logged to. Appliances will ignore this",
							Optional:    true,
							Default:     `%zeushome%/zxtm/log/services/%g.log`,
						},
						"format": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  `%t,%s,%l,%q,%g,%n,%d,%a`,
						},
					},
				},
			},
		},
	}
}

func resourceGLBSet(d *schema.ResourceData, m interface{}) error {

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
		basic["chained_location_order"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("rules"); ok {
		basic["rules"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("domains"); ok {
		basic["domains"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("last_resort_response"); ok {
		basic["last_resort_response"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("location_draining"); ok {
		basic["location_draining"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("location_settings"); ok {
		ls := v.(*schema.Set).List()
		locations := make([]map[string]interface{}, 0)

		for _, item := range ls {
			itemAsMap := item.(map[string]interface{})
			itemAsMap["ips"] = itemAsMap["ips"].(*schema.Set).List()
			itemAsMap["monitors"] = itemAsMap["monitors"].(*schema.Set).List()
			locations = append(locations, itemAsMap)
		}
		basic["location_settings"] = locations
	}

	if v, ok := d.GetOk("dnssec_keys"); ok {
		dks := v.(*schema.Set).List()
		dksAsList := make([]map[string]interface{}, 0)

		for _, item := range dks {
			itemAsMap := item.(map[string]interface{})
			itemAsMap["ssl_key"] = itemAsMap["ssl_key"].(*schema.Set).List()
			dksAsList = append(dksAsList, itemAsMap)
		}

		basic["dnssec_keys"] = dksAsList
	}

	props["basic"] = basic
	logs := d.Get("log").([]interface{})
	props["log"] = logs[0].(map[string]interface{})
	res["properties"] = props

	err := client.Set("glb_services", name, res, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM GLB error whilst creating %s: %v", name, err)
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
		return fmt.Errorf("[ERROR] BrocadeVTM GLB error whilst retrieving %s: %v", d.Id(), err)
	}

	pros := res["properties"].(map[string]interface{})
	basic := pros["basic"].(map[string]interface{})

	for _, attr := range []string{
		"algorithm",
		"all_monitors_needed",
		"autorecovery",
		"chained_auto_failback",
		"chained_location_order",
		"disable_on_failure",
		"dnssec_keys",
		"domains",
		"enabled",
		"geo_effect",
		"last_resort_response",
		"location_draining",
		"location_settings",
		"peer_health_timeout",
		"return_ips_on_fail",
		"rules",
		"ttl",
	} {
		d.Set(attr, basic[attr])
	}

	d.Set("log", pros["log"])

	return nil
}

func resourceGLBDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("glb_services", d, m)
}
