# Terraform Stripe Provider

This provider enables Stripe merchants to manage certain parts of their Stripe infrastructure—products, plans, webhook 
endpoints—via Terraform.

**Example use cases**
* Create and update resources in a repeatable manner
* Migrate resources from test mode to live mode
* Clone resources across multiple Stripe accounts (e.g. different locales or brands)

## Requirements

*	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
*	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)


## Building The Provider

Clone repository to: `$GOPATH/src/github.com/franckverrot/terraform-provider-stripe`

```sh
$ mkdir -p $GOPATH/src/github.com/franckverrot; cd $GOPATH/src/github.com/franckverrot
$ git clone git@github.com:franckverrot/terraform-provider-stripe
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/franckverrot/terraform-provider-stripe
$ make build
```

Or alternatively, to install it as a plugin, run

```sh
$ cd $GOPATH/src/github.com/franckverrot/terraform-provider-stripe
$ make install
```

## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

### Basic Usage

Set an environment variable, `TF_VAR_stripe_api_token` to store your Stripe API token. This helps ensure you 
do not accidentally commit this sensitive token to your repository.

    export TF_VAR_stripe_api_token=<your token>

Your token is now accessible in your Terraform configuration as `var.stripe_api_token`, and can be used to 
configure the provider.

The example below demonstrates the following operations:

  * create a product
  * create a plan for that product
  * create a webhook endpoint for a few events

```hcl
provider "stripe" {
  # NOTE: This is populated from the `TF_VAR_stripe_api_token` environment variable.
  api_token = "${var.stripe_api_token}"
}

resource "stripe_product" "my_product" {
  name = "My Product"
  type = "service"
}

resource "stripe_plan" "my_product_plan1" {
  product  = "${stripe_product.my_product.id}"
  amount   = 12345
  interval = "month"                           # Options: day week month year
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

resource "stripe_coupon" "mlk_day_coupon_25pc_off" {
  code     = "MLK_DAY"
  name     = "King Sales Event"
  duration = "once"

  amount_off = 4200
  currency   = "usd" # lowercase

  metadata = {
    mlk   = "<3"
    sales = "yes"
  }

  max_redemptions = 1024
  redeem_by       = "2019-09-02T12:34:56-08:00" # RFC3339, in the future
}
```

### Supported resources

- [x] [Products](https://stripe.com/docs/api/service_products)
  - [x] name
  - [x] type
  - [x] active (Default: true)
  - [x] attributes (list)
  - [x] metadata (map)
  - [x] statement descriptor
  - [x] unit label
- [x] [Plans](https://stripe.com/docs/api/plans)
  - [x] active (Default: true)
  - [x] aggregate usage
  - [x] amount
  - [x] billing scheme (Default: per_unit)
  - [x] currency
  - [x] interval
  - [x] interval_count (Default: 1)
  - [x] metadata (map)
  - [x] nickname
  - [x] product
  - [ ] tiers
  - [x] tiers mode
  - [ ] transform_usage
  - [x] trial period days
  - [x] usage type (Default: licensed)
- [x] [Webhook Endpoints](https://stripe.com/docs/api/webhook_endpoints)
  - [x] url
  - [x] enabled_events (list)
  - [x] connect
  - Computed:
    - secret
- [x] [Coupons](https://stripe.com/docs/api/coupons)
  - [x] code (aka `id`)
  - [x] name
  - [x] amount off
    - [x] currency
  - [x] percent off
  - [x] duration
    - [x] duration_in_months
  - [x] max redemptions
  - [x] metadata
  - [x] redeem by (should be RC3339-compliant)
  - Computed:
    - [x] valid
    - [x] created
    - [x] livemode
    - [x] times redeemed


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