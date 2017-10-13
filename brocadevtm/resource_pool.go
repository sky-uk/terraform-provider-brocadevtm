package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/pool"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNode,
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validatePoolUnsignedInteger,
						},
						"state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateState,
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateWeight,
						},
					},
				},
			},
			"monitorlist": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_connection_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePoolUnsignedInteger,
				Default:      0,
			},
			"max_idle_connections_pernode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePoolUnsignedInteger,
				Default:      50,
			},
			"max_timed_out_connection_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePoolUnsignedInteger,
				Default:      2,
			},
			"node_close_with_rst": {
				Type:     schema.TypeBool,
				Optional: true,
			},
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
		},
	}

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

// resourcePoolCreate - Creates a  pool resource object
func resourcePoolCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	var createPool pool.Pool
	var poolName string
	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)
	}

	if v, ok := d.GetOk("node"); ok {
		if nodes, ok := v.(*schema.Set); ok {
			nodeList := []pool.MemberNode{}
			for _, value := range nodes.List() {
				nodeObject := value.(map[string]interface{})
				newNode := pool.MemberNode{}
				if nodeValue, ok := nodeObject["node"].(string); ok {
					newNode.Node = nodeValue
				}
				if priorityValue, ok := nodeObject["priority"].(int); ok {
					newNode.Priority = priorityValue
				}
				if stateValue, ok := nodeObject["state"].(string); ok {
					newNode.State = stateValue
				}
				if weightValue, ok := nodeObject["weight"].(int); ok {
					newNode.Weight = weightValue
				}
				nodeList = append(nodeList, newNode)

			}
			createPool.Properties.Basic.NodesTable = nodeList
		}
	}
	if v, ok := d.GetOk("monitorlist"); ok {
		originalMonitors := v.([]interface{})
		monitors := make([]string, len(originalMonitors))
		for i, monitor := range originalMonitors {
			monitors[i] = monitor.(string)
		}
		createPool.Properties.Basic.Monitors = monitors
	}
	if v, ok := d.GetOk("max_connection_attempts"); ok {
		createPool.Properties.Basic.MaxConnectionAttempts = uint(v.(int))
	}
	if v, ok := d.GetOk("max_idle_connections_pernode"); ok {
		createPool.Properties.Basic.MaxIdleConnectionsPerNode = uint(v.(int))
	}
	if v, ok := d.GetOk("max_timed_out_connection_attempts"); ok {
		createPool.Properties.Basic.MaxTimeoutConnectionAttempts = uint(v.(int))
	}
	if v, _ := d.GetOk("node_close_with_rst"); v != nil {
		nodeCloseWithReset := v.(bool)
		createPool.Properties.Basic.NodeCloseWithReset = &nodeCloseWithReset
	}

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

	err := client.Set("pools", poolName, createPool, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Pool error whilst creating %s: %v", poolName, err)
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
			log.Printf("BrocadeVTM Pool Resource %s does not exists", poolName)
			return nil
		}
		return fmt.Errorf("BrocadeVTM Pool error whilst retrieving %s: %v", poolName, err)
	}

	d.Set("name", poolName)
	d.Set("node", poolObj.Properties.Basic.NodesTable)
	d.Set("monitorlist", poolObj.Properties.Basic.Monitors)
	d.Set("max_connection_attempts", poolObj.Properties.Basic.MaxConnectionAttempts)
	d.Set("max_idle_connections_pernode", poolObj.Properties.Basic.MaxIdleConnectionsPerNode)
	d.Set("max_timed_out_connection_attempts", poolObj.Properties.Basic.MaxTimeoutConnectionAttempts)
	d.Set("node_close_with_rst", *poolObj.Properties.Basic.NodeCloseWithReset)
	d.Set("max_connection_timeout", poolObj.Properties.Connection.MaxConnectTime)
	d.Set("max_connections_per_node", poolObj.Properties.Connection.MaxConnectionsPerNode)
	d.Set("max_queue_size", poolObj.Properties.Connection.MaxQueueSize)
	d.Set("max_reply_time", poolObj.Properties.Connection.MaxReplyTime)
	d.Set("queue_timeout", poolObj.Properties.Connection.QueueTimeout)
	d.Set("http_keepalive", *poolObj.Properties.HTTP.HTTPKeepAlive)
	d.Set("http_keepalive_non_idempotent", *poolObj.Properties.HTTP.HTTPKeepAliveNonIdempotent)
	d.Set("load_balancing_priority_enabled", *poolObj.Properties.LoadBalancing.PriorityEnabled)
	d.Set("load_balancing_priority_nodes", poolObj.Properties.LoadBalancing.PriorityNodes)
	d.Set("load_balancing_algorithm", poolObj.Properties.LoadBalancing.Algorithm)
	d.Set("tcp_nagle", *poolObj.Properties.TCP.Nagle)

	return nil
}

