package brocadevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api"
	"log"
	"time"
)

// Provider is a basic structure that describes a provider: the configuration
// keys it takes, the resources it supports, a callback to configure, etc.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_debug": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_CLIENT_DEBUG", false),
				Description: "BrocadeVTM client debug",
			},
			"allow_unverified_ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_ALLOW_UNVERIFIED_SSL", false),
				Description: "If set, BrocadeVTM client will permit unverifiable SSL certificates.",
			},
			"vtm_user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_USERNAME", nil),
				Description: "User to authenticate with BrocadeVTM appliance",
			},
			"vtm_password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_PASSWORD", nil),
				Description: "Password to authenticate with BrocadeVTM appliance",
			},
			"vtm_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADEVTM_SERVER", nil),
				Description: "Server to authenticate with BrocadeVTM appliance",
			},
			"api_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "3.8",
				Description: "BrocadevTM REST API Server version",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"brocadevtm_dns_zone": resourceDNSZone(),
			/*
				"brocadevtm_dns_zone_file":      resourceDNSZoneFile(),
				"brocadevtm_glb":                resourceGLB(),
				"brocadevtm_location":           resourceLocation(),
				"brocadevtm_monitor":            resourceMonitor(),
				"brocadevtm_pool":               resourcePool(),
				"brocadevtm_rule":               resourceRule(),
				"brocadevtm_ssl_server_key":     resourceSSLServerKey(),
				"brocadevtm_traffic_ip_group":   resourceTrafficIPGroup(),
				"brocadevtm_user_authenticator": resourceUserAuthenticator(),
				"brocadevtm_user_group":         resourceUserGroup(),
				"brocadevtm_virtual_server":     resourceVirtualServer(),
			*/
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
		Headers:    map[string]string{"Content-Type": "application/octet-stream"},
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
