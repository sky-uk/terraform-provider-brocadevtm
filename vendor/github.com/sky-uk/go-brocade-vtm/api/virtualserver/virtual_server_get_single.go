package virtualserver

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetSingleVirtualServer - Base Struct
type GetSingleVirtualServer struct {
	*api.BaseAPI
}

// NewGetSingle - Returns a single virtual server
func NewGetSingle(virtualServerName string) *GetSingleVirtualServer {
	this := new(GetSingleVirtualServer)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/virtual_servers/"+virtualServerName, nil, new(VirtualServer))
	return this

}

// GetResponse returns ResponseObject of GetSingleVirtualServer
func (gsv GetSingleVirtualServer) GetResponse() *VirtualServer {
	return gsv.ResponseObject().(*VirtualServer)
}
