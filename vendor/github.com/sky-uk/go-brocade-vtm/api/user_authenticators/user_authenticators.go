package userauthenticators

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const userAuthenticatorsEndpoint = "/api/tm/3.8/config/active/user_authenticators/"

// NewGet : used to retrieve a user authenticator
func NewGet(authenticatorName string) *rest.BaseAPI {
	getUserAuthenticatorAPI := rest.NewBaseAPI(http.MethodGet, userAuthenticatorsEndpoint+authenticatorName, nil, new(UserAuthenticator), new(api.VTMError))
	return getUserAuthenticatorAPI
}

// NewGetAll : used to retrieve a list of user authenticators
func NewGetAll() *rest.BaseAPI {
	getAllUserAuthenticatorsAPI := rest.NewBaseAPI(http.MethodGet, userAuthenticatorsEndpoint, nil, new(UserAuthenticatorList), new(api.VTMError))
	return getAllUserAuthenticatorsAPI
}

// NewPut : used to create or update a user authenticator
func NewPut(authenticatorName string, userAuthenticator UserAuthenticator) *rest.BaseAPI {
	putUserAuthenticatorAPI := rest.NewBaseAPI(http.MethodPut, userAuthenticatorsEndpoint+authenticatorName, userAuthenticator, new(UserAuthenticator), new(api.VTMError))
	return putUserAuthenticatorAPI
}

// NewDelete : used to delete a user authenticator
func NewDelete(authenticatorName string) *rest.BaseAPI {
	deleteUserAuthenticatorAPI := rest.NewBaseAPI(http.MethodDelete, userAuthenticatorsEndpoint+authenticatorName, nil, nil, new(api.VTMError))
	return deleteUserAuthenticatorAPI
}
