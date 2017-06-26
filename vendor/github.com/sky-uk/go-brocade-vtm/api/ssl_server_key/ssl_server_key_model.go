package sslServerKey

// A SSLServerKey is a Brocade trusted certificate.
type SSLServerKey struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall SSLServerKey configuration
type Properties struct {
	Basic Basic `json:"basic,omitempty"`
}

// Basic struct for SSlServerKey
type Basic struct {
	Note    string `json:"note,omitempty"`
	Private string `json:"private,omitempty"`
	Public  string `json:"public,omitempty"`
	Request string `json:"request,omitempty"`
}

// SSLServerKeysList : List of ssl server keys
type SSLServerKeysList struct {
	Children []ChildSSLServerKey `json:"children"`
}

// ChildSSLServerKey : ssl server keys child nodes
type ChildSSLServerKey struct {
	Name string `json:"name"`
	HRef string `json:"href"`
}
