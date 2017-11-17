package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"net/http"
	"regexp"
)

func resourcePool() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolCreate,
		Read:   resourcePoolRead,
		Delete: resourcePoolDelete,
		Update: resourcePoolUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique name of the pool",
			},
			"bandwidth_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the bandwidth management class",
			},
			"failure_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the pool to use when all nodes in this pool have failed",
			},
			"max_connection_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
				Description:  "Maximum numberof nodes an attempt to send a request to befoirce returning an error to the client",
			},
			"max_idle_connections_pernode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
				Default:      50,
				Description:  "Maximum number of unused HTTP keepalive connections",
			},
			"max_timed_out_connection_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
				Default:      2,
				Description:  "Maxiumum failed connection attempts within the max_reply_time.",
			},
			"monitors": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of monitors to associate with this pool",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"node_close_with_rst": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not a connection to a node should be closed with a RST rather than a FIN packet",
				Default:     false,
			},
			"node_connection_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				Description:  "Number of times an attempt to connect to the same node before marking it as failed. Only used when passive_monitoring is enabled",
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"node_delete_behaviour": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "immediate",
				Description: "Node deletion behaviour for this pool",
				ValidateFunc: validation.StringInSlice([]string{
					"drain",
					"immediate",
				}, false),
			},
			"node_drain_to_delete_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The maximum time a node will remain in draining after it has been deleted",
				ValidateFunc: util.ValidateUnsignedInteger,
				Default:      0,
			},
			"nodes_table": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"nodes_list"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateNode,
							Description:  "A node. Combination of IP and port",
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: util.ValidateUnsignedInteger,
							Description:  "Priority assigned to a node. Defaults to 1",
						},
						"state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"active",
								"draining",
								"disabled",
							}, false),
							Description: "State of the node in the pool",
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 100),
							Description:  "Weight assigned to the node. Valid values are between 1 and 100",
						},
					},
				},
			},
			"nodes_list": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"nodes_table"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "Can be used in place of previous table when only the list of ip addresses is known",
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A note assigned to this pool",
			},
			"passive_monitoring": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not the software should check that real requests are working",
				Default:     false,
			},
			"persistence_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The session persistance class to use with this pool",
			},
			"transparent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not connections to the back ends appears to originate from the source client IP",
			},

			"auto_scaling": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addnode_delaytime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: util.ValidateUnsignedInteger,
							Description:  "Time the Traffic Manager should wait before adding node to autoscaled pool",
						},
						"cloud_credentials": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud API Credentials to use for authentication",
						},
						"cluster": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ESX host or cluster to place new VMs",
						},
						"data_center": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of logical vCenter server",
						},
						"data_store": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of VMWare data store",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set if all nodes in this pool are under auto-scaling control",
						},
						"external": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not auto-scaling is handled by an external system",
						},
						"hysteresis": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     20,
							Description: "Time period in seconds for which a change condition must persist prior to instigating the change",
						},
						"imageid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identifier for the image of the instances to create",
						},
						"ips_to_use": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "publicips",
							Description: "Type of IP addresses on the node to use",
							ValidateFunc: validation.StringInSlice([]string{
								"publicips",
								"private_ips",
							}, false),
						},
						"last_node_idle_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3600,
							Description: "Time node must be inactive before being destroyed",
						},
						"max_nodes": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     4,
							Description: "Maximum nodes in auto-scaled pool",
						},
						"min_nodes": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Minimum nodes in auto-scaled pool",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     1,
							Description: "The name prefix of the nodes in the auto-scaling group",
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      80,
							Description:  "Port number to use for each node in auto-scaled pool",
							ValidateFunc: util.ValidateTCPPort,
						},
						"refractory": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     180,
							Description: "Time after instigation of a change before any further changes made to the auto-scaled pool",
						},
						"response_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1000,
							Description: "Expected response time of nodes in ms",
						},
						"scale_down_level": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     95,
							Description: "Percentage of conforming requests above which the pool size is decresed",
						},
						"scale_up_level": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     40,
							Description: "Percentage of conforming requests below which the pool size is increased",
						},

						"securitygroupids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of security group IDs to assciate with a new ec2 instance",
							// When we're able to validate a list we should check each subnet ID starts with 'sg-'
							Elem: &schema.Schema{Type: schema.TypeString},
							Set:  schema.HashString,
						},
						"size_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identifier for the size of the instances to create",
						},
						"subnetids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of VPC subnet IDs where the new ec2 instances will be launched",
							// When we're able to validate a list we should check each subnet ID starts with 'subnet-'
							Elem: &schema.Schema{Type: schema.TypeString},
							Set:  schema.HashString,
						},
					},
				},
			},
			"pool_connection": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_connect_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      4,
							Description:  "How long to wait before giving up when attempting a connection to a node",
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"max_connections_per_node": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Max number of connections allowed to each back end node",
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"max_queue_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Max number connections that can be queued",
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"max_reply_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      30,
							Description:  "How long to wait for a response from a node",
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"queue_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "Max time to keep a connection queued",
							ValidateFunc: util.ValidateUnsignedInteger,
						},
					},
				},
			},
			"dns_autoscale": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the Traffic Manager will periodically resolve the hostnames using a DNS query",
						},
						"hostnames": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of hostnames which will be used for DNS derived autoscaling",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      80,
							Description:  "Port number to use for each node when using DNS dereived autoscaling",
							ValidateFunc: util.ValidateTCPPort,
						},
					},
				},
			},
			"ftp": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"support_rfc_2428": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the backed nodes understand EPRT and EPSV commands",
						},
					},
				},
			},
			"http": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keepalive": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not the pool should maintain HTTP keepalive connections to the nodes",
						},
						"keepalive_non_idempotent": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not the pool should maintain HTTP keepalive connections to the nodes for non-idempotent requests",
						},
					},
				},
			},
			"kerberos_protocol_transition": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Kerberos principle to use when performing Kerberos protocol transition",
						},
						"target": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Kerberos principle name",
						},
					},
				},
			},
			"load_balancing": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Load balancing algorithm to use",
							ValidateFunc: validation.StringInSlice([]string{
								"fastest_response_time",
								"least_connections",
								"perceptive",
								"random",
								"round_robin",
								"weighted_least_connections",
								"weighted_round_robin",
							}, false),
						},
						"priority_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not to enable priority lists",
						},
						"priority_nodes": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "Minimum number of active highest priority nodes",
						},
					},
				},
			},
			"node": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"close_on_death": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Close all connections to a failed node",
						},
						"retry_fail_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     60,
							Description: "Time a traffic manager will wait before retrying a failed node",
						},
					},
				},
			},
			"smtp": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_starttls": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Use when encrypting SMTP traffic",
						},
					},
				},
			},
			"ssl": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not a suitable certificate and private key from the SSL client certificates catalog can be used if the node requests authentication",
						},
						"common_name_match": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of names the common name can be matched against",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"elliptic_curves": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "SSL elliptic curver perference list",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not the pool should encrypt data before sending it to the node",
						},
						"enhance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allows for the traffic manager to prefix each new SSL connection with client information",
						},
						"send_close_alerts": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not to send SSL/TLS closer alert",
						},
						"server_name": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not to use the TLS 1.0 server_name extension",
						},
						"signature_algorithms": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSL signature algorithm preference list",
						},
						"ssl_ciphers": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSL/TLS ciphers to allow for connections to a node",
						},
						"ssl_support_ssl2": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "use_default",
							Description: "Whether or not SSLv2 is enabled",
							ValidateFunc: validation.StringInSlice([]string{
								"disabled",
								"enabled",
								"use_default",
							}, false),
						},
						"ssl_support_ssl3": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "use_default",
							Description: "Whether or not SSLv3 is enabled",
							ValidateFunc: validation.StringInSlice([]string{
								"disabled",
								"enabled",
								"use_default",
							}, false),
						},
						"ssl_support_tls1": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "use_default",
							Description: "Whether or not TLSv1.0 is enabled",
							ValidateFunc: validation.StringInSlice([]string{
								"disabled",
								"enabled",
								"use_default",
							}, false),
						},
						"ssl_support_tls1_1": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "use_default",
							Description: "Whether or not TLSv1.1 is enabled",
							ValidateFunc: validation.StringInSlice([]string{
								"disabled",
								"enabled",
								"use_default",
							}, false),
						},
						"ssl_support_tls1_2": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     false,
							Description: "Whether or not TLSv1.2 is enabled",
							ValidateFunc: validation.StringInSlice([]string{
								"disabled",
								"enabled",
								"use_default",
							}, false),
						},
						"strict_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not strict certificate verification should be performed",
						},
					},
				},
			},
			"tcp": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nagle": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not Nagle's algorithm should be used for connections to nodes",
						},
					},
				},
			},
			"udp": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accept_from": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "dest_only",
							Description: "IP addresses and ports from which responses to UDP requests should be accepted",
							ValidateFunc: validation.StringInSlice([]string{
								"all",
								"dest_ip_only",
								"dest_only",
								"ip_mask",
							}, false),
						},
						"accept_from_mask": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The CIDR mask which matches IPs we wat to receive responses from",
							ValidateFunc: validateAcceptFromMask,
						},
						"response_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Max time a node is permitted to take after receiving a UDP request",
						},
					},
				},
			},
		},
	}

}

