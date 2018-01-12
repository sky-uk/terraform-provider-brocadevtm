package pulsevtm

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-pulse-vtm/api"
)

// Provider is a basic structure that describes a provider: the configuration
// keys it takes, the resources it supports, a callback to configure, etc.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PULSEVTM_CLIENT_DEBUG", false),
				Description: "PulseVTM client debug",
			},
			"allow_unverified_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PULSEVTM_ALLOW_UNVERIFIED_SSL", false),
				Description: "If set, PulseVTM client will permit unverifiable SSL certificates.",
			},
			"vtm_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PULSEVTM_USERNAME", nil),
				Description: "User to authenticate with PulseVTM appliance",
			},
			"vtm_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PULSEVTM_PASSWORD", nil),
				Description: "Password to authenticate with PulseVTM appliance",
			},
			"vtm_server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PULSEVTM_SERVER", nil),
				Description: "Server to authenticate with PulseVTM appliance",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PULSEVTM_API_VERSION", "5.1"),
				Description: "PulsevTM REST API Server version",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"pulsevtm_appliance_nat":      resourceApplianceNat(),
			"pulsevtm_aptimizer_profile":  resourceAptimizerProfile(),
			"pulsevtm_bandwidth":          resourceBandwidth(),
			"pulsevtm_cloud_credentials":  resourceCloudCredentials(),
			"pulsevtm_dns_zone":           resourceDNSZone(),
			"pulsevtm_global_settings":    resourceGlobalSettings(),
			"pulsevtm_dns_zone_file":      resourceDNSZoneFile(),
			"pulsevtm_glb":                resourceGLB(),
			"pulsevtm_location":           resourceLocation(),
			"pulsevtm_monitor":            resourceMonitor(),
			"pulsevtm_persistence":        resourcePersistence(),
			"pulsevtm_pool":               resourcePool(),
			"pulsevtm_rule":               resourceRule(),
			"pulsevtm_ssl_cas_file":       resourceSSLCasFile(),
			"pulsevtm_ssl_client_key":     resourceSSLClientKey(),
			"pulsevtm_ssl_server_key":     resourceSSLServerKey(),
			"pulsevtm_traffic_manager":    resourceTrafficManager(),
			"pulsevtm_traffic_ip_group":   resourceTrafficIPGroup(),
			"pulsevtm_user_authenticator": resourceUserAuthenticator(),
			"pulsevtm_user_group":         resourceUserGroup(),
			"pulsevtm_virtual_server":     resourceVirtualServer(),
			"pulsevtm_ssl_ticket_key":     resourceSSLTicketKey(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	clientDebug := d.Get("client_debug").(bool)
	allowUnverifiedSSL := d.Get("allow_unverified_ssl").(bool)
	vtmUser := d.Get("vtm_user").(string)
	vtmPassword := d.Get("vtm_password").(string)
	vtmServer := d.Get("vtm_server").(string)
	apiVersion := d.Get("api_version").(string)
	var timeout time.Duration = 30

	config := make(map[string]interface{})

	octetHeaders := make(map[string]string)
	octetHeaders["Content-Type"] = "application/octet-stream"
	octetHeaders["Content-Transfer-Encoding"] = "text"

	jsonConfig := api.Params{
		APIVersion: apiVersion,
		Debug:      clientDebug,
		IgnoreSSL:  allowUnverifiedSSL,
		Username:   vtmUser,
		Password:   vtmPassword,
		Server:     vtmServer,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Timeout:    timeout,
	}

	octetConfig := api.Params{
		APIVersion: apiVersion,
		Debug:      clientDebug,
		IgnoreSSL:  allowUnverifiedSSL,
		Username:   vtmUser,
		Password:   vtmPassword,
		Server:     vtmServer,
		Headers:    octetHeaders,
		Timeout:    timeout,
	}

	jsonClient, err := api.Connect(jsonConfig)
	if err != nil {
		log.Println("Error connecting to Pulse REST Server: ", err)
		return nil, err
	}
	octetClient, err := api.Connect(octetConfig)
	if err != nil {
		log.Println("Error connecting to Pulse REST Server: ", err)
		return nil, err
	}

	config["jsonClient"] = jsonClient
	config["octetClient"] = octetClient
	return config, nil
}
