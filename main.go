// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/services/dbaas"
	"terraform-provider-infomaniak/internal/services/domain"
	"terraform-provider-infomaniak/internal/services/kaas"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the kaas with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/Infomaniak/infomaniak",
		Debug:   debug,
	}

	// Register resources
	kaas.Register()
	domain.Register()
	dbaas.Register()

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
