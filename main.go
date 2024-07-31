package main

import (
	"context"
	"flag"
	"log"

	template "github.com/chrismarget/terraform-provider-template/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool
	var printVersion bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
	flag.Parse()

	err := providerserver.Serve(context.Background(), template.NewProvider, providerserver.ServeOpts{
		Address: "registry.terraform.io/chrismarget/template",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
