package monitor

// Monitor : main Monitor data structure
type Monitor struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall monitor configuration
type Properties struct {
	Basic  Basic  `json:"basic,omitempty"`
	HTTP   HTTP   `json:"http,omitempty"`
	RTSP   RTSP   `json:"rtsp,omitempty"`
	SCRIPT SCRIPT `json:"script,omitempty"`
	SIP    SIP    `json:"sip,omitempty"`
	TCP    TCP    `json:"tcp,omitempty"`
	UDP    UDP    `json:"udp,omitempty"`
}

// Basic : Basic monitor configration
type Basic struct {
	Delay        uint     `json:"delay,omitempty"`
	Failures     uint     `json:"failures,omitempty"`
	Type         string   `json:"type,omitempty"`
	Timeout      uint     `json:"timeout,omitempty"`
	UseSSL       *bool    `json:"use_ssl,omitempty"`
	Verbose      *bool    `json:"verbose,omitempty"`
	BackOff      *bool    `json:"back_off,omitempty"`
	Machine      string   `json:"machine,omitempty"`
	Note         string   `json:"note,omitempty"`
	Scope        string   `json:"scope,omitempty"`
	CanEditSSL   *bool    `json:"can_edit_ssl,omitempty"`
	CanUseSSL    *bool    `json:"can_use_ssl,omitempty"`
	EditableKeys []string `json:"editable_keys,omitempty"`
}

// HTTP : HTTP monitor set up
type HTTP struct {
	Authentication string `json:"authentication,omitempty"`
	BodyRegex      string `json:"body_regex,omitempty"`
	HostHeader     string `json:"host_header,omitempty"`
	URIPath        string `json:"path,omitempty"`
	StatusRegex    string `json:"status_regex,omitempty"`
}

// RTSP : RTSP monitor set up
type RTSP struct {
	BodyRegex   string `json:"body_regex,omitempty"`
	URIPath     string `json:"path,omitempty"`
	StatusRegex string `json:"status_regex,omitempty"`
}

// SCRIPT : SCRIPT monitor set up
type SCRIPT struct {
	Arguments []ArgumentIssue `json:"arguments,omitempty"`
	Program   string          `json:"program,omitempty"`
}

// ArgumentIssue : argyments broken down
type ArgumentIssue struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// SIP : SIP monitor set up
type SIP struct {
	BodyRegex   string `json:"body_regex,omitempty"`
	StatusRegex string `json:"status_regex,omitempty"`
	Transport   string `json:"transport,omitempty"`
}

// TCP : TCP monitor set up
type TCP struct {
	CloseString    string `json:"close_string,omitempty"`
	MaxResponseLen uint   `json:"max_response_len,omitempty"`
	ResponseRegex  string `json:"response_regex,omitempty"`
	WriteString    string `json:"write_string,omitempty"`
}

// UDP : UDP monitor set up
type UDP struct {
	AcceptAll *bool `json:"accept_all,omitempty"`
}

// MonitorsList : List of nodes monitored
type MonitorsList struct {
	Children []ChildMonitor `json:"children"`
}

// ChildMonitor : monitored node structure
type ChildMonitor struct {
	Name string `json:"name"`
	HRef string `json:"href"`
}
