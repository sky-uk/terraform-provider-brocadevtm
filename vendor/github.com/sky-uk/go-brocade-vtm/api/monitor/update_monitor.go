package monitor

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// UpdateMonitorAPI : object we use to update a monitor
type UpdateMonitorAPI struct {
	*api.BaseAPI
}

// NewUpdate : creates a new object of type UpdateMonitorAPI
func NewUpdate(name string, monitor Monitor) *UpdateMonitorAPI {
	this := new(UpdateMonitorAPI)
	requestPayLoad := new(Monitor)
	requestPayLoad.Properties.Basic = monitor.Properties.Basic
	requestPayLoad.Properties.HTTP = monitor.Properties.HTTP
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/tm/3.8/config/active/monitors/"+name, requestPayLoad, new(Monitor))
	return this
}

// GetResponse : returns the response object from UpdateMonitorAPI
func (updateMonitorAPI UpdateMonitorAPI) GetResponse() Monitor {
	return *updateMonitorAPI.ResponseObject().(*Monitor)
}
