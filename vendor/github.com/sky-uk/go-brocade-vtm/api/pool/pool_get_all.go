package pool

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetAllPools - Base struct
type GetAllPools struct {
	*api.BaseAPI
}

// NewGetAll - Returns all pools
func NewGetAll() *GetAllPools {
	this := new(GetAllPools)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/pools", nil, new(LBPoolList))
	return this

}

// GetResponse returns ResponseObject of GetAllPools.
func (gap GetAllPools) GetResponse() *LBPoolList {
	return gap.ResponseObject().(*LBPoolList)
}
