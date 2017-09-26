package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
	"github.com/sky-uk/terraform-provider-infoblox/infoblox/util"
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "GLB Algorithm",
				ValidateFunc: validateGLBAlgorithm,
			},
			"all_monitors_needed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether all assigned monitors in a location need to be working",
			},
			"auto_recovery": {
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
			"geo_effect": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				Description:  "How important the client's location is when deciding which location to use",
				ValidateFunc: validateGeoEffect,
			},
			"peer_health_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Reported monitor timeout in seconds",
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"return_ips_on_fail": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to return all IPs or none during a failure of all locations",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "The TTL for the DNS records handled by the GLB service",
			},
		},
	}
}

/*
type Basic struct {
	ChainedLocationOrder []string          `json:"chained_location_order,omitempty"`
	Rules                []string          `json:"rules,omitempty"`
	Domains              []string          `json:"domains,omitempty"`
	LastResortResponse   []string          `json:"last_resort_response,omitempty"`
	LocationDraining     []string          `json:"location_draining,omitempty"`
	LocationSettings     []LocationSetting `json:"location_settings,omitempty"`
	DNSSecKeys           []DNSSecKey       `json:"dnssec_keys,omitempty"`
}

// DNSSecKey : DNS Sec key struct
type DNSSecKey struct {
	Domain  string   `json:"domain,omitempty"`
	SSLKeys []string `json:"ssl_key,omitempty"`
}

// LocationSetting : settings for a location
type LocationSetting struct {
	Location string   `json:"location,omitempty"`
	Weight   uint     `json:"weight"`
	IPS      []string `json:"ips,omitempty"`
	Monitors []string `json:"monitors,omitempty"`
}

// Log : log configuration for a GLB
type Log struct {
	Enabled  bool   `json:"enabled"`
	Filename string `json:"filename,omitempty"`
	Format   string `json:"format,omitempty"`
}

*/

func validateGLBAlgorithm(v interface{}, k string) (ws []string, errors []error) {
	algorithm := v.(string)
	algorithmOptions := regexp.MustCompile(`^(chained|geo|hybrid|load|round_robin|weighted_random)$`)
	if !algorithmOptions.MatchString(algorithm) {
		errors = append(errors, fmt.Errorf("%q must be one of chained, geo, hybrid, load, round_robin or weighted_random", k))
	}
	return
}

func validateGeoEffect(v interface{}, k string) (ws []string, errors []error) {
	geoEffect := v.(int)
	if geoEffect < 0 || geoEffect > 100 {
		errors = append(errors, fmt.Errorf("%q must be a whole number between 0 and 100 (percentage)", k))
	}
	return
}

func resourceGLBCreate(d *schema.ResourceData, m interface{}) error {

	//d.SetId(name)
	return resourceGLBRead(d, m)
}

func resourceGLBRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceGLBUpdate(d *schema.ResourceData, m interface{}) error {

	//d.SetId(name)
	return resourceGLBRead(d, m)
}

func resourceGLBDelete(d *schema.ResourceData, m interface{}) error {

	return nil
}
