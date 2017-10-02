package dnsZoneFile

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

// DNSZoneFileEndpoint : DNS zone file uri endpoint
const DNSZoneFileEndpoint = "/api/tm/3.8/config/active/dns_server/zone_files/"

// NewCreate : Create a new DNS zone file
func NewCreate(dnsZoneFileName string, dnsZoneFile []byte) *rest.BaseAPI {
	createDNSZoneFileAPI := rest.NewBaseAPI(http.MethodPut, DNSZoneFileEndpoint+dnsZoneFileName, dnsZoneFile, nil, new(api.VTMError))
	return createDNSZoneFileAPI
}

// NewGet : returns a DNS zone file
func NewGet(dnsZoneFileName string) *rest.BaseAPI {
	getDNSZoneFileAPI := rest.NewBaseAPI(http.MethodGet, DNSZoneFileEndpoint+dnsZoneFileName, nil, new(string), new(api.VTMError))
	return getDNSZoneFileAPI
}

// NewGetAll : returns a list of DNS zone files {
func NewGetAll() *rest.BaseAPI {
	getAllDNSZoneFileAPI := rest.NewBaseAPI(http.MethodGet, DNSZoneFileEndpoint, nil, new(DNSZoneFiles), new(api.VTMError))
	return getAllDNSZoneFileAPI
}

// NewUpdate : Update a DNS zone file
func NewUpdate(dnsZoneFileName string, dnsZoneFile []byte) *rest.BaseAPI {
	updateDNSZoneFileNameAPI := rest.NewBaseAPI(http.MethodPut, DNSZoneFileEndpoint+dnsZoneFileName, dnsZoneFile, nil, new(api.VTMError))
	return updateDNSZoneFileNameAPI
}

// NewDelete : used to delete a DNS zone file
func NewDelete(dnsZoneFileName string) *rest.BaseAPI {
	deleteDNSZoneFileAPI := rest.NewBaseAPI(http.MethodDelete, DNSZoneFileEndpoint+dnsZoneFileName, nil, nil, new(api.VTMError))
	return deleteDNSZoneFileAPI
}
