package main

import (
	"context"
	"flag"
	"log"

	"github.com/form3tech-oss/terraform-provider-auth0/auth0"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: auth0.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "form3.tech/providers/auth0", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
