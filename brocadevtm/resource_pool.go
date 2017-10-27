package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/pool"
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
				Computed:     true,
				Description:  "Maximum number of unused HTTP keepalive connections",
			},
			"max_timed_out_connection_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
				Computed:     true,
				Description:  "Maxiumum failed connection attempts within the max_reply_time.",
			},
			"monitors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of monitors to associate with this pool",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"node_close_with_rst": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not a connection to a node should be closed with a RST rather than a FIN packet",
				Computed:    true,
			},
			"node_connection_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Number of times an attempt to connect to the same node before marking it as failed. Only used when passive_monitoring is enabled",
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"node_delete_behaviour": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "immediate",
				Description:  "Node deletion behaviour for this pool",
				ValidateFunc: validateNodeDeleteBehaviour,
			},
			"node_drain_to_delete_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The maximum time a node will remain in draining after it has been deleted",
				ValidateFunc: util.ValidateUnsignedInteger,
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
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateState,
							Description:  "State of the node in the pool",
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validateWeight,
							Description:  "Weight assigned to the node. Valid values are between 1 and 100",
						},
						"source_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: util.ValidateIP,
							Description:  "Source IP the Traffic Manager uses to connect to this node",
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
				Computed:    true,
			},
			"persistence_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The session persistance class to use with this pool",
			},
			"transparent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether or not connections to the back ends appears to originate from the source client IP",
			},

			"auto_scaling": {
				Type:     schema.TypeList,
				Optional: true,
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
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "publicips",
							Description:  "Type of IP addresses on the node to use",
							ValidateFunc: validateAutoScalingIPsToUse,
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
						"refactory": {
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
							Elem:        &schema.Schema{Type: schema.TypeString},
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
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"pool_connection": {
				Type:     schema.TypeList,
				Optional: true,
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
						},
					},
				},
			},
			"dns_autoscale": {
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Load balancing algorithm to use",
							ValidateFunc: validatePoolLBAlgo,
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
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
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
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "use_default",
							Description:  "Whether or not SSLv2 is enabled",
							ValidateFunc: validateSSLSupportOptions,
						},
						"ssl_support_ssl3": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "use_default",
							Description:  "Whether or not SSLv3 is enabled",
							ValidateFunc: validateSSLSupportOptions,
						},
						"ssl_support_tls1": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "use_default",
							Description:  "Whether or not TLSv1.0 is enabled",
							ValidateFunc: validateSSLSupportOptions,
						},
						"ssl_support_tls1_1": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "use_default",
							Description:  "Whether or not TLSv1.1 is enabled",
							ValidateFunc: validateSSLSupportOptions,
						},
						"ssl_support_tls1_2": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      false,
							Description:  "Whether or not TLSv1.2 is enabled",
							ValidateFunc: validateSSLSupportOptions,
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
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accept_from": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "dest_only",
							Description:  "IP addresses and ports from which responses to UDP requests should be accepted",
							ValidateFunc: validateUDPAcceptFrom,
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

