package sslServerKey

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// CreateSSLServerKeyAPI : Create Monitor API
type CreateSSLServerKeyAPI struct {
	*api.BaseAPI
}

// NewCreate : Create new SSLServerKey
func NewCreate(name string, requestPayload *SSLServerKey) *CreateSSLServerKeyAPI {
	this := new(CreateSSLServerKeyAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/tm/3.8/config/active/ssl/server_keys/"+name, requestPayload, new(string))
	return this
}

// GetResponse : get response object from create SSLServerKey api call.
func (cma CreateSSLServerKeyAPI) GetResponse() string {
	return cma.ResponseObject().(string)
}
