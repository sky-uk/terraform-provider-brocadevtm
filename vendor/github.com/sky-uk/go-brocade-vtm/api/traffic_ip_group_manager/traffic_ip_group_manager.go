package trafficIpGroupManager

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const trafficIPGroupManagerEndpoint = "/api/tm/3.8/config/active/traffic_managers/"

// NewGetAll : Get a list of traffic managers
func NewGetAll() *rest.BaseAPI {
	getAllTrafficManagerAPI := rest.NewBaseAPI(http.MethodGet, trafficIPGroupManagerEndpoint, nil, new(TrafficManagerChildren), new(api.VTMError))
	return getAllTrafficManagerAPI
}
