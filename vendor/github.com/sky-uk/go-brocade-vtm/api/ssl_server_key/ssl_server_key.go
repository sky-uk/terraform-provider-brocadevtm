package sslServerKey

import (
	"encoding/json"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const sslServerKeyEndpoint = "/api/tm/3.8/config/active/ssl/server_keys/"

// NewCreate : Create new SSLServerKey
func NewCreate(name string, key SSLServerKey) *rest.BaseAPI {
	sslServerKeyCreateAPI := rest.NewBaseAPI(http.MethodPut, sslServerKeyEndpoint+name, key, new(SSLServerKey), new(api.VTMError))
	return sslServerKeyCreateAPI
}

// NewGetAll returns a new object of GetAllSSLServerKeys.
func NewGetAll() *rest.BaseAPI {
	SSLServerKeyGetAllAPI := rest.NewBaseAPI(http.MethodGet, sslServerKeyEndpoint, nil, new(SSLServerKeysList), new(api.VTMError))
	return SSLServerKeyGetAllAPI
}

// NewGet : returns the SSLServerKey details
func NewGet(name string) *rest.BaseAPI {
	SSLServerKeyGetAPI := rest.NewBaseAPI(http.MethodGet, sslServerKeyEndpoint+name, nil, new(SSLServerKey), new(api.VTMError))
	return SSLServerKeyGetAPI
}

// NewUpdate : creates a new object of type UpdateMonitorAPI
func NewUpdate(name string, sslServerKey SSLServerKey) *rest.BaseAPI {
	sslServerKeyUpdateAPI := rest.NewBaseAPI(http.MethodPut, sslServerKeyEndpoint+name, sslServerKey, new(json.RawMessage), new(api.VTMError))
	return sslServerKeyUpdateAPI
}

// NewDelete : returns a new DeleteSSLServerKeyAPI object
func NewDelete(name string) *rest.BaseAPI {
	sslServerKeyDeleteAPI := rest.NewBaseAPI(http.MethodDelete, sslServerKeyEndpoint+name, nil, nil, new(api.VTMError))
	return sslServerKeyDeleteAPI
}
