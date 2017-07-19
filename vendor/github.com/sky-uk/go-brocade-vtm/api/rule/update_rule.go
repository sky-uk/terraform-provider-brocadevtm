package rule

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// UpdateRuleAPI : update rule API object
type UpdateRuleAPI struct {
	*api.BaseAPI
}

// NewUpdate : Update a rule
func NewUpdate(ruleName string, trafficScript []byte) *UpdateRuleAPI {
	this := new(UpdateRuleAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/tm/3.8/config/active/rules/"+ruleName, trafficScript, nil)
	return this
}
