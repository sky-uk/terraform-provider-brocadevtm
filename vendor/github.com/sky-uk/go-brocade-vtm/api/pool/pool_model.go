package pool

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
	Ssl                        Ssl                        `json:"ssl,omitempty"`
	TCP                        TCP                        `json:"tcp,omitempty"`
	UDP                        UDP                        `json:"udp,omitempty"`
}

// Basic - main pool definitions
type Basic struct {
	BandwidthClass               string       `json:"bandwidth_class,omitempty"`
	FailurePool                  string       `json:"failure_pool,omitempty"`
	LARDSize                     uint         `json:"lard_size,omitempty"`
	MaxConnectionAttempts        uint         `json:"max_connection_attempts,omitempty"`
	MaxIdleConnectionsPerNode    uint         `json:"max_idle_connections_pernode,omitempty"`
	MaxTimeoutConnectionAttempts uint         `json:"max_timed_out_connection_attempts,omitempty"`
	Monitors                     []string     `json:"monitors,omitempty"`
	NodeCloseWithReset           *bool        `json:"node_close_with_rst,omitempty"`
	NodeConnectionAttempts       uint         `json:"node_connection_attempts,omitempty"`
	NodeDeleteBehavior           string       `json:"node_delete_behavior,omitempty"`
	NodeDrainDeleteTimeout       uint         `json:"node_drain_to_delete_timeout"`
	NodesTable                   []MemberNode `json:"nodes_table,omitempty"`
	Note                         string       `json:"note,omitempty"`
	PassiveMonitoring            *bool        `json:"passive_monitoring,omitempty"`
	PersistenceClass             string       `json:"persistence_class,omitempty"`
	Transparent                  *bool        `json:"transparent,omitempty"`
}

// AutoScaling - AutoScaling settings
type AutoScaling struct {
	AddNodeDelayTime uint     `json:"transparent"`
	CloudCredentials string   `json:"cloud_credentials,omitempty"`
	Cluster          string   `json:"cluster,omitempty"`
	DataCenter       string   `json:"data_center,omitempty"`
	DataStore        string   `json:"data_store,omitempty"`
	Enabled          bool     `json:"enabled"`
	External         bool     `json:"external"`
	ExtraArgs        string   `json:"extraargs,omitempty"`
	Hysteresis       uint     `json:"hysteresis"`
	ImageID          string   `json:"imageid,omitempty"`
	IPsToUse         string   `json:"ips_to_use,omitempty"`
	LastNodeIdleTime uint     `json:"last_node_idle_time"`
	MaxNodes         uint     `json:"max_nodes"`
	MinNodes         uint     `json:"min_nodes"`
	Name             string   `json:"name,omitempty"`
	Port             uint     `json:"port"`
	Refractory       uint     `json:"refractory"`
	ResponseTime     uint     `json:"response_time"`
	ScaleDownLevel   uint     `json:"scale_down_level"`
	ScaleUpLevel     uint     `json:"scale_up_level"`
	SecurityGroupIDs []string `json:"securitygroupids,omitempty"`
	SizeID           string   `json:"size_id,omitempty"`
	SubnetIDs        []string `json:"subnetids,omitempty"`
}

// Connection - Connection settings
type Connection struct {
	MaxConnectTime        uint `json:"max_connect_time,omitempty"`
	MaxConnectionsPerNode uint `json:"max_connections_per_node,omitempty"`
	MaxQueueSize          uint `json:"max_queue_size,omitempty"`
	MaxReplyTime          uint `json:"max_reply_time,omitempty"`
	QueueTimeout          uint `json:"queue_timeout,omitempty"`
}

// DNSAutoScale - DNSAutoScale settings
type DNSAutoScale struct {
	Enabled   bool     `json:"enabled"`
	Hostnames []string `json:"hostnames"`
	Port      uint     `json:"port"`
}

// FTP - FTP settings
type FTP struct {
	SupportRFC2428 bool `json:"support_rfc_2428"`
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

// LoadBalancing - Pool Load balancing settings
type LoadBalancing struct {
	Algorithm       string `json:"algorithm,omitempty"`
	PriorityEnabled *bool  `json:"priority_enabled,omitempty"`
	PriorityNodes   uint   `json:"priority_nodes,omitempty"`
}

// Node - Node Specific settings
type Node struct {
	CloseOnDeath  *bool `json:"close_on_death,omitempty"`
	RetryFailTime int   `json:"retry_fail_time,omitempty"`
}

// SMTP - SMTP settings
type SMTP struct {
	SendSTARTTLS bool `json:"send_starttls"`
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
	SSLSupportSSL2      string   `json:"ssl_support_ssl2,omitempty"`
	SSLSupportSSL3      string   `json:"ssl_support_ssl3,omitempty"`
	SSLSupportTLS1      string   `json:"ssl_support_tls1,omitempty"`
	SSLSupportTLS2      string   `json:"ssl_support_tls2,omitempty"`
	StrictVerify        bool     `json:"strict_verify"`
}

// TCP - tcp setting
type TCP struct {
	Nagle *bool `json:"nagle,omitempty"`
}

// UDP - UDP setting
type UDP struct {
	AcceptFrom      string `json:"accept_from,omitempty"`
	AcceptFromMask  string `json:"accept_from_mask,omitempty"`
	ResponseTimeout uint   `json:"response_timeout"`
}

// MemberNode - Pool membership details / node /state / weight
type MemberNode struct {
	Node     string `json:"node"`
	Priority int    `json:"priority"`
	State    string `json:"state"`
	Weight   int    `json:"weight"`
	SourceIP string `json:"source_ip"`
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
