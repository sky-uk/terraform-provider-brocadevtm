package glb

// GLB : data structure
type GLB struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall location configuration
type Properties struct {
	Basic Basic `json:"basic,omitempty"`
	Log   Log   `json:"log,omitempty"`
}

// Basic : Basic location
type Basic struct {
	Algorithm            string            `json:"algorithm,omitempty"`
	AllMonitorsNeeded    bool              `json:"all_monitors_needed"`
	AutoRecovery         bool              `json:"autorecovery"`
	ChainedAutoFailback  bool              `json:"chained_auto_failback"`
	ChainedLocationOrder []string          `json:"chained_location_order,omitempty"`
	DisableOnFailure     bool              `json:"disable_on_failure"`
	DNSSecKeys           []DNSSecKey       `json:"dnssec_keys,omitempty"`
	Domains              []string          `json:"domains,omitempty"`
	Enabled              bool              `json:"enabled"`
	GeoEffect            uint              `json:"geo_effect"`
	LastResortResponse   []string          `json:"last_resort_response,omitempty"`
	LocationDraining     []string          `json:"location_draining,omitempty"`
	LocationSettings     []LocationSetting `json:"location_settings,omitempty"`
	ReturnIPSOnFail      bool              `json:"return_ips_on_fail"`
	Rules                []string          `json:"rules,omitempty"`
	TTL                  int               `json:"ttl"`
}

// DNSSecKey : DNS Sec key struct
type DNSSecKey struct {
	Domain  string   `json:"domain,omitempty"`
	SSLKeys []string `json:"ssl_key,omitempty"`
}

// LocationSetting : settings for a location
type LocationSetting struct {
	Location string   `json:"location,omitempty"`
	Weight   uint     `json:"weight"`
	IPS      []string `json:"ips,omitempty"`
	Monitors []string `json:"monitors,omitempty"`
}

// Log : log configuration for a GLB
type Log struct {
	Enabled  bool   `json:"enabled"`
	Filename string `json:"filename,omitempty"`
	Format   string `json:"format,omitempty"`
}

// GlobalLoadBalancers : List of GLBs
type GlobalLoadBalancers struct {
	Children []ChildLocation `json:"children"`
}

// ChildLocation : location name and href
type ChildLocation struct {
	Name string `json:"name"`
	Href string `json:"href"`
}
