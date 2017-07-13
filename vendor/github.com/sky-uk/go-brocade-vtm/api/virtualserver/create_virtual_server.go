// Package virtualserver : Virtual server configuration handling
package virtualserver

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// CreateVirtualServerAPI : Create VirtualServer API
type CreateVirtualServerAPI struct {
	*api.BaseAPI
}

// NewCreate : Create new virtualServer
// Input:
//   virtualServerName : the name of the virtual server
//   virtualServer     : the configMap of the new virtual
//                       server
func NewCreate(virtualServerName string,
	virtualServer VirtualServer) *CreateVirtualServerAPI {

	this := new(CreateVirtualServerAPI)
	this.BaseAPI = api.NewBaseAPI(
		http.MethodPut,
		"/api/tm/3.8/config/active/virtual_servers/"+virtualServerName,
		virtualServer,
		new(interface{}),
	)
	return this
}

// GetResponse : get response object from created virtualServer
func (cvs CreateVirtualServerAPI) GetResponse() string {
	return cvs.ResponseObject().(string)
}
