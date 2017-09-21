package userauthenticators

// UserAuthenticator : The over all UserAuthenticator data structure
type UserAuthenticator struct {
	Properties `json:"properties"`
}

// Properties : Properties contains the overall UserAuthenticator configuration
type Properties struct {
	Basic      Basic      `json:"basic"`
	LDAP       LDAP       `json:"ldap,omitempty"`
	Radius     Radius     `json:"radius,omitempty"`
	TACACSPlus TACACSPlus `json:"tacacs_plus,omitempty"`
}

// Basic : Properties contains the Basic UserAuthenticator configuration
type Basic struct {
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled"`
	Type        string `json:"type,omitempty"`
}

// LDAP : Data structure representing the LDAP access control system
type LDAP struct {
	BaseDN         string `json:"base_dn,omitempty"`
	BindDN         string `json:"bind_dn,omitempty"`
	DNMethod       string `json:"dn_method,omitempty"`
	FallbackGroup  string `json:"fallback_group,omitempty"`
	Filter         string `json:"filter,omitempty"`
	GroupAttribute string `json:"group_attribute,omitempty"`
	GroupField     string `json:"group_field,omitempty"`
	GroupFilter    string `json:"group_filter,omitempty"`
	Port           uint   `json:"port"`
	SearchDN       string `json:"search_dn,omitempty"`
	SearchPassword string `json:"search_password,omitempty"`
	Server         string `json:"server,omitempty"`
	Timeout        uint   `json:"timeout"`
}

// Radius : Data structure representing the Radius access control system
type Radius struct {
	FallbackGroup  string `json:"fallback_group,omitempty"`
	GroupAttribute uint   `json:"group_attribute"`
	GroupVendor    uint   `json:"group_vendor"`
	NasIdentifier  string `json:"nas_identifier,omitempty"`
	NasIPAddress   string `json:"nas_ip_address,omitempty"`
	Port           uint   `json:"port"`
	Secret         string `json:"secret,omitempty"`
	Server         string `json:"server,omitempty"`
	Timeout        uint   `json:"timeout"`
}

// TACACSPlus : Data structure representing the TACACSPlus access control system
type TACACSPlus struct {
	AuthType      string `json:"auth_type,omitempty"`
	FallbackGroup string `json:"fallback_group,omitempty"`
	GroupField    string `json:"group_field,omitempty"`
	GroupService  string `json:"group_service,omitempty"`
	Port          uint   `json:"port"`
	Secret        string `json:"secret,omitempty"`
	Server        string `json:"server,omitempty"`
	Timeout       uint   `json:"timeout"`
}

// UserAuthenticatorList : List of ChildUserAuthenticator
type UserAuthenticatorList struct {
	Children []ChildUserAuthenticator `json:"children"`
}

// ChildUserAuthenticator : UserAuthenticator name/ref pair returned by get all
type ChildUserAuthenticator struct {
	Name string `json:"name"`
	HRef string `json:"href"`
}
