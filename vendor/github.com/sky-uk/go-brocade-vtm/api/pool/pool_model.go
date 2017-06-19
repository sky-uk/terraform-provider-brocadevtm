package pool

// Pool - main pool struct
type Pool struct {
	Properties Properties `json:"properties"`
	Name       string     `json:"name,omitempty"`
}

// Properties - General Properties for the pool
type Properties struct {
	Basic         Basic         `json:"basic,omitempty"`
	Connection    Connection    `json:"connection,omitempty"`
	HTTP          HTTP          `json:"http,omitempty"`
	LoadBalancing LoadBalancing `json:"load_balancing,omitempty"`
	Node          Node          `json:"node,omitempty"`
	Ssl           Ssl           `json:"ssl,omitempty"`
	TCP           TCP           `json:"tcp,omitempty"`
}

// Basic - main pool definitions
type Basic struct {
	BandwidthClass               string       `json:"bandwidth_class,omitempty"`
	FailurePool                  string       `json:"failure_pool,omitempty"`
	MaxConnectionAttempts        int          `json:"max_connection_attempts,omitempty"`
	MaxIdleConnectionsPerNode    int          `json:"max_idle_connections_pernode,omitempty"`
	MaxTimeoutConnectionAttempts int          `json:"max_timed_out_connection_attempts,omitempty"`
	Monitors                     []string     `json:"monitors,omitempty"`
	NodeCloseWithReset           *bool        `json:"node_close_with_rst,omitempty"`
	NodeConnectionAttempts       int          `json:"node_connection_attempts,omitempty"`
	NodeDeleteBehavior           string       `json:"node_delete_behavior,omitempty"`
	NodeDrainDeleteTimeout       int          `json:"node_drain_to_delete_timeout"`
	NodesTable                   []MemberNode `json:"nodes_table,omitempty"`
	Note                         string       `json:"note,omitempty"`
	PassiveMonitoring            *bool        `json:"passive_monitoring,omitempty"`
	PersistenceClass             string       `json:"persistence_class,omitempty"`
	Transparent                  *bool        `json:"transparent,omitempty"`
}

// Connection - Connection setting
type Connection struct {
	MaxConnectTime        int `json:"max_connect_time,omitempty"`
	MaxConnectionsPerNode int `json:"max_connections_per_node,omitempty"`
	MaxQueueSize          int `json:"max_queue_size,omitempty"`
	MaxReplyTime          int `json:"max_reply_time,omitempty"`
	QueueTimeout          int `json:"queue_timeout,omitempty"`
}

// HTTP - http settings
type HTTP struct {
	HTTPKeepAlive              *bool `json:"keepalive,omitempty"`
	HTTPKeepAliveNonIdempotent *bool `json:"keepalive_non_idempotent,omitempty"`
}

// LoadBalancing - Pool Load balancing settings
type LoadBalancing struct {
	Algorithm       string `json:"algorithm,omitempty"`
	PriorityEnabled *bool  `json:"priority_enabled,omitempty"`
	PriorityNodes   int    `json:"priority_nodes,omitempty"`
}

// Node - Node Specific settings
type Node struct {
	CloseOnDeath  *bool `json:"close_on_death,omitempty"`
	RetryFailTime int   `json:"retry_fail_time,omitempty"`
}

// Ssl - SSL related settings
type Ssl struct {
	ClientAuth          *bool    `json:"client_auth,omitempty"`
	CommonNameMatch     []string `json:"common_name_match,omitempty"`
	ElipticCurves       []string `json:"eliptic_curves,omitempty"`
	Enabled             *bool    `json:"enabled,omitempty"`
	Enhance             *bool    `json:"enhance,omitempty"`
	SendCloseAlerts     *bool    `json:"send_close_alerts,omitempty"`
	ServerName          *bool    `json:"server_name,omitempty"`
	SignatureAlgorithms string   `json:"signature_algorithms,omitempty"`
	SslCiphers          string   `json:"ssl_ciphers,omitempty"`
}

// TCP - tcp setting
type TCP struct {
	Nagle *bool `json:"enabled,omitempty"`
}

// MemberNode - Pool membership details / node /state / weight
type MemberNode struct {
	Node     string `json:"node"`
	Priority int    `json:"priority"`
	State    string `json:"state"`
	Weight   int    `json:"weight"`
}

// LBPoolList - Used to return all pools
type LBPoolList struct {
	ChildPools []ChildPools `json:"children,omitempty"`
}

// ChildPools - Used to display data about all pools ie name and link
type ChildPools struct {
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
}
