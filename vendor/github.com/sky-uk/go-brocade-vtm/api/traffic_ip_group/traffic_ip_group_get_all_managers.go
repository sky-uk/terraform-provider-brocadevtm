package trafficIpGroups

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetAllTrafficManagersAPI : base object.
type GetAllTrafficManagersAPI struct {
	*api.BaseAPI
}

// NewGetTrafficManagerList : Get a list of traffic managers
func NewGetTrafficManagerList() *GetAllTrafficManagersAPI {
	this := new(GetAllTrafficManagersAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/traffic_managers", nil, new(TrafficManagerChildren))
	return this
}

// GetResponse : get the response for GetAllTrafficManagersAPI
func (getAllTrafficManagers GetAllTrafficManagersAPI) GetResponse() TrafficManagerChildren {
	return *getAllTrafficManagers.ResponseObject().(*TrafficManagerChildren)
}
