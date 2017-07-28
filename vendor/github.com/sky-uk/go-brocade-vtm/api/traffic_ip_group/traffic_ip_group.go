package trafficIpGroups

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const trafficIPGroupEndpoint = "/api/tm/3.8/config/active/traffic_ip_groups/"

// NewCreate : used to create a new traffic IP group
func NewCreate(name string, trafficIPGroup TrafficIPGroup) *rest.BaseAPI {
	trafficIPGroupCreateAPI := rest.NewBaseAPI(http.MethodPut, trafficIPGroupEndpoint+name, trafficIPGroup, new(TrafficIPGroup), new(api.VTMError))
	return trafficIPGroupCreateAPI
}

// NewGetAll : used to get a list of all traffic IP groups and their associated HRefs
func NewGetAll() *rest.BaseAPI {
	trafficIPGroupGetAllAPI := rest.NewBaseAPI(http.MethodGet, trafficIPGroupEndpoint, nil, new(TrafficIPGroupList), new(api.VTMError))
	return trafficIPGroupGetAllAPI
}

// NewGet : used to get a traffic IP group
func NewGet(name string) *rest.BaseAPI {
	trafficIPGroupGetAPI := rest.NewBaseAPI(http.MethodGet, trafficIPGroupEndpoint+name, nil, new(TrafficIPGroup), new(api.VTMError))
	return trafficIPGroupGetAPI
}

// NewUpdate : used to update an existing traffic IP group
func NewUpdate(name string, trafficIPGroup TrafficIPGroup) *rest.BaseAPI {
	trafficIPGroupUpdateAPI := rest.NewBaseAPI(http.MethodPut, trafficIPGroupEndpoint+name, trafficIPGroup, new(TrafficIPGroup), new(api.VTMError))
	return trafficIPGroupUpdateAPI
}

// NewDelete : used to delete an existing traffic IP group
func NewDelete(name string) *rest.BaseAPI {
	trafficIPGroupDeleteAPI := rest.NewBaseAPI(http.MethodDelete, trafficIPGroupEndpoint+name, nil, nil, new(api.VTMError))
	return trafficIPGroupDeleteAPI
}
