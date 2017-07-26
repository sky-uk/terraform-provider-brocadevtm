package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-rest-api"
	"github.com/sky-uk/go-brocade-vtm/api/pool"
	"log"
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
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNode,
						},
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
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
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_connection_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"max_idle_connections_pernode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"max_timed_out_connection_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"node_close_with_rst": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"max_connection_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"max_connections_per_node": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"max_queue_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"max_reply_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"queue_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"http_keepalive": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"http_keepalive_non_idempotent": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"load_balancing_priority_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"load_balancing_priority_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"tcp_nagle": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
		},
	}

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

	vtmClient := m.(*rest.Client)
	var createPool pool.Pool
	var poolName string
	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)
	} else {
		return fmt.Errorf("Pool name argument required")
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
	if v, ok := d.GetOk("max_connection_attempts"); ok {
		createPool.Properties.Basic.MaxConnectionAttempts = v.(int)
	}
	if v, ok := d.GetOk("max_idle_connections_pernode"); ok {
		createPool.Properties.Basic.MaxIdleConnectionsPerNode = v.(int)
	}
	if v, ok := d.GetOk("max_timed_out_connection_attempts"); ok {
		createPool.Properties.Basic.MaxTimeoutConnectionAttempts = v.(int)
	}
	if v, ok := d.GetOk("monitorlist"); ok {
		originalMonitors := v.([]interface{})
		monitors := make([]string, len(originalMonitors))
		for i, monitor := range originalMonitors {
			monitors[i] = monitor.(string)
		}
		createPool.Properties.Basic.Monitors = monitors
	}
	if v, ok := d.GetOk("node_close_with_rst"); ok {
		nodeCloseWithReset := v.(bool)
		createPool.Properties.Basic.NodeCloseWithReset = &nodeCloseWithReset
	}

	if v, ok := d.GetOk("max_connection_timeout"); ok {
		createPool.Properties.Connection.MaxConnectTime = v.(int)
	}
	if v, ok := d.GetOk("max_connections_per_node"); ok {
		createPool.Properties.Connection.MaxConnectionsPerNode = v.(int)
	}
	if v, ok := d.GetOk("max_queue_size"); ok {
		createPool.Properties.Connection.MaxQueueSize = v.(int)
	}
	if v, ok := d.GetOk("max_reply_time"); ok {
		createPool.Properties.Connection.MaxReplyTime = v.(int)
	}
	if v, ok := d.GetOk("queue_timeout"); ok {
		createPool.Properties.Connection.QueueTimeout = v.(int)
	}
	if v, ok := d.GetOk("http_keepalive"); ok {
		httpKeepAlive := v.(bool)
		createPool.Properties.HTTP.HTTPKeepAlive = &httpKeepAlive
	}
	if v, ok := d.GetOk("http_keepalive_non_idempotent"); ok {
		httpKeepAliveNonIdempotent := v.(bool)
		createPool.Properties.HTTP.HTTPKeepAliveNonIdempotent = &httpKeepAliveNonIdempotent
	}
	if v, ok := d.GetOk("load_balancing_priority_enabled"); ok {
		loadBalancingPriorityEnabled := v.(bool)
		createPool.Properties.LoadBalancing.PriorityEnabled = &loadBalancingPriorityEnabled
	}
	if v, ok := d.GetOk("load_balancing_priority_nodes"); ok {
		createPool.Properties.LoadBalancing.PriorityNodes = v.(int)
	}
	if v, ok := d.GetOk("tcp_nagle"); ok {
		tcpNagle := v.(bool)
		createPool.Properties.TCP.Nagle = &tcpNagle
	}

	createAPI := pool.NewCreate(poolName, createPool)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Could not create pool: %+v", err)
	}
	if createAPI.StatusCode() != 201 && createAPI.StatusCode() != 200 {
		return fmt.Errorf("Invalid HTTP response code %+v returned. Response object was %+v", createAPI.StatusCode(), createAPI.ResponseObject())
	}

	d.SetId(poolName)
	return resourcePoolRead(d, m)
}

