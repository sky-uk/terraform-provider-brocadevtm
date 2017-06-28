package trafficIpGroups

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// DeleteTrafficIPGroupAPI base object.
type DeleteTrafficIPGroupAPI struct {
	*api.BaseAPI
}

// NewDelete returns a new object of DeleteTrafficIPGroupAPI.
func NewDelete(tipg string) *DeleteTrafficIPGroupAPI {
	this := new(DeleteTrafficIPGroupAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/tm/3.8/config/active/traffic_ip_groups/"+tipg, nil, nil)
	return this
}
