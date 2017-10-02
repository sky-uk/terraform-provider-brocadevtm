package dnsZone

import (
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

// DNSZoneEndpoint : DNS zone uri endpoint
const DNSZoneEndpoint = "/api/tm/3.8/config/active/dns_server/zones/"

// NewCreate : Create a new DNS zone
func NewCreate(dnsZoneName string, dnsZone DNSZone) *rest.BaseAPI {
	createDNSZoneAPI := rest.NewBaseAPI(http.MethodPut, DNSZoneEndpoint+dnsZoneName, dnsZone, new(DNSZone), new(api.VTMError))
	return createDNSZoneAPI
}

// NewGet : returns a DNS zone
func NewGet(dnsZoneName string) *rest.BaseAPI {
	getDNSZoneAPI := rest.NewBaseAPI(http.MethodGet, DNSZoneEndpoint+dnsZoneName, nil, new(DNSZone), new(api.VTMError))
	return getDNSZoneAPI
}

// NewGetAll : returns a list of DNS zones {
func NewGetAll() *rest.BaseAPI {
	getAllDNSZoneAPI := rest.NewBaseAPI(http.MethodGet, DNSZoneEndpoint, nil, new(DNSZones), new(api.VTMError))
	return getAllDNSZoneAPI
}

// NewUpdate : Update a DNS zone
func NewUpdate(dnsZoneName string, dnsZone DNSZone) *rest.BaseAPI {
	updateDNSZoneNameAPI := rest.NewBaseAPI(http.MethodPut, DNSZoneEndpoint+dnsZoneName, dnsZone, new(DNSZone), new(api.VTMError))
	return updateDNSZoneNameAPI
}

// NewDelete : used to delete a DNS zone
func NewDelete(dnsZoneName string) *rest.BaseAPI {
	deleteDNSZoneAPI := rest.NewBaseAPI(http.MethodDelete, DNSZoneEndpoint+dnsZoneName, nil, nil, new(api.VTMError))
	return deleteDNSZoneAPI
}
