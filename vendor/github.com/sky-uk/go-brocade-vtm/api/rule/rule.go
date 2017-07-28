package rule

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const ruleEndpoint = "/api/tm/3.8/config/active/rules/"

// NewCreate : Create a new rule
func NewCreate(ruleName string, trafficScript []byte) *rest.BaseAPI {
	createRuleAPI := rest.NewBaseAPI(http.MethodPut, ruleEndpoint+ruleName, trafficScript, nil, new(api.VTMError))
	return createRuleAPI
}

// NewGet : returns a rule
func NewGet(ruleName string) *rest.BaseAPI {
	getRuleAPI := rest.NewBaseAPI(http.MethodGet, ruleEndpoint+ruleName, nil, new(string), new(api.VTMError))
	return getRuleAPI
}

// NewGetAll : returns a list of rules {
func NewGetAll() *rest.BaseAPI {
	this := rest.NewBaseAPI(http.MethodGet, ruleEndpoint, nil, new(Rules), new(api.VTMError))
	return this
}

// NewUpdate : Update a rule
func NewUpdate(ruleName string, trafficScript []byte) *rest.BaseAPI {
	updateRuleAPI := rest.NewBaseAPI(http.MethodPut, ruleEndpoint+ruleName, trafficScript, nil, new(api.VTMError))
	return updateRuleAPI
}

// NewDelete : used to delete a rule
func NewDelete(ruleName string) *rest.BaseAPI {
	this := rest.NewBaseAPI(http.MethodDelete, ruleEndpoint+ruleName, nil, nil, new(api.VTMError))
	return this
}
