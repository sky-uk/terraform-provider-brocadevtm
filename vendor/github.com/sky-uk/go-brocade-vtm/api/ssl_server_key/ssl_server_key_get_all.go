package sslServerKey

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetAllSSLServerKeys base object.
type GetAllSSLServerKeys struct {
	*api.BaseAPI
}

// NewGetAll returns a new object of GetAllSSLServerKeys.
func NewGetAll() *GetAllSSLServerKeys {
	this := new(GetAllSSLServerKeys)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/ssl/server_keys", nil, new(SSLServerKeysList))
	return this
}

// GetResponse returns ResponseObject of SSLServerKeysList.
func (gam GetAllSSLServerKeys) GetResponse() *SSLServerKeysList {
	return gam.ResponseObject().(*SSLServerKeysList)
}
