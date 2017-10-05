package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/virtualserver"
	"github.com/sky-uk/go-rest-api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"net/http"
	"regexp"
)

func resourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualServerCreate,
		Read:   resourceVirtualServerRead,
		Update: resourceVirtualServerUpdate,
		Delete: resourceVirtualServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the virtual server",
				Required:    true,
				ForceNew:    true,
			},
			"add_cluster_ip": {
				Type:        schema.TypeBool,
				Description: "Whether or not the virtual server should add an 'X-Cluster-Client-Ip' header to the request that contains the remote client's IP address.",
				Optional:    true,
				Default:     false,
			},
			"add_x_forwarded_for": {
				Type:        schema.TypeBool,
				Description: "Whether or not the virtual server should append the remote client's IP address to the 'X-Forwarded-For header'. If the header does not exist, it will be added.",
				Optional:    true,
				Default:     false,
			},
			"add_x_forwarded_proto": {
				Type:        schema.TypeBool,
				Description: "Whether or not the virtual server should add an 'X-Forwarded-Proto' header to the request that contains the original protocol used by the client to connect to the traffic manager.",
				Optional:    true,
				Default:     false,
			},
			"autodetect_upgrade_headers": {
				Type:        schema.TypeBool,
				Description: "Whether the traffic manager should check for HTTP responses that confirm an HTTP connection is transitioning to the WebSockets protocol. ",
				Optional:    true,
				Default:     false,
			},
			"bandwidth_class": {
				Type:        schema.TypeString,
				Description: "The bandwidth management class that this server should use, if any.",
				Optional:    true,
			},
			"close_with_rst": {
				Type:        schema.TypeBool,
				Description: "Whether or not connections from clients should be closed with a RST packet, rather than a FIN packet.",
				Optional:    true,
				Default:     false,
			},
			"completionrules": {
				Type:        schema.TypeList,
				Description: "Rules that are run at the end of a transaction, in order, comma separated.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"connect_timeout": {
				Type:         schema.TypeInt,
				Description:  "The time, in seconds, to wait for data from a new connection. If no data isreceived within this time, the connection will be closed. A value of 0 (zero) will disable the timeout.",
				Optional:     true,
				Default:      10,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the virtual server is enabled.",
				Optional:    true,
				Default:     false,
			},
			"ftp_force_server_secure": {
				Type:        schema.TypeBool,
				Description: "Whether or not the virtual server should require that incoming FTP data connections from the nodes originate from the same IP address as the node",
				Optional:    true,
				Default:     false,
			},
			"glb_services": {
				Type:        schema.TypeList,
				Description: "The associated GLB services for this DNS virtual server.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"listen_on_any": {
				Type:        schema.TypeBool,
				Description: "Whether to listen on all IP addresses",
				Optional:    true,
				Default:     false,
			},
			"listen_on_hosts": {
				Type:        schema.TypeList,
				Description: "Hostnames and IP addresses to listen on",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"listen_on_traffic_ips": {
				Type:        schema.TypeList,
				Description: "List of traffic IPs to listen on",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"note": {
				Type:        schema.TypeString,
				Description: " A description for the virtual server.",
				Optional:    true,
			},
			"pool": {
				Type:        schema.TypeString,
				Description: "The default pool to use for traffic.",
				Required:    true,
			},
			"port": {
				Type:         schema.TypeInt,
				Description:  "The port on which to listen for incoming connections",
				Required:     true,
				ValidateFunc: util.ValidatePortNumber,
			},
			"protection_class": {
				Type:        schema.TypeString,
				Description: "The service protection class that should be used to protect this server, if any.",
				Optional:    true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "The protocol that the virtual server is using.",
				Optional:     true,
				Default:      "http",
				ValidateFunc: validateVirtualServerProtocol,
			},
			"request_rules": {
				Type:        schema.TypeList,
				Description: "Rules to be applied to incoming requests, in order, comma separated.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"response_rules": {
				Type:        schema.TypeList,
				Description: "Rules to be applied to responses, in order, comma separated.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"slm_class": {
				Type:         schema.TypeString,
				Description:  "The service level monitoring class that this server should use, if any.",
				Optional:     true,
				Default:      "http",
				ValidateFunc: validateVirtualServerProtocol,
			},
			"so_nagle": {
				Type:        schema.TypeBool,
				Description: "Whether or not Nagle's algorithm should be used for TCP connections.",
				Optional:    true,
				Default:     false,
			},
			"ssl_client_cert_headers": {
				Type:         schema.TypeString,
				Description:  "What HTTP headers the virtual server should add to each request to show the data in the client certificate.",
				Optional:     true,
				Default:      "none",
				ValidateFunc: validateSSLClientCertHeaders,
			},
			"ssl_decrypt": {
				Type:        schema.TypeBool,
				Description: "Whether or not the virtual server should decrypt incoming SSL traffic.",
				Optional:    true,
				Default:     false,
			},
			"ssl_honor_fallback_scsv": {
				Type:         schema.TypeString,
				Description:  " Whether or not the Fallback SCSV sent by TLS clients is honored by this virtual server. ",
				Optional:     true,
				Default:      "use_default",
				ValidateFunc: validateSSLHonorFallbackSCSV,
			},
			"transparent": {
				Type:        schema.TypeBool,
				Description: "Whether or not bound sockets should be configured for transparent proxying",
				Optional:    true,
				Default:     false,
			},

			"error_file": {
				Type:        schema.TypeString,
				Description: "The error message to be sent to the client when the traffic manager detects an internal or backend error for the virtual server.",
				Optional:    true,
				Default:     "Default",
			},

			"expect_starttls": {
				Type:        schema.TypeBool,
				Description: "Whether or not the traffic manager should expect the connection to start off in plain text and then upgrade to SSL using STARTTLS when handling SMTP traffic",
				Optional:    true,
				Default:     false,
			},

			"proxy_close": {
				Type:        schema.TypeBool,
				Description: "If set to Yes the traffic manager will send the client FIN to the back-end server and wait for a server response instead of closing the connection immediately. ",
				Optional:    true,
				Default:     false,
			},

			"aptimizer": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: " Whether the virtual server should optimize web content",
							Optional:    true,
							Default:     false,
						},
						"profile": {
							Type:        schema.TypeList,
							Description: "A table of Aptimizer profiles and the application scopes that apply to them.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "The name of an Aptimizer acceleration profile.",
										Required:    true,
									},
									"urls": {
										Type:        schema.TypeList,
										Description: "The application scopes which apply to the acceleration profile.",
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},

			"vs_connection": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keepalive": {
							Type:        schema.TypeBool,
							Description: "Whether or not the virtual server should use keepalive connections with the remote clients.",
							Optional:    true,
							Default:     false,
						},
						"keepalive_timeout": {
							Type:         schema.TypeInt,
							Description:  "The length of time that the virtual server should keep an idle keepalive connection before discarding it. A value of 0 (zero) will mean that the keepalives are never closed by the traffic manager.",
							Optional:     true,
							Default:      10,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"max_client_buffer": {
							Type:         schema.TypeInt,
							Description:  "The amount of memory, in bytes, that the virtual server should use to store data sent by the client.",
							Optional:     true,
							Default:      65536,
							ValidateFunc: validateMaxBuffer,
						},
						"max_server_buffer": {
							Type:         schema.TypeInt,
							Description:  "The amount of memory, in bytes, that the virtual server should use to store data returned by the server.",
							Optional:     true,
							Default:      65536,
							ValidateFunc: validateMaxBuffer,
						},
						"max_transaction_duration": {
							Type:         schema.TypeInt,
							Description:  " The total amount of time a transaction can take, counted from the first byte being received until the transaction is complete. ",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},

						"server_first_banner": {
							Type:        schema.TypeString,
							Description: "If specified, the traffic manager will use the value as the banner to send for server-first protocols such as POP, SMTP and IMAP. ",
							Optional:    true,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Description:  "A connection should be closed if no additional data has been received for this period of time. A value of 0 (zero) will disable this timeout.",
							Optional:     true,
							Default:      300,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
					},
				},
			},

			"cookie": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:         schema.TypeString,
							Description:  "The way in which the traffic manager should rewrite the domain portion of any cookies set by a back-end web server.",
							Optional:     true,
							Default:      "no_rewrite",
							ValidateFunc: validateCookieDomain,
						},
						"new_domain": {
							Type:        schema.TypeString,
							Description: "The domain to use when rewriting a cookie's domain to a named value.",
							Optional:    true,
						},
						"path_regex": {
							Type:        schema.TypeString,
							Description: "If you wish to rewrite the path portion of any cookies set by a back-end web server, provide a regular expression to match the path.",
							Optional:    true,
						},
						"path_replace": {
							Type:        schema.TypeString,
							Description: "If cookie path regular expression matches, it will be replaced by this substitution.",
							Optional:    true,
						},
						"secure": {
							Type:         schema.TypeString,
							Description:  "Whether or not the traffic manager should modify the 'secure' tag of anycookies set by a back-end web server.",
							Optional:     true,
							Default:      "no_modify",
							ValidateFunc: validateCookieSecure,
						},
					},
				},
			},
			"dns": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"edns_client_subnet": {
							Type:        schema.TypeBool,
							Description: "Enable/Disable use of EDNS client subnet option.",
							Optional:    true,
							Default:     false,
						},
						"edns_udpsize": {
							Type:         schema.TypeInt,
							Description:  "EDNS UDP size advertised in responses",
							Optional:     true,
							Default:      4096,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"max_udpsize": {
							Type:         schema.TypeInt,
							Description:  "Maximum UDP answer size",
							Optional:     true,
							Default:      4096,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"rrset_order": {
							Type:         schema.TypeString,
							Description:  "Response record ordering.",
							Optional:     true,
							Default:      "fixed",
							ValidateFunc: validateDNSrrsetOrder,
						},
						"verbose": {
							Type:        schema.TypeBool,
							Description: "Whether or not the DNS Server should emit verbose logging.",
							Optional:    true,
							Default:     false,
						},
						"zones": {
							Type:        schema.TypeList,
							Description: "The DNS zones.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
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
						"data_source_port": {
							Type:         schema.TypeInt,
							Description:  "The source port to be used for active-mode FTP data connections. If 0, a random high port will be used.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"force_client_secure": {
							Type:        schema.TypeBool,
							Description: "Whether or not the virtual server should require that incoming FTP dataconnections from the client originate from the same IP address as the corresponding client control connection.",
							Optional:    true,
							Default:     false,
						},
						"port_range_high": {
							Type:         schema.TypeInt,
							Description:  "If non-zero, then this controls the upper bound of the port range to use for FTP data connections.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"port_range_low": {
							Type:         schema.TypeInt,
							Description:  "If non-zero, then this controls the lower bound of the port range to use for FTP data connections.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"ssl_data": {
							Type:        schema.TypeBool,
							Description: "Use SSL on the data connection as well as the control connection (if not enabled it is left to the client and server to negotiate this).",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			"gzip": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compress_level": {
							Type:         schema.TypeInt,
							Description:  "The source port to be used for active-mode FTP data connections. If 0, a random high port will be used.",
							Optional:     true,
							Default:      1,
							ValidateFunc: validateCompressLevel,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Compress web pages sent back by the server",
							Optional:    true,
							Default:     false,
						},
						"etag_rewrite": {
							Type:         schema.TypeString,
							Description:  "How the ETag header should be manipulated when compressing content.",
							Optional:     true,
							Default:      "wrap",
							ValidateFunc: validateETagRewrite,
						},
						"include_mime": {
							Type:        schema.TypeList,
							Description: "MIME types to compress. Complete MIME types can be used, or a type can end in a '*' to match multiple types.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"max_size": {
							Type:         schema.TypeInt,
							Description:  "Maximum document size to compress (0 means unlimited).",
							Optional:     true,
							Default:      10000000,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"min_size": {
							Type:         schema.TypeInt,
							Description:  "Minimum document size to compress.",
							Optional:     true,
							Default:      1000,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"no_size": {
							Type:        schema.TypeBool,
							Description: "Compress documents with no given size.",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			"http": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chunk_overhead_forwarding": {
							Type:         schema.TypeString,
							Description:  "Handling of HTTP chunk overhead.",
							Optional:     true,
							Default:      "lazy",
							ValidateFunc: validateChunkOverheadForwarding,
						},
						"location_regex": {
							Type:        schema.TypeString,
							Description: "If the 'Location' header matches this regular expression, rewrite the header using the 'location_replace' pattern.",
							Optional:    true,
						},
						"location_replace": {
							Type:        schema.TypeString,
							Description: "If the 'Location' header matches the 'location_regex' regular expression, rewrite the header with this pattern",
							Optional:    true,
						},
						"location_rewrite": {
							Type:         schema.TypeString,
							Description:  "If the 'Location' header matches the 'location_regex' regular expression, rewrite the header with this pattern",
							Optional:     true,
							Default:      "if_host_matches",
							ValidateFunc: validateLocationRewrite,
						},
						"mime_default": {
							Type:        schema.TypeString,
							Description: "Auto-correct MIME types if the server sends the 'default' MIME type for files.",
							Optional:    true,
							Default:     "text/plain",
						},
						"mime_detect": {
							Type:        schema.TypeBool,
							Description: "Auto-detect MIME types if the server does not provide them.",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			"http2": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_timeout": {
							Type:         schema.TypeInt,
							Description:  "The time, in seconds, to wait for a request on a new HTTP/2 connection.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"data_frame_size": {
							Type:         schema.TypeInt,
							Description:  "This setting controls the preferred frame size used when sending body data to the client.",
							Optional:     true,
							Default:      4096,
							ValidateFunc: validateDataFrameSize,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Description: "This setting allows the HTTP/2 protocol to be used by a HTTP virtual server.",
							Optional:    true,
							Default:     false,
						},
						"header_table_size": {
							Type:         schema.TypeInt,
							Description:  "This setting controls the amount of memory allowed for header compression on each HTTP/2 connection.",
							Optional:     true,
							Default:      4096,
							ValidateFunc: validateHeaderTableSize,
						},
						"headers_index_blacklist": {
							Type:        schema.TypeList,
							Description: "A list of header names that should never be compressed using indexing.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"headers_index_default": {
							Type:        schema.TypeBool,
							Description: "The HTTP/2 HPACK compression scheme allows for HTTP headers to be compressed using indexing. If this is true only hraders included in headers_index_blacklist are marked as 'never index', if false all headers are marked as never index unless in whitelist",
							Optional:    true,
							Default:     false,
						},
						"headers_index_whitelist": {
							Type:        schema.TypeList,
							Description: "A list of header names that can be compressed using indexing when the value of headers_index_default is set to False.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"idle_timeout_no_streams": {
							Type:         schema.TypeInt,
							Description:  "The time, in seconds, to wait for a new HTTP/2 request on a previously used HTTP/2 connection that has no open HTTP/2 streams. A value of 0 disables the timeout.",
							Optional:     true,
							Default:      120,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"idle_timeout_open_streams": {
							Type:         schema.TypeInt,
							Description:  "The time, in seconds, to wait for data on an idle HTTP/2 connection, which has open streams, when no data has been sent recently. A value of 0 disables the timeout.",
							Optional:     true,
							Default:      600,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"max_concurrent_streams": {
							Type:         schema.TypeInt,
							Description:  "This setting controls the number of streams a client is permitted to open concurrently on a single connection.",
							Optional:     true,
							Default:      200,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"max_frame_size": {
							Type:         schema.TypeInt,
							Description:  "This setting controls the maximum HTTP/2 frame size clients are permitted to send to the traffic manager.",
							Optional:     true,
							Default:      16384,
							ValidateFunc: validateMaxFrameSize,
						},
						"max_header_padding": {
							Type:         schema.TypeInt,
							Description:  "The maximum size, in bytes, of the random-length padding to add to HTTP/2 header frames. The padding, a random number of zero bytes up to the maximum specified.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"merge_cookie_headers": {
							Type:        schema.TypeBool,
							Description: "Whether Cookie headers received from an HTTP/2 client should be merged into a single Cookie header using RFC6265 rules before forwarding to anHTTP/1.1 server.",
							Optional:    true,
							Default:     false,
						},
						"stream_window_size": {
							Type:         schema.TypeInt,
							Description:  "This setting controls the flow control window for each HTTP/2 stream.",
							Optional:     true,
							Default:      65535,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
					},
				},
			},

			"kerberos_protocol_transition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether or not the virtual server should use Kerberos Protocol Transition.",
							Optional:    true,
							Default:     false,
						},
						"principal": {
							Type:        schema.TypeString,
							Description: "The Kerberos principal this virtual server should use to perform Kerberos Protocol Transition.",
							Optional:    true,
						},
						"target": {
							Type:        schema.TypeString,
							Description: "The Kerberos principal name of the service this virtual server targets.",
							Optional:    true,
						},
					},
				},
			},

			"log": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"always_flush": {
							Type:        schema.TypeBool,
							Description: "Write log data to disk immediately, rather than buffering data.",
							Optional:    true,
							Default:     false,
						},
						"client_connection_failures": {
							Type:        schema.TypeBool,
							Description: "Should the virtual server log failures occurring on connections to clients.",
							Optional:    true,
							Default:     false,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether or not to log connections to the virtual server to a disk on the file system.",
							Optional:    true,
							Default:     false,
						},
						"filename": {
							Type:        schema.TypeString,
							Description: "The name of the file in which to store the request logs. The filename can contain macros which will be expanded by the traffic manager to generate the full filename.",
							Optional:    true,
							Default:     "%zeushome%/zxtm/log/%v.log",
						},
						"format": {
							Type:        schema.TypeString,
							Description: "The log file format. This specifies the line of text that will be written to the log file when a connection to the traffic manager is completed. Many parameters from the connection can be recorded using macros.",
							Optional:    true,
							Default:     "%h %l %u %t \"%r\" %s %b \"%{Referer}i\"\"%{User-agent}i\"",
						},
						"save_all": {
							Type:        schema.TypeBool,
							Description: "Whether to log all connections by default, or log no connections by default.",
							Optional:    true,
							Default:     false,
						},
						"server_connection_failures": {
							Type:        schema.TypeBool,
							Description: "Should the virtual server log failures occurring on connections to nodes.",
							Optional:    true,
							Default:     false,
						},
						"session_persistence_verbose": {
							Type:        schema.TypeBool,
							Description: "Should the virtual server log session persistence events.",
							Optional:    true,
							Default:     false,
						},
						"ssl_failures": {
							Type:        schema.TypeBool,
							Description: "Should the virtual server log failures occurring on SSL secure negotiation.",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			"recent_connections": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether or not connections handled by this virtual server should be shown on the Activity Connections page.",
							Optional:    true,
							Default:     false,
						},
						"save_all": {
							Type:        schema.TypeBool,
							Description: "Whether or not all connections handled by this virtual server should be shown on the Connections page.",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			"request_tracing": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Record a trace of major connection processing events for each request and response.",
							Optional:    true,
							Default:     false,
						},
						"trace_io": {
							Type:        schema.TypeBool,
							Description: " Include details of individual I/O events in request and response traces. Requires request tracing to be enabled.",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			"rtsp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"streaming_port_range_high": {
							Type:         schema.TypeInt,
							Description:  "If non-zero this controls the upper bound of the port range to use for streaming data connections.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"streaming_port_range_low": {
							Type:         schema.TypeInt,
							Description:  "If non-zero this controls the lower bound of the port range to use for streaming data connections.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"streaming_timeout": {
							Type:         schema.TypeInt,
							Description:  "If non-zero data-streams associated with RTSP connections will timeout if no data is transmitted for this many seconds",
							Optional:     true,
							Default:      30,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
					},
				},
			},

			"sip": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dangerous_requests": {
							Type:         schema.TypeString,
							Description:  " The action to take when a SIP request with body data arrives that should be routed to an external IP.",
							Optional:     true,
							Default:      "node",
							ValidateFunc: validateSIPDangerousRequestsAction,
						},
						"follow_route": {
							Type:        schema.TypeBool,
							Description: "Should the virtual server follow routing information contained in SIP requests.",
							Optional:    true,
							Default:     false,
						},
						"max_connection_mem": {
							Type:         schema.TypeInt,
							Description:  "this setting limits the amount of memory each SIP client can use. When the limit is reached new requests will be sent a 413 response. If the value is set to 0 (zero) the memory limit is disabled.",
							Optional:     true,
							Default:      65536,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"mode": {
							Type:         schema.TypeString,
							Description:  " The action to take when a SIP request with body data arrives that should be routed to an external IP.",
							Optional:     true,
							Default:      "sip_gateway",
							ValidateFunc: validateSIPMode,
						},
						"rewrite_uri": {
							Type:        schema.TypeBool,
							Description: " Replace the Request-URI of SIP requests with the address of the selected backend node.",
							Optional:    true,
							Default:     false,
						},
						"streaming_port_range_high": {
							Type:         schema.TypeInt,
							Description:  "If non-zero this controls the upper bound of the port range to use for streaming data connections.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"streaming_port_range_low": {
							Type:         schema.TypeInt,
							Description:  "If non-zero this controls the lower bound of the port range to use for streaming data connections.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"streaming_timeout": {
							Type:         schema.TypeInt,
							Description:  "If non-zero a UDP stream will timeout when no data has been seen within this time.",
							Optional:     true,
							Default:      60,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"timeout_messages": {
							Type:        schema.TypeBool,
							Description: "When timing out a SIP transaction, send a 'timed out' response to the client and, in the case of an INVITE transaction, a CANCEL request to the server.",
							Optional:    true,
							Default:     false,
						},
						"transaction_timeout": {
							Type:         schema.TypeInt,
							Description:  "The virtual server should discard a SIP transaction when no further messages have been seen within this time.",
							Optional:     true,
							Default:      30,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
					},
				},
			},

			"ssl": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add_http_headers": {
							Type:        schema.TypeBool,
							Description: "Whether or not the virtual server should add HTTP headers to each request to show the SSL connection parameters.",
							Optional:    true,
							Default:     false,
						},
						"client_cert_cas": {
							Type:        schema.TypeList,
							Description: "The certificate authorities that this virtual server should trust to validate client certificates. If no certificate authorities are selected, and client certificates are requested, then all client certificates will be accepted.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"elliptic_curves": {
							Type:        schema.TypeList,
							Description: "The SSL elliptic curve preference list for SSL connections to this virtual server using TLS version 1.0 or higher. Leaving this empty will make the virtual server use the globally configured curve preference list.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"issued_certs_never_expire": {
							Type:        schema.TypeList,
							Description: "When the virtual server verifies certificates signed by these certificate authorities, it doesn't check the 'not after' date, i.e., they are considered valid even after their expiration date has passed  (but not if they have been revoked)",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ocsp_enable": {
							Type:        schema.TypeBool,
							Description: "Whether or not the traffic manager should use OCSP to check the revocation status of client certificates.",
							Optional:    true,
							Default:     false,
						},
						"ocsp_issuers": {
							Type:        schema.TypeList,
							Description: "A table of certificate issuer specific OCSP settings",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issuer": {
										Type:        schema.TypeString,
										Description: "The name of an issuer",
										Optional:    true,
										Default:     "DEFAULT",
									},
									"aia": {
										Type:        schema.TypeBool,
										Description: "Whether the traffic manager should use AIA information",
										Optional:    true,
										Default:     false,
									},
									"nonce": {
										Type:         schema.TypeString,
										Description:  "How to use the OCSP nonce extension, which protects against OCSP replay attacks",
										Optional:     true,
										ValidateFunc: validateVirtualServerOCSPNonce,
									},
									"required": {
										Type:         schema.TypeString,
										Description:  "Whether we should do an OCSP check for this issuer",
										Optional:     true,
										ValidateFunc: validateVirtualServerOCSPRequired,
									},
									"responder_cert": {
										Type:        schema.TypeString,
										Description: "The expected responder certificate",
										Optional:    true,
									},
									"signer": {
										Type:        schema.TypeString,
										Description: "The certificate with which to sign the request",
										Optional:    true,
									},
									"url": {
										Type:        schema.TypeString,
										Description: "Which OCSP responders this virtual server should use to verify client certificates",
										Optional:    true,
									},
								},
							},
						},
						"ocsp_max_response_age": {
							Type:         schema.TypeInt,
							Description:  "The number of seconds for which an OCSP response is considered valid if it has not yet exceeded the time specified in the 'nextUpdate' field. If set to 0 (zero) then OCSP responses are considered valid until the time specified in their 'nextUpdate' field.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"ocsp_stapling": {
							Type:        schema.TypeBool,
							Description: "If OCSP URIs are present in certificates used by this virtual server, allow the traffic manager to provide OCSP responses for these certificates as part of the handshake.",
							Optional:    true,
							Default:     false,
						},
						"ocsp_time_tolerance": {
							Type:         schema.TypeInt,
							Description:  "The number of seconds outside the permitted range for which the 'thisUpdate' and 'nextUpdate' fields of an OCSP response are still considered valid.",
							Optional:     true,
							Default:      30,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"ocsp_timeout": {
							Type:         schema.TypeInt,
							Description:  "The number of seconds after which OCSP requests will be timed out.",
							Optional:     true,
							Default:      10,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"prefer_sslv3": {
							Type:        schema.TypeBool,
							Description: "Deprecated. Formerly allowed a preference for SSLv3 for performance reasons.",
							Optional:    true,
							Default:     false,
						},
						"request_client_cert": {
							Type:         schema.TypeString,
							Description:  "Whether or not the virtual server should request an identifying SSL certificate from each client.",
							Optional:     true,
							Default:      "dont_request",
							ValidateFunc: validateSSLRequestClientCert,
						},
						"send_close_alerts": {
							Type:        schema.TypeBool,
							Description: "Whether or not to send an SSL/TLS 'close alert' when the traffic manager is initiating an SSL socket disconnection.",
							Optional:    true,
							Default:     false,
						},
						"server_cert_alt_certificates": {
							Type:        schema.TypeList,
							Description: "The SSL certificates and corresponding private keys.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"server_cert_default": {
							Type:        schema.TypeString,
							Description: "The default SSL certificate to use for this virtual server.",
							Optional:    true,
						},
						"ssl_server_cert_host_mapping": {
							Type:        schema.TypeList,
							Description: "Host specific SSL server certificate mappings",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Description: "Which host the SSL certificate refers to",
										Optional:    true,
									},
									"certificate": {
										Type:        schema.TypeString,
										Description: "The SSL server certificate for a particular destination",
										Optional:    true,
									},
									"alt_certificates": {
										Type:        schema.TypeList,
										Description: "SSL server certificates for a particular destination IP",
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"signature_algorithms": {
							Type:        schema.TypeString,
							Description: "The SSL signature algorithms preference list for SSL connections to this virtual server using TLS version 1.2 or higher.",
							Optional:    true,
						},
						"ssl_ciphers": {
							Type:        schema.TypeString,
							Description: "The SSL/TLS ciphers to allow for connections to this virtual server. ",
							Optional:    true,
						},
						"ssl_support_ssl2": {
							Type:         schema.TypeString,
							Description:  "Whether or not SSLv2 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validateVirtualServerUseSSLSupport,
						},
						"ssl_support_ssl3": {
							Type:         schema.TypeString,
							Description:  "Whether or not SSLv3 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validateVirtualServerUseSSLSupport,
						},
						"ssl_support_tls1": {
							Type:         schema.TypeString,
							Description:  "Whether or not TLSv1.0 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validateVirtualServerUseSSLSupport,
						},
						"ssl_support_tls1_1": {
							Type:         schema.TypeString,
							Description:  "Whether or not TLSv1.1 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validateVirtualServerUseSSLSupport,
						},
						"ssl_support_tls1_2": {
							Type:         schema.TypeString,
							Description:  "Whether or not TLSv1.2 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validateVirtualServerUseSSLSupport,
						},
						"trust_magic": {
							Type:        schema.TypeBool,
							Description: "If the traffic manager is receiving traffic sent from another traffic manager, then enabling this option will allow it to decode extra information on the true origin of the SSL connection.",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
			"syslog": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether or not to log connections to the virtual server to a remote syslog host.",
							Optional:    true,
							Default:     false,
						},
						"format": {
							Type:        schema.TypeString,
							Description: "The log format for the remote syslog.",
							Optional:    true,
							Default:     "%h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\"",
						},
						"ip_end_point": {
							Type:        schema.TypeString,
							Description: "The remote host and port (default is 514) to send request log lines to.",
							Optional:    true,
						},
						"msg_len_limit": {
							Type:         schema.TypeInt,
							Description:  "Maximum length in bytes of a message sent to the remote syslog. Messages longer than this will be truncated before they are sent.",
							Optional:     true,
							Default:      1024,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
					},
				},
			},

			"udp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_point_persistence": {
							Type:        schema.TypeBool,
							Description: "Whether or not UDP datagrams from the same IP and port are sent to the same node in the pool if there's an existing UDP transaction.",
							Optional:    true,
							Default:     false,
						},
						"port_smp": {
							Type:        schema.TypeBool,
							Description: "Whether or not UDP datagrams should be distributed across all traffic manager processes.",
							Optional:    true,
							Default:     false,
						},

						"response_datagrams_expected": {
							Type:        schema.TypeInt,
							Description: "The virtual server should discard any UDP connection and reclaim resourceswhen the node has responded with this number of datagrams. If set to -1, the connection will not be discarded until the timeout is reached.",
							Optional:    true,
							Default:     1,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Description:  "The virtual server should discard any UDP connection and reclaim resources when no further UDP traffic has been seen within this time.",
							Optional:     true,
							Default:      7,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
					},
				},
			},

			"web_cache": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_out": {
							Type:        schema.TypeString,
							Description: "The 'Cache-Control' header to add to every cached HTTP response, no-cache or max-age=600 for example.",
							Optional:    true,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Description: "If set to true the traffic manager will attempt to cache web server responses.",
							Optional:    true,
							Default:     false,
						},
						"error_page_time": {
							Type:         schema.TypeInt,
							Description:  "Time period to cache error pages for.",
							Optional:     true,
							Default:      30,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"max_time": {
							Type:         schema.TypeInt,
							Description:  "Maximum time period to cache web pages for.",
							Optional:     true,
							Default:      600,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
						"refresh_time": {
							Type:         schema.TypeInt,
							Description:  "If a cached page is about to expire within this time, the traffic manager will start to forward some new requests on to the web servers. Setting this value to 0 will stop the traffic manager updating the cache before it expires.",
							Optional:     true,
							Default:      2,
							ValidateFunc: util.ValidateUnsignedInteger,
						},
					},
				},
			},
		},
	}

}

