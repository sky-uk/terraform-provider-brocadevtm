package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/pool"
	"github.com/sky-uk/go-rest-api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"log"
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
				Type:        schema.TypeList,
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
				Type:     schema.TypeSet,
				Required: true,
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
						/* extraargs is in doco, but causes errors as it doesn't appear to be in API
						"extraargs": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Extra comma separated arguments to send the auto-scaling API",
						},
						*/
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
					Schema: map[string]*schema.Schema{},
				},
			},
			"kerberos_protocol_transition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"load_balancing": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"node": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"smtp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"ssl": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"tcp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"udp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},

			/*
				"max_connection_timeout": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validatePoolUnsignedInteger,
					Default:      4,
				},
				"max_connections_per_node": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validatePoolUnsignedInteger,
				},
				"max_queue_size": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validatePoolUnsignedInteger,
					Default:      0,
				},
				"max_reply_time": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validatePoolUnsignedInteger,
					Default:      30,
				},
				"queue_timeout": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validatePoolUnsignedInteger,
					Default:      10,
				},
				"http_keepalive": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"http_keepalive_non_idempotent": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"load_balancing_priority_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"load_balancing_priority_nodes": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validatePoolUnsignedInteger,
					Default:      1,
				},
				"load_balancing_algorithm": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validatePoolLBAlgo,
				},
				"tcp_nagle": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			*/
		},
	}

}

func validateAutoScalingIPsToUse(v interface{}, k string) (ws []string, errors []error) {
	ipType := v.(string)
	ipTypeOptions := regexp.MustCompile(`^(publicips|private_ips)`)
	if !ipTypeOptions.MatchString(ipType) {
		errors = append(errors, fmt.Errorf("%q must be one of publicips or private_ips", k))
	}
	return
}

func validateNodeDeleteBehaviour(v interface{}, k string) (ws []string, errors []error) {
	behaviour := v.(string)
	behvaiourOptions := regexp.MustCompile(`^(immediate|drain)$`)
	if !behvaiourOptions.MatchString(behaviour) {
		errors = append(errors, fmt.Errorf("%q must be one of immediate or drain", k))
	}
	return
}

func validatePoolLBAlgo(v interface{}, k string) (ws []string, errors []error) {
	algo := v.(string)
	algoOptions := regexp.MustCompile(`^(fastest_response_time|least_connections|perceptive|random|round_robin|weighted_least_connections|weighted_round_robin)$`)
	if !algoOptions.MatchString(algo) {
		errors = append(errors, fmt.Errorf("%q must be one of fastest_response_time, least_connections, perceptive, random, round_robin, weighted_least_connections, weighted_round_robin", k))
	}
	return
}

