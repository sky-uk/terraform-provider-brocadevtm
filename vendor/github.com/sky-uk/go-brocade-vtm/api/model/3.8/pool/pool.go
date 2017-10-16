package pool

// ChildPools - Used to display data about all pools ie name and link
type ChildPools struct {
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
}

// Pool - main pool struct
type Pool struct {
	Properties Properties `json:"properties"`
}

// Properties - General Properties for the pool
type Properties struct {
	Basic                      Basic                      `json:"basic,omitempty"`
	AutoScaling                AutoScaling                `json:"auto_scaling,omitempty"`
	Connection                 Connection                 `json:"connection,omitempty"`
	DNSAutoScale               DNSAutoScale               `json:"dns_autoscale,omitempty"`
	FTP                        FTP                        `json:"ftp,omitempty"`
	HTTP                       HTTP                       `json:"http,omitempty"`
	KerberosProtocolTransition KerberosProtocolTransition `json:"kerberos_protocol_transition,omitempty"`
	LoadBalancing              LoadBalancing              `json:"load_balancing,omitempty"`
	Node                       Node                       `json:"node,omitempty"`
	SMTP                       SMTP                       `json:"smtp,omitempty"`
	Ssl                        Ssl                        `json:"ssl,omitempty"`
	TCP                        TCP                        `json:"tcp,omitempty"`
	UDP                        UDP                        `json:"udp,omitempty"`
}

// Basic - main pool definitions
type Basic struct {
	BandwidthClass               string       `json:"bandwidth_class,omitempty"`
	FailurePool                  string       `json:"failure_pool,omitempty"`
	MaxConnectionAttempts        *uint        `json:"max_connection_attempts,omitempty"`
	MaxIdleConnectionsPerNode    *uint        `json:"max_idle_connections_pernode,omitempty"`
	MaxTimeoutConnectionAttempts *uint        `json:"max_timed_out_connection_attempts,omitempty"`
	Monitors                     []string     `json:"monitors,omitempty"`
	NodeCloseWithReset           *bool        `json:"node_close_with_rst,omitempty"`
	NodeConnectionAttempts       *uint        `json:"node_connection_attempts,omitempty"`
	NodeDeleteBehavior           string       `json:"node_delete_behavior,omitempty"`
	NodeDrainDeleteTimeout       *uint        `json:"node_drain_to_delete_timeout,omitempty"`
	NodesTable                   []MemberNode `json:"nodes_table,omitempty"`
	Note                         string       `json:"note,omitempty"`
	PassiveMonitoring            *bool        `json:"passive_monitoring,omitempty"`
	PersistenceClass             string       `json:"persistence_class,omitempty"`
	Transparent                  *bool        `json:"transparent,omitempty"`
}

// AutoScaling - AutoScaling settings
type AutoScaling struct {
	AddNodeDelayTime *uint    `json:"addnode_delaytime,omitempty"`
	CloudCredentials string   `json:"cloud_credentials,omitempty"`
	Cluster          string   `json:"cluster,omitempty"`
	DataCenter       string   `json:"data_center,omitempty"`
	DataStore        string   `json:"data_store,omitempty"`
	Enabled          *bool    `json:"enabled,omitempty"`
	External         *bool    `json:"external,omitempty"`
	ExtraArgs        string   `json:"extraargs,omitempty"`
	Hysteresis       *uint    `json:"hysteresis,omitempty"`
	ImageID          string   `json:"imageid,omitempty"`
	IPsToUse         string   `json:"ips_to_use,omitempty"`
	LastNodeIdleTime *uint    `json:"last_node_idle_time,omitempty"`
	MaxNodes         *uint    `json:"max_nodes,omitempty"`
	MinNodes         *uint    `json:"min_nodes,omitempty"`
	Name             string   `json:"name,omitempty,omitempty"`
	Port             *uint    `json:"port,omitempty,omitempty"`
	Refractory       *uint    `json:"refractory,omitempty"`
	ResponseTime     *uint    `json:"response_time,omitempty"`
	ScaleDownLevel   *uint    `json:"scale_down_level,omitempty"`
	ScaleUpLevel     *uint    `json:"scale_up_level,omitempty"`
	SecurityGroupIDs []string `json:"securitygroupids,omitempty"`
	SizeID           string   `json:"size_id,omitempty"`
	SubnetIDs        []string `json:"subnetids,omitempty"`
}

