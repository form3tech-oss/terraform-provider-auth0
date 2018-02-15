package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/terraform-provider-auth0/auth0"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: auth0.Provider})
}
