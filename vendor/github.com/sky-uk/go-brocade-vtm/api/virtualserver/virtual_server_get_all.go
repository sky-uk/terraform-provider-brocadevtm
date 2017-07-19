package virtualserver

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetAllVirtualServers : returns the list of all virtual servers
type GetAllVirtualServers struct {
	*api.BaseAPI
}

// NewGetAll returns a new object of GetAllVirtualServers.
func NewGetAll() *GetAllVirtualServers {
	this := new(GetAllVirtualServers)
	this.BaseAPI = api.NewBaseAPI(
		http.MethodGet,
		"/api/tm/3.8/config/active/virtual_servers",
		nil,
		new(VirtualServersList),
	)
	return this
}

// GetResponse returns ResponseObject of GetAllVirtualServers.
func (gav GetAllVirtualServers) GetResponse() VirtualServersList {
	return *gav.ResponseObject().(*VirtualServersList)
}