// resourcePoolRead - Reads a  pool resource
func resourcePoolRead(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*rest.Client)
	//var readPool pool.Pool
	var poolName string
	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)
	}
	/*
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
			readPool.Properties.Basic.NodesTable = nodeList
		}

	}
	if v, ok := d.GetOk("max_connection_attempts"); ok {
		readPool.Properties.Basic.MaxConnectionAttempts = v.(int)
	}
	if v, ok := d.GetOk("max_idle_connections_pernode"); ok {
		readPool.Properties.Basic.MaxIdleConnectionsPerNode = v.(int)
	}
	if v, ok := d.GetOk("max_timed_out_connection_attempts"); ok {
		readPool.Properties.Basic.MaxTimeoutConnectionAttempts = v.(int)
	}
	if v, ok := d.GetOk("monitorlist"); ok {
		originalMonitors := v.([]interface{})
		monitors := make([]string, len(originalMonitors))
		for i, monitor := range originalMonitors {
			monitors[i] = monitor.(string)
		}
		readPool.Properties.Basic.Monitors = monitors
	}
	if v, ok := d.GetOk("node_close_with_rst"); ok {
		nodeCloseWithRst := v.(bool)
		readPool.Properties.Basic.NodeCloseWithReset = &nodeCloseWithRst
	}
	if v, ok := d.GetOk("max_connection_timeout"); ok {
		readPool.Properties.Connection.MaxConnectTime = v.(int)
	}
	if v, ok := d.GetOk("max_connections_per_node"); ok {
		readPool.Properties.Connection.MaxConnectionsPerNode = v.(int)
	}
	if v, ok := d.GetOk("max_queue_size"); ok {
		readPool.Properties.Connection.MaxQueueSize = v.(int)
	}
	if v, ok := d.GetOk("max_reply_time"); ok {
		readPool.Properties.Connection.MaxReplyTime = v.(int)
	}
	if v, ok := d.GetOk("queue_timeout"); ok {
		readPool.Properties.Connection.QueueTimeout = v.(int)
	}
	if v, ok := d.GetOk("http_keepalive"); ok {
		httpKeepAlive := v.(bool)
		readPool.Properties.HTTP.HTTPKeepAlive = &httpKeepAlive
	}
	if v, ok := d.GetOk("http_keepalive_non_idempotent"); ok {
		httpKeepAliveNonIdempotent := v.(bool)
		readPool.Properties.HTTP.HTTPKeepAliveNonIdempotent = &httpKeepAliveNonIdempotent
	}
	if v, ok := d.GetOk("load_balancing_priority_enabled"); ok {
		loadBalancingPriorityEnabled := v.(bool)
		readPool.Properties.LoadBalancing.PriorityEnabled = &loadBalancingPriorityEnabled
	}
	if v, ok := d.GetOk("load_balancing_priority_nodes"); ok {
		readPool.Properties.LoadBalancing.PriorityNodes = v.(int)
	}
	if v, ok := d.GetOk("tcp_nagle"); ok {
		tcpNagle := v.(bool)
		readPool.Properties.TCP.Nagle = &tcpNagle
	}
	*/

	getAPI := pool.NewGet(poolName)
	readErr := vtmClient.Do(getAPI)
	if readErr != nil {
		return fmt.Errorf("Error reading pool:", readErr)
	}
	response := getAPI.ResponseObject().(*pool.Pool)

	d.Set("name", poolName)
	d.Set("node", response.Properties.Basic.NodesTable)
	d.Set("max_connection_attempts", response.Properties.Basic.MaxConnectionAttempts)
	d.Set("max_idle_connections_pernode", response.Properties.Basic.MaxIdleConnectionsPerNode)
	d.Set("max_timed_out_connection_attempts", response.Properties.Basic.MaxTimeoutConnectionAttempts)
	d.Set("monitorlist", response.Properties.Basic.Monitors)
	d.Set("node_close_with_rst", response.Properties.Basic.NodeCloseWithReset)
	d.Set("max_connection_timeout", response.Properties.Connection.MaxConnectTime)
	d.Set("max_connections_per_node", response.Properties.Connection.MaxConnectionsPerNode)
	d.Set("max_queue_size", response.Properties.Connection.MaxQueueSize)
	d.Set("max_reply_time", response.Properties.Connection.MaxReplyTime)
	d.Set("queue_timeout", response.Properties.Connection.QueueTimeout)
	d.Set("http_keepalive", response.Properties.HTTP.HTTPKeepAlive)
	d.Set("http_keepalive_non_idempotent", response.Properties.HTTP.HTTPKeepAliveNonIdempotent)
	d.Set("load_balancing_priority_enabled", response.Properties.LoadBalancing.PriorityEnabled)
	d.Set("load_balancing_priority_nodes", response.Properties.LoadBalancing.PriorityNodes)
	d.Set("tcp_nagle", response.Properties.TCP.Nagle)

	return nil
}

