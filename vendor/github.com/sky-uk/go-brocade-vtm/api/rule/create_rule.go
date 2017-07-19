package rule

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// CreateRuleAPI : create rule API object
type CreateRuleAPI struct {
	*api.BaseAPI
}

// NewCreate : Create a new rule
func NewCreate(ruleName string, trafficScript []byte) *CreateRuleAPI {
	this := new(CreateRuleAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/tm/3.8/config/active/rules/"+ruleName, trafficScript, nil)
	return this
}