func validateVirtualServerOCSPRequired(v interface{}, k string) (ws []string, errors []error) {
	ocspRequired := v.(string)
	ocspRequiredOptions := regexp.MustCompile(`^(none|optional|strict)$`)
	if !ocspRequiredOptions.MatchString(ocspRequired) {
		errors = append(errors, fmt.Errorf("%q must be one of none, optional, strict", k))
	}
	return
}

func validateVirtualServerOCSPNonce(v interface{}, k string) (ws []string, errors []error) {
	nonce := v.(string)
	nonceOptions := regexp.MustCompile(`^(off|on|strict)$`)
	if !nonceOptions.MatchString(nonce) {
		errors = append(errors, fmt.Errorf("%q must be one of off, on or strict", k))
	}
	return
}

func validateSSLClientCertHeaders(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"all",
		"none",
		"simple":
		return
	}
	errors = append(errors, fmt.Errorf("SSL Client Cert Header must be one of all, none or simple"))
	return
}

func validateSSLHonorFallbackSCSV(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"disabled",
		"enabled",
		"use_default":
		return
	}
	errors = append(errors, fmt.Errorf("SSL Honor Fallback SCSV must be one of disabled, enabled or use_default"))
	return
}

func validateCookieDomain(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"no_rewrite",
		"set_to_named",
		"set_to_request":
		return
	}
	errors = append(errors, fmt.Errorf("Cookie Domain must be one of no_rewrite, set_to_named or set_to_request"))
	return
}

