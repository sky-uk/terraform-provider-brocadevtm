package virtualserver

// VirtualServer data structure
type VirtualServer struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall monitor configuration
type Properties struct {
	Basic                      Basic                      `json:"basic"`
	Aptimizer                  Aptimizer                  `json:"aptimizer,omitempty"`
	Connection                 Connection                 `json:"connection,omitempty"`
	ConnectionErrors           ConnectionErrors           `json:"connection_errors,omitempty"`
	Cookie                     Cookie                     `json:"cookie,omitempty"`
	DNS                        DNS                        `json:"dns,omitempty"`
	Ftp                        Ftp                        `json:"ftp,omitempty"`
	Gzip                       Gzip                       `json:"gzip,omitempty"`
	HTTP                       HTTP                       `json:"http,omitempty"`
	HTTP2                      HTTP2                      `json:"http2,omitempty"`
	KerberosProtocolTransition KerberosProtocolTransition `json:"kerberos_protocol_transition"`
	Log                        Log                        `json:"log,omitempty"`
	RecentConnections          RecentConnections          `json:"recent_connections,omitempty"`
	RequestTracing             RequestTracing             `json:"request_tracing,omitempty"`
	RTSP                       RTSP                       `json:"rtsp,omitempty"`
	SIP                        SIP                        `json:"sip,omitempty"`
	SMTP                       SMTP                       `json:"smtp,omitempty"`
	Ssl                        Ssl                        `json:"ssl,omitempty"`
	SysLog                     SysLog                     `json:"sys_log,omitempty"`
	TCP                        TCP                        `json:"tcp,omitempty"`
	UDP                        UDP                        `json:"udp,omitempty"`
	WebCache                   WebCache                   `json:"web_cache,omitempty"`
}

// Basic : Basic virtual server configration
type Basic struct {
	AddClusterIP             *bool    `json:"add_cluster_ip,omitempty"`        // true
	AddXForwarded            *bool    `json:"add_x_forwarded_for,omitempty"`   // false
	AddXForwardedProto       *bool    `json:"add_x_forwarded_proto,omitempty"` // false
	AutoUpgradeProtocols     []string `json:"auto_upgrade_protocols,omitempty"`
	AutoDetectUpgradeHeaders bool     `json:"autodetect_upgrade_headers"`
	BandwidthClass           string   `json:"bandwidth_class,omitempty"` // ""
	CloseWithRst             *bool    `json:"close_with_rst,omitempty"`  // false
	CompletionRules          []string `json:"completionrules,omitempty"` // []
	ConnectTimeout           uint     `json:"connect_timeout"`           // 10
	Enabled                  *bool    `json:"enabled,omitempty"`
	FtpForceServerSecure     *bool    `json:"ftp_force_server_secure,omitempty"` // true
	GlbServices              []string `json:"glb_services,omitempty"`            // []
	ListenOnAny              *bool    `json:"listen_on_any,omitempty"`           // false
	ListenOnHosts            []string `json:"listen_on_hosts,omitempty"`         // []
	ListenOnTrafficIps       []string `json:"listen_on_traffic_ips,omitempty"`   // []
	MSS                      uint     `json:"mss"`
	Note                     string   `json:"note,omitempty"` // ""
	Pool                     string   `json:"pool,omitempty"`
	Port                     uint     `json:"port,omitempty"`
	ProtectionClass          string   `json:"protection_class,omitempty"` // ""
	Protocol                 string   `json:"protocol,omitempty"`         // "http"
	RequestRules             []string `json:"request_rules,omitempty"`    // []
	ResponseRules            []string `json:"response_rules,omitempty"`   // []
	SlmClass                 string   `json:"slm_class,omitempty"`        // ""
	SoNagle                  *bool    `json:"so_nagle,omitempty"`         // false,
	// SslClientCertHeaders : enum can accept these strings("all", "none"
	// (default), "none"
	SslClientCertHeaders string `json:"ssl_client_cert_headers,omitempty"` // "none"
	SslDecrypt           *bool  `json:"ssl_decrypt,omitempty"`             // false
	SslHonorFallbackScsv string `json:"ssl_honor_fallback_scsv,omitempty"` // "use_default"
	Transparent          *bool  `json:"transparent,omitempty"`             // false
}

