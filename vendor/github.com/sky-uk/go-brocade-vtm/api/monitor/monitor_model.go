package monitor

// Monitor : main Monitor data structure
type Monitor struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall monitor configuration
type Properties struct {
	Basic Basic `json:"basic,omitempty"`
	HTTP  HTTP  `json:"http,omitempty"`
}

// Basic : Basic monitor configration
type Basic struct {
	Delay    uint   `json:"delay,omitempty"`
	Failures uint   `json:"failures,omitempty"`
	Type     string `json:"type,omitempty"`
	Timeout  uint   `json:"timeout,omitempty"`
	UseSSL   *bool  `json:"use_ssl,omitempty"`
	Verbose  *bool  `json:"verbose,omitempty"`
}

// HTTP : HTTP monitor set up
type HTTP struct {
	Authentication string `json:"authentication,omitempty"`
	BodyRegex      string `json:"body_regex,omitempty"`
	HostHeader     string `json:"host_header,omitempty"`
	URIPath        string `json:"path,omitempty"`
	StatusRegex    string `json:"status_regex,omitempty"`
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