func validateCookieSecure(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"no_modify",
		"set_secure",
		"unset_secure":
		return
	}
	errors = append(errors, fmt.Errorf("Cookie Secure must be one of no_modify, set_secure or unset_secure"))
	return
}

func validateDNSrrsetOrder(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"cyclic",
		"fixed":
		return
	}
	errors = append(errors, fmt.Errorf("DNS RRSET Order must be one of cyclic or fixed"))
	return
}

func validateCompressLevel(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 1 || v.(int) > 9 {
		errors = append(errors, fmt.Errorf("Compression level must be a value within 1-9"))
	}
	return
}

func validateDataFrameSize(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 100 || v.(int) > 16777206 {
		errors = append(errors, fmt.Errorf("data_frame_size must be a value within 100-16777206"))
	}
	return
}

func validateMaxFrameSize(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 16384 || v.(int) > 16777215 {
		errors = append(errors, fmt.Errorf("max_frame_size must be a value within 16384-16777215"))
	}
	return
}

func validateETagRewrite(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"delete",
		"ignore",
		"weaken",
		"wrap":
		return
	}
	errors = append(errors, fmt.Errorf("ETag Rewrite must be one of wrap, delete, ignore, weaken or wrap"))
	return
}

func validateMaxBuffer(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 1024 || v.(int) > 16777216 {
		errors = append(errors, fmt.Errorf("%q must be within 1024-16777216", k))
	}
	return
}

