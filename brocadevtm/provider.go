package brocadevtm

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
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
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_CLIENT_DEBUG", false),
				Description: "BrocadeVTM client debug",
			},
			"allow_unverified_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_ALLOW_UNVERIFIED_SSL", false),
				Description: "If set, BrocadeVTM client will permit unverifiable SSL certificates.",
			},
			"vtm_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_USERNAME", nil),
				Description: "User to authenticate with BrocadeVTM appliance",
			},
			"vtm_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_PASSWORD", nil),
				Description: "Password to authenticate with BrocadeVTM appliance",
			},
			"vtm_server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_SERVER", nil),
				Description: "Server to authenticate with BrocadeVTM appliance",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "3.8",
				Description: "BrocadevTM REST API Server version",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"brocadevtm_bandwidth":       resourceBandwidth(),
			"brocadevtm_dns_zone":        resourceDNSZone(),
			"brocadevtm_global_settings": resourceGlobalSettings(),
			"brocadevtm_dns_zone_file":   resourceDNSZoneFile(),
			"brocadevtm_glb":             resourceGLB(),
			"brocadevtm_location":        resourceLocation(),
			"brocadevtm_monitor":         resourceMonitor(),
			"brocadevtm_persistence":     resourcePersistence(),
			"brocadevtm_pool":            resourcePool(),
			"brocadevtm_rule":               resourceRule(),
			"brocadevtm_ssl_server_key": resourceSSLServerKey(),
			//	"brocadevtm_traffic_ip_group":   resourceTrafficIPGroup(),in
			"brocadevtm_user_authenticator": resourceUserAuthenticator(),
			"brocadevtm_user_group":         resourceUserGroup(),
			"brocadevtm_virtual_server":     resourceVirtualServer(),
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
		log.Println("Error connecting to Brocade REST Server: ", err)
		return nil, err
	}
	octetClient, err := api.Connect(octetConfig)
	if err != nil {
		log.Println("Error connecting to Brocade REST Server: ", err)
		return nil, err
	}

	config["jsonClient"] = jsonClient
	config["octetClient"] = octetClient
	return config, nil
}
