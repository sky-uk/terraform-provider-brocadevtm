package pool

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// CreatePoolAPI - Base Struct
type CreatePoolAPI struct {
	*api.BaseAPI
}

//NewCreate - Creates a new pool
func NewCreate(poolName string, pool Pool) *CreatePoolAPI {
	return execCreateUpdate(poolName, pool)
}

//NewUpdate - Placeholder to create
func NewUpdate(poolName string, pool Pool) *CreatePoolAPI {
	return execCreateUpdate(poolName, pool)
}

func execCreateUpdate(poolName string, pool Pool) *CreatePoolAPI {
	this := new(CreatePoolAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/tm/3.8/config/active/pools/"+poolName, pool, new(Pool))
	return this
}

// GetResponse - Returns the http call response
func (cp CreatePoolAPI) GetResponse() Pool {
	return *cp.ResponseObject().(*Pool)
}
