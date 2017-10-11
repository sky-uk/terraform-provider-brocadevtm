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
	KerberosProtocolTransition KerberosProtocolTransition `json:"kerberos_protocol_transition,omitempty"`
	Log                        Log                        `json:"log,omitempty"`
	RecentConnections          RecentConnections          `json:"recent_connections,omitempty"`
	RequestTracing             RequestTracing             `json:"request_tracing,omitempty"`
	RTSP                       RTSP                       `json:"rtsp,omitempty"`
	SIP                        SIP                        `json:"sip,omitempty"`
	SMTP                       SMTP                       `json:"smtp,omitempty"`
	Ssl                        Ssl                        `json:"ssl,omitempty"`
	SysLog                     SysLog                     `json:"syslog,omitempty"`
	TCP                        TCP                        `json:"tcp,omitempty"`
	UDP                        UDP                        `json:"udp,omitempty"`
	WebCache                   WebCache                   `json:"web_cache,omitempty"`
}

// Basic : Basic virtual server configration
type Basic struct {
	AddClusterIP             bool     `json:"add_cluster_ip"`
	AddXForwarded            bool     `json:"add_x_forwarded_for"`
	AddXForwardedProto       bool     `json:"add_x_forwarded_proto"`
	AutoUpgradeProtocols     []string `json:"auto_upgrade_protocols,omitempty"`
	AutoDetectUpgradeHeaders bool     `json:"autodetect_upgrade_headers"`
	BandwidthClass           string   `json:"bandwidth_class,omitempty"`
	CloseWithRst             bool     `json:"close_with_rst"`
	CompletionRules          []string `json:"completionrules,omitempty"`
	ConnectTimeout           uint     `json:"connect_timeout"`
	Enabled                  bool     `json:"enabled"`
	FtpForceServerSecure     bool     `json:"ftp_force_server_secure"`
	GlbServices              []string `json:"glb_services,omitempty"`
	ListenOnAny              bool     `json:"listen_on_any"`
	ListenOnHosts            []string `json:"listen_on_hosts,omitempty"`
	ListenOnTrafficIps       []string `json:"listen_on_traffic_ips,omitempty"`
	MSS                      uint     `json:"mss"`
	Note                     string   `json:"note,omitempty"`
	Pool                     string   `json:"pool,omitempty"`
	Port                     uint     `json:"port"`
	ProtectionClass          string   `json:"protection_class,omitempty"`
	Protocol                 string   `json:"protocol,omitempty"`
	RequestRules             []string `json:"request_rules,omitempty"`
	ResponseRules            []string `json:"response_rules,omitempty"`
	SlmClass                 string   `json:"slm_class,omitempty"`
	SoNagle                  bool     `json:"so_nagle"`
	SslClientCertHeaders     string   `json:"ssl_client_cert_headers,omitempty"`
	SslDecrypt               bool     `json:"ssl_decrypt"`
	SslHonorFallbackScsv     string   `json:"ssl_honor_fallback_scsv,omitempty"`
	Transparent              bool     `json:"transparent"`
}

// Aptimizer : whether virtual server should aptimize web content
type Aptimizer struct {
	Enabled *bool              `json:"enabled,omitempty"`
	Profile []AptimizerProfile `json:"profile,omitempty"`
}

// AptimizerProfile : Aptimizer profiles and the application scopes that apply to them.
type AptimizerProfile struct {
	Name string   `json:"name,omitempty"`
	URLs []string `json:"urls,omitempty"`
}

// Connection : connection parameters
type Connection struct {
	Keepalive              *bool  `json:"keepalive,omitempty"`
	KeepaliveTimeout       *uint  `json:"keepalive_timeout,omitempty"`
	MaxClientBuffer        *uint  `json:"max_client_buffer,omitempty"`
	MaxServerBuffer        *uint  `json:"max_server_buffer,omitempty"`
	MaxTransactionDuration *uint  `json:"max_transaction_duration,omitempty"`
	ServerFirstBanner      string `json:"server_first_banner,omitempty"`
	Timeout                *uint  `json:"timeout,omitempty"`
}

// ConnectionErrors : error file params
type ConnectionErrors struct {
	ErrorFile string `json:"error_file,omitempty"`
}

// Cookie : how cookies are handled
type Cookie struct {
	Domain      string `json:"domain,omitempty"`
	NewDomain   string `json:"new_domain,omitempty"`
	PathRegex   string `json:"path_regex,omitempty"`
	PathReplace string `json:"path_replace,omitempty"`
	Secure      string `json:"secure,omitempty"`
}

