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
				ValidateFunc: validatePoolUnsignedInteger,
				Description:  "Maximum numberof nodes an attempt to send a request to befoirce returning an error to the client",
			},
			"max_idle_connections_pernode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePoolUnsignedInteger,
				Default:      50,
				Description:  "Maximum number of unused HTTP keepalive connections",
			},
			"max_timed_out_connection_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePoolUnsignedInteger,
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
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "Number of times an attempt to connect to the same node before marking it as failed. Only used when passive_monitoring is enabled",
			},
			"node_delete_behaviour": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "immediate",
				Description:  "Node deletion behaviour for this pool",
				ValidateFunc: validateNodeDeleteBehaviour,
			},
			"node_drain_to_delete_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum time a node will remain in draining after it has been deleted",
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
							ValidateFunc: validatePoolUnsignedInteger,
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

func validatePoolUnsignedInteger(v interface{}, k string) (ws []string, errors []error) {
	checkNumber := v.(int)
	if checkNumber < 0 {
		errors = append(errors, fmt.Errorf("%q can't be negative", k))
	}
	return
}

func validateNode(v interface{}, k string) (ws []string, errors []error) {
	node := v.(string)
	validateNode := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+$`)
	if !validateNode.MatchString(node) {
		errors = append(errors, fmt.Errorf("Must be a valid IP and port seperated by a colon. i.e 127.0.0.1:80"))
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
		if node, ok := nodeItem["node"].(string); ok && node != "" {
			memberNode.Node = node
		}
		if priority, ok := nodeItem["priority"].(int); ok {
			nodePriority := uint(priority)
			memberNode.Priority = &nodePriority
		}
		if state, ok := nodeItem["state"].(string); ok && state != "" {
			memberNode.State = state
		}
		if weight, ok := nodeItem["weight"].(int); ok {
			memberNode.Weight = &weight
		}
		if sourceIP, ok := nodeItem["source_ip"].(string); ok && sourceIP != "" {
			memberNode.SourceIP = sourceIP
		}
		memberNodes = append(memberNodes, memberNode)
	}
	return memberNodes
}

// resourcePoolCreate - Creates a  pool resource object
func resourcePoolCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)

	var createPool pool.Pool
	var poolName string
	if v, ok := d.GetOk("name"); ok && v != "" {
		poolName = v.(string)
	}
	if v, ok := d.GetOk("bandwidth_class"); ok && v != "" {
		createPool.Properties.Basic.BandwidthClass = v.(string)
	}
	if v, ok := d.GetOk("failure_pool"); ok && v != "" {
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
	if v, ok := d.GetOk("node_delete_behaviour"); ok && v != "" {
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
	if v, ok := d.GetOk("note"); ok && v != "" {
		createPool.Properties.Basic.Note = v.(string)
	}
	createPool.Properties.Basic.PassiveMonitoring = d.Get("passive_monitoring").(bool)
	if v, ok := d.GetOk("persistence_class"); ok && v != "" {
		createPool.Properties.Basic.PersistenceClass = v.(string)
	}
	createPool.Properties.Basic.Transparent = d.Get("transparent").(bool)

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
		if v, ok := d.GetOk("bandwidth_class"); ok && v != "" {
			updatePool.Properties.Basic.BandwidthClass = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("failure_pool") {
		if v, ok := d.GetOk("failure_pool"); ok && v != "" {
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
		if v, ok := d.GetOk("node_delete_behaviour"); ok && v != "" {
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
		if v, ok := d.GetOk("note"); ok && v != "" {
			updatePool.Properties.Basic.Note = v.(string)
		}
		hasChanges = true
	}
	updatePool.Properties.Basic.PassiveMonitoring = d.Get("passive_monitoring").(bool)
	if d.HasChange("passive_monitoring") {
		hasChanges = true
	}
	if d.HasChange("persistence_class") {
		if v, ok := d.GetOk("persistence_class"); ok && v != "" {
			updatePool.Properties.Basic.PersistenceClass = v.(string)
		}
		hasChanges = true
	}
	updatePool.Properties.Basic.Transparent = d.Get("transparent").(bool)
	if d.HasChange("transparent") {
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
