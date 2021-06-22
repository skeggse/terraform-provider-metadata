# Terraform metadata provider

This extremely simple provider implements a state-based storage option for storing values that
Terraform should durably retain, without having them associated with real resources. This is a
workaround for use-cases during full Terraform rollout, such as when you have resources that live in
an unsupported provider but need to depend on that data prior to implementing that provider.

# Terraform Provider Scaffolding

TODO: Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://www.terraform.io/docs/registry/providers/publishing.html) so that others can use it.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Using the provider

```tf
resource "metadata_value" "some-data" {
    update = var.tag_value != nil
    inputs = {
        # Ignored value if `update` is false.
        tag_value = var.tag_value
    }
}

resource "aws_vpc" "example" {
    ...

    # You want to sometimes update this value, but otherwise have it stay the same.
    tags = {
        example = metadata_value.some-data.outputs.tag_value
    }
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
$ make testacc
```
