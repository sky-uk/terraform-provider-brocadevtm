package virtualserver

// VirtualServer data structure
type VirtualServer struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall monitor configuration
type Properties struct {
	Basic            Basic            `json:"basic"`
	Aptimizer        Aptimizer        `json:"aptimizer,omitempty"`
	Connection       Connection       `json:"connection,omitempty"`
	ConnectionErrors ConnectionErrors `json:"connection_errors,omitempty"`
	Cookie           Cookie           `json:"cookie,omitempty"`
	DNS              DNS              `json:"dns,omitempty"`
	Ftp              Ftp              `json:"ftp,omitempty"`
	Gzip             Gzip             `json:"gzip,omitempty"`
	Ssl              Ssl              `json:"ssl,omitempty"`
}

// Basic : Basic virtual server configration
type Basic struct {
	AddClusterIP         *bool    `json:"add_cluster_ip,omitempty"`        // true
	AddXForwarded        *bool    `json:"add_x_forwarded_for,omitempty"`   // false
	AddXForwardedProto   *bool    `json:"add_x_forwarded_proto,omitempty"` // false
	BandwidthClass       string   `json:"bandwidth_class,omitempty"`       // ""
	CloseWithRst         *bool    `json:"close_with_rst,omitempty"`        // false
	CompletionRules      []string `json:"completionrules,omitempty"`       // []
	ConnectTimeout       uint     `json:"connect_timeout,omitempty"`       // 10
	Enabled              *bool    `json:"enabled,omitempty"`
	FtpForceServerSecure *bool    `json:"ftp_force_server_secure,omitempty"` // true
	GlbServices          []string `json:"glb_services,omitempty"`            // []
	ListenOnAny          *bool    `json:"listen_on_any,omitempty"`           // false
	ListenOnHosts        []string `json:"listen_on_hosts,omitempty"`         // []
	ListenOnTrafficIps   []string `json:"listen_on_traffic_ips,omitempty"`   // []
	Note                 string   `json:"note,omitempty"`                    // ""
	Pool                 string   `json:"pool,omitempty"`
	Port                 uint     `json:"port,omitempty"`
	ProtectionClass      string   `json:"protection_class,omitempty"` // ""
	Protocol             string   `json:"protocol,omitempty"`         // "http"
	RequestRules         []string `json:"request_rules,omitempty"`    // []
	ResponseRules        []string `json:"response_rules,omitempty"`   // []
	SlmClass             string   `json:"slm_class,omitempty"`        // ""
	SoNagle              *bool    `json:"so_nagle,omitempty"`         // false,
	// SslClientCertHeaders : enum can accept these strings("all", "none"
	// (default), "none"
	SslClientCertHeaders string `json:"ssl_client_cert_headers,omitempty"` // "none"
	SslDecrypt           *bool  `json:"ssl_decrypt,omitempty"`             // false
	SslHonorFallbackScsv string `json:"ssl_honor_fallback_scsv,omitempty"` // "use_default"
	Transparent          *bool  `json:"transparent,omitempty"`             // false
}

// Aptimizer : whether virtual server should aptimize web content
type Aptimizer struct {
	Enabled *bool    `json:"enabled,omitempty"` // false
	Profile []string `json:"profile,omitempty"` // [], NOTE: actually more complex
}

// Connection : connection parameters
type Connection struct {
	Keepalive              *bool  `json:"keepalive,omitempty"`                // false
	KeepaliveTimeout       uint   `json:"keepalive_timeout,omitempty"`        // 10
	MaxClientBuffer        uint   `json:"max_client_buffer,omitempty"`        // 65536
	MaxServerBuffer        uint   `json:"max_server_buffer,omitempty"`        // 65536
	MaxTransactionDuration uint   `json:"max_transaction_duration,omitempty"` // none
	ServerFirstBanner      string `json:"server_first_banner,omitempty"`      // none
	Timeout                uint   `json:"timeout,omitempty"`                  // 300
}

