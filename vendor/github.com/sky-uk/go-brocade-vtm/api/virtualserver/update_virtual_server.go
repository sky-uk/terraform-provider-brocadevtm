package virtualserver

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// UpdateVirtualServerAPI : object we use to update a virtual server
type UpdateVirtualServerAPI struct {
	*api.BaseAPI
}

// NewUpdate : creates a new object of type UpdateVirtualServerAPI
func NewUpdate(name string, virtualServer VirtualServer) *UpdateVirtualServerAPI {
	this := new(UpdateVirtualServerAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/tm/3.8/config/active/virtual_servers/"+name, virtualServer, new(VirtualServer))
	return this
}

// GetResponse : returns the response object from UpdateVirtualServerAPI
func (updateVirtualServerAPI UpdateVirtualServerAPI) GetResponse() VirtualServer {
	return *updateVirtualServerAPI.ResponseObject().(*VirtualServer)
}
