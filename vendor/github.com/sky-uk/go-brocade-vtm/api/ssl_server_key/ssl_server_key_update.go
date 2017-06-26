package sslServerKey

import (
	"encoding/json"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// UpdateSSLServerKeyAPI  : object we use to update a monitor
type UpdateSSLServerKeyAPI struct {
	*api.BaseAPI
}

// NewUpdate : creates a new object of type UpdateMonitorAPI
func NewUpdate(name string, sslServerKey SSLServerKey) *UpdateSSLServerKeyAPI {
	this := new(UpdateSSLServerKeyAPI)
	requestPayLoad := new(SSLServerKey)
	requestPayLoad.Properties.Basic = sslServerKey.Properties.Basic
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/tm/3.8/config/active/ssl/server_keys/"+name, requestPayLoad, new(json.RawMessage))
	return this
}

// GetResponse : returns the response object from UpdateMonitorAPI
func (updateMonitorAPI UpdateSSLServerKeyAPI) GetResponse() string {
	return updateMonitorAPI.ResponseObject().(string)
}