// resourcePoolUpdate - Updates an existing pool resource
func resourcePoolUpdate(d *schema.ResourceData, m interface{}) error {

	var poolName string
	var updatePool pool.Pool
	hasChanges := false

	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)
	}

	if d.HasChange("node") {
		if v, ok := d.GetOk("node"); ok {
			if nodes, ok := v.(*schema.Set); ok {
				nodeList := []pool.MemberNode{}
				for _, value := range nodes.List() {
					nodeObject := value.(map[string]interface{})
					newNode := pool.MemberNode{}
					if nodeValue, ok := nodeObject["node"].(string); ok {
						newNode.Node = nodeValue
					}
					if priorityValue, ok := nodeObject["priority"].(int); ok {
						newNode.Priority = priorityValue
					}
					if stateValue, ok := nodeObject["state"].(string); ok {
						newNode.State = stateValue
					}
					if weightValue, ok := nodeObject["weight"].(int); ok {
						newNode.Weight = weightValue
					}
					nodeList = append(nodeList, newNode)

				}
				updatePool.Properties.Basic.NodesTable = nodeList
			}
		}
		hasChanges = true
	}

	if d.HasChange("monitorlist") {
		if v, ok := d.GetOk("monitorlist"); ok {
			originalMonitors := v.([]interface{})
			monitors := make([]string, len(originalMonitors))
			for i, monitor := range originalMonitors {
				monitors[i] = monitor.(string)
			}
			updatePool.Properties.Basic.Monitors = monitors
		}
		hasChanges = true
	}

	if d.HasChange("max_connection_attempts") {
		if v, ok := d.GetOk("max_connection_attempts"); ok {
			updatePool.Properties.Basic.MaxConnectionAttempts = uint(v.(int))
		}
		hasChanges = true
	}

	if d.HasChange("max_idle_connections_pernode") {
		if v, ok := d.GetOk("max_idle_connections_pernode"); ok {
			updatePool.Properties.Basic.MaxIdleConnectionsPerNode = uint(v.(int))
		}
		hasChanges = true
	}

	if d.HasChange("max_timed_out_connection_attempts") {
		if v, ok := d.GetOk("max_timed_out_connection_attempts"); ok {
			updatePool.Properties.Basic.MaxTimeoutConnectionAttempts = uint(v.(int))
		}
		hasChanges = true
	}

	if d.HasChange("node_close_with_rst") {
		nodeCloseWithRst := d.Get("node_close_with_rst").(bool)
		updatePool.Properties.Basic.NodeCloseWithReset = &nodeCloseWithRst
		hasChanges = true
	}

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

	if hasChanges {
		config := m.(map[string]interface{})
		client := config["jsonClient"].(*api.Client)

		err := client.Set("pools", poolName, updatePool, nil)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Pool error whilst updating %s: %v", poolName, err)
		}
		d.SetId(poolName)
	}
	return resourcePoolRead(d, m)

}

// resourcePoolDelete - Deletes a pool resource
func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	var poolName string
	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)
	}

	err := client.Delete("pools", poolName)
	if client.StatusCode == http.StatusNoContent || client.StatusCode == http.StatusNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("BrocadeVTM Pool error whilst deleting %s: %v", poolName, err)
	}
	return nil
}