// resourcePoolDelete - Deletes a pool resource

func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*rest.Client)
	var poolName string
	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)
	}
	deleteAPI := pool.NewDelete(poolName)
	deleteErr := vtmClient.Do(deleteAPI)
	if deleteErr != nil {
		log.Println("Error Deleting the pool:", deleteErr)
	}
	d.SetId("")
	return nil
}

// resourcePoolUpdate - Updates an existing pool resource
func resourcePoolUpdate(d *schema.ResourceData, m interface{}) error {
	vtmClient := m.(*rest.Client)
	var poolName string
	var updatePool pool.Pool
	if v, ok := d.GetOk("name"); ok {
		poolName = v.(string)

	} else {
		return fmt.Errorf("Pool name argument required")
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
	}

	if d.HasChange("max_connection_attempts") {
		if v, ok := d.GetOk("max_connection_attempts"); ok {
			updatePool.Properties.Basic.MaxConnectionAttempts = v.(int)
		}
	}

	if d.HasChange("max_idle_connections_pernode") {
		if v, ok := d.GetOk("max_idle_connections_pernode"); ok {
			updatePool.Properties.Basic.MaxIdleConnectionsPerNode = v.(int)
		}
	}

	if d.HasChange("max_timed_out_connection_attempts") {
		if v, ok := d.GetOk("max_timed_out_connection_attempts"); ok {
			updatePool.Properties.Basic.MaxTimeoutConnectionAttempts = v.(int)
		}
	}

	if d.HasChange("node_close_with_rst") {
		if v, ok := d.GetOk("node_close_with_rst"); ok {
			nodeCloseWithRst := v.(bool)
			updatePool.Properties.Basic.NodeCloseWithReset = &nodeCloseWithRst
		}
	}

	if d.HasChange("max_connection_timeout") {
		if v, ok := d.GetOk("max_connection_timeout"); ok {
			updatePool.Properties.Connection.MaxConnectTime = v.(int)
		}
	}

	if d.HasChange("max_connections_per_node") {
		if v, ok := d.GetOk("max_connections_per_node"); ok {
			updatePool.Properties.Connection.MaxConnectionsPerNode = v.(int)
		}
	}

	if d.HasChange("max_queue_size") {
		if v, ok := d.GetOk("max_queue_size"); ok {
			updatePool.Properties.Connection.MaxQueueSize = v.(int)
		}
	}

	if d.HasChange("max_reply_time") {
		if v, ok := d.GetOk("max_reply_time"); ok {
			updatePool.Properties.Connection.MaxReplyTime = v.(int)
		}
	}

	if d.HasChange("queue_timeout") {
		if v, ok := d.GetOk("queue_timeout"); ok {
			updatePool.Properties.Connection.QueueTimeout = v.(int)
		}
	}

	if d.HasChange("http_keepalive") {
		if v, ok := d.GetOk("http_keepalive"); ok {
			httpKeepAlive := v.(bool)
			updatePool.Properties.HTTP.HTTPKeepAlive = &httpKeepAlive
		}
	}

	if d.HasChange("http_keepalive_non_idempotent") {
		if v, ok := d.GetOk("http_keepalive_non_idempotent"); ok {
			httpKeepAliveNonIdempotent := v.(bool)
			updatePool.Properties.HTTP.HTTPKeepAliveNonIdempotent = &httpKeepAliveNonIdempotent
		}
	}

	if d.HasChange("load_balancing_priority_enabled") {
		if v, ok := d.GetOk("load_balancing_priority_enabled"); ok {
			loadBalancingPriorityEnabled := v.(bool)
			updatePool.Properties.LoadBalancing.PriorityEnabled = &loadBalancingPriorityEnabled
		}
	}

	if d.HasChange("load_balancing_priority_nodes") {
		if v, ok := d.GetOk("load_balancing_priority_nodes"); ok {
			updatePool.Properties.LoadBalancing.PriorityNodes = v.(int)
		}
	}

	if d.HasChange("tcp_nagle") {
		if v, ok := d.GetOk("tcp_nagle"); ok {
			tcpNagle := v.(bool)
			updatePool.Properties.TCP.Nagle = &tcpNagle
		}
	}

	updatePoolAPI := pool.NewUpdate(poolName, updatePool)
	updatePoolErr := vtmClient.Do(updatePoolAPI)
	if updatePoolErr != nil {
		return fmt.Errorf("Error updating pool %s", updatePoolErr)
	}

	return resourcePoolRead(d, m)

}
