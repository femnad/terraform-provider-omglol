package main

import (
	"context"
	"flag"
	"github.com/femnad/terraform-provider-omglol/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "https://registry.terraform.io/providers/femnad/omglol",
	}

	err := providerserver.Serve(context.Background(), provider.New(), opts)
	if err != nil {
		log.Fatal(err)
	}
}
