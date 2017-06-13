package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: brocadevtm.Provider})
}
