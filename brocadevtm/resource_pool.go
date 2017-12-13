package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"log"
	"net/http"
	"regexp"
)

func resourcePool() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolSet,
		Read:   resourcePoolRead,
		Update: resourcePoolSet,
		Delete: resourcePoolDelete,

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
			"node_delete_behavior": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "immediate",
				Description: "Node deletion behavior for this pool",
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
							Default:     1,
							Description: "Minimum nodes in auto-scaled pool",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
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
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      10,
							Description:  "Max time to keep a connection queued",
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
			"l4accel": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snat": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether connections to the back-end nodes should appear to originate from an IP address raised on the traffic manager, rather than the IP address from which they were received by the traffic manager.",
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
						"cipher_suites": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The SSL/TLS cipher suites to allow for connections to a back-end node",
						},
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
						"session_cache_enabled": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "use_default",
							Description: "Whether or not the SSL client cache will be used for this pool",
							ValidateFunc: validation.StringInSlice([]string{
								"disabled",
								"enabled",
								"use_default",
							}, false),
						},
						"signature_algorithms": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSL signature algorithm preference list",
						},
						"strict_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not strict certificate verification should be performed",
						},
						"support_ssl3": {
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
						"support_tls1": {
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
						"support_tls1_1": {
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
						"support_tls1_2": {
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
		errors = append(errors, fmt.Errorf("[ERROR] %q must be in the format xxx.xxx.xxx.xxx/xx e.g. 10.0.0.0/8", k))
	}
	return
}

// validateNode : check a node is given in the correct format
func validateNode(v interface{}, k string) (ws []string, errors []error) {
	node := v.(string)
	validateNode := regexp.MustCompile(`[\w.-]+:\d{1,5}$`)
	if !validateNode.MatchString(node) {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be a valid IP/Hostname and port seperated by a colon. i.e 127.0.0.1:80", k))
	}
	return
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

func basicPoolKeys() []string {
	return []string{
		"bandwidth_class",
		"failure_pool",
		"max_connection_attempts",
		"max_idle_connections_pernode",
		"max_timed_out_connection_attempts",
		"monitors",
		"node_close_with_rst",
		"node_connection_attempts",
		"node_delete_behavior",
		"node_drain_to_delete_timeout",
		"note",
		"passive_monitoring",
		"persistence_class",
		"transparent",
	}
}

func poolSectionName(name string) string {
	if name == "pool_connection" {
		return "connection"
	}
	if name == "connection" {
		return "pool_connection"
	}
	return name
}

func resourcePoolSet(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	var nodesTableDefined bool

	name := d.Get("name").(string)

	poolRequest := make(map[string]interface{})
	poolProperties := make(map[string]interface{})

	util.GetSection(d, "basic", poolProperties, basicPoolKeys())

	for _, section := range []string{
		"auto_scaling",
		"dns_autoscale",
		"ftp",
		"http",
		"kerberos_protocol_transition",
		"load_balancing",
		"node",
		"pool_connection",
		"smtp",
		"ssl",
		"tcp",
		"udp",
		"l4accel",
	} {
		if d.HasChange(section) {
			poolProperties[poolSectionName(section)] = d.Get(section).(*schema.Set).List()[0]
		}
	}

	nodesTable := d.Get("nodes_table")
	nodesList := d.Get("nodes_list")
	if d.HasChange("nodes_table") {
		poolProperties["basic"].(map[string]interface{})["nodes_table"] = nodesTable.(*schema.Set).List()
		nodesTableDefined = true
		log.Printf(fmt.Sprintf("[DEBUG] Nodes table is %+v", nodesTable.(*schema.Set).List()))
	}
	// We only want to use nodes_list when nodes_table hasn't been defined.
	if nodesTableDefined == false {
		if d.HasChange("nodes_list") {
			poolProperties["basic"].(map[string]interface{})["nodes_table"] = buildNodesTableFromList(nodesList)
			//nodesListDefined = true
			log.Printf(fmt.Sprintf("[DEBUG] Nodes list is %+v", buildNodesTableFromList(nodesList)))
		}
	}

	if len(nodesTable.(*schema.Set).List()) == 0 && len(nodesList.(*schema.Set).List()) == 0 {
		return fmt.Errorf("[ERROR] creating/updating resource: one of nodes_list or nodes_table must be defined")
	}

	poolRequest["properties"] = poolProperties
	util.TraverseMapTypes(poolRequest)
	err := client.Set("pools", name, poolRequest, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM Pool error whilst creating/updating %s: %s", name, err)
	}
	d.SetId(name)

	return resourcePoolRead(d, m)
}

func resourcePoolRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	poolResponse := make(map[string]interface{})

	client.WorkWithConfigurationResources()
	err := client.GetByName("pools", d.Id(), &poolResponse)
	if err != nil {
		if client.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] BrocadeVTM Pools error whilst retrieving %s: %v", d.Id(), err)
	}

	poolsProperties := poolResponse["properties"].(map[string]interface{})
	poolsBasic := poolsProperties["basic"].(map[string]interface{})

	for _, key := range basicPoolKeys() {
		err := d.Set(key, poolsBasic[key])
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM Pools error whilst setting attribute %s in state", key)
		}
	}

	if _, ok := d.GetOk("nodes_list"); ok {
		var nodesList []string
		for _, item := range poolsBasic["nodes_table"].([]interface{}) {
			node := item.(map[string]interface{})
			nodesList = append(nodesList, node["node"].(string))
		}
		err := d.Set("nodes_list", nodesList)
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM Pools error whilst setting attribute nodes_list in state")
		}
	}
	err = d.Set("nodes_table", poolsBasic["nodes_table"])
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM Pools error whilst setting attribute nodes_table in state")
	}

	for _, section := range []string{
		"auto_scaling",
		"dns_autoscale",
		"ftp",
		"http",
		"kerberos_protocol_transition",
		"load_balancing",
		"node",
		"connection",
		"smtp",
		"ssl",
		"tcp",
		"udp",
		"l4accel",
	} {
		set := make([]map[string]interface{}, 0)
		//set = append(set, poolsProperties[section].(map[string]interface{}))
		readSectionMap, err := util.BuildReadMap(poolsProperties[section].(map[string]interface{}))
		if err != nil {
			return err
		}
		set = append(set, readSectionMap)
		err = d.Set(poolSectionName(section), set)
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM Pools error whilst setting attribute %s in state", section)
		}
	}
	return nil
}

// resourcePoolDelete - Deletes a pool resource
func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("pools", d, m)
}
