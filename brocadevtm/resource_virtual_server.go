package brocadevtm

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
)

func resourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualServerSet,
		Read:   resourceVirtualServerRead,
		Update: resourceVirtualServerSet,
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
			/*
				"mss": {
					Type:         schema.TypeInt,
					Description:  "The maximum TCP segment size. This will place a maximum on the size of TCP segments that are sent by this machine, and will advertise to the client this value as the maximum size of TCP segment to send to this machine. Setting this to zero causes the default maximum TCP segment size to be advertised and used.",
					Optional:     true,
					ValidateFunc: util.ValidateUnsignedInteger,
				},
			*/
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
				Type:        schema.TypeString,
				Description: "The protocol that the virtual server is using.",
				Optional:    true,
				Default:     "http",
				ValidateFunc: validation.StringInSlice([]string{
					"client_first", "dns", "dns_tcp", "ftp", "http", "https", "imaps", "imapv2",
					"imapv3", "imapv4", "ldap", "ldaps", "pop3", "pop3s", "rtsp", "server_first",
					"siptcp", "sipudp", "smtp", "ssl", "stream", "telnet", "udp", "udpstreaming",
				}, false),
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
				Type:        schema.TypeString,
				Description: "The service level monitoring class that this server should use, if any.",
				Optional:    true,
			},
			"so_nagle": {
				Type:        schema.TypeBool,
				Description: "Whether or not Nagle's algorithm should be used for TCP connections.",
				Optional:    true,
				Default:     false,
			},
			"ssl_client_cert_headers": {
				Type:        schema.TypeString,
				Description: "What HTTP headers the virtual server should add to each request to show the data in the client certificate.",
				Optional:    true,
				Default:     "none",
				ValidateFunc: validation.StringInSlice([]string{
					"all",
					"none",
					"simple",
				}, false),
			},
			"ssl_decrypt": {
				Type:        schema.TypeBool,
				Description: "Whether or not the virtual server should decrypt incoming SSL traffic.",
				Optional:    true,
				Default:     false,
			},
			"ssl_honor_fallback_scsv": {
				Type:        schema.TypeString,
				Description: " Whether or not the Fallback SCSV sent by TLS clients is honored by this virtual server. ",
				Optional:    true,
				Default:     "use_default",
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
					"use_default",
				}, false),
			},
			"transparent": {
				Type:        schema.TypeBool,
				Description: "Whether or not bound sockets should be configured for transparent proxying",
				Optional:    true,
				Default:     false,
			},

			"aptimizer": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
							Type:        schema.TypeSet,
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
				Computed: true,
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
							ValidateFunc: validation.IntBetween(1024, 16777216),
						},
						"max_server_buffer": {
							Type:         schema.TypeInt,
							Description:  "The amount of memory, in bytes, that the virtual server should use to store data returned by the server.",
							Optional:     true,
							Default:      65536,
							ValidateFunc: validation.IntBetween(1024, 16777216),
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

			"connection_errors": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_file": {
							Type:        schema.TypeString,
							Description: "The error message to be sent to the client when the traffic manager detects an internal or backend error for the virtual server.",
							Optional:    true,
							Default:     "Default",
						},
					},
				},
			},

			"cookie": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Description: "The way in which the traffic manager should rewrite the domain portion of any cookies set by a back-end web server.",
							Optional:    true,
							Default:     "no_rewrite",
							ValidateFunc: validation.StringInSlice([]string{
								"no_rewrite", "set_to_named", "set_to_request"}, false),
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
							Type:        schema.TypeString,
							Description: "Whether or not the traffic manager should modify the 'secure' tag of anycookies set by a back-end web server.",
							Optional:    true,
							Default:     "no_modify",
							ValidateFunc: validation.StringInSlice([]string{
								"no_modify", "set_secure", "unset_secure"}, false),
						},
					},
				},
			},

			"dns": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
							ValidateFunc: util.ValidateUDPSize,
						},
						"max_udpsize": {
							Type:         schema.TypeInt,
							Description:  "Maximum UDP answer size",
							Optional:     true,
							Default:      4096,
							ValidateFunc: util.ValidateUDPSize,
						},
						"rrset_order": {
							Type:         schema.TypeString,
							Description:  "Response record ordering.",
							Optional:     true,
							Default:      "fixed",
							ValidateFunc: validation.StringInSlice([]string{"cyclic", "fixed"}, false),
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
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source_port": {
							Type:         schema.TypeInt,
							Description:  "The source port to be used for active-mode FTP data connections. If 0, a random high port will be used.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidatePortNumber,
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
							ValidateFunc: util.ValidatePortNumber,
						},
						"port_range_low": {
							Type:         schema.TypeInt,
							Description:  "If non-zero, then this controls the lower bound of the port range to use for FTP data connections.",
							Optional:     true,
							Default:      0,
							ValidateFunc: util.ValidatePortNumber,
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
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compress_level": {
							Type:         schema.TypeInt,
							Description:  "The source port to be used for active-mode FTP data connections. If 0, a random high port will be used.",
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 9),
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
							ValidateFunc: validation.StringInSlice([]string{"delete", "ignore", "weaken", "wrap"}, false),
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
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chunk_overhead_forwarding": {
							Type:         schema.TypeString,
							Description:  "Handling of HTTP chunk overhead.",
							Optional:     true,
							Default:      "lazy",
							ValidateFunc: validation.StringInSlice([]string{"lazy", "eager"}, false),
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
							Type:        schema.TypeString,
							Description: "If the 'Location' header matches the 'location_regex' regular expression, rewrite the header with this pattern",
							Optional:    true,
							Default:     "if_host_matches",
							ValidateFunc: validation.StringInSlice([]string{
								"always", "if_host_matches", "never"}, false),
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
				Computed: true,
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
							ValidateFunc: validation.IntBetween(100, 16777206),
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
							ValidateFunc: validation.IntBetween(4096, 1048576),
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
							ValidateFunc: validation.IntBetween(16384, 16777215),
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
				Computed: true,
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
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						//always_flush field is stated to exist in brocade API documentation but does not appear to exist in corresponding API version
						/*
							"always_flush": {
								Type:        schema.TypeBool,
								Description: "Write log data to disk immediately, rather than buffering data.",
								Optional:    true,
								Default:     false,
							},
						*/
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
							Default:     `%h %l %u %t "%r" %s %b "%{Referer}i""%{User-agent}i""`,
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
				Computed: true,
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
				Computed: true,
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
				Computed: true,
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
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dangerous_requests": {
							Type:         schema.TypeString,
							Description:  " The action to take when a SIP request with body data arrives that should be routed to an external IP.",
							Optional:     true,
							Default:      "node",
							ValidateFunc: validation.StringInSlice([]string{"forbid", "forward", "node"}, false),
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
							ValidateFunc: validation.StringInSlice([]string{"full_gateway", "route", "sip_gateway"}, false),
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

			"smtp": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expect_starttls": {
							Type:        schema.TypeBool,
							Description: "Whether or not the traffic manager should expect the connection to start off in plain text and then upgrade to SSL using STARTTLS when handling SMTP traffic",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			"ssl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
							Type:        schema.TypeSet,
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
										ValidateFunc: validation.StringInSlice([]string{"off", "on", "strict"}, false),
									},
									"required": {
										Type:         schema.TypeString,
										Description:  "Whether we should do an OCSP check for this issuer",
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"none", "optional", "strict"}, false),
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
							ValidateFunc: validation.StringInSlice([]string{"dont_request", "request", "require"}, false),
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
						"server_cert_host_mapping": {
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
							ValidateFunc: validation.StringInSlice([]string{"use_default", "disabled", "enabled"}, false),
						},
						"ssl_support_ssl3": {
							Type:         schema.TypeString,
							Description:  "Whether or not SSLv3 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validation.StringInSlice([]string{"use_default", "disabled", "enabled"}, false),
						},
						"ssl_support_tls1": {
							Type:         schema.TypeString,
							Description:  "Whether or not TLSv1.0 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validation.StringInSlice([]string{"use_default", "disabled", "enabled"}, false),
						},
						"ssl_support_tls1_1": {
							Type:         schema.TypeString,
							Description:  "Whether or not TLSv1.1 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validation.StringInSlice([]string{"use_default", "disabled", "enabled"}, false),
						},
						"ssl_support_tls1_2": {
							Type:         schema.TypeString,
							Description:  "Whether or not TLSv1.2 is enabled for this virtual server",
							Optional:     true,
							Default:      "use_default",
							ValidateFunc: validation.StringInSlice([]string{"use_default", "disabled", "enabled"}, false),
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
				Computed: true,
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
							Default:     `%h %l %u %t "%r" %s %b "%{Referer}i" "%{User-agent}i"`,
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
							ValidateFunc: validation.IntBetween(480, 65535),
						},
					},
				},
			},

			"tcp": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_close": {
							Type:        schema.TypeBool,
							Description: "If set to Yes the traffic manager will send the client FIN to the back-end server and wait for a server response instead of closing the connection immediately. ",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			"udp": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
				Computed: true,
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