func validateHeaderTableSize(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 4096 || v.(int) > 1048576 {
		errors = append(errors, fmt.Errorf("header_table_size must be a value within 4096-1048576"))
	}
	return
}

func validateChunkOverheadForwarding(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"lazy",
		"eager":
		return
	}
	errors = append(errors, fmt.Errorf("Chunk Overhead Forwarding must be one of lazy or eager"))
	return
}

func validateLocationRewrite(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"always",
		"if_host_matches",
		"never":
		return
	}
	errors = append(errors, fmt.Errorf("Location Rewrite must be one of always, if_host_matches or never"))
	return
}

func validateSIPDangerousRequestsAction(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"forbid",
		"forward",
		"node":
		return
	}
	errors = append(errors, fmt.Errorf("Dangerous requests action must be one of forbid, forward or node"))
	return
}

func validateSIPMode(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"full_gateway",
		"route",
		"sip_gateway":
		return
	}
	errors = append(errors, fmt.Errorf("SIP mode must be one of full_gateway, route or sip_gateway"))
	return
}

func validateSSLRequestClientCert(v interface{}, k string) (ws []string, errors []error) {
	switch v.(string) {
	case
		"dont_request",
		"request",
		"require":
		return
	}
	errors = append(errors, fmt.Errorf("SSL Request Client Cert must be one of dont_request, request or require"))
	return
}

func validateVirtualServerProtocol(v interface{}, k string) (ws []string, errors []error) {
	protocol := v.(string)
	protocolOptions := regexp.MustCompile(`^(client_first|dns|dns_tcp|ftp|http|https|imaps|imapv2|imapv3|imapv4|ldap|ldaps|pop3|pop3s|rtsp|server_first|siptcp|sipudp|smtp|ssl|stream|telnet|udp|udpstreaming)$`)
	if !protocolOptions.MatchString(protocol) {
		errors = append(errors, fmt.Errorf("%q must be one of client_first, dns, dns_tcp, ftp, http, https, imaps, imapv2, imapv3, imapv4, ldap, ldaps, pop3, pop3s, rtsp, server_first, siptcp, sipudp, smtp, ssl, stream, telnet, udp or udpstreaming", k))
	}
	return
}

func validateVirtualServerUseSSLSupport(v interface{}, k string) (ws []string, errors []error) {
	sslUseSSLSupport := v.(string)
	sslUseSSLSupportOptions := regexp.MustCompile(`^(use_default|disabled|enabled)$`)
	if !sslUseSSLSupportOptions.MatchString(sslUseSSLSupport) {
		errors = append(errors, fmt.Errorf("%q must be one of use_default, disabled or enabled", k))
	}
	return
}

func assignAptimizerValues(aptmizerMap map[string]interface{}) (aptimizerStruct virtualserver.Aptimizer) {
	enabled := aptmizerMap["enabled"].(bool)
	aptimizerStruct.Enabled = &enabled

	profileList := []virtualserver.AptimizerProfile{}
	var profile virtualserver.AptimizerProfile
	for _, value := range aptmizerMap["profile"].([]interface{}) {
		profileItem := value.(map[string]interface{})
		profile.Name = profileItem["name"].(string)
		profile.URLs = util.BuildStringArrayFromInterface(profileItem["urls"])
		profileList = append(profileList, profile)
	}
	aptimizerStruct.Profile = profileList
	return
}

func assignConnectionValues(connectionMap map[string]interface{}) (connectionStruct virtualserver.Connection) {
	keepAlive := connectionMap["keepalive"].(bool)
	connectionStruct.Keepalive = &keepAlive

	keepAliveTimeout := uint(connectionMap["keepalive_timeout"].(int))
	connectionStruct.KeepaliveTimeout = &keepAliveTimeout

	maxClientBuffer := uint(connectionMap["max_client_buffer"].(int))
	connectionStruct.MaxClientBuffer = &maxClientBuffer

	maxServerBuffer := uint(connectionMap["max_server_buffer"].(int))
	connectionStruct.MaxServerBuffer = &maxServerBuffer

	maxTransactionDuration := uint(connectionMap["max_transaction_duration"].(int))
	connectionStruct.MaxTransactionDuration = &maxTransactionDuration

	connectionStruct.ServerFirstBanner = connectionMap["server_first_banner"].(string)

	timeout := uint(connectionMap["timeout"].(int))
	connectionStruct.Timeout = &timeout
	return
}

func assignCookieValues(cookieMap map[string]interface{}) (cookieStruct virtualserver.Cookie) {
	cookieStruct.Domain = cookieMap["domain"].(string)
	cookieStruct.NewDomain = cookieMap["new_domain"].(string)
	cookieStruct.PathRegex = cookieMap["path_regex"].(string)
	cookieStruct.PathReplace = cookieMap["path_replace"].(string)
	cookieStruct.Secure = cookieMap["secure"].(string)
	return
}

func assignDNSValues(dnsMap map[string]interface{}) (dnsStruct virtualserver.DNS) {
	ednsClientSubnet := dnsMap["edns_client_subnet"].(bool)
	dnsStruct.EDNSClientSubnet = &ednsClientSubnet

	ednsUDPsize := uint(dnsMap["edns_udpsize"].(int))
	dnsStruct.EdnsUdpsize = &ednsUDPsize

	maxUDPSize := uint(dnsMap["max_udpsize"].(int))
	dnsStruct.MaxUdpsize = &maxUDPSize

	dnsStruct.RrsetOrder = dnsMap["rrset_order"].(string)

	verbose := dnsMap["verbose"].(bool)
	dnsStruct.Verbose = &verbose

	dnsStruct.Zones = util.BuildStringArrayFromInterface(dnsMap["zones"])
	return
}

