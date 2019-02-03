# Terraform Stripe Provider


## Requirements

*	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
*	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)


## Building The Provider

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-stripe`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-stripe
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-stripe
$ make build
```

## Using the provider

### Basic Usage

In order for Terraform to pick up your API token, start by exporting the
special environment variable that will be loaded automatically for you:

    export TF_VAR_stripe_api_token=<your token>


Now that your token is available in HCL files as `var.stripe_api_token`, you
can use it to configure the provider.  Here's an example that will:

  * Create a product
  * Create a plan for that product
  * Create a webhook endpoint for a few events

```hcl
provider "stripe" {
  api_token = "${var.stripe_api_token}"
}

resource "stripe_product" "my_product" {
  name = "My Product"
  type = "service"
}

resource "stripe_plan" "my_product_plan1" {
  product  = "${stripe_product.my_product.id}"
  amount   = 12345
  interval = "month"                           # day week month year
  currency = "usd"
}

resource "stripe_webhook_endpoint" "my_endpoint" {
  url = "https://mydomain.example.com/webhook"

  enabled_events = [
    "charge.succeeded",
    "charge.failed",
    "source.chargeable",
  ]
}
```

### Supported resources

- [Products](https://stripe.com/docs/api/service_products)
- [Plans](https://stripe.com/docs/api/plans)
- [Webhook Endpoints](https://stripe.com/docs/api/webhook_endpoints)


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-stripe
...
```


## License

Mozilla Public License Version 2.0 – Franck Verrot – Copyright 2018