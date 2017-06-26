package sslServerKey

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// DeleteSSLServerKeyAPI : object used to call delete on monitor
type DeleteSSLServerKeyAPI struct {
	*api.BaseAPI
}

// NewDelete : returns a new DeleteSSLServerKeyAPI object
func NewDelete(name string) *DeleteSSLServerKeyAPI {
	this := new(DeleteSSLServerKeyAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/tm/3.8/config/active/ssl/server_keys/"+name, nil, nil)
	return this
}
