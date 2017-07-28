package trafficIpGroupManager

// TrafficManagerChildren : A list of traffic manager children
type TrafficManagerChildren struct {
	Children []TrafficMangerChild `json:"children"`
}

// TrafficMangerChild : A traffic manager name and reference
type TrafficMangerChild struct {
	Name string `json:"name"`
	HREF string `json:"href"`
}
