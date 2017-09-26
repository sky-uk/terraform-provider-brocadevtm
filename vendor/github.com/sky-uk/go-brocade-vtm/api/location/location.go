package location

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

const locationEndpoint = "/api/tm/3.8/config/active/locations/"

// NewCreate : used to create a new location
func NewCreate(locationName string, location Location) *rest.BaseAPI {
	createLocationAPI := rest.NewBaseAPI(http.MethodPut, locationEndpoint+locationName, location, new(Location), new(api.VTMError))
	return createLocationAPI
}

// NewGetAll : used to retrieve a list of locations and their HRef
func NewGetAll() *rest.BaseAPI {
	getAllLocationAPI := rest.NewBaseAPI(http.MethodGet, locationEndpoint, nil, new(Locations), new(api.VTMError))
	return getAllLocationAPI
}

// NewGet - used to retrieve a location
func NewGet(locationName string) *rest.BaseAPI {
	getLocationAPI := rest.NewBaseAPI(http.MethodGet, locationEndpoint+locationName, nil, new(Location), new(api.VTMError))
	return getLocationAPI
}

// NewUpdate : used to update a location
func NewUpdate(locationName string, location Location) *rest.BaseAPI {
	updateLocationAPI := rest.NewBaseAPI(http.MethodPut, locationEndpoint+locationName, location, new(Location), new(api.VTMError))
	return updateLocationAPI
}

// NewDelete : used to delete a location
func NewDelete(locationName string) *rest.BaseAPI {
	deleteLocationAPI := rest.NewBaseAPI(http.MethodDelete, locationEndpoint+locationName, nil, nil, new(api.VTMError))
	return deleteLocationAPI
}