func validateNode(v interface{}, k string) (ws []string, errors []error) {
	node := v.(string)
	validateNode := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+$`)
	if !validateNode.MatchString(node) {
		errors = append(errors, fmt.Errorf("%q must be a valid IP and port seperated by a colon. i.e 127.0.0.1:80", k))
	}
	return
}

func validateWeight(v interface{}, k string) (ws []string, errors []error) {
	weight := v.(int)

	if weight < 1 || weight > 100 {
		errors = append(errors, fmt.Errorf("%q must be between 1-100", k))
	}
	return
}

func validateState(v interface{}, k string) (ws []string, errors []error) {
	state := v.(string)
	stateOptions := regexp.MustCompile(`^(active|draining|disabled)$`)
	if !stateOptions.MatchString(state) {
		errors = append(errors, fmt.Errorf("%q must be one of active, draining, disabled", k))
	}
	return
}

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

func checkStringPrefix(prefix string, list []string) error {
	for _, item := range list {
		checkFormat := regexp.MustCompile(`^` + prefix)
		if !checkFormat.MatchString(item) {
			return fmt.Errorf("one or more items in the list of strings doesn't match the prefix %s", prefix)
		}
	}
	return nil
}

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

func buildFTP(ftpBlock interface{}) pool.FTP {

	ftpObject := pool.FTP{}
	ftpList := ftpBlock.([]interface{})
	ftpItem := ftpList[0].(map[string]interface{})

	if supportRFC2428, ok := ftpItem["support_rfc_2428"].(bool); ok {
		ftpObject.SupportRFC2428 = &supportRFC2428
	}
	return ftpObject
}

// resourcePoolCreate - Creates a  pool resource object
func resourcePoolCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)

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
		createPool.Properties.Basic.MaxIdleConnectionsPerNode = uint(v.(int))
	}
	if v, ok := d.GetOk("max_timed_out_connection_attempts"); ok {
		createPool.Properties.Basic.MaxTimeoutConnectionAttempts = uint(v.(int))
	}
	if v, ok := d.GetOk("monitors"); ok {
		createPool.Properties.Basic.Monitors = util.BuildStringArrayFromInterface(v)
	}
	createPool.Properties.Basic.NodeCloseWithReset = d.Get("node_close_with_rst").(bool)
	if v, ok := d.GetOk("node_connection_attempts"); ok {
		createPool.Properties.Basic.NodeConnectionAttempts = uint(v.(int))
	}
	if v, ok := d.GetOk("node_delete_behaviour"); ok {
		createPool.Properties.Basic.NodeDeleteBehavior = v.(string)
	}
	if v, ok := d.GetOk("node_drain_to_delete_timeout"); ok {
		log.Println("The node drain to delete timeout is: ", v.(int))
		nodeDrainDeleteTimeout := uint(v.(int))
		createPool.Properties.Basic.NodeDrainDeleteTimeout = &nodeDrainDeleteTimeout
	}
	if v, ok := d.GetOk("nodes_table"); ok {
		createPool.Properties.Basic.NodesTable = buildNodesTable(v.(*schema.Set))
	}
	if v, ok := d.GetOk("note"); ok {
		createPool.Properties.Basic.Note = v.(string)
	}
	createPool.Properties.Basic.PassiveMonitoring = d.Get("passive_monitoring").(bool)
	if v, ok := d.GetOk("persistence_class"); ok {
		createPool.Properties.Basic.PersistenceClass = v.(string)
	}
	createPool.Properties.Basic.Transparent = d.Get("transparent").(bool)
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

	/*

		if v, ok := d.GetOk("max_connection_timeout"); ok {
			createPool.Properties.Connection.MaxConnectTime = uint(v.(int))
		}
		if v, ok := d.GetOk("max_connections_per_node"); ok {
			createPool.Properties.Connection.MaxConnectionsPerNode = uint(v.(int))
		}
		if v, ok := d.GetOk("max_queue_size"); ok {
			createPool.Properties.Connection.MaxQueueSize = uint(v.(int))
		}
		if v, ok := d.GetOk("max_reply_time"); ok {
			createPool.Properties.Connection.MaxReplyTime = uint(v.(int))
		}
		if v, ok := d.GetOk("queue_timeout"); ok {
			createPool.Properties.Connection.QueueTimeout = uint(v.(int))
		}
		if v, _ := d.GetOk("http_keepalive"); v != nil {
			httpKeepAlive := v.(bool)
			createPool.Properties.HTTP.HTTPKeepAlive = &httpKeepAlive
		}
		if v, _ := d.GetOk("http_keepalive_non_idempotent"); v != nil {
			httpKeepAliveNonIdempotent := v.(bool)
			createPool.Properties.HTTP.HTTPKeepAliveNonIdempotent = &httpKeepAliveNonIdempotent
		}
		if v, _ := d.GetOk("load_balancing_priority_enabled"); v != nil {
			loadBalancingPriorityEnabled := v.(bool)
			createPool.Properties.LoadBalancing.PriorityEnabled = &loadBalancingPriorityEnabled
		}
		if v, ok := d.GetOk("load_balancing_priority_nodes"); ok {
			createPool.Properties.LoadBalancing.PriorityNodes = uint(v.(int))
		}
		if v, ok := d.GetOk("load_balancing_algorithm"); ok && v != "" {
			createPool.Properties.LoadBalancing.Algorithm = v.(string)
		}
		if v, _ := d.GetOk("tcp_nagle"); v != nil {
			tcpNagle := v.(bool)
			createPool.Properties.TCP.Nagle = &tcpNagle
		}
	*/
	createAPI := pool.NewCreate(poolName, createPool)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Pool error whilst creating %s: %v", poolName, createAPI.ErrorObject())
	}

	d.SetId(poolName)
	return resourcePoolRead(d, m)
}

// resourcePoolRead - Reads a  pool resource
func resourcePoolRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var poolName string

	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)
	}

	getAPI := pool.NewGet(poolName)
	err := vtmClient.Do(getAPI)
	if err != nil {
		if getAPI.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM Pool error whilst retrieving %s: %v", poolName, getAPI.ErrorObject())
	}
	response := getAPI.ResponseObject().(*pool.Pool)

	d.Set("name", poolName)
	d.Set("bandwidth_class", response.Properties.Basic.BandwidthClass)
	d.Set("failure_pool", response.Properties.Basic.FailurePool)
	d.Set("max_connection_attempts", *response.Properties.Basic.MaxConnectionAttempts)
	d.Set("max_idle_connections_pernode", response.Properties.Basic.MaxIdleConnectionsPerNode)
	d.Set("max_timed_out_connection_attempts", response.Properties.Basic.MaxTimeoutConnectionAttempts)
	d.Set("monitors", response.Properties.Basic.Monitors)
	d.Set("node_close_with_rst", response.Properties.Basic.NodeCloseWithReset)
	d.Set("node_connection_attempts", response.Properties.Basic.NodeConnectionAttempts)
	d.Set("node_delete_behaviour", response.Properties.Basic.NodeDeleteBehavior)
	d.Set("node_drain_to_delete_timeout", *response.Properties.Basic.NodeDrainDeleteTimeout)
	d.Set("nodes_table", response.Properties.Basic.NodesTable)
	d.Set("note", response.Properties.Basic.Note)
	d.Set("passive_monitoring", response.Properties.Basic.PassiveMonitoring)
	d.Set("persistence_class", response.Properties.Basic.PersistenceClass)
	d.Set("transparent", response.Properties.Basic.Transparent)

	d.Set("auto_scaling", []pool.AutoScaling{response.Properties.AutoScaling})
	d.Set("pool_connection", []pool.Connection{response.Properties.Connection})
	d.Set("dns_autoscale", []pool.DNSAutoScale{response.Properties.DNSAutoScale})
	d.Set("ftp", []pool.FTP{response.Properties.FTP})

	/*
		d.Set("node", response.Properties.Basic.NodesTable)
		d.Set("monitorlist", response.Properties.Basic.Monitors)
		d.Set("max_connection_attempts", *response.Properties.Basic.MaxConnectionAttempts)
		d.Set("max_idle_connections_pernode", response.Properties.Basic.MaxIdleConnectionsPerNode)
		d.Set("max_timed_out_connection_attempts", response.Properties.Basic.MaxTimeoutConnectionAttempts)
		d.Set("node_close_with_rst", *response.Properties.Basic.NodeCloseWithReset)
		d.Set("max_connection_timeout", response.Properties.Connection.MaxConnectTime)
		d.Set("max_connections_per_node", response.Properties.Connection.MaxConnectionsPerNode)
		d.Set("max_queue_size", response.Properties.Connection.MaxQueueSize)
		d.Set("max_reply_time", response.Properties.Connection.MaxReplyTime)
		d.Set("queue_timeout", response.Properties.Connection.QueueTimeout)
		d.Set("http_keepalive", *response.Properties.HTTP.HTTPKeepAlive)
		d.Set("http_keepalive_non_idempotent", *response.Properties.HTTP.HTTPKeepAliveNonIdempotent)
		d.Set("load_balancing_priority_enabled", *response.Properties.LoadBalancing.PriorityEnabled)
		d.Set("load_balancing_priority_nodes", response.Properties.LoadBalancing.PriorityNodes)
		d.Set("load_balancing_algorithm", response.Properties.LoadBalancing.Algorithm)
		d.Set("tcp_nagle", *response.Properties.TCP.Nagle)
	*/

	return nil
}

// resourcePoolUpdate - Updates an existing pool resource
func resourcePoolUpdate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
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
		updatePool.Properties.Basic.MaxIdleConnectionsPerNode = uint(d.Get("max_idle_connections_pernode").(int))
		hasChanges = true
	}
	if d.HasChange("max_timed_out_connection_attempts") {
		updatePool.Properties.Basic.MaxTimeoutConnectionAttempts = uint(d.Get("max_timed_out_connection_attempts").(int))
		hasChanges = true
	}
	if d.HasChange("monitors") {
		updatePool.Properties.Basic.Monitors = util.BuildStringArrayFromInterface(d.Get("monitors"))
		hasChanges = true
	}
	updatePool.Properties.Basic.NodeCloseWithReset = d.Get("node_close_with_rst").(bool)
	if d.HasChange("node_close_with_rst") {
		hasChanges = true
	}
	if d.HasChange("node_connection_attempts") {
		updatePool.Properties.Basic.NodeConnectionAttempts = uint(d.Get("node_connection_attempts").(int))
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
	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok {
			updatePool.Properties.Basic.Note = v.(string)
		}
		hasChanges = true
	}
	updatePool.Properties.Basic.PassiveMonitoring = d.Get("passive_monitoring").(bool)
	if d.HasChange("passive_monitoring") {
		hasChanges = true
	}
	if d.HasChange("persistence_class") {
		if v, ok := d.GetOk("persistence_class"); ok {
			updatePool.Properties.Basic.PersistenceClass = v.(string)
		}
		hasChanges = true
	}
	updatePool.Properties.Basic.Transparent = d.Get("transparent").(bool)
	if d.HasChange("transparent") {
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

	/*


		if d.HasChange("max_connection_timeout") {
			if v, ok := d.GetOk("max_connection_timeout"); ok {
				updatePool.Properties.Connection.MaxConnectTime = uint(v.(int))
			}
			hasChanges = true
		}

		if d.HasChange("max_connections_per_node") {
			if v, ok := d.GetOk("max_connections_per_node"); ok {
				updatePool.Properties.Connection.MaxConnectionsPerNode = uint(v.(int))
			}
			hasChanges = true
		}

		if d.HasChange("max_queue_size") {
			if v, ok := d.GetOk("max_queue_size"); ok {
				updatePool.Properties.Connection.MaxQueueSize = uint(v.(int))
			}
			hasChanges = true
		}

		if d.HasChange("max_reply_time") {
			if v, ok := d.GetOk("max_reply_time"); ok {
				updatePool.Properties.Connection.MaxReplyTime = uint(v.(int))
			}
			hasChanges = true
		}

		if d.HasChange("queue_timeout") {
			if v, ok := d.GetOk("queue_timeout"); ok {
				updatePool.Properties.Connection.QueueTimeout = uint(v.(int))
			}
			hasChanges = true
		}

		if d.HasChange("http_keepalive") {
			httpKeepAlive := d.Get("http_keepalive").(bool)
			updatePool.Properties.HTTP.HTTPKeepAlive = &httpKeepAlive
			hasChanges = true
		}

		if d.HasChange("http_keepalive_non_idempotent") {
			httpKeepAliveNonIdempotent := d.Get("http_keepalive_non_idempotent").(bool)
			updatePool.Properties.HTTP.HTTPKeepAliveNonIdempotent = &httpKeepAliveNonIdempotent
			hasChanges = true
		}

		if d.HasChange("load_balancing_priority_enabled") {
			loadBalancingPriorityEnabled := d.Get("load_balancing_priority_enabled").(bool)
			updatePool.Properties.LoadBalancing.PriorityEnabled = &loadBalancingPriorityEnabled
			hasChanges = true
		}

		if d.HasChange("load_balancing_priority_nodes") {
			if v, ok := d.GetOk("load_balancing_priority_nodes"); ok {
				updatePool.Properties.LoadBalancing.PriorityNodes = uint(v.(int))
			}
			hasChanges = true
		}

		if d.HasChange("load_balancing_algorithm") {
			if v, ok := d.GetOk("load_balancing_algorithm"); ok && v != "" {
				updatePool.Properties.LoadBalancing.Algorithm = v.(string)
			}
			hasChanges = true
		}

		if d.HasChange("tcp_nagle") {
			tcpNagle := d.Get("tcp_nagle").(bool)
			updatePool.Properties.TCP.Nagle = &tcpNagle
			hasChanges = true
		}
	*/

	if hasChanges {
		updatePoolAPI := pool.NewUpdate(poolName, updatePool)
		err := vtmClient.Do(updatePoolAPI)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Pool error whilst updating %s: %v", poolName, updatePoolAPI.ErrorObject())
		}
		d.SetId(poolName)
	}

	return resourcePoolRead(d, m)

}

// resourcePoolDelete - Deletes a pool resource
func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	poolName := d.Id()

	deleteAPI := pool.NewDelete(poolName)
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf(fmt.Sprintf("BrocadeVTM Pool error whilst deleting %s: %v", poolName, deleteAPI.ErrorObject()))
	}
	d.SetId("")
	return nil
}
