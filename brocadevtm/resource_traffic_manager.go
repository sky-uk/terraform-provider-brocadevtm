package brocadevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"fmt"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"net/http"
	"regexp"
)

func resourceTrafficManager() *schema.Resource {
	return &schema.Resource{

		Create: resourceTrafficManagerCreate,
		Read:   resourceTrafficManagerRead,
		Update: resourceTrafficManagerUpdate,
		Delete: resourceTrafficManagerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the traffic manager",
				Required:    true,
			},
			"adminMasterXMLIP": {
				Type:        schema.TypeString,
				Description: "The Application Firewall master XML IP",
				Optional:    true,
				Default:     "0.0.0.0",
			},
			"adminSlaveXMLIP": {
				Type:        schema.TypeString,
				Description: "The Application Firewall master XML IP",
				Optional:    true,
				Default:     "0.0.0.0",
			},
			"appliance_card": {
				Type:        schema.TypeSet,
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
				Type:        schema.TypeSet,
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
			"authenticationServerIP": {
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
			"numberOfCPUs": {
				Type:         schema.TypeInt,
				Description:  "The number of Application Firewall decider process to run.",
				Optional:     true,
				Default:      0,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"restServerPort": {
				Type:         schema.TypeInt,
				Description:  "The Application Firewall REST Internal API port, this port should not be accessed directly",
				Optional:     true,
				Default:      0,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"start_sysd": {
				Type:        schema.TypeBool,
				Description: "Whether or not to start the sysd process on software installations. Appliance and EC2 will always run sysd regardless of this config key",
				Computed:    true,
			},
			"trafficip": {
				Type:        schema.TypeSet,
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
			"updaterIP": {
				Type:        schema.TypeString,
				Description: "The Application Firewall Updater IP.",
				Optional:    true,
				Default:     "0.0.0.0",
			},
			"appliance": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
							Description: "Name (hostname.domainname) of the appliance",
							Optional:    true,
						},
						"hosts": {
							Type:        schema.TypeSet,
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
							Type:        schema.TypeSet,
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
							Type:        schema.TypeSet,
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
						"managereturnpath": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages return path routing. If disabled, the appliance won't modify iptables / rules / routes for this feature",
							Optional:    true,
							Default:     true,
						},
						"managesysctl": {
							Type:        schema.TypeBool,
							Description: "Whether or not the software manages user specified sysctl keys",
							Optional:    true,
							Computed:    true,
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
							Type:        schema.TypeSet,
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
						"shim_client_id": {
							Type:        schema.TypeString,
							Description: "The client ID provided by the portal for this server",
							Optional:    true,
						},
						"shim_client_key": {
							Type:        schema.TypeString,
							Description: "The client key provided by the portal for this server",
							Optional:    true,
						},
						"shim_enabled": {
							Type:        schema.TypeBool,
							Description: "Enable the Riverbed Cloud SteelHead discovery agent on this appliance",
							Optional:    true,
							Default:     false,
						},
						"shim_ips": {
							Type:        schema.TypeString,
							Description: "The IP addresses of the Riverbed Cloud SteelHeads to use, as a space or comma separated list. If using priority load balancing this should be in ascending order of priority (highest priority last)",
							Optional:    true,
						},
						"shim_load_balance": {
							Type:         schema.TypeString,
							Description:  "The load balancing method for selecting a Riverbed Cloud SteelHead appliance",
							Optional:     true,
							Default:      "round_robin",
							ValidateFunc: validation.StringInSlice([]string{"priority", "round_robin"}, false),
						},
						"shim_log_level": {
							Type:         schema.TypeString,
							Description:  "The minimum severity that the discovery agent will record to its log",
							Optional:     true,
							Default:      "notice",
							ValidateFunc: validation.StringInSlice([]string{"critical", "debug", "info", "notice", "serious", "warning"}, false),
						},
						"shim_mode": {
							Type:         schema.TypeString,
							Description:  "The mode used to discover Riverbed Cloud SteelHeads in the local cloud or data center",
							Optional:     true,
							Default:      "portal",
							ValidateFunc: validation.StringInSlice([]string{"local", "manual", "portal"}, false),
						},
						"shim_portal_url": {
							Type:        schema.TypeString,
							Description: "The hostname or IP address of the local portal to use",
							Optional:    true,
						},
						"shim_proxy_host": {
							Type:        schema.TypeString,
							Description: "The IP or hostname of the proxy server to use to connect to the portal. Leave blank to not use a proxy server",
							Optional:    true,
						},
						"shim_proxy_port": {
							Type:        schema.TypeString,
							Description: "The port of the proxy server, must be set if a proxy server has been configured",
							Optional:    true,
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
		},
	}
}

func validateApplianceCardLabel(v interface{}, k string) (ws []string, errors []error) {
	r := regexp.MustCompile(`^[\w.:@\-]{1,64}$`)
	if r.MatchString(v.(string)) {
		return
	}
	errors = append(errors, fmt.Errorf("Label must be a valid network interface card label"))
	return
}

func validateApplianceIFBond(v interface{}, k string) (ws []string, errors []error) {
	r := regexp.MustCompile(`^(bond\d+)?$`)
	if r.MatchString(v.(string)) {
		return
	}
	errors = append(errors, fmt.Errorf("Bond must match regex '^(bond\\d+)?$'"))
	return
}

func assignApplianceValues(v []interface{}) map[string]interface{} {
	values := v[0].(map[string]interface{})
	applianceValuesMap := make(map[string]interface{})

	tableNameList := []string{"hosts", "if", "ip", "routes"}
	for _, element := range tableNameList {
		if len(values[element].(*schema.Set).List()) > 0 {
			applianceValuesMap[element] = values[element].(*schema.Set).List()
		}
	}

	attributeNameList := []string{"gateway_ipv4", "gateway_ipv6", "ipmi_lan_access", "ipmi_lan_addr", "ipmi_lan_gateway", "ipmi_lan_ipsrc",
		"ipmi_lan_mask", "ipv4_forwarding", "ipv6_forwarding", "licence_agreed", "manageazureroutes", "manageec2conf", "manageiptrans",
		"managereturnpath", "managevpcconf", "name_servers", "ntpservers", "search_domains", "shim_client_id", "shim_client_key",
		"shim_enabled", "shim_ips", "shim_load_balance", "shim_log_level", "shim_mode", "shim_portal_url", "shim_proxy_host", "shim_proxy_port",
		"ssh_enabled", "ssh_password_allowed", "ssh_port", "timezone", "vlans"}

	for _, element := range attributeNameList {
		applianceValuesMap[element] = values[element]
	}

	return applianceValuesMap
}

func resourceTrafficManagerCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	trafficManagerConfiguration := make(map[string]interface{})
	trafficManagerPropertiesConfiguration := make(map[string]interface{})
	trafficManagerBasicConfiguration := make(map[string]interface{})

	name := d.Get("name").(string)

	tableAttributeList := []string{"appliance_card", "appliance_sysctl", "trafficip"}
	for _, element := range tableAttributeList {
		if v, ok := d.GetOk(element); ok {
			trafficManagerBasicConfiguration[element] = v.(*schema.Set).List()
		}
	}

	trafficManagerBasicConfiguration["adminMasterXMLIP"] = d.Get("adminMasterXMLIP").(string)
	trafficManagerBasicConfiguration["adminSlaveXMLIP"] = d.Get("adminSlaveXMLIP").(string)
	trafficManagerBasicConfiguration["authenticationServerIP"] = d.Get("authenticationServerIP").(string)
	trafficManagerBasicConfiguration["num_aptimizer_threads"] = d.Get("num_aptimizer_threads").(int)
	trafficManagerBasicConfiguration["num_children"] = d.Get("num_children").(int)
	trafficManagerBasicConfiguration["numberOfCPUs"] = d.Get("numberOfCPUs").(int)
	trafficManagerBasicConfiguration["restServerPort"] = d.Get("restServerPort").(int)
	trafficManagerBasicConfiguration["updaterIP"] = d.Get("updaterIP").(string)
	trafficManagerBasicConfiguration["location"] = d.Get("location").(string)
	trafficManagerBasicConfiguration["nameip"] = d.Get("nameip")

	if v, ok := d.GetOk("appliance"); ok {
		trafficManagerPropertiesConfiguration["appliance"] = assignApplianceValues(v.([]interface{}))
	}

	trafficManagerPropertiesConfiguration["basic"] = trafficManagerBasicConfiguration
	trafficManagerConfiguration["properties"] = trafficManagerPropertiesConfiguration

	err := client.Set("traffic_managers", name, &trafficManagerConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst creating Traffic Manager %s: %v", name, err)
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
		return fmt.Errorf("BrocadeVTM error whilst retrieving Traffic Manager %s: %v", name, err)
	}

	trafficManagerPropertiesConfig := trafficManagerConfiguration["properties"].(map[string]interface{})

	for i, element := range trafficManagerPropertiesConfig["basic"].(map[string]interface{}) {
		d.Set(i, element)
	}

	d.Set("appliance", trafficManagerPropertiesConfig["appliance"])

	return nil
}

func resourceTrafficManagerUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	trafficManagerConfiguration := make(map[string]interface{})
	trafficManagerPropertiesConfiguration := make(map[string]interface{})
	trafficManagerBasicConfiguration := make(map[string]interface{})

	tableAttributeList := []string{"appliance_card", "appliance_sysctl", "trafficip"}
	for _, element := range tableAttributeList {
		if d.HasChange(element) {
			trafficManagerBasicConfiguration[element] = d.Get(element).(*schema.Set).List()
		}
	}

	if d.HasChange("adminMasterXMLIP") {
		trafficManagerBasicConfiguration["adminMasterXMLIP"] = d.Get("adminMasterXMLIP")
	}
	if d.HasChange("adminSlaveXMLIP") {
		trafficManagerBasicConfiguration["adminSlaveXMLIP"] = d.Get("adminSlaveXMLIP").(string)
	}
	if d.HasChange("authenticationServerIP") {
		trafficManagerBasicConfiguration["authenticationServerIP"] = d.Get("authenticationServerIP").(string)
	}
	if d.HasChange("num_aptimizer_threads") {
		trafficManagerBasicConfiguration["num_aptimizer_threads"] = d.Get("num_aptimizer_threads").(int)
	}
	if d.HasChange("num_children") {
		trafficManagerBasicConfiguration["num_children"] = d.Get("num_children").(int)
	}
	if d.HasChange("numberOfCPUs") {
		trafficManagerBasicConfiguration["numberOfCPUs"] = d.Get("numberOfCPUs").(int)
	}
	if d.HasChange("restServerPort") {
		trafficManagerBasicConfiguration["restServerPort"] = d.Get("restServerPort").(int)
	}
	if d.HasChange("updaterIP") {
		trafficManagerBasicConfiguration["updaterIP"] = d.Get("updaterIP").(string)
	}
	if d.HasChange("location") {
		trafficManagerBasicConfiguration["location"] = d.Get("location").(string)
	}
	if d.HasChange("nameip") {
		trafficManagerBasicConfiguration["nameip"] = d.Get("nameip").(string)
	}

	if d.HasChange("appliance") {
		trafficManagerPropertiesConfiguration["appliance"] = assignApplianceValues(d.Get("appliance").([]interface{}))
	}

	trafficManagerPropertiesConfiguration["basic"] = trafficManagerBasicConfiguration
	trafficManagerConfiguration["properties"] = trafficManagerPropertiesConfiguration

	err := client.Set("traffic_managers", d.Id(), &trafficManagerConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst updating Traffic Manager %s: %v", d.Id(), err)
	}

	return resourceTrafficManagerRead(d, m)
}

func resourceTrafficManagerDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("traffic_managers", d, m)
}
