package main

import (
	"github.com/form3tech-oss/terraform-provider-auth0/auth0"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: auth0.Provider})
}
