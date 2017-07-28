package rule

// Rules : List of rules
type Rules struct {
	Children []ChildRule `json:"children"`
}

// ChildRule : contains a rule's name and href
type ChildRule struct {
	Name string `json:"name"`
	HRef string `json:"href"`
}

// TrafficScriptRule : contains the traffic script and the name of the rule
type TrafficScriptRule struct {
	Name   string
	Script string
}
