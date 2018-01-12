package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sky-uk/terraform-provider-pulsevtm/pulsevtm"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pulsevtm.Provider})
}
