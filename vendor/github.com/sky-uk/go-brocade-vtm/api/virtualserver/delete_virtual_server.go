package virtualserver

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// DeleteVirtualServerAPI : object used to call delete on monitor
type DeleteVirtualServerAPI struct {
	*api.BaseAPI
}

// NewDelete : returns a new DeleteVirtualServerAPI object
func NewDelete(name string) *DeleteVirtualServerAPI {
	this := new(DeleteVirtualServerAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/tm/3.8/config/active/virtual_servers/"+name, nil, nil)
	return this
}