// validateSSLSupportOptions : check the assigned SSL support choice is valid
func validateSSLSupportOptions(v interface{}, k string) (ws []string, errors []error) {
	ssl2Support := v.(string)
	ssl2SupportOptions := regexp.MustCompile(`^(disabled|enabled|use_default)$`)
	if !ssl2SupportOptions.MatchString(ssl2Support) {
		errors = append(errors, fmt.Errorf("%q must be one of disabled, enabled or use_default", k))
	}
	return
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

// validateUDPAcceptFrom : checks the assigned UDP accept from choice is valid
func validateUDPAcceptFrom(v interface{}, k string) (ws []string, errors []error) {
	acceptFrom := v.(string)
	acceptFromOptions := regexp.MustCompile(`^(all|dest_ip_only|dest_only|ip_mask)`)
	if !acceptFromOptions.MatchString(acceptFrom) {
		errors = append(errors, fmt.Errorf("%q must be one of all, dest_ip_only, dest_only or ip_mask", k))
	}
	return
}

// validateAutoScalingIPsToUse : check the assigned auto scaling IPs to use is a valid choice
func validateAutoScalingIPsToUse(v interface{}, k string) (ws []string, errors []error) {
	ipType := v.(string)
	ipTypeOptions := regexp.MustCompile(`^(publicips|private_ips)`)
	if !ipTypeOptions.MatchString(ipType) {
		errors = append(errors, fmt.Errorf("%q must be one of publicips or private_ips", k))
	}
	return
}

// validateNodeDeleteBehaviour : check the assigned node delete behaviour is a valid choice
func validateNodeDeleteBehaviour(v interface{}, k string) (ws []string, errors []error) {
	behaviour := v.(string)
	behvaiourOptions := regexp.MustCompile(`^(immediate|drain)$`)
	if !behvaiourOptions.MatchString(behaviour) {
		errors = append(errors, fmt.Errorf("%q must be one of immediate or drain", k))
	}
	return
}

// validatePoolLBAlgo : check the assigned algorithm is a valid choice
func validatePoolLBAlgo(v interface{}, k string) (ws []string, errors []error) {
	algo := v.(string)
	algoOptions := regexp.MustCompile(`^(fastest_response_time|least_connections|perceptive|random|round_robin|weighted_least_connections|weighted_round_robin)$`)
	if !algoOptions.MatchString(algo) {
		errors = append(errors, fmt.Errorf("%q must be one of fastest_response_time, least_connections, perceptive, random, round_robin, weighted_least_connections, weighted_round_robin", k))
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

// validateWeight : check the assigned weight is a valid choice
func validateWeight(v interface{}, k string) (ws []string, errors []error) {
	weight := v.(int)

	if weight < 1 || weight > 100 {
		errors = append(errors, fmt.Errorf("%q must be between 1-100", k))
	}
	return
}

// validateState : check the assigned state is a valid choice
func validateState(v interface{}, k string) (ws []string, errors []error) {
	state := v.(string)
	stateOptions := regexp.MustCompile(`^(active|draining|disabled)$`)
	if !stateOptions.MatchString(state) {
		errors = append(errors, fmt.Errorf("%q must be one of active, draining, disabled", k))
	}
	return
}

// checkStringPrefix : check a string starts with a given prefix
func checkStringPrefix(prefix string, list []string) error {
	for _, item := range list {
		checkFormat := regexp.MustCompile(`^` + prefix)
		if !checkFormat.MatchString(item) {
			return fmt.Errorf("one or more items in the list of strings doesn't match the prefix %s", prefix)
		}
	}
	return nil
}

// buildNodesTable : builds the nodes table
func buildNodesTable(nodesTable *schema.Set) []pool.MemberNode {

	memberNodes := make([]pool.MemberNode, 0)
	for _, item := range nodesTable.List() {
		nodeItem := item.(map[string]interface{})
		memberNode := pool.MemberNode{}
		if node, ok := nodeItem["node"].(string); ok {
			memberNode.Node = node
		}
		if priority, ok := nodeItem["priority"].(int); ok {
			nodePriority := uint(priority)
			memberNode.Priority = &nodePriority
		}
		if state, ok := nodeItem["state"].(string); ok {
			memberNode.State = state
		}
		if weight, ok := nodeItem["weight"].(int); ok {
			memberNode.Weight = &weight
		}
		if sourceIP, ok := nodeItem["source_ip"].(string); ok {
			memberNode.SourceIP = sourceIP
		}
		memberNodes = append(memberNodes, memberNode)
	}
	return memberNodes
}

// buildAutoScaling : build the auto scaling object
func buildAutoScaling(autoScalingBlock interface{}) (pool.AutoScaling, error) {

	autoScalingObject := pool.AutoScaling{}
	autoScalingList := autoScalingBlock.([]interface{})
	autoScalingItem := autoScalingList[0].(map[string]interface{})

	if addNodeDelayTime, ok := autoScalingItem["addnode_delaytime"].(int); ok {
		autoScaleAddNodeDelayTime := uint(addNodeDelayTime)
		autoScalingObject.AddNodeDelayTime = &autoScaleAddNodeDelayTime
	}
	if cloudCredentials, ok := autoScalingItem["cloud_credentials"].(string); ok {
		autoScalingObject.CloudCredentials = cloudCredentials
	}
	if cluster, ok := autoScalingItem["cluster"].(string); ok {
		autoScalingObject.Cluster = cluster
	}
	if dataCentre, ok := autoScalingItem["data_center"].(string); ok {
		autoScalingObject.DataCenter = dataCentre
	}
	if enabled, ok := autoScalingItem["enabled"].(bool); ok {
		autoScalingObject.Enabled = &enabled
	}
	if external, ok := autoScalingItem["external"].(bool); ok {
		autoScalingObject.External = &external
	}
	if extraArgs, ok := autoScalingItem["extraargs"].(string); ok {
		autoScalingObject.ExtraArgs = extraArgs
	}
	if hysteresis, ok := autoScalingItem["hysteresis"].(int); ok {
		uintHysteresis := uint(hysteresis)
		autoScalingObject.Hysteresis = &uintHysteresis
	}
	if imageID, ok := autoScalingItem["imageid"].(string); ok {
		autoScalingObject.ImageID = imageID
	}
	if ipsToUse, ok := autoScalingItem["ips_to_use"].(string); ok {
		autoScalingObject.IPsToUse = ipsToUse
	}
	if lastNodeIdleTime, ok := autoScalingItem["last_node_idle_time"].(int); ok {
		uintLastNodeIdleTime := uint(lastNodeIdleTime)
		autoScalingObject.LastNodeIdleTime = &uintLastNodeIdleTime
	}
	if maxNodes, ok := autoScalingItem["max_nodes"].(int); ok {
		uintMaxNodes := uint(maxNodes)
		autoScalingObject.MaxNodes = &uintMaxNodes
	}
	if minNodes, ok := autoScalingItem["min_nodes"].(int); ok {
		uintMinNodes := uint(minNodes)
		autoScalingObject.MinNodes = &uintMinNodes
	}
	if name, ok := autoScalingItem["name"].(string); ok {
		autoScalingObject.Name = name
	}
	if port, ok := autoScalingItem["port"].(int); ok {
		uintPort := uint(port)
		autoScalingObject.Port = &uintPort
	}
	if refactory, ok := autoScalingItem["refactory"].(int); ok {
		uintRefactory := uint(refactory)
		autoScalingObject.Refractory = &uintRefactory
	}
	if responseTime, ok := autoScalingItem["response_time"].(int); ok {
		uintResponseTime := uint(responseTime)
		autoScalingObject.ResponseTime = &uintResponseTime
	}
	if scaleDownLevel, ok := autoScalingItem["scale_down_level"].(int); ok {
		uintScaleDownLevel := uint(scaleDownLevel)
		autoScalingObject.ScaleDownLevel = &uintScaleDownLevel
	}
	if scaleUpLevel, ok := autoScalingItem["scale_up_level"].(int); ok {
		uintScaleUpLevel := uint(scaleUpLevel)
		autoScalingObject.ScaleUpLevel = &uintScaleUpLevel
	}
	if securityGroupIDs, ok := autoScalingItem["securitygroupids"].(*schema.Set); ok {
		securityGroupIDList := util.BuildStringListFromSet(securityGroupIDs)
		err := checkStringPrefix("sg-", securityGroupIDList)
		if err != nil {
			return autoScalingObject, err
		}
		autoScalingObject.SecurityGroupIDs = securityGroupIDList
	}
	if sizeID, ok := autoScalingItem["size_id"].(string); ok {
		autoScalingObject.SizeID = sizeID
	}
	if subnetIDs, ok := autoScalingItem["subnetids"].(*schema.Set); ok {
		subnetIDList := util.BuildStringListFromSet(subnetIDs)
		err := checkStringPrefix("subnet-", subnetIDList)
		if err != nil {
			return autoScalingObject, err
		}
		autoScalingObject.SubnetIDs = subnetIDList
	}

	return autoScalingObject, nil
}

// buildConnection : build the connection object
func buildConnection(connectionBlock interface{}) pool.Connection {

	connectionObject := pool.Connection{}
	connectionList := connectionBlock.([]interface{})
	connectionItem := connectionList[0].(map[string]interface{})

	if maxConnectTime, ok := connectionItem["max_connect_time"].(int); ok {
		maxConnectTimeUint := uint(maxConnectTime)
		connectionObject.MaxConnectTime = &maxConnectTimeUint
	}
	if maxConnectionsPerNode, ok := connectionItem["max_connections_per_node"].(int); ok {
		maxConnectionsPerNodeUint := uint(maxConnectionsPerNode)
		connectionObject.MaxConnectionsPerNode = &maxConnectionsPerNodeUint
	}

	if maxQueueSize, ok := connectionItem["max_queue_size"].(int); ok {
		maxQueueSizeUint := uint(maxQueueSize)
		connectionObject.MaxQueueSize = &maxQueueSizeUint
	}
	if maxReplyTime, ok := connectionItem["max_reply_time"].(int); ok {
		maxReplyTimeUint := uint(maxReplyTime)
		connectionObject.MaxReplyTime = &maxReplyTimeUint
	}
	if queueTimeout, ok := connectionItem["queue_timeout"].(int); ok {
		queueTimeoutUint := uint(queueTimeout)
		connectionObject.QueueTimeout = &queueTimeoutUint
	}
	return connectionObject
}

// buildDNSAutoScale : build the DNS auto scale object
func buildDNSAutoScale(dnsAutoScaleBlock interface{}) pool.DNSAutoScale {

	dnsAutoScaleObject := pool.DNSAutoScale{}
	dnsAutoScaleList := dnsAutoScaleBlock.([]interface{})
	dnsAutoScaleItem := dnsAutoScaleList[0].(map[string]interface{})

	if enabled, ok := dnsAutoScaleItem["enabled"].(bool); ok {
		dnsAutoScaleObject.Enabled = &enabled
	}
	if hostnames, ok := dnsAutoScaleItem["hostnames"]; ok {
		dnsAutoScaleObject.Hostnames = util.BuildStringListFromSet(hostnames.(*schema.Set))
	}
	if port, ok := dnsAutoScaleItem["port"].(int); ok {
		portUint := uint(port)
		dnsAutoScaleObject.Port = &portUint
	}
	return dnsAutoScaleObject
}

// buildFTP : build the FTP object
func buildFTP(ftpBlock interface{}) pool.FTP {

	ftpObject := pool.FTP{}
	ftpList := ftpBlock.([]interface{})
	ftpItem := ftpList[0].(map[string]interface{})

	if supportRFC2428, ok := ftpItem["support_rfc_2428"].(bool); ok {
		ftpObject.SupportRFC2428 = &supportRFC2428
	}
	return ftpObject
}

// buildHTTP : build the HTTP object
func buildHTTP(httpBlock interface{}) pool.HTTP {

	httpObject := pool.HTTP{}
	httpList := httpBlock.([]interface{})
	httpItem := httpList[0].(map[string]interface{})

	if keepalive, ok := httpItem["keepalive"].(bool); ok {
		httpObject.HTTPKeepAlive = &keepalive
	}
	if keepaliveNonIdempotent, ok := httpItem["keepalive_non_idempotent"].(bool); ok {
		httpObject.HTTPKeepAlive = &keepaliveNonIdempotent
	}
	return httpObject
}

// buildKerberosProtocolTransition : build the kerberos protocol transitition object
func buildKerberosProtocolTransition(kerberosBlock interface{}) pool.KerberosProtocolTransition {

	kerberosObject := pool.KerberosProtocolTransition{}
	kerberosList := kerberosBlock.([]interface{})
	kerberosItem := kerberosList[0].(map[string]interface{})

	if principle, ok := kerberosItem["principal"].(string); ok {
		kerberosObject.Principal = principle
	}
	if target, ok := kerberosItem["target"].(string); ok {
		kerberosObject.Target = target
	}
	return kerberosObject
}

// buildLoadBalancing : build the load balancing object
func buildLoadBalancing(loadBalancingBlock interface{}) pool.LoadBalancing {

	loadBalancingObject := pool.LoadBalancing{}
	loadBalancingList := loadBalancingBlock.([]interface{})
	loadBalancingItem := loadBalancingList[0].(map[string]interface{})

	if algorithm, ok := loadBalancingItem["algorithm"].(string); ok {
		loadBalancingObject.Algorithm = algorithm
	}
	if priorityEnabled, ok := loadBalancingItem["priority_enabled"].(bool); ok {
		loadBalancingObject.PriorityEnabled = &priorityEnabled
	}
	if priorityNodes, ok := loadBalancingItem["priority_nodes"].(int); ok {
		priorityNodesUint := uint(priorityNodes)
		loadBalancingObject.PriorityNodes = &priorityNodesUint
	}
	return loadBalancingObject
}

// buildNode : build the Node object
func buildNode(nodeBlock interface{}) pool.Node {

	nodeObject := pool.Node{}
	nodeList := nodeBlock.([]interface{})
	nodeItem := nodeList[0].(map[string]interface{})

	if closeOnDeath, ok := nodeItem["close_on_death"].(bool); ok {
		nodeObject.CloseOnDeath = &closeOnDeath
	}
	if retryFailTime, ok := nodeItem["retry_fail_time"].(int); ok {
		retryFailTimeUint := uint(retryFailTime)
		nodeObject.RetryFailTime = &retryFailTimeUint
	}
	return nodeObject
}

// buildSMTP : build the SMTP object
func buildSMTP(smtpBlock interface{}) pool.SMTP {

	smtpObject := pool.SMTP{}
	smtpList := smtpBlock.([]interface{})
	smtpItem := smtpList[0].(map[string]interface{})

	if sendStartTLS, ok := smtpItem["send_starttls"].(bool); ok {
		smtpObject.SendSTARTTLS = &sendStartTLS
	}
	return smtpObject
}

// buildSSL : build the SSL object
func buildSSL(sslBlock interface{}) pool.Ssl {

	sslObject := pool.Ssl{}
	sslList := sslBlock.([]interface{})
	sslItem := sslList[0].(map[string]interface{})

	if clientAuth, ok := sslItem["client_auth"].(bool); ok {
		sslObject.ClientAuth = &clientAuth
	}
	if commonNameMatch, ok := sslItem["common_name_match"]; ok {
		sslObject.CommonNameMatch = util.BuildStringListFromSet(commonNameMatch.(*schema.Set))
	}
	if ellipticCurves, ok := sslItem["elliptic_curves"].(*schema.Set); ok {
		sslObject.EllipticCurves = util.BuildStringListFromSet(ellipticCurves)
	}
	if enable, ok := sslItem["enable"].(bool); ok {
		sslObject.Enable = &enable
	}
	if sendCloseAlerts, ok := sslItem["send_close_alerts"].(bool); ok {
		sslObject.SendCloseAlerts = &sendCloseAlerts
	}
	if serverName, ok := sslItem["server_name"].(bool); ok {
		sslObject.ServerName = &serverName
	}
	if signatureAlgorithms, ok := sslItem["signature_algorithms"].(string); ok {
		sslObject.SignatureAlgorithms = signatureAlgorithms
	}
	if sslCiphers, ok := sslItem["ssl_ciphers"].(string); ok {
		sslObject.SslCiphers = sslCiphers
	}
	if sslSupportSSL2, ok := sslItem["ssl_support_ssl2"].(string); ok {
		sslObject.SSLSupportSSL2 = sslSupportSSL2
	}
	if sslSupportSSL3, ok := sslItem["ssl_support_ssl3"].(string); ok {
		sslObject.SSLSupportSSL3 = sslSupportSSL3
	}
	if sslSupportTLS1, ok := sslItem["ssl_support_tls1"].(string); ok {
		sslObject.SSLSupportTLS1 = sslSupportTLS1
	}
	if sslSupportTLS1_1, ok := sslItem["ssl_support_tls1_1"].(string); ok {
		sslObject.SSLSupportTLS1_1 = sslSupportTLS1_1
	}
	if sslSupportTLS1_2, ok := sslItem["ssl_support_tls1_2"].(string); ok {
		sslObject.SSLSupportTLS1_2 = sslSupportTLS1_2
	}
	if strictVerify, ok := sslItem["strict_verify"].(bool); ok {
		sslObject.StrictVerify = &strictVerify
	}
	return sslObject
}

// buildTCP : build the TCP object
func buildTCP(tcpBlock interface{}) pool.TCP {

	tcpObject := pool.TCP{}
	tcpList := tcpBlock.([]interface{})
	tcpItem := tcpList[0].(map[string]interface{})

	if nagle, ok := tcpItem["nagle"].(bool); ok {
		tcpObject.Nagle = &nagle
	}
	return tcpObject
}

// buildUDP : build the UDP object
func buildUDP(udpBlock interface{}) pool.UDP {

	udpObject := pool.UDP{}
	udpList := udpBlock.([]interface{})
	udpItem := udpList[0].(map[string]interface{})

	if acceptFrom, ok := udpItem["accept_from"].(string); ok {
		udpObject.AcceptFrom = acceptFrom
	}
	if acceptFromMask, ok := udpItem["accept_from_mask"].(string); ok {
		udpObject.AcceptFromMask = acceptFromMask
	}
	if responseTimeout, ok := udpItem["response_timeout"].(int); ok {
		responseTimeoutUint := uint(responseTimeout)
		udpObject.ResponseTimeout = &responseTimeoutUint
	}
	return udpObject
}

// resourcePoolCreate - Creates a  pool resource object
func resourcePoolCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	var nodesTableDefined, nodesListDefined bool

	var createPool pool.Pool
	var poolName string
	if v, ok := d.GetOk("name"); ok && v != "" {
		poolName = v.(string)
	}
	if v, ok := d.GetOk("bandwidth_class"); ok {
		createPool.Properties.Basic.BandwidthClass = v.(string)
	}
	if v, ok := d.GetOk("failure_pool"); ok {
		createPool.Properties.Basic.FailurePool = v.(string)
	}
	if v, ok := d.GetOk("max_connection_attempts"); ok {
		maxConnectionAttempts := uint(v.(int))
		createPool.Properties.Basic.MaxConnectionAttempts = &maxConnectionAttempts
	}
	if v, ok := d.GetOk("max_idle_connections_pernode"); ok {
		maxIdleConnectionsPerNode := uint(v.(int))
		createPool.Properties.Basic.MaxIdleConnectionsPerNode = &maxIdleConnectionsPerNode
	}
	if v, ok := d.GetOk("max_timed_out_connection_attempts"); ok {
		maxTimedOutConnectionAttempts := uint(v.(int))
		createPool.Properties.Basic.MaxTimeoutConnectionAttempts = &maxTimedOutConnectionAttempts
	}
	if v, ok := d.GetOk("monitors"); ok {
		createPool.Properties.Basic.Monitors = util.BuildStringArrayFromInterface(v)
	}
	if v, ok := d.GetOk("node_close_with_rst"); ok {
		nodeCloseWithRst := v.(bool)
		createPool.Properties.Basic.NodeCloseWithReset = &nodeCloseWithRst
	}
	if v, ok := d.GetOk("node_connection_attempts"); ok {
		nodeConnectionAttempts := uint(v.(int))
		createPool.Properties.Basic.NodeConnectionAttempts = &nodeConnectionAttempts
	}
	if v, ok := d.GetOk("node_delete_behaviour"); ok {
		createPool.Properties.Basic.NodeDeleteBehavior = v.(string)
	}
	if v, ok := d.GetOk("node_drain_to_delete_timeout"); ok {
		nodeDrainDeleteTimeout := uint(v.(int))
		createPool.Properties.Basic.NodeDrainDeleteTimeout = &nodeDrainDeleteTimeout
	}
	if v, ok := d.GetOk("nodes_table"); ok {
		createPool.Properties.Basic.NodesTable = buildNodesTable(v.(*schema.Set))
		nodesTableDefined = true
	} else {
		if v, ok := d.GetOk("nodes_list"); ok {
			addresses := v.(*schema.Set).List()
			nodesTable := make([]pool.MemberNode, 0)
			for _, ip_addr := range addresses {
				var rec pool.MemberNode
				rec.Node = ip_addr.(string)
				nodesTable = append(nodesTable, rec)
			}
			createPool.Properties.Basic.NodesTable = nodesTable
			nodesListDefined = true
		}
	}
	if nodesTableDefined == false && nodesListDefined == false {
		return fmt.Errorf("Error creating resource: no one of nodes_table or nodes_list attr has been defined")
	}
	if v, ok := d.GetOk("note"); ok {
		createPool.Properties.Basic.Note = v.(string)
	}
	if v, ok := d.GetOk("passive_monitoring"); ok {
		passiveMonitoring := v.(bool)
		createPool.Properties.Basic.PassiveMonitoring = &passiveMonitoring
	}
	if v, ok := d.GetOk("persistence_class"); ok {
		createPool.Properties.Basic.PersistenceClass = v.(string)
	}
	if v, ok := d.GetOk("transparent"); ok {
		transparent := v.(bool)
		createPool.Properties.Basic.Transparent = &transparent
	}
	if v, ok := d.GetOk("auto_scaling"); ok {
		autoScaling, err := buildAutoScaling(v)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Pool - auto_scaling error whilst creating %s: %v", poolName, err)
		}
		createPool.Properties.AutoScaling = autoScaling
	}
	if v, ok := d.GetOk("pool_connection"); ok {
		createPool.Properties.Connection = buildConnection(v)
	}
	if v, ok := d.GetOk("dns_autoscale"); ok {
		createPool.Properties.DNSAutoScale = buildDNSAutoScale(v)
	}
	if v, ok := d.GetOk("ftp"); ok {
		createPool.Properties.FTP = buildFTP(v)
	}
	if v, ok := d.GetOk("http"); ok {
		createPool.Properties.HTTP = buildHTTP(v)
	}
	if v, ok := d.GetOk("kerberos_protocol_transition"); ok {
		createPool.Properties.KerberosProtocolTransition = buildKerberosProtocolTransition(v)
	}
	if v, ok := d.GetOk("load_balancing"); ok {
		createPool.Properties.LoadBalancing = buildLoadBalancing(v)
	}
	if v, ok := d.GetOk("node"); ok {
		createPool.Properties.Node = buildNode(v)
	}
	if v, ok := d.GetOk("smtp"); ok {
		createPool.Properties.SMTP = buildSMTP(v)
	}
	if v, ok := d.GetOk("ssl"); ok {
		createPool.Properties.Ssl = buildSSL(v)
	}
	if v, ok := d.GetOk("tcp"); ok {
		createPool.Properties.TCP = buildTCP(v)
	}
	if v, ok := d.GetOk("udp"); ok {
		createPool.Properties.UDP = buildUDP(v)
	}

	err := client.Set("pools", poolName, createPool, nil)
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
	var poolName string

	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)
	}

	var poolObj pool.Pool
	client.WorkWithConfigurationResources()
	err := client.GetByName("pools", poolName, &poolObj)
	if err != nil {
		if client.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM Pool error whilst retrieving %s: %s", poolName, err)
	}

	d.Set("name", poolName)
	d.Set("bandwidth_class", poolObj.Properties.Basic.BandwidthClass)
	d.Set("failure_pool", poolObj.Properties.Basic.FailurePool)
	d.Set("max_connection_attempts", *poolObj.Properties.Basic.MaxConnectionAttempts)
	d.Set("max_idle_connections_pernode", *poolObj.Properties.Basic.MaxIdleConnectionsPerNode)
	d.Set("max_timed_out_connection_attempts", *poolObj.Properties.Basic.MaxTimeoutConnectionAttempts)
	d.Set("monitors", poolObj.Properties.Basic.Monitors)
	d.Set("node_close_with_rst", *poolObj.Properties.Basic.NodeCloseWithReset)
	d.Set("node_connection_attempts", *poolObj.Properties.Basic.NodeConnectionAttempts)
	d.Set("node_delete_behaviour", poolObj.Properties.Basic.NodeDeleteBehavior)
	d.Set("node_drain_to_delete_timeout", *poolObj.Properties.Basic.NodeDrainDeleteTimeout)
	d.Set("nodes_table", poolObj.Properties.Basic.NodesTable)
	d.Set("note", poolObj.Properties.Basic.Note)
	d.Set("passive_monitoring", *poolObj.Properties.Basic.PassiveMonitoring)
	d.Set("persistence_class", poolObj.Properties.Basic.PersistenceClass)
	d.Set("transparent", *poolObj.Properties.Basic.Transparent)
	d.Set("auto_scaling", []pool.AutoScaling{poolObj.Properties.AutoScaling})
	d.Set("pool_connection", []pool.Connection{poolObj.Properties.Connection})
	d.Set("dns_autoscale", []pool.DNSAutoScale{poolObj.Properties.DNSAutoScale})
	d.Set("ftp", []pool.FTP{poolObj.Properties.FTP})
	d.Set("http", []pool.HTTP{poolObj.Properties.HTTP})
	d.Set("kerberos_protocol_transition", []pool.KerberosProtocolTransition{poolObj.Properties.KerberosProtocolTransition})
	d.Set("load_balancing", []pool.LoadBalancing{poolObj.Properties.LoadBalancing})
	d.Set("node", []pool.Node{poolObj.Properties.Node})
	d.Set("smtp", []pool.SMTP{poolObj.Properties.SMTP})
	d.Set("ssl", []pool.Ssl{poolObj.Properties.Ssl})
	d.Set("tcp", []pool.TCP{poolObj.Properties.TCP})
	d.Set("udp", []pool.UDP{poolObj.Properties.UDP})
	return nil
}