// DNS configuration section
type DNS struct {
	EDNSClientSubnet *bool    `json:"edns_client_subnet,omitempty"`
	EdnsUdpsize      *uint    `json:"edns_udpsize,omitempty"`
	MaxUdpsize       *uint    `json:"max_udpsize,omitempty"`
	RrsetOrder       string   `json:"rrset_order,omitempty"`
	Verbose          *bool    `json:"verbose,omitempty"`
	Zones            []string `json:"zones,omitempty"`
}

// Ftp configuration section
type Ftp struct {
	DataSourcePort    *uint `json:"data_source_port,omitempty"`
	ForceClientSecure *bool `json:"force_client_secure,omitempty"`
	PortRangeHigh     *uint `json:"port_range_high,omitempty"`
	PortRangeLow      *uint `json:"port_range_low,omitempty"`
	SslData           *bool `json:"ssl_data,omitempty"`
}

// Gzip configuration section
type Gzip struct {
	CompressLevel *uint    `json:"compress_level,omitempty"`
	Enabled       *bool    `json:"enabled,omitempty"`
	EtagRewrite   string   `json:"etag_rewrite,omitempty"`
	IncludeMime   []string `json:"include_mime,omitempty"`
	MaxSize       *uint    `json:"max_size,omitempty"`
	MinSize       *uint    `json:"min_size,omitempty"`
	NoSize        *bool    `json:"no_size,omitempty"`
}

// HTTP configuration section
type HTTP struct {
	ChunkOverheadForwarding string `json:"chunk_overhead_forwarding,omitempty"`
	LocationRegex           string `json:"location_regex,omitempty"`
	LocationReplace         string `json:"location_replace,omitempty"`
	LocationRewrite         string `json:"location_rewrite,omitempty"`
	MIMEDefault             string `json:"mime_default,omitempty"`
	MIMEDetect              *bool  `json:"mime_detect,omitempty"`
}

// HTTP2 configuration section
type HTTP2 struct {
	ConnectTimeout         *uint    `json:"connect_timeout,omitempty"`
	DataFrameSize          *uint    `json:"data_frame_size,omitempty"`
	Enabled                *bool    `json:"enabled,omitempty"`
	HeaderTableSize        *uint    `json:"header_table_size,omitempty"`
	HeadersIndexBlacklist  []string `json:"headers_index_blacklist,omitempty"`
	HeadersIndexDefault    *bool    `json:"headers_index_default,omitempty"`
	HeadersIndexWhitelist  []string `json:"headers_index_whitelist,omitempty"`
	IdleTimeoutNoStreams   *uint    `json:"idle_timeout_no_streams,omitempty"`
	IdleTimeoutOpenStreams *uint    `json:"idle_timeout_open_streams,omitempty"`
	MaxConcurrentStreams   *uint    `json:"max_concurrent_streams,omitempty"`
	MaxFrameSize           *uint    `json:"max_frame_size,omitempty"`
	MaxHeaderPadding       *uint    `json:"max_header_padding,omitempty"`
	MergeCookieHeaders     *bool    `json:"merge_cookie_headers,omitempty"`
	StreamWindowSize       *uint    `json:"stream_window_size,omitempty"`
}

// KerberosProtocolTransition configuration section
type KerberosProtocolTransition struct {
	Enabled   *bool  `json:"enabled,omitempty"`
	Principal string `json:"principal,omitempty"`
	Target    string `json:"target,omitempty"`
}

// Log configuration section
type Log struct {
	AlwaysFlush               *bool  `json:"always_flush,omitempty"`
	ClientConnectionFailures  *bool  `json:"client_connection_failures,omitempty"`
	Enabled                   *bool  `json:"enabled,omitempty"`
	Filename                  string `json:"filename,omitempty"`
	Format                    string `json:"format,omitempty"`
	SaveAll                   *bool  `json:"save_all,omitempty"`
	ServerConnectionFailures  *bool  `json:"server_connection_failures,omitempty"`
	SessionPersistenceVerbose *bool  `json:"session_persistence_verbose,omitempty"`
	SSLFailures               *bool  `json:"ssl_failures,omitempty"`
}

// RecentConnections configuration section
type RecentConnections struct {
	Enabled *bool `json:"enabled,omitempty"`
	SaveAll *bool `json:"save_all,omitempty"`
}

// RequestTracing configuration section
type RequestTracing struct {
	Enabled *bool `json:"enabled,omitempty"`
	TraceIO *bool `json:"trace_io,omitempty"`
}

// RTSP configuration section
type RTSP struct {
	StreamingPortRangeHigh *uint `json:"streaming_port_range_high,omitempty"`
	StreamingPortRangeLow  *uint `json:"streaming_port_range_low,omitempty"`
	StreamingTimeout       *uint `json:"streaming_timeout,omitempty"`
}

