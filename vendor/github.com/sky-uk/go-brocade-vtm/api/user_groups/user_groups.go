package usergroups

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const userGroupsEndpoint = "/api/tm/3.8/config/active/user_groups/"

// NewGet : used to retrieve a user group
func NewGet(groupName string) *rest.BaseAPI {
	getUserGroupAPI := rest.NewBaseAPI(http.MethodGet, userGroupsEndpoint+groupName, nil, new(UserGroup), new(api.VTMError))
	return getUserGroupAPI
}

// NewGetAll : used to retrieve a list of user group
func NewGetAll() *rest.BaseAPI {
	getAllUserGroupAPI := rest.NewBaseAPI(http.MethodGet, userGroupsEndpoint, nil, new(UserGroupList), new(api.VTMError))
	return getAllUserGroupAPI
}

// NewPut : used to create or update a user group
func NewPut(groupName string, userGroup UserGroup) *rest.BaseAPI {
	putUserGroupAPI := rest.NewBaseAPI(http.MethodPut, userGroupsEndpoint+groupName, userGroup, new(UserGroup), new(api.VTMError))
	return putUserGroupAPI
}

// NewDelete : used to delete a user group
func NewDelete(groupName string) *rest.BaseAPI {
	deleteUserGroupAPI := rest.NewBaseAPI(http.MethodDelete, userGroupsEndpoint+groupName, nil, nil, new(api.VTMError))
	return deleteUserGroupAPI
}