// Connection - Connection setting
type Connection struct {
	MaxConnectTime        *uint `json:"max_connect_time,omitempty"`
	MaxConnectionsPerNode *uint `json:"max_connections_per_node,omitempty"`
	MaxQueueSize          *uint `json:"max_queue_size,omitempty"`
	MaxReplyTime          *uint `json:"max_reply_time,omitempty"`
	QueueTimeout          *uint `json:"queue_timeout,omitempty"`
}

// DNSAutoScale - DNSAutoScale settings
type DNSAutoScale struct {
	Enabled   *bool    `json:"enabled,omitempty"`
	Hostnames []string `json:"hostnames,omitempty"`
	Port      *uint    `json:"port,omitempty"`
}

// FTP - FTP settings
type FTP struct {
	SupportRFC2428 *bool `json:"support_rfc_2428,omitempty"`
}

// HTTP - http settings
type HTTP struct {
	HTTPKeepAlive              *bool `json:"keepalive,omitempty"`
	HTTPKeepAliveNonIdempotent *bool `json:"keepalive_non_idempotent,omitempty"`
}

// KerberosProtocolTransition - KerberosProtocolTransition settings
type KerberosProtocolTransition struct {
	Principal string `json:"principal,omitempty"`
	Target    string `json:"target,omitempty"`
}

// LBPoolList - Used to return all pools
type LBPoolList struct {
	ChildPools []ChildPools `json:"children,omitempty"`
}

// LoadBalancing - Pool Load balancing settings
type LoadBalancing struct {
	Algorithm       string `json:"algorithm,omitempty"`
	PriorityEnabled *bool  `json:"priority_enabled,omitempty"`
	PriorityNodes   *uint  `json:"priority_nodes,omitempty"`
}

// MemberNode - Pool membership details / node /state / weight
type MemberNode struct {
	Node     string `json:"node,omitempty"`
	Priority *uint  `json:"priority,omitempty"`
	State    string `json:"state,omitempty"`
	Weight   *int   `json:"weight,omitempty"`
	SourceIP string `json:"source_ip,omitempty"`
}

// Node - Node Specific settings
type Node struct {
	CloseOnDeath  *bool `json:"close_on_death,omitempty"`
	RetryFailTime *uint `json:"retry_fail_time,omitempty"`
}

// SMTP - SMTP settings
type SMTP struct {
	SendSTARTTLS *bool `json:"send_starttls,omitempty"`
}

// Ssl - SSL related settings
type Ssl struct {
	ClientAuth          *bool    `json:"client_auth,omitempty"`
	CommonNameMatch     []string `json:"common_name_match,omitempty"`
	EllipticCurves      []string `json:"elliptic_curves,omitempty"`
	Enable              *bool    `json:"enable,omitempty"`
	Enhance             *bool    `json:"enhance,omitempty"`
	SendCloseAlerts     *bool    `json:"send_close_alerts,omitempty"`
	ServerName          *bool    `json:"server_name,omitempty"`
	SignatureAlgorithms string   `json:"signature_algorithms,omitempty"`
	SslCiphers          string   `json:"ssl_ciphers,omitempty"`
	SSLSupportSSL2      string   `json:"ssl_support_ssl2,omitempty"`
	SSLSupportSSL3      string   `json:"ssl_support_ssl3,omitempty"`
	SSLSupportTLS1      string   `json:"ssl_support_tls1,omitempty"`
	SSLSupportTLS1_1    string   `json:"ssl_support_tls1_1,omitempty"`
	SSLSupportTLS1_2    string   `json:"ssl_support_tls1_2,omitempty"`
	StrictVerify        *bool    `json:"strict_verify,omitempty"`
}

// TCP - tcp setting
type TCP struct {
	Nagle *bool `json:"nagle,omitempty"`
}

// UDP - UDP setting
type UDP struct {
	AcceptFrom      string `json:"accept_from,omitempty"`
	AcceptFromMask  string `json:"accept_from_mask,omitempty"`
	ResponseTimeout *uint  `json:"response_timeout,omitempty"`
}