// SIP configuration section
type SIP struct {
	DangerousRequests      string `json:"dangerous_requests,omitempty"`
	FollowRoute            *bool  `json:"follow_route,omitempty"`
	MaxConnectionMem       *uint  `json:"max_connection_mem,omitempty"`
	Mode                   string `json:"mode,omitempty"`
	RewriteURI             *bool  `json:"rewrite_uri,omitempty"`
	StreamingPortRangeHigh *uint  `json:"streaming_port_range_high,omitempty"`
	StreamingPortRangeLow  *uint  `json:"streaming_port_range_low,omitempty"`
	StreamingTimeout       *uint  `json:"streaming_timeout,omitempty"`
	TimeoutMessages        *bool  `json:"timeout_messages,omitempty"`
	TransactionTimeout     *uint  `json:"transaction_timeout,omitempty"`
}

// SMTP configuration section
type SMTP struct {
	ExpectSTARTTLS *bool `json:"expect_starttls,omitempty"`
}

// Ssl configuration section
type Ssl struct {
	AddHTTPHeaders            *bool        `json:"add_http_headers,omitempty"`
	ClientCertCAS             []string     `json:"client_cert_cas,omitempty"`
	EllipticCurves            []string     `json:"elliptic_curves,omitempty"`
	IssuedCertsNeverExpire    []string     `json:"issued_certs_never_expire,omitempty"`
	OCSPEnable                *bool        `json:"ocsp_enable,omitempty"`
	OCSPIssuers               []OCSPIssuer `json:"ocsp_issuers,omitempty"`
	OCSPMaxResponseAge        *uint        `json:"ocsp_max_response_age,omitempty"`
	OCSPStapling              *bool        `json:"ocsp_stapling,omitempty"`
	OCSPTimeTolerance         *uint        `json:"ocsp_time_tolerance,omitempty"`
	OCSPTimeout               *uint        `json:"ocsp_timeout,omitempty"`
	PreferSSLv3               *bool        `json:"prefer_sslv3,omitempty"`
	RequestClientCert         string       `json:"request_client_cert,omitempty"`
	SendCloseAlerts           *bool        `json:"send_close_alerts,omitempty"`
	ServerCertAltCertificates []string     `json:"server_cert_alt_certificates,omitempty"`
	ServerCertDefault         string       `json:"server_cert_default,omitempty"`
	ServerCertHostMap         []CertItem   `json:"server_cert_host_mapping,omitempty"`
	SignatureAlgorithms       string       `json:"signature_algorithms,omitempty"`
	SSLCiphers                string       `json:"ssl_ciphers,omitempty"`
	SslSupportSsl2            string       `json:"ssl_support_ssl2,omitempty"`
	SslSupportSsl3            string       `json:"ssl_support_ssl3,omitempty"`
	SslSupportTLS1            string       `json:"ssl_support_tls1,omitempty"`
	SslSupportTLS1_1          string       `json:"ssl_support_tls1_1,omitempty"`
	SslSupportTLS1_2          string       `json:"ssl_support_tls1_2,omitempty"`
	TrustMagic                *bool        `json:"trust_magic,omitempty"`
}

// SysLog configuration section
type SysLog struct {
	Enabled     *bool  `json:"enabled,omitempty"`
	Format      string `json:"format,omitempty"`
	IPEndpoint  string `json:"ip_end_point,omitempty"`
	MsgLenLimit *uint  `json:"msg_len_limit,omitempty"`
}

// TCP configuration section
type TCP struct {
	ProxyClose *bool `json:"proxy_close,omitempty"`
}

//UDP configuration section
type UDP struct {
	EndPointPersistence       *bool `json:"end_point_persistence,omitempty"`
	PortSMP                   *bool `json:"port_smp,omitempty"`
	ResponseDatagramsExpected *int  `json:"response_datagrams_expected,omitempty"`
	Timeout                   *uint `json:"timeout,omitempty"`
}

// WebCache configuration section
type WebCache struct {
	ControlOut    string `json:"control_out,omitempty"`
	Enabled       *bool  `json:"enabled,omitempty"`
	ErrorPageTime *uint  `json:"error_page_time,omitempty"`
	MaxTime       *uint  `json:"max_time,omitempty"`
	RefreshTime   *uint  `json:"refresh_time,omitempty"`
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
	Children []ChildVirtualServer `json:"children,omitempty"`
}

// ChildVirtualServer : monitored node structure
type ChildVirtualServer struct {
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
}
