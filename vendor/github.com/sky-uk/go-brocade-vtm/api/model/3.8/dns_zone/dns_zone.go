package dnsZone

// DNSZone : contains the DNS zone
type DNSZone struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall location configuration
type Properties struct {
	Basic Basic `json:"basic,omitempty"`
}

// Basic : Basic contains the attributes for a DNS zone
type Basic struct {
	Origin   string `json:"origin,omitempty"`
	ZoneFile string `json:"zonefile,omitempty"`
}

// DNSZones : List of DNS Zones
type DNSZones struct {
	Children []ChildLocation `json:"children"`
}

// ChildLocation : location name and href
type ChildLocation struct {
	Name string `json:"name"`
	Href string `json:"href"`
}
