package sslServerKey

import (
	"fmt"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetSSLServerKey base object.
type GetSSLServerKey struct {
	*api.BaseAPI
}

// String returns a string representation of the monitor
func (sslServerKey SSLServerKey) String() string {
	return fmt.Sprintf("SSLServerKey: %+v", sslServerKey.Properties)
}

// NewGet : returns the SSLServerKey details
func NewGet(name string) *GetSSLServerKey {
	this := new(GetSSLServerKey)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/ssl/server_keys/"+name, nil, new(SSLServerKey))
	return this
}

// GetResponse returns ResponseObject of GetOneSSLServerKey.
func (reqObj GetSSLServerKey) GetResponse() *SSLServerKey {
	return reqObj.ResponseObject().(*SSLServerKey)
}