func assignFTPValues(ftpMap map[string]interface{}) (ftpStruct virtualserver.Ftp) {
	dataSourcePort := uint(ftpMap["data_source_port"].(int))
	ftpStruct.DataSourcePort = &dataSourcePort

	forceClientSecure := ftpMap["force_client_secure"].(bool)
	ftpStruct.ForceClientSecure = &forceClientSecure

	portRangeHigh := uint(ftpMap["port_range_high"].(int))
	ftpStruct.PortRangeHigh = &portRangeHigh

	portRangeLow := uint(ftpMap["port_range_low"].(int))
	ftpStruct.PortRangeLow = &portRangeLow

	sslData := ftpMap["ssl_data"].(bool)
	ftpStruct.SslData = &sslData
	return
}

func assignGZIPValues(gzipMap map[string]interface{}) (gzipStruct virtualserver.Gzip) {
	compressLevel := uint(gzipMap["compress_level"].(int))
	gzipStruct.CompressLevel = &compressLevel

	enabled := gzipMap["enabled"].(bool)
	gzipStruct.Enabled = &enabled

	gzipStruct.EtagRewrite = gzipMap["etag_rewrite"].(string)
	gzipStruct.IncludeMime = util.BuildStringArrayFromInterface(gzipMap["include_mime"])

	maxSize := uint(gzipMap["max_size"].(int))
	gzipStruct.MaxSize = &maxSize

	minSize := uint(gzipMap["min_size"].(int))
	gzipStruct.MinSize = &minSize

	noSize := gzipMap["no_size"].(bool)
	gzipStruct.NoSize = &noSize
	return
}

func assignHTTPValues(httpMap map[string]interface{}) (httpStruct virtualserver.HTTP) {
	httpStruct.ChunkOverheadForwarding = httpMap["chunk_overhead_forwarding"].(string)
	httpStruct.LocationRegex = httpMap["location_regex"].(string)
	httpStruct.LocationReplace = httpMap["location_replace"].(string)
	httpStruct.LocationRewrite = httpMap["location_rewrite"].(string)
	httpStruct.MIMEDefault = httpMap["mime_default"].(string)

	mimeDetect := httpMap["mime_detect"].(bool)
	httpStruct.MIMEDetect = &mimeDetect
	return
}

func assignHTTP2Values(http2Map map[string]interface{}) (http2Struct virtualserver.HTTP2) {
	connectTimeout := uint(http2Map["connect_timeout"].(int))
	http2Struct.ConnectTimeout = &connectTimeout

	dataFrameSize := uint(http2Map["data_frame_size"].(int))
	http2Struct.DataFrameSize = &dataFrameSize

	enabled := http2Map["enabled"].(bool)
	http2Struct.Enabled = &enabled

	headerTableSize := uint(http2Map["header_table_size"].(int))
	http2Struct.HeaderTableSize = &headerTableSize

	http2Struct.HeadersIndexBlacklist = util.BuildStringArrayFromInterface(http2Map["headers_index_blacklist"])

	headersIndexDefault := http2Map["headers_index_default"].(bool)
	http2Struct.HeadersIndexDefault = &headersIndexDefault
	http2Struct.HeadersIndexWhitelist = util.BuildStringArrayFromInterface(http2Map["headers_index_whitelist"])

	idleTimeoutNoStreams := uint(http2Map["idle_timeout_no_streams"].(int))
	http2Struct.IdleTimeoutNoStreams = &idleTimeoutNoStreams

	idleTimeoutOpenStreams := uint(http2Map["idle_timeout_open_streams"].(int))
	http2Struct.IdleTimeoutOpenStreams = &idleTimeoutOpenStreams

	maxConcurrentStreams := uint(http2Map["max_concurrent_streams"].(int))
	http2Struct.MaxConcurrentStreams = &maxConcurrentStreams

	maxFrameSize := uint(http2Map["max_frame_size"].(int))
	http2Struct.MaxFrameSize = &maxFrameSize

	maxHeaderPadding := uint(http2Map["max_header_padding"].(int))
	http2Struct.MaxHeaderPadding = &maxHeaderPadding

	mergeCookieHeaders := http2Map["merge_cookie_headers"].(bool)
	http2Struct.MergeCookieHeaders = &mergeCookieHeaders

	streamWindowSize := uint(http2Map["stream_window_size"].(int))
	http2Struct.StreamWindowSize = &streamWindowSize
	return
}

func assignKerberosProtocolTransitionValues(kptMap map[string]interface{}) (kptStruct virtualserver.KerberosProtocolTransition) {
	enabled := kptMap["enabled"].(bool)
	kptStruct.Enabled = &enabled

	kptStruct.Principal = kptMap["principal"].(string)
	kptStruct.Target = kptMap["target"].(string)
	return
}

func assignLogValues(logMap map[string]interface{}) (logStruct virtualserver.Log) {
	alwaysFlush := logMap["always_flush"].(bool)
	logStruct.AlwaysFlush = &alwaysFlush

	clientConnectionFailures := logMap["client_connection_failures"].(bool)
	logStruct.ClientConnectionFailures = &clientConnectionFailures

	enabled := logMap["enabled"].(bool)
	logStruct.Enabled = &enabled

	logStruct.Filename = logMap["filename"].(string)
	logStruct.Format = logMap["format"].(string)

	saveAll := logMap["save_all"].(bool)
	logStruct.SaveAll = &saveAll

	serverConnectionFailures := logMap["server_connection_failures"].(bool)
	logStruct.ServerConnectionFailures = &serverConnectionFailures

	sessionPersistenceVerbose := logMap["session_persistence_verbose"].(bool)
	logStruct.SessionPersistenceVerbose = &sessionPersistenceVerbose

	sslFailures := logMap["ssl_failures"].(bool)
	logStruct.SSLFailures = &sslFailures
	return
}

func assignRecentConnectionsValues(recentConnectionsMap map[string]interface{}) (recentConnectionsStruct virtualserver.RecentConnections) {
	enabled := recentConnectionsMap["enabled"].(bool)
	recentConnectionsStruct.Enabled = &enabled

	saveAll := recentConnectionsMap["save_all"].(bool)
	recentConnectionsStruct.SaveAll = &saveAll
	return
}

func assignRequestTracingValues(requestTracingMap map[string]interface{}) (requestTracingStruct virtualserver.RequestTracing) {
	enabled := requestTracingMap["enabled"].(bool)
	requestTracingStruct.Enabled = &enabled

	traceIO := requestTracingMap["trace_io"].(bool)
	requestTracingStruct.TraceIO = &traceIO
	return
}

func assignRTSPValues(rtspMap map[string]interface{}) (rtspStruct virtualserver.RTSP) {
	streamingPortRangeHigh := uint(rtspMap["streaming_port_range_high"].(int))
	rtspStruct.StreamingPortRangeHigh = &streamingPortRangeHigh

	streamingPortRangeLow := uint(rtspMap["streaming_port_range_low"].(int))
	rtspStruct.StreamingPortRangeLow = &streamingPortRangeLow

	streamingTimeout := uint(rtspMap["streaming_timeout"].(int))
	rtspStruct.StreamingTimeout = &streamingTimeout
	return
}

func assignSIPValues(sipMap map[string]interface{}) (sipStruct virtualserver.SIP) {
	sipStruct.DangerousRequests = sipMap["dangerous_requests"].(string)

	followRoute := sipMap["follow_route"].(bool)
	sipStruct.FollowRoute = &followRoute

	maxConnectionMem := uint(sipMap["max_connection_mem"].(int))
	sipStruct.MaxConnectionMem = &maxConnectionMem

	sipStruct.Mode = sipMap["mode"].(string)

	rewriteURI := sipMap["rewrite_uri"].(bool)
	sipStruct.RewriteURI = &rewriteURI

	streamingPortRangeHigh := uint(sipMap["streaming_port_range_high"].(int))
	sipStruct.StreamingPortRangeHigh = &streamingPortRangeHigh

	streamingPortRangeLow := uint(sipMap["streaming_port_range_low"].(int))
	sipStruct.StreamingPortRangeLow = &streamingPortRangeLow

	streamingTimeout := uint(sipMap["streaming_timeout"].(int))
	sipStruct.StreamingTimeout = &streamingTimeout

	timeoutMessages := sipMap["timeout_messages"].(bool)
	sipStruct.TimeoutMessages = &timeoutMessages

	transactionTimeout := uint(sipMap["transaction_timeout"].(int))
	sipStruct.TransactionTimeout = &transactionTimeout
	return
}

