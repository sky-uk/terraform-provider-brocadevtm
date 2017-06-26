package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
		},

		ResourcesMap: map[string]*schema.Resource{
			"brocadevtm_monitor":          resourceMonitor(),
			"brocadevtm_pool":             resourcePool(),
			"brocadevtm_traffic_ip_group": resourceTrafficIPGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	clientDebug := d.Get("client_debug").(bool)
	allowUnverifiedSSL := d.Get("allow_unverified_ssl").(bool)

	vtmUser := d.Get("vtm_user").(string)
	if vtmUser == "" {
		return nil, fmt.Errorf("vtm_user must be provided")
	}

	vtmPassword := d.Get("vtm_password").(string)

	if vtmPassword == "" {
		return nil, fmt.Errorf("vtm_password must be provided")
	}

	vtmServer := d.Get("vtm_server").(string)

	if vtmServer == "" {
		return nil, fmt.Errorf("vtm_server must be provided")
	}

	config := Config{
		Debug:       clientDebug,
		Insecure:    allowUnverifiedSSL,
		VTMUser:     vtmUser,
		VTMPassword: vtmPassword,
		VTMServer:   vtmServer,
	}

	return config.Client()
}

// This is a global MutexKV for use within this plugin.
var vtmMutexKV = mutexkv.NewMutexKV()
