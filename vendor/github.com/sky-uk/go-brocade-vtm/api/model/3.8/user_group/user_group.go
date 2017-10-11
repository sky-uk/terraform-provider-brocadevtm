package userGroup

// UserGroup : main UserGroup data structure
type UserGroup struct {
	Properties Properties `json:"properties"`
}

// Properties : Properties contains the overall UserGroup configuration
type Properties struct {
	Basic Basic `json:"basic"`
}

// Basic : Properties contains the Basic UserGroup configuration
type Basic struct {
	Description        string       `json:"description,omitempty"`
	PasswordExpireTime uint         `json:"password_expire_time,omitempty"`
	Timeout            uint         `json:"timeout,omitempty"`
	Permissions        []Permission `json:"permissions,omitempty"`
}

// Permission : Permission data structure, linking a configuration element to an access level
type Permission struct {
	Name        string `json:"name"`
	AccessLevel string `json:"access_level"`
}

// UserGroups : List of UserGroups
type UserGroups struct {
	Children []ChildUserGroup `json:"children"`
}

// ChildUserGroup : UserGroup name/ref pair returned by get all
type ChildUserGroup struct {
	Name string `json:"name"`
	HRef string `json:"href"`
}
