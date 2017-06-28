package trafficIpGroups

// TrafficIPGroupList : List of ChildTrafficIPGroup data structure
type TrafficIPGroupList struct {
	Children []ChildTrafficIPGroup `json:"children"`
}

// ChildTrafficIPGroup : Child of TrafficIPGroupList, contains name and href of an Traffic IP Group
type ChildTrafficIPGroup struct {
	Name string `json:"name"`
	HRef string `json:"href"`
}

// TrafficIPGroup : Main traffic IP group data structure
type TrafficIPGroup struct {
	Properties Properties `json:"properties"`
}

// Properties : Stored within TrafficIPGroup, contains Basic data structure
type Properties struct {
	Basic Basic `json:"basic"`
}

// IPMapping : IPMapping data structure for use within Basic data structure
type IPMapping struct {
	IP             string `json:"ip"`
	TrafficManager string `json:"traffic_manager"`
}

// Basic : Stored within Properties, contains attributes of an TrafficIPGroup
type Basic struct {
	Enabled                      *bool       `json:"enabled,omitempty"`
	HashSourcePort               *bool       `json:"hash_source_port,omitempty"`
	IPAssignmentMode             string      `json:"ip_assignment_mode,omitempty"`
	IPMapping                    []IPMapping `json:"ip_mapping,omitempty"`
	IPAddresses                  []string    `json:"ipaddresses,omitempty"`
	KeepTogether                 *bool       `json:"keeptogether,omitempty"`
	Location                     int         `json:"location,omitempty"`
	Machines                     []string    `json:"machines,omitempty"`
	Mode                         string      `json:"mode,omitempty"`
	Multicast                    string      `json:"multicast,omitempty"`
	Note                         string      `json:"note,omitempty"`
	RhiBgpMetricBase             uint        `json:"rhi_bgp_metric_base,omitempty"`
	RhiBgpPassiveMetricOffset    uint        `json:"rhi_bgp_passive_metric_offset,omitempty"`
	RhiOspfv2MetricBase          uint        `json:"rhi_ospfv2_metric_base,omitempty"`
	RhiOspfv2PassiveMetricOffset uint        `json:"rhi_ospfv2_passive_metric_offset,omitempty"`
	RhiProtocols                 string      `json:"rhi_protocols,omitempty"`
	Slaves                       []string    `json:"slaves,omitempty"`
}

// TrafficManagerChildren : A list of traffic manager children
type TrafficManagerChildren struct {
	Children []TrafficMangerChild `json:"children"`
}

// TrafficMangerChild : A traffic manager name and reference
type TrafficMangerChild struct {
	Name string `json:"name"`
	HREF string `json:"href"`
}
