package rule

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetAllRules : used for retrieving a list of rules and their href
type GetAllRules struct {
	*api.BaseAPI
}

// NewGetAll : returns a list of rules {
func NewGetAll() *GetAllRules {
	this := new(GetAllRules)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/rules/", nil, new(Rules))
	return this
}

// GetResponse returns the response object of GetAllRules
func (getAllRules GetAllRules) GetResponse() *Rules {
	return getAllRules.ResponseObject().(*Rules)
}
