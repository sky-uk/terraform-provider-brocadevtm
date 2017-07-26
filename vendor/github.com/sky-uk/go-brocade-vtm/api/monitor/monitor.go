package monitor

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const monitorEndpoint = "/api/tm/3.8/config/active/monitors/"

// NewCreate : Create new monitor
func NewCreate(name string, monitor Monitor) *rest.BaseAPI {
	createMonitorAPI := rest.NewBaseAPI(http.MethodPut, monitorEndpoint+name, monitor, new(Monitor), new(api.VTMError))
	return createMonitorAPI
}

// NewDelete : returns a new DeleteMonitorAPI object
func NewDelete(name string) *rest.BaseAPI {
	deleteMonitorAPI := rest.NewBaseAPI(http.MethodDelete, monitorEndpoint+name, nil, nil, new(api.VTMError))
	return deleteMonitorAPI
}

// NewGetAll : returns a new object of GetAllMonitors.
func NewGetAll() *rest.BaseAPI {
	getAllMonitorsAPI := rest.NewBaseAPI(http.MethodGet, monitorEndpoint, nil, new(MonitorsList), new(api.VTMError))
	return getAllMonitorsAPI
}

// NewGet : returns the monitor details
func NewGet(name string) *rest.BaseAPI {
	getMonitorAPI := rest.NewBaseAPI(http.MethodGet, monitorEndpoint+name, nil, new(Monitor), new(api.VTMError))
	return getMonitorAPI
}

// NewUpdate : creates a new object of type UpdateMonitorAPI
func NewUpdate(name string, monitor Monitor) *rest.BaseAPI {
	monitorUpdateAPI := rest.NewBaseAPI(http.MethodPut, monitorEndpoint+name, monitor, new(Monitor), new(api.VTMError))
	return monitorUpdateAPI
}
