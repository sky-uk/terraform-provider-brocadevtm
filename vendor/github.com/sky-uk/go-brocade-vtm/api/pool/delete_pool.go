package pool

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// DeletePoolAPI - Base Struct
type DeletePoolAPI struct {
	*api.BaseAPI
}

// NewDelete - Deletes a pool
func NewDelete(poolName string) *DeletePoolAPI {
	this := new(DeletePoolAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/tm/3.8/config/active/pools/"+poolName, nil, new(Pool))
	return this
}