// validateAcceptFromMask : check the assigned accept from mask is valid
func validateAcceptFromMask(v interface{}, k string) (ws []string, errors []error) {
	acceptFromMask := v.(string)
	acceptFromPattern := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+$`)
	if !acceptFromPattern.MatchString(acceptFromMask) {
		errors = append(errors, fmt.Errorf("%q must be in the format xxx.xxx.xxx.xxx/xx e.g. 10.0.0.0/8", k))
	}
	return
}

// validateNode : check a node is given in the correct format
func validateNode(v interface{}, k string) (ws []string, errors []error) {
	node := v.(string)
	validateNode := regexp.MustCompile(`[\w.-]+:\d{1,5}$`)
	if !validateNode.MatchString(node) {
		errors = append(errors, fmt.Errorf("%q must be a valid IP/Hostname and port seperated by a colon. i.e 127.0.0.1:80", k))
	}
	return
}

func getPoolMapAttributeList(mapName string) []string {

	var attributes []string

	switch mapName {
	case "basic":
		attributes = []string{"bandwidth_class",
			"failure_pool",
			"max_connection_attempts",
			"max_idle_connections_pernode",
			"max_timed_out_connection_attempts",
			"monitors",
			"node_close_with_rst",
			"node_connection_attempts",
			"node_drain_to_delete_timeout",
			"note",
			"passive_monitoring",
			"persistence_class",
			"transparent",
		}
	case "nodes_table":
		attributes = []string{"node",
			"priority",
			"state",
			"weight",
			"source_ip",
		}
	case "auto_scaling":
		attributes = []string{"addnode_delaytime",
			"addnode_delaytime",
			"cloud_credentials",
			"cluster",
			"data_center",
			"data_store",
			"enabled",
			"external",
			"hysteresis",
			"imageid",
			"ips_to_use",
			"last_node_idle_time",
			"max_nodes",
			"min_nodes",
			"name",
			"port",
			"refractory",
			"response_time",
			"scale_down_level",
			"scale_up_level",
			"securitygroupids",
			"size_id",
			"subnetids",
		}
	case "pool_connection":
		attributes = []string{"max_connect_time",
			"max_connections_per_node",
			"max_queue_size",
			"max_reply_time",
			"queue_timeout",
		}
	case "dns_autoscale":
		attributes = []string{"enabled", "hostnames", "port"}
	case "ftp":
		attributes = []string{"support_rfc_2428"}
	case "http":
		attributes = []string{"keepalive", "keepalive_non_idempotent"}
	case "kerberos_protocol_transition":
		attributes = []string{"principal", "target"}
	case "load_balancing":
		attributes = []string{"algorithm", "priority_enabled", "priority_nodes"}
	case "node":
		attributes = []string{"close_on_death", "retry_fail_time"}
	case "smtp":
		attributes = []string{"send_starttls"}
	case "ssl":
		attributes = []string{"client_auth",
			"common_name_match",
			"elliptic_curves",
			"enable",
			"enhance",
			"send_close_alerts",
			"server_name",
			"signature_algorithms",
			"ssl_ciphers",
			"ssl_support_ssl2",
			"ssl_support_ssl3",
			"ssl_support_tls1",
			"ssl_support_tls1_1",
			"ssl_support_tls1_2",
			"strict_verify",
		}
	case "sub_sections":
		attributes = []string{"auto_scaling",
			"dns_autoscale",
			"ftp",
			"http",
			"kerberos_protocol_transition",
			"load_balancing",
			"node",
			"smtp",
			"ssl",
			"tcp",
			"udp",
		}
	case "tcp":
		attributes = []string{"nagle"}
	case "udp":
		attributes = []string{"accept_from", "accept_from_mask", "response_timeout"}
	default:
		attributes = []string{}
	}
	return attributes
}

func buildNodesTableFromList(nodes interface{}) []map[string]interface{} {

	addresses := nodes.(*schema.Set).List()
	nodesTable := make([]map[string]interface{}, 0)

	for _, address := range addresses {
		node := make(map[string]interface{})
		node["node"] = address
		nodesTable = append(nodesTable, node)
	}
	return nodesTable
}

// resourcePoolCreate - Creates a  pool resource object
func resourcePoolCreate(d *schema.ResourceData, m interface{}) error {

	var nodesTableDefined, nodesListDefined bool
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	poolConfiguration := make(map[string]interface{})
	poolPropertiesConfiguration := make(map[string]interface{})

	poolName := d.Get("name").(string)

	// basic section
	poolBasicConfiguration := make(map[string]interface{})
	poolBasicConfiguration = util.AddSimpleGetAttributesToMap(d, poolBasicConfiguration, "", getPoolMapAttributeList("basic"))
	poolBasicConfiguration["node_delete_behavior"] = d.Get("node_delete_behaviour").(string)

	if v, ok := d.GetOk("nodes_table"); ok {
		poolBasicConfiguration["nodes_table"] = v.(*schema.Set).List()
		nodesTableDefined = true
	} else {
		if v, ok := d.GetOk("nodes_list"); ok {
			poolBasicConfiguration["nodes_table"] = buildNodesTableFromList(v)
			nodesListDefined = true
		}
	}
	if nodesTableDefined == false && nodesListDefined == false {
		return fmt.Errorf("Error creating resource: no one of nodes_table or nodes_list attr has been defined")
	}
	poolPropertiesConfiguration["basic"] = poolBasicConfiguration

	// connection section - we can't use "connection" as an attribute in the schema as it's reserved
	if v, ok := d.GetOk("pool_connection"); ok {
		poolPropertiesConfiguration["connection"] = v.(*schema.Set).List()[0]
	}

	// all other sections
	for _, section := range getPoolMapAttributeList("sub_sections") {
		if v, ok := d.GetOk(section); ok {
			builtList, err := util.BuildListMaps(v.(*schema.Set), getPoolMapAttributeList(section))
			if err != nil {
				return err
			}
			poolPropertiesConfiguration[section] = builtList[0]
		}
	}

	poolConfiguration["properties"] = poolPropertiesConfiguration
	err := client.Set("pools", poolName, poolConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Pool error whilst creating %s: %s", poolName, err)
	}

	d.SetId(poolName)
	return resourcePoolRead(d, m)
}

// resourcePoolRead - Reads a  pool resource
func resourcePoolRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	poolName := d.Id()
	poolConfiguration := make(map[string]interface{})

	client.WorkWithConfigurationResources()
	err := client.GetByName("pools", poolName, &poolConfiguration)
	if err != nil {
		if client.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM Pool error whilst retrieving %s: %s", poolName, err)
	}

	d.Set("name", poolName)
	poolPropertiesConfiguration := poolConfiguration["properties"].(map[string]interface{})
	poolBasicConfiguration := poolPropertiesConfiguration["basic"].(map[string]interface{})
	util.SetSimpleAttributesFromMap(d, poolBasicConfiguration, "", getPoolMapAttributeList("basic"))
	d.Set("node_delete_behaviour", poolBasicConfiguration["node_delete_behavior"])

	if _, ok := d.GetOk("nodes_list"); ok {
		var nodesList []string
		for _, item := range poolBasicConfiguration["nodes_table"].([]interface{}) {
			node := item.(map[string]interface{})
			nodesList = append(nodesList, node["node"].(string))
		}
		d.Set("nodes_list", nodesList)
	}
	d.Set("nodes_table", poolBasicConfiguration["nodes_table"])

	// pool_connection section
	poolSection := make([]map[string]interface{}, 0)
	poolSection = append(poolSection, poolPropertiesConfiguration["connection"].(map[string]interface{}))
	d.Set("pool_connection", poolSection)

	// all other sections
	for _, sectionName := range getPoolMapAttributeList("sub_sections") {
		section := make([]map[string]interface{}, 0)
		// sections with more complex structures need to be handled differently
		if sectionName == "auto_scaling" || sectionName == "ssl" || sectionName == "dns_autoscale" {
			autoScalingMapList, err := util.BuildReadListMaps(poolPropertiesConfiguration[sectionName].(map[string]interface{}), sectionName)
			if err != nil {
				return err
			}
			section = append(section, autoScalingMapList)
		} else {
			section = append(section, poolPropertiesConfiguration[sectionName].(map[string]interface{}))
		}
		d.Set(sectionName, section)
	}
	return nil
}

// resourcePoolUpdate - Updates an existing pool resource
func resourcePoolUpdate(d *schema.ResourceData, m interface{}) error {

	poolName := d.Id()
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	poolConfiguration := make(map[string]interface{})
	poolPropertiesConfiguration := make(map[string]interface{})

	// basic section
	poolBasicConfiguration := make(map[string]interface{})
	poolBasicConfiguration = util.AddChangedSimpleAttributesToMap(d, poolBasicConfiguration, "", getPoolMapAttributeList("basic"))
	poolBasicConfiguration["node_delete_behavior"] = d.Get("node_delete_behaviour").(string)

	if d.HasChange("nodes_table") {
		poolBasicConfiguration["nodes_table"] = d.Get("nodes_table").(*schema.Set).List()
	} else {
		if v, ok := d.GetOk("nodes_list"); ok {
			poolBasicConfiguration["nodes_table"] = buildNodesTableFromList(v)
		}
	}
	poolPropertiesConfiguration["basic"] = poolBasicConfiguration

	// connection section
	if d.HasChange("pool_connection") {
		poolPropertiesConfiguration["connection"] = d.Get("pool_connection").(*schema.Set).List()[0]
	}

	for _, section := range getPoolMapAttributeList("sub_sections") {
		if v, ok := d.GetOk(section); ok {
			builtList, err := util.BuildListMaps(v.(*schema.Set), getPoolMapAttributeList(section))
			if err != nil {
				return err
			}
			poolPropertiesConfiguration[section] = builtList[0]
		}
	}

	// all other sections
	for _, section := range getPoolMapAttributeList("sub_sections") {
		if d.HasChange(section) {
			builtList, err := util.BuildListMaps(d.Get(section).(*schema.Set), getPoolMapAttributeList(section))
			if err != nil {
				return err
			}
			poolPropertiesConfiguration[section] = builtList[0]
		}
	}

	poolConfiguration["properties"] = poolPropertiesConfiguration
	err := client.Set("pools", poolName, poolConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Pool error whilst updating %s: %s", poolName, err)
	}
	d.SetId(poolName)
	return resourcePoolRead(d, m)
}

// resourcePoolDelete - Deletes a pool resource
func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("pools", d, m)
}
