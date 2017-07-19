package monitor

import (
	"fmt"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// GetSingleMonitor base object.
type GetSingleMonitor struct {
	*api.BaseAPI
}

// FilterByName : returns a monitor object if the monitor name matches.
func (monitorList MonitorsList) FilterByName(name string) *ChildMonitor {
	var foundMonitor ChildMonitor
	for _, childMonitor := range monitorList.Children {
		if childMonitor.Name == name {
			foundMonitor = childMonitor
			break
		}
	}
	return &foundMonitor
}

// String returns a string representation of the monitor
func (monitor Monitor) String() string {
	return fmt.Sprintf("Monitor: %+v", monitor.Properties)
}

// NewGetSingle : returns the monitor details
func NewGetSingle(name string) *GetSingleMonitor {
	this := new(GetSingleMonitor)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/tm/3.8/config/active/monitors/"+name, nil, new(Monitor))
	return this
}

// GetResponse returns ResponseObject of GetSingleMonitor.
func (gsm GetSingleMonitor) GetResponse() Monitor {
	return *gsm.ResponseObject().(*Monitor)
}