// resourcePoolUpdate - Updates an existing pool resource
func resourcePoolUpdate(d *schema.ResourceData, m interface{}) error {

	var updatePool pool.Pool
	hasChanges := false
	poolName := d.Id()

	if d.HasChange("bandwidth_class") {
		if v, ok := d.GetOk("bandwidth_class"); ok {
			updatePool.Properties.Basic.BandwidthClass = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("failure_pool") {
		if v, ok := d.GetOk("failure_pool"); ok {
			updatePool.Properties.Basic.FailurePool = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("max_connection_attempts") {
		maxConnectionAttempts := uint(d.Get("max_connection_attempts").(int))
		updatePool.Properties.Basic.MaxConnectionAttempts = &maxConnectionAttempts
		hasChanges = true
	}
	if d.HasChange("max_idle_connections_pernode") {
		maxIdleConnectionsPerNode := uint(d.Get("max_idle_connections_pernode").(int))
		updatePool.Properties.Basic.MaxIdleConnectionsPerNode = &maxIdleConnectionsPerNode
		hasChanges = true
	}
	if d.HasChange("max_timed_out_connection_attempts") {
		maxTimedOutConnectionAttempts := uint(d.Get("max_timed_out_connection_attempts").(int))
		updatePool.Properties.Basic.MaxTimeoutConnectionAttempts = &maxTimedOutConnectionAttempts
		hasChanges = true
	}
	if d.HasChange("monitors") {
		updatePool.Properties.Basic.Monitors = util.BuildStringArrayFromInterface(d.Get("monitors"))
		hasChanges = true
	}
	if d.HasChange("node_close_with_rst") {
		nodeCloseWithRst := d.Get("node_close_with_rst").(bool)
		updatePool.Properties.Basic.NodeCloseWithReset = &nodeCloseWithRst
	}
	if d.HasChange("node_connection_attempts") {
		nodeConnectionAttempts := uint(d.Get("node_connection_attempts").(int))
		updatePool.Properties.Basic.NodeConnectionAttempts = &nodeConnectionAttempts
		hasChanges = true
	}
	if d.HasChange("node_delete_behaviour") {
		if v, ok := d.GetOk("node_delete_behaviour"); ok {
			updatePool.Properties.Basic.NodeDeleteBehavior = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("node_drain_to_delete_timeout") {
		nodeDrainTimeout := uint(d.Get("node_drain_to_delete_timeout").(int))
		updatePool.Properties.Basic.NodeDrainDeleteTimeout = &nodeDrainTimeout
		hasChanges = true
	}
	if d.HasChange("nodes_table") {
		updatePool.Properties.Basic.NodesTable = buildNodesTable(d.Get("nodes_table").(*schema.Set))
		hasChanges = true
	}
	if d.HasChange("nodes_list") {
		if v, ok := d.GetOk("nodes_list"); ok {
			addresses := v.(*schema.Set).List()
			nodesTable := make([]pool.MemberNode, 0)
			for _, ip_addr := range addresses {
				var rec pool.MemberNode
				rec.Node = ip_addr.(string)
				nodesTable = append(nodesTable, rec)
			}
			updatePool.Properties.Basic.NodesTable = nodesTable
		}
	}
	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok {
			updatePool.Properties.Basic.Note = v.(string)
		}
		hasChanges = true
	}

	if d.HasChange("passive_monitoring") {
		passiveMonitoring := d.Get("passive_monitoring").(bool)
		updatePool.Properties.Basic.PassiveMonitoring = &passiveMonitoring
		hasChanges = true
	}
	if d.HasChange("persistence_class") {
		if v, ok := d.GetOk("persistence_class"); ok {
			updatePool.Properties.Basic.PersistenceClass = v.(string)
		}
		hasChanges = true
	}

	if d.HasChange("transparent") {
		transparent := d.Get("transparent").(bool)
		updatePool.Properties.Basic.Transparent = &transparent
		hasChanges = true
	}
	if d.HasChange("auto_scaling") {
		if v, ok := d.GetOk("auto_scaling"); ok {
			autoScaling, err := buildAutoScaling(v)
			if err != nil {
				return fmt.Errorf("BrocadeVTM Pool - auto_scaling error whilst updating %s: %v", poolName, err)
			}
			updatePool.Properties.AutoScaling = autoScaling
		}
		hasChanges = true
	}
	if d.HasChange("pool_connection") {
		if v, ok := d.GetOk("pool_connection"); ok {
			updatePool.Properties.Connection = buildConnection(v)
		}
		hasChanges = true
	}
	if d.HasChange("dns_autoscale") {
		if v, ok := d.GetOk("dns_autoscale"); ok {
			updatePool.Properties.DNSAutoScale = buildDNSAutoScale(v)
		}
		hasChanges = true
	}
	if d.HasChange("ftp") {
		if v, ok := d.GetOk("ftp"); ok {
			updatePool.Properties.FTP = buildFTP(v)
		}
		hasChanges = true
	}
	if d.HasChange("http") {
		if v, ok := d.GetOk("http"); ok {
			updatePool.Properties.HTTP = buildHTTP(v)
		}
		hasChanges = true
	}
	if d.HasChange("kerberos_protocol_transition") {
		if v, ok := d.GetOk("kerberos_protocol_transition"); ok {
			updatePool.Properties.KerberosProtocolTransition = buildKerberosProtocolTransition(v)
		}
		hasChanges = true
	}
	if d.HasChange("load_balancing") {
		if v, ok := d.GetOk("load_balancing"); ok {
			updatePool.Properties.LoadBalancing = buildLoadBalancing(v)
		}
		hasChanges = true
	}
	if d.HasChange("node") {
		if v, ok := d.GetOk("node"); ok {
			updatePool.Properties.Node = buildNode(v)
		}
		hasChanges = true
	}
	if d.HasChange("smtp") {
		if v, ok := d.GetOk("smtp"); ok {
			updatePool.Properties.SMTP = buildSMTP(v)
		}
		hasChanges = true
	}
	if d.HasChange("ssl") {
		if v, ok := d.GetOk("ssl"); ok {
			updatePool.Properties.Ssl = buildSSL(v)
		}
		hasChanges = true
	}
	if d.HasChange("tcp") {
		if v, ok := d.GetOk("tcp"); ok {
			updatePool.Properties.TCP = buildTCP(v)
		}
		hasChanges = true
	}
	if d.HasChange("udp") {
		if v, ok := d.GetOk("udp"); ok {
			updatePool.Properties.UDP = buildUDP(v)
		}
		hasChanges = true
	}

	if hasChanges {
		config := m.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)
		err := client.Set("pools", poolName, updatePool, nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Pool error whilst updating %s: %s", poolName, err)
		}
		d.SetId(poolName)
	}

	return resourcePoolRead(d, m)
}

// resourcePoolDelete - Deletes a pool resource
func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("pools", d, m)
}