// ConnectionErrors : error file params
type ConnectionErrors struct {
	ErrorFile string `json:"error_file,omitempty"` // "Default"
}

// Cookie : how cookies are handled
type Cookie struct {
	Domain      string `json:"domain,omitempty"`       // "no_rewrite"
	NewDomain   string `json:"new_domain,omitempty"`   // ""
	PathRegex   string `json:"path_regex,omitempty"`   // ""
	PathReplace string `json:"path_replace,omitempty"` // ""
	Secure      string `json:"secure,omitempty"`       // "no_modify"
}

// DNS configuration section
type DNS struct {
	EdnsUdpsize uint     `json:"edns_udpsize,omitempty"` // 4096
	MaxUdpsize  uint     `json:"max_udpsize,omitempty"`  // 4096
	RrsetOrder  string   `json:"rrset_order,omitempty"`  // "fixed"
	Verbose     *bool    `json:"verbose,omitempty"`      // false
	Zones       []string `json:"zones,omitempty"`        // []
}

// Ftp configuration section
type Ftp struct {
	DataSourcePort    uint  `json:"data_source_port,omitempty"`    // 0
	ForceClientSecure *bool `json:"force_client_secure,omitempty"` // true
	PortRangeHigh     uint  `json:"port_range_high,omitempty"`     // 0
	PortRangeLow      uint  `json:"port_range_low,omitempty"`      // 0
	SslData           *bool `json:"ssl_data,omitempty"`            // true
}

// Gzip configuration section
type Gzip struct {
	CompressLevel uint     `json:"compress_level,omitempty"` // 1
	Enabled       *bool    `json:"enabled,omitempty"`        // false
	EtagRewrite   string   `json:"etag_rewrite,omitempty"`   // "wrap"
	IncludeMime   []string `json:"include_mime,omitempty"`   // [
	// "text/html"
	// "text/plain"
	// ]
	MaxSize uint  `json:"max_size,omitempty"` // 10000000
	MinSize uint  `json:"min_size,omitempty"` // 1000
	NoSize  *bool `json:"no_size,omitempty"`  // true
}

// Ssl configuration section
type Ssl struct {
	ServerCertDefault string       `json:"server_cert_default,omitempty"` // ""
	ServerCertHostMap []CertItem   `json:"server_cert_host_mapping,omitempty"`
	SslSupportSsl2    string       `json:"ssl_support_ssl2,omitempty"`
	SslSupportSsl3    string       `json:"ssl_support_ssl3,omitempty"`
	SslSupportTLS1    string       `json:"ssl_support_tls1,omitempty"`
	SslSupportTLS1_1  string       `json:"ssl_support_tls1_1,omitempty"`
	SslSupportTLS1_2  string       `json:"ssl_support_tls1_2,omitempty"`
	OCSPEnable        *bool        `json:"ocsp_enable,omitempty"`
	OCSPIssuers       []OCSPIssuer `json:"ocsp_issuers,omitempty"`
}

// OCSPIssuer configuration serction
type OCSPIssuer struct {
	Issuer        string `json:"issuer,omitempty"`
	AIA           *bool  `json:"aia,omitempty"`
	Nonce         string `json:"nonce,omitempty"`
	Required      string `json:"required,omitempty"`
	ResponderCert string `json:"responder_cert,omitempty"`
	Signer        string `json:"signer,omitempty"`
	URL           string `json:"url,omitempty"`
}

// CertItem : a single certificate item in the cert map
type CertItem struct {
	Host            string   `json:"host,omitempty"`
	AltCertificates []string `json:"alt_certificates,omitempty"`
	Certificate     string   `json:"certificate,omitempty"`
}

// VirtualServersList : List of nodes monitored
type VirtualServersList struct {
	Children []ChildVirtualServer `json:"children"`
}

// ChildVirtualServer : monitored node structure
type ChildVirtualServer struct {
	Name string `json:"name"`
	Href string `json:"href"`
}