// Aptimizer : whether virtual server should aptimize web content
type Aptimizer struct {
	Enabled *bool              `json:"enabled,omitempty"` // false
	Profile []AptimizerProfile `json:"profile,omitempty"` // [], NOTE: actually more complex
}

// AptimizerProfile : Aptimizer profiles and the application scopes that apply to them.
type AptimizerProfile struct {
	Name string   `json:"name"`
	URLs []string `json:"urls"`
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
	EDNSClientSubnet bool     `json:"edns_client_subnet"`
	EdnsUdpsize      uint     `json:"edns_udpsize,omitempty"` // 4096
	MaxUdpsize       uint     `json:"max_udpsize,omitempty"`  // 4096
	RrsetOrder       string   `json:"rrset_order,omitempty"`  // "fixed"
	Verbose          *bool    `json:"verbose,omitempty"`      // false
	Zones            []string `json:"zones,omitempty"`        // []
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

// HTTP configuration section
type HTTP struct {
	ChunkOverheadForwarding string `json:"chunk_overhead_forwarding,omitempty"`
	LocationRegex           string `json:"location_regex,omitempty"`
	LocationReplace         string `json:"location_replace,omitempty"`
	LocationRewrite         string `json:"location_rewrite,omitempty"`
	MIMEDefault             string `json:"mime_default,omitempty"`
	MIMEDetect              bool   `json:"mime_detect"`
}

// HTTP2 configuration section
type HTTP2 struct {
	ConnectTimeout         uint     `json:"connect_timeout"`
	DataFrameSize          uint     `json:"data_frame_size"`
	Enabled                bool     `json:"enabled"`
	HeaderTableSize        uint     `json:"header_table_size"`
	HeadersIndexBlacklist  []string `json:"headers_index_blacklist,omitempty"`
	HeadersIndexDefault    bool     `json:"headers_index_default"`
	HeadersIndexWhitelist  []string `json:"headers_index_whitelist,omitempty"`
	IdleTimeoutNoStreams   uint     `json:"idle_timeout_no_streams"`
	IdleTimeoutOpenStreams uint     `json:"idle_timeout_open_streams"`
	MaxConcurrentStreams   uint     `json:"max_concurrent_streams"`
	MaxFrameSize           uint     `json:"max_frame_size"`
	MaxHeaderPadding       uint     `json:"max_header_padding"`
	MergeCookieHeaders     bool     `json:"merge_cookie_headers"`
	StreamWindowSize       uint     `json:"stream_window_size"`
}

// KerberosProtocolTransition configuration section
type KerberosProtocolTransition struct {
	Enabled   bool   `json:"enabled"`
	Principal string `json:"principal,omitempty"`
	Target    string `json:"target,omitempty"`
}

// Log configuration section
type Log struct {
	AlwaysFlush               bool   `json:"always_flush"`
	ClientConnectionFailures  bool   `json:"client_connection_failures"`
	Enabled                   bool   `json:"enabled"`
	Filename                  string `json:"filename,omitempty"`
	Format                    string `json:"format,omitempty"`
	SaveAll                   bool   `json:"save_all"`
	ServerConnectionFailures  bool   `json:"server_connection_failures"`
	SessionPersistenceVerbose bool   `json:"session_persistence_verbose"`
	SSLFailures               bool   `json:"ssl_failures"`
}

// RecentConnections configuration section
type RecentConnections struct {
	Enabled bool `json:"enabled"`
	SaveAll bool `json:"save_all"`
}

// RequestTracing configuration section
type RequestTracing struct {
	Enabled bool `json:"enabled"`
	TraceIO bool `json:"trace_io"`
}

// RTSP configuration section
type RTSP struct {
	StreamingPortRangeHigh uint `json:"streaming_port_range_high"`
	StreamingPortRangeLow  uint `json:"streaming_port_range_low"`
	StreamingTimeout       uint `json:"streaming_timeout"`
}

// SIP configuration section
type SIP struct {
	DangerousRequests      string `json:"dangerous_requests,omitempty"`
	FollowRoute            bool   `json:"follow_route"`
	MaxConnectionMem       uint   `json:"max_connection_mem"`
	Mode                   string `json:"mode,omitempty"`
	RewriteURI             bool   `json:"rewrite_uri"`
	StreamingPortRangeHigh uint   `json:"streaming_port_range_high"`
	StreamingPortRangeLow  uint   `json:"streaming_port_range_low"`
	StreamingTimeout       uint   `json:"streaming_timeout"`
	TimeoutMessages        bool   `json:"timeout_messages"`
	TransactionTimeout     uint   `json:"transaction_timeout"`
}

// SMTP configuration section
type SMTP struct {
	ExpectSTARTTLS bool `json:"expect_starttls"`
}

// Ssl configuration section
type Ssl struct {
	AddHTTPHeaders            bool         `json:"add_http_headers"`
	ClientCertCAS             []string     `json:"client_cert_cas,omitempty"`
	EllipticCurves            []string     `json:"elliptic_curves,omitempty"`
	IssuedCertsNeverExpire    []string     `json:"issued_certs_never_expire,omitempty"`
	OCSPEnable                *bool        `json:"ocsp_enable,omitempty"`
	OCSPIssuers               []OCSPIssuer `json:"ocsp_issuers,omitempty"`
	OSCPMaxResponseAge        uint         `json:"oscp_max_response_age"`
	OSCPStapling              bool         `json:"oscp_stapling"`
	OSCPTimeTolerance         uint         `json:"oscp_time_tolerance"`
	OSCPTimeout               uint         `json:"oscp_timeout"`
	PreferSSLv3               bool         `json:"prefer_sslv3"`
	RequestClientCert         string       `json:"request_client_cert,omitempty"`
	SendCloseAlerts           bool         `json:"send_close_alerts"`
	ServerCertAltCertificates []string     `json:"server_cert_alt_certificates,omitempty"`
	ServerCertDefault         string       `json:"server_cert_default,omitempty"` // ""
	ServerCertHostMap         []CertItem   `json:"server_cert_host_mapping,omitempty"`
	SignatureAlgorithms       string       `json:"signature_algorithms,omitempty"`
	SSLCiphers                string       `json:"ssl_ciphers,omitempty"`
	SslSupportSsl2            string       `json:"ssl_support_ssl2,omitempty"`
	SslSupportSsl3            string       `json:"ssl_support_ssl3,omitempty"`
	SslSupportTLS1            string       `json:"ssl_support_tls1,omitempty"`
	SslSupportTLS1_1          string       `json:"ssl_support_tls1_1,omitempty"`
	SslSupportTLS1_2          string       `json:"ssl_support_tls1_2,omitempty"`
	TrustMagic                bool         `json:"trust_magic"`
}

// SysLog configuration section
type SysLog struct {
	Enabled     bool   `json:"enabled"`
	Format      string `json:"format,omitempty"`
	IPEndpoint  string `json:"ip_endpoint,omitempty"`
	MsgLenLimit uint   `json:"msg_len_limit"`
}

// TCP configuration section
type TCP struct {
	ProxyClose bool `json:"proxy_close"`
}

//UDP configuration section
type UDP struct {
	EndPointPersistence       bool `json:"end_point_persistence"`
	PortSMP                   bool `json:"port_smp"`
	ResponseDatagramsExpected int  `json:"response_datagrams_expected"`
	Timeout                   uint `json:"timeout"`
}

// WebCache configuration section
type WebCache struct {
	ControlOut    string `json:"control_out,omitempty"`
	Enabled       bool   `json:"enabled"`
	ErrorPageTime uint   `json:"error_page_time"`
	MaxTime       uint   `json:"max_time"`
	RefreshTime   uint   `json:"refresh_time"`
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
