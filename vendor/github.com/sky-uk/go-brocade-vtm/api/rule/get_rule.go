package rule

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetRuleAPI base object.
type GetRuleAPI struct {
	*api.BaseAPI
}

// NewGetRule : returns a rule
func NewGetRule(ruleName string) *GetRuleAPI {
	this := new(GetRuleAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/rules/"+ruleName, nil, new(string))
	return this
}

// GetResponse returns the string representation of the traffic script
func (getRule *GetRuleAPI) GetResponse() string {
	return string(getRule.RawResponse())
}
