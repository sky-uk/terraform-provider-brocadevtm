package glb

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const glbEndpoint = "/api/tm/3.8/config/active/glb_services/"

// NewCreate : used to create a new glb
func NewCreate(glbName string, glb GLB) *rest.BaseAPI {
	createGLBAPI := rest.NewBaseAPI(http.MethodPut, glbEndpoint+glbName, glb, new(GLB), new(api.VTMError))
	return createGLBAPI
}

// NewGetAll : used to retrieve a list of GLBs and their HRef
func NewGetAll() *rest.BaseAPI {
	getAllGLBAPI := rest.NewBaseAPI(http.MethodGet, glbEndpoint, nil, new(GlobalLoadBalancers), new(api.VTMError))
	return getAllGLBAPI
}

// NewGet - used to retrieve a GLB
func NewGet(glbName string) *rest.BaseAPI {
	getGLBAPI := rest.NewBaseAPI(http.MethodGet, glbEndpoint+glbName, nil, new(GLB), new(api.VTMError))
	return getGLBAPI
}

// NewUpdate : used to update a GLB
func NewUpdate(glbName string, glb GLB) *rest.BaseAPI {
	updateGLBAPI := rest.NewBaseAPI(http.MethodPut, glbEndpoint+glbName, glb, new(GLB), new(api.VTMError))
	return updateGLBAPI
}

// NewDelete : used to delete a GLB
func NewDelete(glbName string) *rest.BaseAPI {
	deleteGLBAPI := rest.NewBaseAPI(http.MethodDelete, glbEndpoint+glbName, nil, nil, new(api.VTMError))
	return deleteGLBAPI
}
