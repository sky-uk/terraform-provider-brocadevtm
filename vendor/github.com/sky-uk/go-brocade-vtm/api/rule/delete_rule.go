package rule

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// DeleteRuleAPI : object used to call delete on a rule
type DeleteRuleAPI struct {
	*api.BaseAPI
}

// NewDelete : used to delete a rule
func NewDelete(ruleName string) *DeleteRuleAPI {
	this := new(DeleteRuleAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/tm/3.8/config/active/rules/"+ruleName, nil, nil)
	return this
}
