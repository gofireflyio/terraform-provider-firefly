package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/gofireflyio/terraform-provider-firefly/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// version will be set by the goreleaser configuration to appropriate value for the compiled binary
var version string = "dev"

// commit will be set by the goreleaser configuration to the appropriate value for the compiled binary
var commit string = ""

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/firefly/firefly",
		Debug:   debug,
	}

	// Setup logging for debug mode
	if debug {
		tflog.Info(ctx, "Starting Firefly provider in debug mode", map[string]any{
			"version": version,
			"commit":  commit,
		})
	}

	err := providerserver.Serve(ctx, provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
