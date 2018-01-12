package pulsevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/sky-uk/go-pulse-vtm/api"
	"github.com/sky-uk/terraform-provider-pulsevtm/pulsevtm/util"
	"log"
	"net/http"
	"regexp"
)

func resourceTrafficManager() *schema.Resource {
	return &schema.Resource{

		Create: resourceTrafficManagerSet,
		Read:   resourceTrafficManagerRead,
		Update: resourceTrafficManagerSet,
		Delete: resourceTrafficManagerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the traffic manager",
				Required:    true,
			},
			"admin_master_xmlip": {
				Type:        schema.TypeString,
				Description: "The Application Firewall master XML IP",
				Optional:    true,
				Default:     "0.0.0.0",
			},
			"admin_slave_xmlip": {
				Type:        schema.TypeString,
				Description: "The Application Firewall master XML IP",
				Optional:    true,
				Default:     "0.0.0.0",
			},
			"appliance_card": {
				Type:        schema.TypeList,
				Description: "The table of network cards of a hardware appliance",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Network card PCI ID",
							Required:    true,
						},
						"interfaces": {
							Type:        schema.TypeList,
							Description: "The order of the interfaces of a network card",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"label": {
							Type:         schema.TypeString,
							Description:  "The labels of the installed network cards",
							Required:     true,
							ValidateFunc: validateApplianceCardLabel,
						},
					},
				},
			},
			"appliance_sysctl": {
				Type:        schema.TypeList,
				Description: "Custom kernel parameters applied by the user with sysctl interface",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysctl": {
							Type:        schema.TypeString,
							Description: "The name of the kernel parameter, e.g. net.ipv4.forward",
							Required:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Associated optional description for the sysctl",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the kernel parameter",
							Optional:    true,
						},
					},
				},
			},
			"authentication_server_ip": {
				Type:        schema.TypeString,
				Description: "The Application Firewall Authentication Server IP.",
				Optional:    true,
				Default:     "0.0.0.0",
			},
			"cloud_platform": {
				Type:        schema.TypeString,
				Description: "Cloud platform where the traffic manager is running",
				Optional:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "This is the location of the local traffic manager is in",
				Optional:    true,
			},
			"nameip": {
				Type:        schema.TypeString,
				Description: "Replace Traffic Manager name with an IP address",
				Optional:    true,
			},
			"num_aptimizer_threads": {
				Type:         schema.TypeInt,
				Description:  "How many worker threads the Web Accelerator process should create to optimise content. By default, one thread will be created for each CPU on the system.",
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"num_children": {
				Type:         schema.TypeInt,
				Description:  "The number of worker processes the software will run. By default, one child process will be created for each CPU on the system.",
				Optional:     true,
				Default:      0,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"number_of_cpus": {
				Type:         schema.TypeInt,
				Description:  "The number of Application Firewall decider process to run.",
				Optional:     true,
				Default:      0,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"rest_server_port": {
				Type:         schema.TypeInt,
				Description:  "The Application Firewall REST Internal API port, this port should not be accessed directly",
				Optional:     true,
				Default:      0,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"trafficip": {
				Type:        schema.TypeList,
				Description: "Custom kernel parameters applied by the user with sysctl interface",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "A network interface.",
							Required:    true,
						},
						"networks": {
							Type:        schema.TypeList,
							Description: "A set of IP/masks to which the network interface maps.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"updater_ip": {
				Type:        schema.TypeString,
				Description: "The Application Firewall Updater IP.",
				Optional:    true,
				Default:     "0.0.0.0",
			},
			"appliance": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_ipv4": {
							Type:        schema.TypeString,
							Description: "The default gateway",
							Optional:    true,
						},
						"gateway_ipv6": {
							Type:        schema.TypeString,
							Description: "The default IPv6 gateway",
							Optional:    true,
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "Name (hostname.domainname) of the appliance. This value is Read Only",
							Computed:    true,
						},
						"hosts": {
							Type:        schema.TypeList,
							Description: "A table of hostname to static ip address mappings, to be placed in the /etc/ hosts file",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "The name of a host",
										Required:    true,
									},
									"ip_address": {
										Type:        schema.TypeString,
										Description: "The static IP address of the host",
										Required:    true,
									},
								},
							},
						},
						"if": {
							Type:        schema.TypeList,
							Description: "A table of network interface specific settings",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "A network interface name",
										Required:    true,
									},
									"autoneg": {
										Type:        schema.TypeBool,
										Description: "Whether auto-negotiation should be enabled for the interface",
										Optional:    true,
										Default:     true,
									},
									"bmode": {
										Type:         schema.TypeString,
										Description:  "The trunk of which the interface should be a member",
										Optional:     true,
										Default:      "802_3ad",
										ValidateFunc: validation.StringInSlice([]string{"802_3ad"}, true),
									},
									"bond": {
										Type:         schema.TypeString,
										Description:  "The trunk of which the interface should be a member",
										Optional:     true,
										ValidateFunc: validateApplianceIFBond,
									},
									"duplex": {
										Type:        schema.TypeBool,
										Description: "Whether full-duplex should be enabled for the interface",
										Optional:    true,
										Default:     true,
									},
									"mode": {
										Type:         schema.TypeString,
										Description:  "Set the configuriation mode of an interface, the interface name is used in place of the * (asterisk).",
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"dhcp", "static"}, true),
									},
									"mtu": {
										Type:         schema.TypeInt,
										Description:  "The maximum transmission unit (MTU) of the interface",
										Optional:     true,
										Default:      1500,
										ValidateFunc: util.ValidateUnsignedInteger,
									},
									"speed": {
										Type:         schema.TypeString,
										Description:  "The speed of the interface",
										Optional:     true,
										Default:      "1000",
										ValidateFunc: validation.StringInSlice([]string{"10", "100", "1000"}, false),
									},
								},
							},
						},
						"ip": {
							Type:        schema.TypeList,
							Description: "A table of network interfaces and their network settings.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "A network interface name",
										Required:    true,
									},
									"addr": {
										Type:        schema.TypeString,
										Description: "The IP address for the interface",
										Required:    true,
									},
									"isexternal": {
										Type:        schema.TypeBool,
										Description: "Whether the interface is externally facing",
										Optional:    true,
										Default:     false,
									},
									"mask": {
										Type:        schema.TypeString,
										Description: "The IP mask (netmask) for the interface",
										Required:    true,
									},
								},
							},
						},
						"ipmi_lan_access": {
							Type:        schema.TypeBool,
							Description: "Whether IPMI LAN access should be enabled or not",
							Optional:    true,
							Default:     false,
						},
						"ipmi_lan_addr": {
							Type:        schema.TypeString,
							Description: "The IP address of the appliance IPMI LAN channel",
							Optional:    true,
						},
						"ipmi_lan_gateway": {
							Type:        schema.TypeString,
							Description: "The default gateway of the IPMI LAN channel",
							Optional:    true,
						},
						"ipmi_lan_ipsrc": {
							Type:         schema.TypeString,
							Description:  "The default gateway of the IPMI LAN channel",
							Optional:     true,
							Default:      "static",
							ValidateFunc: validation.StringInSlice([]string{"dhcp", "static"}, false),
						},
						"ipmi_lan_mask": {
							Type:        schema.TypeString,
							Description: "Set the IP netmask for the IPMI LAN channel",
							Optional:    true,
						},
						"ipv4_forwarding": {
							Type:        schema.TypeBool,
							Description: "Whether or not IPv4 forwarding is enabled",
							Optional:    true,
							Default:     false,
						},
						"ipv6_forwarding": {
							Type:        schema.TypeBool,
							Description: "Whether or not IPv6 forwarding is enabled",
							Optional:    true,
							Default:     false,
						},
						"licence_agreed": {
							Type:        schema.TypeBool,
							Description: "Whether or not the license agreement has been accepted.",
							Optional:    true,
							Default:     false,
						},
						"manageazureroutes": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages the Azure policy routing",
							Optional:    true,
							Default:     true,
						},
						"managedpa": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages the system configuration based on Data Plane Acceleration mode",
							Optional:    true,
							Default:     false,
						},
						"manageec2conf": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages the EC2 config",
							Optional:    true,
							Default:     true,
						},
						"manageiptrans": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages the IP transparency",
							Optional:    true,
							Default:     true,
						},
						"managereservedports": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages the system configuration for reserved ports.",
							Optional:    true,
							Default:     true,
						},
						"managereturnpath": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages return path routing. If disabled, the appliance won't modify iptables / rules / routes for this feature",
							Optional:    true,
							Default:     true,
						},
						"manageservices": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages the system services.",
							Optional:    true,
							Default:     true,
						},
						"managevpcconf": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages the EC2-VPC secondary IPs",
							Optional:    true,
							Default:     true,
						},
						"name_servers": {
							Type:        schema.TypeList,
							Description: "The IP addresses of the nameservers the appliance should use and place in /etc/resolv.conf",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ntpservers": {
							Type:        schema.TypeList,
							Description: "The NTP servers the appliance should use to synchronize its clock",
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"routes": {
							Type:        schema.TypeList,
							Description: "A table of destination IP addresses and routing details to reach them.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "A destination IP address",
										Required:    true,
									},
									"gw": {
										Type:        schema.TypeString,
										Description: "The gateway IP to configure for the route",
										Required:    true,
									},
									"if": {
										Type:        schema.TypeString,
										Description: "The network interface to configure for the route",
										Required:    true,
									},
									"mask": {
										Type:        schema.TypeString,
										Description: "The netmask to apply to the IP address",
										Required:    true,
									},
								},
							},
						},
						"search_domains": {
							Type:        schema.TypeList,
							Description: "The search domains the appliance should use and place in /etc/resolv.conf",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ssh_enabled": {
							Type:        schema.TypeBool,
							Description: "Whether or not the SSH server is enabled on the appliance",
							Optional:    true,
							Default:     true,
						},
						"ssh_password_allowed": {
							Type:        schema.TypeBool,
							Description: "Whether or not the SSH server allows password based login",
							Optional:    true,
							Default:     true,
						},
						"ssh_port": {
							Type:         schema.TypeInt,
							Description:  "The port that the SSH server should listen on",
							Optional:     true,
							Default:      22,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"timezone": {
							Type:        schema.TypeString,
							Description: "The timezone the appliance should use. This must be a path to a timezone file that exists under /usr/share/zoneinfo/",
							Optional:    true,
							Default:     "US/Pacific",
						},
						"vlans": {
							Type:        schema.TypeList,
							Description: "The VLANs the software should raise. ",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"cluster_comms": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_update": {
							Type:        schema.TypeBool,
							Description: "Whether or not this instance of the software can send configuration updates to other members of the cluster. When not clustered this key is ignored.",
							Optional:    true,
							Default:     false,
						},
						"bind_ip": {
							Type:        schema.TypeString,
							Description: "The IP address that the software should bind to for internal administration communications.",
							Optional:    true,
							Default:     "*",
						},
						"external_ip": {
							Type:        schema.TypeString,
							Description: "This is the optional external ip of the traffic manager, which is used to circumvent natting when traffic managers in a cluster span different networks.",
							Optional:    true,
						},
						"port": {
							Type:         schema.TypeInt,
							Description:  "The port that the software should listen on for internal administration communications.",
							Optional:     true,
							Default:      9080,
							ValidateFunc: util.ValidatePortNumber,
						},
					},
				},
			},
			"ec2": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trafficips_public_enis": {
							Type:        schema.TypeList,
							Description: "List of MAC addresses of interfaces which the traffic manager can use to associate the EC2 elastic IPs (Traffic IPs) to the instance.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"vpcid": {
							Type:        schema.TypeString,
							Description: "The ID of the VPC the instance is in, should be set when the appliance is first booted. Not required for non-VPC EC2 or non-EC2 systems.",
							Optional:    true,
						},
					},
				},
			},
			"fault_tolerance": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bgp_router_id": {
							Type:        schema.TypeString,
							Description: "The BGP router id If set to empty, then the IPv4 address used to communicate with the default IPv4 gateway is used instead. Specifying 0.0.0.0 will stop the traffic manager routing software from running the BGP protocol.",
							Optional:    true,
						},
						"lss_dedicated_ips": {
							Type:        schema.TypeList,
							Description: "IP addresses associated with the links dedicated by the user for receiving L4 state sync messages from other peers in a cluster.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ospfv2_ip": {
							Type:        schema.TypeString,
							Description: "The traffic manager's permanent IPv4 address which the routing software will use for peering and transit traffic, and as its OSPF router ID. If set to empty, then the address used to communicate with the default IPv4 gateway is used instead. Specifying 0.0.0.0 will stop the traffic manager routing software from running the OSPF protocol.",
							Optional:    true,
						},
						"ospfv2_neighbor_addrs": {
							Type:        schema.TypeList,
							Description: "The IP addresses of routers which are expected to be found as OSPFv2 neighbors of the traffic manager. The special  value %gateway% is a placeholder for the default gateway",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"iptables": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_enabled": {
							Type:        schema.TypeBool,
							Description: "This key overrides the product ID used by traffic manager instances to discover each other when clustering. Traffic managers will only discover each other if their product IDs are the same and their versions are compatible.",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
			"iptrans": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fwmark": {
							Type:         schema.TypeInt,
							Description:  "The netfilter forwarding mark to use for IP transparency rules",
							Optional:     true,
							Default:      320,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"iptables_enabled": {
							Type:        schema.TypeBool,
							Description: "Whether IP transparency may be used via netfilter/iptables. This require Linux 2.6.24 and the iptables socket extension",
							Optional:    true,
							Default:     false,
						},
						"routing_table": {
							Type:         schema.TypeInt,
							Description:  "The special routing table ID to use for IP transparency rules",
							Optional:     true,
							Default:      320,
							ValidateFunc: validation.IntBetween(256, 2147483647),
						},
					},
				},
			},
			"java": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:         schema.TypeInt,
							Description:  "The port the Java Extension handler process should listen on",
							Optional:     true,
							Default:      9060,
							ValidateFunc: validation.IntBetween(1024, 65535),
						},
					},
				},
			},
			"remote_licensing": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_address": {
							Type:        schema.TypeString,
							Description: "The e-mail address sent as part of a remote licensing request",
							Optional:    true,
						},
						"message": {
							Type:        schema.TypeString,
							Description: "A free-text field sent as part of a remote licensing request",
							Optional:    true,
						},
					},
				},
			},
			"rest_api": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_ips": {
							Type:        schema.TypeList,
							Description: "A list of IP Addresses which the REST API will listen on for connections. Read only",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"port": {
							Type:         schema.TypeInt,
							Description:  "The port on which the REST API should listen for requests",
							Optional:     true,
							Default:      9070,
							ValidateFunc: util.ValidatePortNumber,
						},
					},
				},
			},
			"snmp": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow": {
							Type:        schema.TypeList,
							Description: "restrict which IP addresses can access the SNMP command responder service. The value can be all, localhost, or a list of IP CIDR subnet masks",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"auth_password": {
							Type:        schema.TypeString,
							Description: "The authentication password. Required (minimum length 8 characters) if security_level includes authentication",
							Optional:    true,
							Sensitive:   true,
						},
						"bind_ip": {
							Type:        schema.TypeString,
							Description: "The IP address the SNMP service should bind its listen port to.  The value * (asterisk) means SNMP will listen on all IP addresses",
							Optional:    true,
							Default:     "*",
						},
						"community": {
							Type:        schema.TypeString,
							Description: "The community string required for SNMPv1 and SNMPv2c commands.  (If empty, all SNMPv1 and SNMPv2c commands will be rejected)",
							Optional:    true,
							Default:     "public",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether or not the SNMP command responder service should be enabled on this traffic manager",
							Optional:    true,
							Default:     false,
						},
						"hash_algorithm": {
							Type:         schema.TypeString,
							Description:  "The hash algorithm for authenticated SNMPv3 communications.",
							Optional:     true,
							Default:      "md5",
							ValidateFunc: validation.StringInSlice([]string{"md5", "sha1"}, false),
						},
						"port": {
							Type:        schema.TypeString,
							Description: "The port the SNMP command responder service should listen on. The value default denotes port 161 if the software is running with root privileges, and 1161 otherwise",
							Optional:    true,
							Default:     "default",
						},
						"priv_password": {
							Type:        schema.TypeString,
							Description: "The privacy password. Required (minimum length 8 characters) if security_level includes privacy (message encryption)",
							Optional:    true,
							Sensitive:   true,
						},
						"security_level": {
							Type:         schema.TypeString,
							Description:  "The security level for SNMPv3 communications",
							Optional:     true,
							Default:      "noauthnopriv",
							ValidateFunc: validation.StringInSlice([]string{"noauthnopriv", "authpriv", "authnopriv"}, false),
						},
						"username": {
							Type:        schema.TypeString,
							Description: "The username required for SNMP v3 commands.  (If empty, all SNMPv3 commands will be rejected)",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func validateApplianceCardLabel(v interface{}, k string) (ws []string, errors []error) {
	r := regexp.MustCompile(`^[\w.:@\-]{1,64}$`)
	if r.MatchString(v.(string)) {
		return
	}
	errors = append(errors, fmt.Errorf("[ERROR] Label must be a valid network interface card label"))
	return
}

func validateApplianceIFBond(v interface{}, k string) (ws []string, errors []error) {
	r := regexp.MustCompile(`^(bond\d+)?$`)
	if r.MatchString(v.(string)) {
		return
	}
	errors = append(errors, fmt.Errorf("[ERROR] Bond must match regex '^(bond\\d+)?$'"))
	return
}

func getTrafficManagerAttributeName(attribute string) string {
	switch attribute {
	case "adminMasterXMLIP":
		return "admin_master_xmlip"
	case "admin_master_xmlip":
		return "adminMasterXMLIP"
	case "adminSlaveXMLIP":
		return "admin_slave_xmlip"
	case "admin_slave_xmlip":
		return "adminSlaveXMLIP"
	case "authenticationServerIP":
		return "authentication_server_ip"
	case "authentication_server_ip":
		return "authenticationServerIP"
	case "number_of_cpus":
		return "numberOfCPUs"
	case "numberOfCPUs":
		return "number_of_cpus"
	case "restServerPort":
		return "rest_server_port"
	case "rest_server_port":
		return "restServerPort"
	case "updaterIP":
		return "updater_ip"
	case "updater_ip":
		return "updaterIP"
	default:
		return attribute
	}

}

func resourceTrafficManagerSet(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	trafficManagerConfiguration := make(map[string]interface{})
	trafficManagerPropertiesConfiguration := make(map[string]interface{})
	trafficManagerBasicConfiguration := make(map[string]interface{})

	name := d.Get("name").(string)

	for _, section := range []string{
		"appliance", "cluster_comms", "ec2", "fault_tolerance", "iptables", "iptrans", "java", "remote_licensing", "rest_api", "snmp",
	} {
		if d.HasChange(section) {
			trafficManagerPropertiesConfiguration[section] = d.Get(section).([]interface{})[0]
		}
	}

	for _, attribute := range []string{"appliance_card", "appliance_sysctl", "trafficip", "num_children", "location", "nameip", "admin_master_xmlip", "admin_slave_xmlip", "authentication_server_ip", "num_aptimizer_threads", "number_of_cpus", "rest_server_port", "updater_ip"} {
		if d.HasChange(attribute) {
			trafficManagerBasicConfiguration[getTrafficManagerAttributeName(attribute)] = d.Get(attribute)
		}
	}

	trafficManagerPropertiesConfiguration["basic"] = trafficManagerBasicConfiguration
	trafficManagerConfiguration["properties"] = trafficManagerPropertiesConfiguration

	err := client.Set("traffic_managers", name, &trafficManagerConfiguration, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM error whilst creating Traffic Manager %s: %v", name, err)
	}
	d.SetId(name)

	return resourceTrafficManagerRead(d, m)
}

func resourceTrafficManagerRead(d *schema.ResourceData, m interface{}) error {
	name := d.Id()
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	trafficManagerConfiguration := make(map[string]interface{})

	err := client.GetByName("traffic_managers", name, &trafficManagerConfiguration)

	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM error whilst retrieving Traffic Manager %s: %v", name, err)
	}

	trafficManagerPropertiesConfig := trafficManagerConfiguration["properties"].(map[string]interface{})

	basicTables := map[string]string{
		"appliance_card":   "name",
		"appliance_sysctl": "sysctl",
		"trafficip":        "name",
	}

	trafficManagerPropertiesConfig["basic"] = util.ReorderTablesInSection(trafficManagerPropertiesConfig, basicTables, "basic", d)
	for key, value := range trafficManagerPropertiesConfig["basic"].(map[string]interface{}) {
		err = d.Set(getTrafficManagerAttributeName(key), value)
		if err != nil {
			return fmt.Errorf("[ERROR] PulseVTM error whilst setting Traffic Manager attribute %s: %v", getTrafficManagerAttributeName(key), err)
		}
	}

	applianceTables := map[string]string{
		"if":     "name",
		"ip":     "name",
		"hosts":  "name",
		"routes": "name",
	}

	trafficManagerPropertiesConfig["appliance"] = util.ReorderTablesInSection(trafficManagerPropertiesConfig, applianceTables, "appliance", d)

	applianceSection := make([]interface{}, 0)
	applianceSection = append(applianceSection, trafficManagerPropertiesConfig["appliance"].(map[string]interface{}))

	err = d.Set("appliance", applianceSection)

	if err != nil {
		log.Println("[ERROR]  Response we're trying to set")
		spew.Dump(applianceSection)

		return fmt.Errorf("[ERROR] PulseVTM error whilst setting Traffic Manager attribute appliance: %v", err)
	}

	sectionNames := []string{"cluster_comms", "ec2", "fault_tolerance", "iptables", "iptrans", "java", "remote_licensing", "rest_api", "snmp"}

	for _, section := range sectionNames {
		sectionAsSliceOfInterfaces := make([]interface{}, 0)
		sectionAsSliceOfInterfaces = append(sectionAsSliceOfInterfaces, trafficManagerPropertiesConfig[section].(map[string]interface{}))

		err := d.Set(section, sectionAsSliceOfInterfaces)
		if err != nil {
			log.Println("[ERROR]  Response we're trying to set")
			spew.Dump(section)
			return fmt.Errorf("[ERROR] PulseVTM error whilst setting Traffic Manager attribute %s: %v", section, err)
		}
	}

	return nil
}

func resourceTrafficManagerDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("traffic_managers", d, m)
}
