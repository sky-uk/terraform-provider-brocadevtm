package trafficIpGroups

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetAllTrafficIPGroupsAPI base object.
type GetAllTrafficIPGroupsAPI struct {
	*api.BaseAPI
}

// NewGetAll returns a new object of GetAllTrafficIPGroupsAPI.
func NewGetAll() *GetAllTrafficIPGroupsAPI {
	this := new(GetAllTrafficIPGroupsAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/traffic_ip_groups", nil, new(TrafficIPGroupList))
	return this
}

// GetResponse returns ResponseObject of GetAllTrafficIPGroupsAPI.
func (ga GetAllTrafficIPGroupsAPI) GetResponse() *TrafficIPGroupList {
	return ga.ResponseObject().(*TrafficIPGroupList)
}

// FilterByName : returns a monitor object if the monitor name matches.
func (trafficIPGroups TrafficIPGroupList) FilterByName(name string) *ChildTrafficIPGroup {
	var trafficIPGroup ChildTrafficIPGroup
	for _, childTrafficIPGroup := range trafficIPGroups.Children {
		if childTrafficIPGroup.Name == name {
			trafficIPGroup = childTrafficIPGroup
			break
		}
	}
	return &trafficIPGroup
}
