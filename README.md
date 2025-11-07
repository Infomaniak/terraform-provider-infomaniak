[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Infomaniak_terraform-provider-infomaniak&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Infomaniak_terraform-provider-infomaniak)
[![Go Report Card](https://goreportcard.com/badge/github.com/Infomaniak/terraform-provider-infomaniak)](https://goreportcard.com/report/github.com/Infomaniak/terraform-provider-infomaniak)
[![License](https://img.shields.io/github/license/Infomaniak/terraform-provider-infomaniak)](LICENSE)
[![Release](https://img.shields.io/github/release/Infomaniak/terraform-provider-infomaniak.svg)](https://github.com/Infomaniak/terraform-provider-infomaniak/releases)
[![Terraform Registry](https://img.shields.io/badge/Terraform-Registry-blueviolet?logo=terraform)](https://registry.terraform.io/providers/Infomaniak/infomaniak/latest/docs)
[![OpenTofu Registry](https://img.shields.io/badge/OpenTofu-Registry-yellow?logo=opentofu)](https://search.opentofu.org/provider/infomaniak/infomaniak/latest)
[![Coverage](https://img.shields.io/endpoint?url=https://infomaniak.github.io/terraform-provider-infomaniak/main/coverage_mocked_acceptance_badge.json)](https://infomaniak.github.io/terraform-provider-infomaniak/main/coverage_mocked_acceptance.html)


# Infomaniak Terraform Provider

The Infomaniak Provider allows Terraform to manage Infomaniak's resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To override Terraform's provider installation and use the in-dev provider, fill your `~/.terraformrc` with the following :
You need to replace `$GOPATH/bin` with the actual location of the binary on your system.

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/infomaniak/infomaniak" = "$GOPATH/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Test the Provider

Two types of tests are available for the provider:

- Unit tests: Test specific functions
- Acceptance mocked tests: End-to-end tests using real Terraform files with mocked infrastructure calls

Run tests using the [makefile](./makefile):
```bash
make prepare-go
make prepare-env
make test-unit
```