func assignSSLValues(sslMap map[string]interface{}) (sslStruct virtualserver.Ssl) {
	addHTTPHeaders := sslMap["add_http_headers"].(bool)
	sslStruct.AddHTTPHeaders = &addHTTPHeaders

	sslStruct.ClientCertCAS = util.BuildStringArrayFromInterface(sslMap["client_cert_cas"])
	sslStruct.EllipticCurves = util.BuildStringArrayFromInterface(sslMap["elliptic_curves"])
	sslStruct.IssuedCertsNeverExpire = util.BuildStringArrayFromInterface(sslMap["issued_certs_never_expire"])

	enable := sslMap["ocsp_enable"].(bool)
	sslStruct.OCSPEnable = &enable

	ocspIssuerList := []virtualserver.OCSPIssuer{}
	var ocspIssuer virtualserver.OCSPIssuer

	for _, value := range sslMap["ocsp_issuers"].([]interface{}) {
		issuerElement := value.(map[string]interface{})
		ocspIssuer.Issuer = issuerElement["issuer"].(string)
		aia := issuerElement["aia"].(bool)
		ocspIssuer.AIA = &aia
		ocspIssuer.Nonce = issuerElement["nonce"].(string)
		ocspIssuer.Required = issuerElement["required"].(string)
		ocspIssuer.ResponderCert = issuerElement["responder_cert"].(string)
		ocspIssuer.Signer = issuerElement["signer"].(string)
		ocspIssuer.URL = issuerElement["url"].(string)
		ocspIssuerList = append(ocspIssuerList, ocspIssuer)
	}
	sslStruct.OCSPIssuers = ocspIssuerList

	ocspMaxResponseAge := uint(sslMap["ocsp_max_response_age"].(int))
	sslStruct.OCSPMaxResponseAge = &ocspMaxResponseAge

	ocspStapling := sslMap["ocsp_stapling"].(bool)
	sslStruct.OCSPStapling = &ocspStapling

	ocspTimeTolerance := uint(sslMap["ocsp_time_tolerance"].(int))
	sslStruct.OCSPTimeTolerance = &ocspTimeTolerance

	ocspTimeout := uint(sslMap["ocsp_timeout"].(int))
	sslStruct.OCSPTimeout = &ocspTimeout

	preferSSLv3 := sslMap["prefer_sslv3"].(bool)
	sslStruct.PreferSSLv3 = &preferSSLv3

	sslStruct.RequestClientCert = sslMap["request_client_cert"].(string)

	sendCloseAlerts := sslMap["send_close_alerts"].(bool)
	sslStruct.SendCloseAlerts = &sendCloseAlerts

	sslStruct.ServerCertAltCertificates = util.BuildStringArrayFromInterface(sslMap["server_cert_alt_certificates"])
	sslStruct.ServerCertDefault = sslMap["server_cert_default"].(string)

	certHostMappingList := []virtualserver.CertItem{}
	var certHostMapping virtualserver.CertItem

	for _, value := range sslMap["ssl_server_cert_host_mapping"].([]interface{}) {
		indivdualHostMapping := value.(map[string]interface{})
		certHostMapping.Host = indivdualHostMapping["host"].(string)
		certHostMapping.Certificate = indivdualHostMapping["certificate"].(string)
		certHostMapping.AltCertificates = util.BuildStringArrayFromInterface(indivdualHostMapping["alt_certificates"])
		certHostMappingList = append(certHostMappingList, certHostMapping)
	}
	sslStruct.ServerCertHostMap = certHostMappingList

	sslStruct.SignatureAlgorithms = sslMap["signature_algorithms"].(string)
	sslStruct.SSLCiphers = sslMap["ssl_ciphers"].(string)
	sslStruct.SslSupportSsl2 = sslMap["ssl_support_ssl2"].(string)
	sslStruct.SslSupportSsl3 = sslMap["ssl_support_ssl3"].(string)
	sslStruct.SslSupportTLS1 = sslMap["ssl_support_tls1"].(string)
	sslStruct.SslSupportTLS1_1 = sslMap["ssl_support_tls1_1"].(string)
	sslStruct.SslSupportTLS1_2 = sslMap["ssl_support_tls1_2"].(string)

	trustMagic := sslMap["trust_magic"].(bool)
	sslStruct.TrustMagic = &trustMagic

	return
}

func assignSysLogValues(sysLogMap map[string]interface{}) (sysLogStruct virtualserver.SysLog) {
	enabled := sysLogMap["enabled"].(bool)
	sysLogStruct.Enabled = &enabled

	sysLogStruct.Format = sysLogMap["format"].(string)
	sysLogStruct.IPEndpoint = sysLogMap["ip_end_point"].(string)

	msgLenLimit := uint(sysLogMap["msg_len_limit"].(int))
	sysLogStruct.MsgLenLimit = &msgLenLimit
	return
}

func assignUDPValues(udpMap map[string]interface{}) (udpStruct virtualserver.UDP) {
	endPointPersistence := udpMap["end_point_persistence"].(bool)
	udpStruct.EndPointPersistence = &endPointPersistence

	portSMP := udpMap["port_smp"].(bool)
	udpStruct.PortSMP = &portSMP

	responseDatagramsExpected := udpMap["response_datagrams_expected"].(int)
	udpStruct.ResponseDatagramsExpected = &responseDatagramsExpected

	timeout := uint(udpMap["timeout"].(int))
	udpStruct.Timeout = &timeout
	return
}

func assignWebCacheValues(webCacheMap map[string]interface{}) (webCacheStruct virtualserver.WebCache) {
	webCacheStruct.ControlOut = webCacheMap["control_out"].(string)

	enabled := webCacheMap["enabled"].(bool)
	webCacheStruct.Enabled = &enabled

	errorPageTime := uint(webCacheMap["error_page_time"].(int))
	webCacheStruct.ErrorPageTime = &errorPageTime

	maxTime := uint(webCacheMap["max_time"].(int))
	webCacheStruct.MaxTime = &maxTime

	refreshTime := uint(webCacheMap["refresh_time"].(int))
	webCacheStruct.RefreshTime = &refreshTime
	return
}

func resourceVirtualServerCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var virtualServer virtualserver.VirtualServer

	virtualServerName := d.Get("name").(string)
	virtualServer.Properties.Basic.AddClusterIP = d.Get("add_cluster_ip").(bool)
	virtualServer.Properties.Basic.AddXForwarded = d.Get("add_x_forwarded_for").(bool)
	virtualServer.Properties.Basic.AddXForwardedProto = d.Get("add_x_forwarded_proto").(bool)
	virtualServer.Properties.Basic.AutoDetectUpgradeHeaders = d.Get("autodetect_upgrade_headers").(bool)
	virtualServer.Properties.Basic.BandwidthClass = d.Get("bandwidth_class").(string)
	virtualServer.Properties.Basic.CloseWithRst = d.Get("close_with_rst").(bool)
	virtualServer.Properties.Basic.CompletionRules = util.BuildStringArrayFromInterface(d.Get("completionrules"))
	virtualServer.Properties.Basic.ConnectTimeout = uint(d.Get("connect_timeout").(int))
	virtualServer.Properties.Basic.Enabled = d.Get("enabled").(bool)
	virtualServer.Properties.Basic.FtpForceServerSecure = d.Get("ftp_force_server_secure").(bool)
	virtualServer.Properties.Basic.GlbServices = util.BuildStringArrayFromInterface(d.Get("glb_services"))
	virtualServer.Properties.Basic.ListenOnAny = d.Get("listen_on_any").(bool)
	virtualServer.Properties.Basic.ListenOnHosts = util.BuildStringArrayFromInterface(d.Get("listen_on_hosts"))
	virtualServer.Properties.Basic.ListenOnTrafficIps = util.BuildStringArrayFromInterface(d.Get("listen_on_traffic_ips"))
	virtualServer.Properties.Basic.Note = d.Get("note").(string)
	virtualServer.Properties.Basic.Pool = d.Get("pool").(string)
	virtualServer.Properties.Basic.Port = uint(d.Get("port").(int))
	virtualServer.Properties.Basic.ProtectionClass = d.Get("protection_class").(string)
	virtualServer.Properties.Basic.Protocol = d.Get("protocol").(string)
	virtualServer.Properties.Basic.RequestRules = util.BuildStringArrayFromInterface(d.Get("request_rules"))
	virtualServer.Properties.Basic.ResponseRules = util.BuildStringArrayFromInterface(d.Get("response_rules"))
	virtualServer.Properties.Basic.SlmClass = d.Get("slm_class").(string)
	virtualServer.Properties.Basic.SoNagle = d.Get("so_nagle").(bool)
	virtualServer.Properties.Basic.SslClientCertHeaders = d.Get("ssl_client_cert_headers").(string)
	virtualServer.Properties.Basic.SslDecrypt = d.Get("ssl_decrypt").(bool)
	virtualServer.Properties.Basic.SslHonorFallbackScsv = d.Get("ssl_honor_fallback_scsv").(string)
	virtualServer.Properties.Basic.Transparent = d.Get("transparent").(bool)
	virtualServer.Properties.ConnectionErrors.ErrorFile = d.Get("error_file").(string)
	expectSTARTTLS := d.Get("expect_starttls").(bool)
	virtualServer.Properties.SMTP.ExpectSTARTTLS = &expectSTARTTLS
	proxyClose := d.Get("proxy_close").(bool)
	virtualServer.Properties.TCP.ProxyClose = &proxyClose

	if v, ok := d.GetOk("aptimizer"); ok {
		aptimizerList := v.([]interface{})
		virtualServer.Properties.Aptimizer = assignAptimizerValues(aptimizerList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("vs_connection"); ok {
		connectionList := v.([]interface{})
		virtualServer.Properties.Connection = assignConnectionValues(connectionList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("cookie"); ok {
		cookieList := v.([]interface{})
		virtualServer.Properties.Cookie = assignCookieValues(cookieList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("dns"); ok {
		dnsList := v.([]interface{})
		virtualServer.Properties.DNS = assignDNSValues(dnsList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("ftp"); ok {
		ftpList := v.([]interface{})
		virtualServer.Properties.Ftp = assignFTPValues(ftpList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("gzip"); ok {
		gzipList := v.([]interface{})
		virtualServer.Properties.Gzip = assignGZIPValues(gzipList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("http"); ok {
		httpList := v.([]interface{})
		virtualServer.Properties.HTTP = assignHTTPValues(httpList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("http2"); ok {
		http2List := v.([]interface{})
		virtualServer.Properties.HTTP2 = assignHTTP2Values(http2List[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("kerberos_protocol_transition"); ok {
		kptList := v.([]interface{})
		virtualServer.Properties.KerberosProtocolTransition = assignKerberosProtocolTransitionValues(kptList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("log"); ok {
		logList := v.([]interface{})
		virtualServer.Properties.Log = assignLogValues(logList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("recent_connections"); ok {
		recentConnectionsList := v.([]interface{})
		virtualServer.Properties.RecentConnections = assignRecentConnectionsValues(recentConnectionsList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("request_tracing"); ok {
		requestTracingList := v.([]interface{})
		virtualServer.Properties.RequestTracing = assignRequestTracingValues(requestTracingList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("rtsp"); ok {
		rtspList := v.([]interface{})
		virtualServer.Properties.RTSP = assignRTSPValues(rtspList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("sip"); ok {
		sipList := v.([]interface{})
		virtualServer.Properties.SIP = assignSIPValues(sipList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("ssl"); ok {
		sslList := v.([]interface{})
		virtualServer.Properties.Ssl = assignSSLValues(sslList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("syslog"); ok {
		sysLogList := v.([]interface{})
		virtualServer.Properties.SysLog = assignSysLogValues(sysLogList[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("udp"); ok {
		udpList := v.([]interface{})
		virtualServer.Properties.UDP = assignUDPValues(udpList[0].(map[string]interface{}))
	}
	if v, ok := d.GetOk("web_cache"); ok {
		webCacheList := v.([]interface{})
		virtualServer.Properties.WebCache = assignWebCacheValues(webCacheList[0].(map[string]interface{}))
	}

	createAPI := virtualserver.NewCreate(virtualServerName, virtualServer)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("BrocadeVTM Virtual Server error whilst creating %s: %v", virtualServerName, err))
	}
	d.SetId(virtualServerName)

	return resourceVirtualServerRead(d, m)
}

func resourceVirtualServerRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)

	getSingleAPI := virtualserver.NewGet(d.Id())
	err := vtmClient.Do(getSingleAPI)
	if err != nil {
		if getSingleAPI.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("BrocadeVTM Virtual Server error whilst retrieving %s: %v", d.Id(), err))
	}

	returnedVirtualServer := *getSingleAPI.ResponseObject().(*virtualserver.VirtualServer)

	d.Set("add_cluster_ip", returnedVirtualServer.Properties.Basic.AddClusterIP)
	d.Set("add_x_forwarded_for", returnedVirtualServer.Properties.Basic.AddXForwarded)
	d.Set("add_x_forwarded_proto", returnedVirtualServer.Properties.Basic.AddXForwardedProto)
	d.Set("autodetect_upgrade_headers", returnedVirtualServer.Properties.Basic.AutoDetectUpgradeHeaders)
	d.Set("bandwidth_class", returnedVirtualServer.Properties.Basic.BandwidthClass)
	d.Set("close_with_rst", returnedVirtualServer.Properties.Basic.CloseWithRst)
	d.Set("completionrules", returnedVirtualServer.Properties.Basic.CompletionRules)
	d.Set("connect_timeout", returnedVirtualServer.Properties.Basic.ConnectTimeout)
	d.Set("enabled", returnedVirtualServer.Properties.Basic.Enabled)
	d.Set("ftp_force_server_secure", returnedVirtualServer.Properties.Basic.FtpForceServerSecure)
	d.Set("glb_services", returnedVirtualServer.Properties.Basic.GlbServices)
	d.Set("listen_on_any", returnedVirtualServer.Properties.Basic.ListenOnAny)
	d.Set("listen_on_hosts", returnedVirtualServer.Properties.Basic.ListenOnHosts)
	d.Set("listen_on_traffic_ips", returnedVirtualServer.Properties.Basic.ListenOnTrafficIps)
	d.Set("note", returnedVirtualServer.Properties.Basic.Note)
	d.Set("pool", returnedVirtualServer.Properties.Basic.Pool)
	d.Set("port", returnedVirtualServer.Properties.Basic.Port)
	d.Set("protection_class", returnedVirtualServer.Properties.Basic.ProtectionClass)
	d.Set("protocol", returnedVirtualServer.Properties.Basic.Protocol)
	d.Set("request_rules", returnedVirtualServer.Properties.Basic.RequestRules)
	d.Set("response_rules", returnedVirtualServer.Properties.Basic.ResponseRules)
	d.Set("slm_class", returnedVirtualServer.Properties.Basic.SlmClass)
	d.Set("so_nagle", returnedVirtualServer.Properties.Basic.SoNagle)
	d.Set("ssl_client_cert_headers", returnedVirtualServer.Properties.Basic.SslClientCertHeaders)
	d.Set("ssl_decrypt", returnedVirtualServer.Properties.Basic.SslDecrypt)
	d.Set("ssl_honor_fallback_scsv", returnedVirtualServer.Properties.Basic.SslHonorFallbackScsv)
	d.Set("transparent", returnedVirtualServer.Properties.Basic.Transparent)
	d.Set("error_file", returnedVirtualServer.Properties.ConnectionErrors.ErrorFile)
	d.Set("expect_starttls", returnedVirtualServer.Properties.SMTP.ExpectSTARTTLS)
	d.Set("proxy_close", returnedVirtualServer.Properties.TCP.ProxyClose)

	d.Set("aptimizer", []virtualserver.Aptimizer{returnedVirtualServer.Properties.Aptimizer})
	d.Set("vs_connection", []virtualserver.Connection{returnedVirtualServer.Properties.Connection})
	d.Set("cookie", []virtualserver.Cookie{returnedVirtualServer.Properties.Cookie})
	d.Set("dns", []virtualserver.DNS{returnedVirtualServer.Properties.DNS})
	d.Set("ftp", []virtualserver.Ftp{returnedVirtualServer.Properties.Ftp})
	d.Set("gzip", []virtualserver.Gzip{returnedVirtualServer.Properties.Gzip})
	d.Set("http", []virtualserver.HTTP{returnedVirtualServer.Properties.HTTP})
	d.Set("http2", []virtualserver.HTTP2{returnedVirtualServer.Properties.HTTP2})
	d.Set("kerberos_protocol_transition", []virtualserver.KerberosProtocolTransition{returnedVirtualServer.Properties.KerberosProtocolTransition})
	d.Set("log", []virtualserver.Log{returnedVirtualServer.Properties.Log})
	d.Set("recent_connections", []virtualserver.RecentConnections{returnedVirtualServer.Properties.RecentConnections})
	d.Set("request_tracing", []virtualserver.RequestTracing{returnedVirtualServer.Properties.RequestTracing})
	d.Set("rtsp", []virtualserver.RTSP{returnedVirtualServer.Properties.RTSP})
	d.Set("sip", []virtualserver.SIP{returnedVirtualServer.Properties.SIP})
	d.Set("ssl", []virtualserver.Ssl{returnedVirtualServer.Properties.Ssl})
	d.Set("syslog", []virtualserver.SysLog{returnedVirtualServer.Properties.SysLog})
	d.Set("udp", []virtualserver.UDP{returnedVirtualServer.Properties.UDP})
	d.Set("web_cache", []virtualserver.WebCache{returnedVirtualServer.Properties.WebCache})

	return nil
}

func resourceVirtualServerUpdate(d *schema.ResourceData, m interface{}) error {
	var virtualServer virtualserver.VirtualServer
	var hasChanges bool

	virtualServer.Properties.Basic.Pool = d.Get("pool").(string)
	virtualServer.Properties.Basic.Port = uint(d.Get("port").(int))

	oldAddClusterIP, newAddClusterIP := d.GetChange("add_cluster_ip")
	if oldAddClusterIP.(bool) != newAddClusterIP.(bool) {
		virtualServer.Properties.Basic.AddClusterIP = newAddClusterIP.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.AddClusterIP = oldAddClusterIP.(bool)
	}

	oldAddXForwardedFor, newAddForwardedFor := d.GetChange("add_x_forwarded_for")
	if oldAddXForwardedFor.(bool) != newAddForwardedFor.(bool) {
		virtualServer.Properties.Basic.AddXForwarded = newAddForwardedFor.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.AddXForwarded = oldAddXForwardedFor.(bool)
	}

	oldAddXForwardedProto, newAddXForwardedProto := d.GetChange("add_x_forwarded_proto")
	if oldAddXForwardedProto.(bool) != newAddXForwardedProto.(bool) {
		virtualServer.Properties.Basic.AddXForwardedProto = newAddXForwardedProto.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.AddXForwardedProto = oldAddXForwardedProto.(bool)
	}

	oldAutodetectUpgradeHeaders, newAutodetectUpgradeHeaders := d.GetChange("autodetect_upgrade_headers")
	if oldAutodetectUpgradeHeaders.(bool) != newAutodetectUpgradeHeaders.(bool) {
		virtualServer.Properties.Basic.AutoDetectUpgradeHeaders = newAutodetectUpgradeHeaders.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.AutoDetectUpgradeHeaders = oldAutodetectUpgradeHeaders.(bool)
	}

	if d.HasChange("bandwidth_class") {
		virtualServer.Properties.Basic.BandwidthClass = d.Get("bandwidth_class").(string)
		hasChanges = true
	}

	oldCloseWithRST, newCloseWithRST := d.GetChange("close_with_rst")
	if oldCloseWithRST.(bool) != newCloseWithRST.(bool) {
		virtualServer.Properties.Basic.CloseWithRst = newCloseWithRST.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.CloseWithRst = oldCloseWithRST.(bool)
	}

	if d.HasChange("completionrules") {
		virtualServer.Properties.Basic.CompletionRules = util.BuildStringArrayFromInterface(d.Get("completionrules"))
		hasChanges = true
	}

	oldConnectTimeout, newConnectTimeout := d.GetChange("connect_timeout")
	if uint(oldConnectTimeout.(int)) != uint(newConnectTimeout.(int)) {
		virtualServer.Properties.Basic.ConnectTimeout = uint(newConnectTimeout.(int))
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.ConnectTimeout = uint(oldConnectTimeout.(int))
	}

	oldEnable, newEnable := d.GetChange("enable")
	if oldEnable.(bool) != newEnable.(bool) {
		virtualServer.Properties.Basic.Enabled = newEnable.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.Enabled = oldEnable.(bool)
	}

	oldForceServerSecure, newForceServerSecure := d.GetChange("ftp_force_server_secure")
	if oldForceServerSecure.(bool) != newForceServerSecure.(bool) {
		virtualServer.Properties.Basic.FtpForceServerSecure = newForceServerSecure.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.FtpForceServerSecure = oldForceServerSecure.(bool)
	}

	if d.HasChange("glb_services") {
		virtualServer.Properties.Basic.GlbServices = util.BuildStringArrayFromInterface(d.Get("glb_services"))
		hasChanges = true
	}

	oldListenOnAny, newListenOnAny := d.GetChange("listen_on_any")
	if oldListenOnAny.(bool) != newListenOnAny.(bool) {
		virtualServer.Properties.Basic.ListenOnAny = newListenOnAny.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.ListenOnAny = oldListenOnAny.(bool)
	}

	if d.HasChange("listen_on_hosts") {
		virtualServer.Properties.Basic.ListenOnHosts = util.BuildStringArrayFromInterface(d.Get("listen_on_hosts"))
		hasChanges = true
	}

	if d.HasChange("listen_on_traffic_ips") {
		virtualServer.Properties.Basic.ListenOnTrafficIps = util.BuildStringArrayFromInterface(d.Get("listen_on_traffic_ips"))
		hasChanges = true
	}

	if d.HasChange("note") {
		virtualServer.Properties.Basic.Note = d.Get("note").(string)
		hasChanges = true
	}

	if d.HasChange("protection_class") {
		virtualServer.Properties.Basic.ProtectionClass = d.Get("protection_class").(string)
		hasChanges = true
	}

	if d.HasChange("protection_class") {
		virtualServer.Properties.Basic.Protocol = d.Get("protocol").(string)
		hasChanges = true
	}

	if d.HasChange("request_rules") {
		virtualServer.Properties.Basic.RequestRules = util.BuildStringArrayFromInterface(d.Get("request_rules"))
		hasChanges = true
	}

	if d.HasChange("response_rules") {
		virtualServer.Properties.Basic.ResponseRules = util.BuildStringArrayFromInterface(d.Get("response_rules"))
		hasChanges = true
	}

	if d.HasChange("slm_class") {
		virtualServer.Properties.Basic.SlmClass = d.Get("slm_class").(string)
		hasChanges = true
	}

	oldSoNagle, newSoNagle := d.GetChange("so_nagle")
	if oldSoNagle.(bool) != newSoNagle.(bool) {
		virtualServer.Properties.Basic.SoNagle = newSoNagle.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.SoNagle = oldSoNagle.(bool)
	}

	if d.HasChange("ssl_client_cert_headers") {
		virtualServer.Properties.Basic.SslClientCertHeaders = d.Get("ssl_client_cert_headers").(string)
		hasChanges = true
	}

	oldSSLDecrypt, newSSLDecrypt := d.GetChange("ssl_decrypt")
	if oldSSLDecrypt.(bool) != newSSLDecrypt.(bool) {
		virtualServer.Properties.Basic.SslDecrypt = newSSLDecrypt.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.SslDecrypt = oldSSLDecrypt.(bool)
	}

	if d.HasChange("ssl_honor_fallback_scsv") {
		virtualServer.Properties.Basic.SslHonorFallbackScsv = d.Get("ssl_honor_fallback_scsv").(string)
		hasChanges = true
	}

	oldTransparent, newTransparent := d.GetChange("transparent")
	if oldTransparent.(bool) != newTransparent.(bool) {
		virtualServer.Properties.Basic.Transparent = newTransparent.(bool)
		hasChanges = true
	} else {
		virtualServer.Properties.Basic.Transparent = oldTransparent.(bool)
	}

	if d.HasChange("error_file") {
		virtualServer.Properties.ConnectionErrors.ErrorFile = d.Get("error_file").(string)
		hasChanges = true
	}

	oldExpectSTARTTLS, newExpectSTARTTLS := d.GetChange("expect_starttls")
	if oldExpectSTARTTLS.(bool) != newExpectSTARTTLS.(bool) {
		starTTL := newExpectSTARTTLS.(bool)
		virtualServer.Properties.SMTP.ExpectSTARTTLS = &starTTL
		hasChanges = true
	} else {
		starTTL := oldExpectSTARTTLS.(bool)
		virtualServer.Properties.SMTP.ExpectSTARTTLS = &starTTL
	}

	oldProxyClose, newProxyClose := d.GetChange("proxy_close")
	if oldProxyClose.(bool) != newProxyClose.(bool) {
		proxyClose := newProxyClose.(bool)
		virtualServer.Properties.TCP.ProxyClose = &proxyClose
		hasChanges = true
	} else {
		proxyClose := oldProxyClose.(bool)
		virtualServer.Properties.TCP.ProxyClose = &proxyClose
	}

	if d.HasChange("aptimizer") {
		aptimizerList := d.Get("aptimizer").([]interface{})
		virtualServer.Properties.Aptimizer = assignAptimizerValues(aptimizerList[0].(map[string]interface{}))
		hasChanges = true
	}

	if d.HasChange("vs_connection") {
		connectionList := d.Get("vs_connection").([]interface{})
		virtualServer.Properties.Connection = assignConnectionValues(connectionList[0].(map[string]interface{}))
		hasChanges = true
	}

	if d.HasChange("cookie") {
		cookieList := d.Get("cookie").([]interface{})
		virtualServer.Properties.Cookie = assignCookieValues(cookieList[0].(map[string]interface{}))
	}

	if d.HasChange("dns") {
		dnsList := d.Get("dns").([]interface{})
		virtualServer.Properties.DNS = assignDNSValues(dnsList[0].(map[string]interface{}))
	}

	if d.HasChange("ftp") {
		ftpList := d.Get("ftp").([]interface{})
		virtualServer.Properties.Ftp = assignFTPValues(ftpList[0].(map[string]interface{}))
	}

	if d.HasChange("gzip") {
		gzipList := d.Get("gzip").([]interface{})
		virtualServer.Properties.Gzip = assignGZIPValues(gzipList[0].(map[string]interface{}))
	}

	if d.HasChange("http") {
		httpList := d.Get("http").([]interface{})
		virtualServer.Properties.HTTP = assignHTTPValues(httpList[0].(map[string]interface{}))
	}

	if d.HasChange("http2") {
		http2List := d.Get("http2").([]interface{})
		virtualServer.Properties.HTTP2 = assignHTTP2Values(http2List[0].(map[string]interface{}))
	}

	if d.HasChange("kerberos_protocol_transition") {
		kptList := d.Get("kerberos_protocol_transition").([]interface{})
		virtualServer.Properties.KerberosProtocolTransition = assignKerberosProtocolTransitionValues(kptList[0].(map[string]interface{}))
	}

	if d.HasChange("log") {
		logList := d.Get("log").([]interface{})
		virtualServer.Properties.Log = assignLogValues(logList[0].(map[string]interface{}))
	}

	if d.HasChange("recent_connections") {
		recentConnectionsList := d.Get("recent_connections").([]interface{})
		virtualServer.Properties.RecentConnections = assignRecentConnectionsValues(recentConnectionsList[0].(map[string]interface{}))
	}

	if d.HasChange("request_tracing") {
		requestTracingList := d.Get("request_tracing").([]interface{})
		virtualServer.Properties.RequestTracing = assignRequestTracingValues(requestTracingList[0].(map[string]interface{}))
	}

	if d.HasChange("rtsp") {
		rtspList := d.Get("rtsp").([]interface{})
		virtualServer.Properties.RTSP = assignRTSPValues(rtspList[0].(map[string]interface{}))
	}

	if d.HasChange("sip") {
		sipList := d.Get("sip").([]interface{})
		virtualServer.Properties.SIP = assignSIPValues(sipList[0].(map[string]interface{}))
	}

	if d.HasChange("ssl") {
		sslList := d.Get("ssl").([]interface{})
		virtualServer.Properties.Ssl = assignSSLValues(sslList[0].(map[string]interface{}))
	}

	if d.HasChange("syslog") {
		sysLogList := d.Get("syslog").([]interface{})
		virtualServer.Properties.SysLog = assignSysLogValues(sysLogList[0].(map[string]interface{}))
	}

	if d.HasChange("udp") {
		udpList := d.Get("udp").([]interface{})
		virtualServer.Properties.UDP = assignUDPValues(udpList[0].(map[string]interface{}))
	}

	if d.HasChange("web_cache") {
		webCacheList := d.Get("web_cache").([]interface{})
		virtualServer.Properties.WebCache = assignWebCacheValues(webCacheList[0].(map[string]interface{}))
	}

	if hasChanges {
		headers := make(map[string]string)
		client := m.(*rest.Client)
		vtmClient := *client
		headers["Content-Type"] = "application/json"
		vtmClient.Headers = headers

		updateAPI := virtualserver.NewUpdate(d.Id(), virtualServer)
		err := vtmClient.Do(updateAPI)

		if err != nil {
			return fmt.Errorf("BrocadeVTM error whilst updating virtual server %s: %vv", d.Id(), err)
		}
		return resourceVirtualServerRead(d, m)
	}

	return nil
}

func resourceVirtualServerDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var virtualServerName string

	if v, ok := d.GetOk("name"); ok && v != "" {
		virtualServerName = v.(string)
	}

	deleteAPI := virtualserver.NewDelete(virtualServerName)
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf(fmt.Sprintf("BrocadeVTM Virtual Server error whilst deleting %s: %v", virtualServerName, err))
	}

	d.SetId("")
	return nil
}
