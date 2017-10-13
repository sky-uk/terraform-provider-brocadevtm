package location

// Location data structure
type Location struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall location configuration
type Properties struct {
	Basic Basic `json:"basic,omitempty"`
}

// Basic : Basic location
type Basic struct {
	ID        uint    `json:"id,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Note      string  `json:"note,omitempty"`
	Type      string  `json:"type,omitempty"`
}

// Locations : List of locations
type Locations struct {
	Children []ChildLocation `json:"children"`
}

// ChildLocation : location name and href
type ChildLocation struct {
	Name string `json:"name"`
	Href string `json:"href"`
}
