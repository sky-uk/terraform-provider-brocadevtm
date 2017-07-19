package trafficIpGroups

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// CreateUpdateTrafficIPGroupAPI base object.
type CreateUpdateTrafficIPGroupAPI struct {
	*api.BaseAPI
}

// execCreateUpdate returns a new object of CreateTrafficIPGroupA.
func execCreateUpdate(name string, tipg TrafficIPGroup) *CreateUpdateTrafficIPGroupAPI {
	this := new(CreateUpdateTrafficIPGroupAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/tm/3.8/config/active/traffic_ip_groups/"+name, tipg, new(TrafficIPGroup))
	return this
}

// NewCreate returns a new object of CreateUpdateTrafficIPGroupAPI.
func NewCreate(name string, requestPayload TrafficIPGroup) *CreateUpdateTrafficIPGroupAPI {
	return execCreateUpdate(name, requestPayload)
}

// NewUpdate returns a new object of CreateUpdateTrafficIPGroupAPI.
func NewUpdate(name string, requestPayload TrafficIPGroup) *CreateUpdateTrafficIPGroupAPI {
	return execCreateUpdate(name, requestPayload)
}

// GetResponse returns ResponseObject of CreateTrafficIPGroupAPI.
func (ga CreateUpdateTrafficIPGroupAPI) GetResponse() TrafficIPGroup {
	return *ga.ResponseObject().(*TrafficIPGroup)
}