func basicKeys() []string {
	return []string{
		"add_cluster_ip", "add_x_forwarded_for", "add_x_forwarded_proto", "autodetect_upgrade_headers",
		"bandwidth_class", "close_with_rst", "completionrules", "connect_timeout", "enabled",
		"ftp_force_server_secure", "glb_services", "listen_on_any", "listen_on_hosts", "listen_on_traffic_ips",
		"note", "pool", "port", "protection_class", "protocol", "request_rules", "response_rules",
		"slm_class", "so_nagle", "ssl_client_cert_headers", "ssl_decrypt", "ssl_honor_fallback_scsv",
		"transparent",
	}
}

func getSection(d *schema.ResourceData, sectionName string, properties map[string]interface{}, keys []string) error {
	m, err := util.GetAttributesToMap(d, keys)
	if err != nil {
		log.Println("Error getting section ", sectionName, err)
		return err
	}
	properties[sectionName] = m
	return nil
}

func sectionName(name string) string {
	if name == "vs_connection" {
		return "connection"
	}
	if name == "connection" {
		return "vs_connection"
	}
	return name
}

func resourceVirtualServerSet(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	name := d.Get("name").(string)

	res := make(map[string]interface{})
	pros := make(map[string]interface{})

	getSection(d, "basic", pros, basicKeys())

	for _, section := range []string{
		"aptimizer", "vs_connection", "connection_errors", "cookie",
		"dns", "ftp", "gzip", "http", "http2", "kerberos_protocol_transition",
		"log", "recent_connections", "request_tracing", "rtsp", "sip", "smtp",
		"ssl", "syslog", "tcp", "udp", "web_cache",
	} {
		if d.HasChange(section) {
			pros[sectionName(section)] = d.Get(section).([]interface{})[0]
		}
	}

	res["properties"] = pros
	util.TraverseMapTypes(res)
	err := client.Set("virtual_servers", name, res, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Virtual Server error whilst creating %s: %s", name, err)
	}
	d.SetId(name)

	return resourceVirtualServerRead(d, m)
}

func resourceVirtualServerRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	res := make(map[string]interface{})

	client.WorkWithConfigurationResources()
	err := client.GetByName("virtual_servers", d.Id(), &res)
	if err != nil {
		if client.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM Virtual Server error whilst retrieving %s: %v", d.Id(), err)
	}

	props := res["properties"].(map[string]interface{})
	basic := props["basic"].(map[string]interface{})

	for _, key := range basicKeys() {
		d.Set(key, basic[key])
	}

	for _, section := range []string{
		"aptimizer", "connection", "connection_errors", "cookie", "dns", "ftp",
		"gzip", "http", "http2", "kerberos_protocol_transition", "log",
		"recent_connections", "request_tracing", "rtsp", "sip", "smtp",
		"ssl", "syslog", "tcp", "udp", "web_cache",
	} {
		set := make([]map[string]interface{}, 0)
		set = append(set, props[section].(map[string]interface{}))
		d.Set(sectionName(section), set)
	}

	return nil
}

func resourceVirtualServerDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("virtual_servers", d, m)
}
