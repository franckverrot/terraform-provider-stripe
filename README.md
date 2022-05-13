# Terraform Stripe Provider

This provider enables Stripe merchants to manage certain parts of their
Stripe infrastructure—products, plans, webhook endpoints—via Terraform.

**Example use cases**
* Create and update resources in a repeatable manner
* Clone resources across multiple Stripe accounts (e.g. different locales or brands)

Since Terraform 13 and the Terraform Registry, no need to download a release
or compile this on your own machine (unless you really want to.)  Just jump
to the [Basic Usage](#basic-usage) section and get going!

## Requirements

*	[Terraform](https://www.terraform.io/downloads.html) 0.11.x to 0.13.x
*	[Go](https://golang.org/doc/install) 1.8 to 1.14 (to build the provider plugin)


## Building The Provider

Clone repository anywhere:

```sh
$ git clone https://github.com/franckverrot/terraform-provider-stripe.git
```

Enter the provider directory and build the provider

```sh
$ cd terraform-provider-stripe
$ make compile
```

Or alternatively, to install it as a plugin, run

```sh
$ cd terraform-provider-stripe
$ make install
```

## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

### Basic Usage

Set an environment variable, `TF_VAR_stripe_api_token` to store your Stripe
API token. This helps ensure you do not accidentally commit this sensitive
token to your repository.

    export TF_VAR_stripe_api_token=<your token>

Your token is now accessible in your Terraform configuration as
`var.stripe_api_token`, and can be used to configure the provider.

The example below demonstrates the following operations:

  * create a product
  * create a plan for that product
  * create a webhook endpoint for a few events

```hcl
terraform {
  required_providers {
    stripe = {
      source = "franckverrot/stripe"
      version = "1.8.0"
    }
  }
}

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

- [x] [Products](https://stripe.com/docs/api/products)
  - [x] name
  - [x] type
  - [x] active (Default: true)
  - [x] attributes (list)
  - [x] metadata (map)
  - [x] statement descriptor
  - [x] unit label
- [x] [Prices](https://stripe.com/docs/api/prices)
  - [x] active (Default: true)
  - [x] currency
  - [x] metadata (map)
  - [x] nickname
  - [x] product
  - [x] recurring
  - [x] unit_amount
  - [x] billing_scheme
  - [x] unit_amount_decimal
  - [x] tiers (Stripe API doesn't provide the API to update this at the moment, so the deletion should be done via dashboard page)
  - [x] tiers mode
- [x] [Plans](https://stripe.com/docs/api/plans)
  - [x] active (Default: true)
  - [x] aggregate usage
  - [x] amount
  - [x] amount_decimal
  - [x] billing scheme (Default: per_unit)
  - [x] currency
  - [x] interval
  - [x] interval_count (Default: 1)
  - [x] metadata (map)
  - [x] nickname
  - [x] product
  - [x] tiers
  - [x] tiers mode
  - [x] transform_usage
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
- [x] [TaxRates](https://stripe.com/docs/api/tax_rates)
  - [x] code (aka `id`)
  - [x] active
  - [x] description
  - [x] display_name
  - [x] inclusive
  - [x] jurisdiction
  - [ ] DELETE API (Stripe API doesn't provide the API at the moment, so the deletion should be done via dashboard page)
  - Computed:
    - [x] created
    - [x] livemode

- [x] [Customer Portal](https://stripe.com/docs/api/customer_portal)
  - [x] business_profile
    - [x] headline
    - [x] privacy_policy_url
    - [x] terms_of_service_url
  - [x] features
    - [x] customer_update
      - [x] allowed_updates
    - [x] invoice_history
    - [x] payment_method_update
    - [x] subscription_cancel
    - [x] subscription_pause
    - [x] subscription_update
  - [x] default_return_url
  - [x] metadata


### Importing existing resources

Scenario: you create something manually and would like to start managing it
with Terraform instead.

This provider support a straightforward/naive import procedure, here's how
you could do it for a coupon.

First, import the resource:

```
$ terraform import stripe_coupon.mlk_day_coupon_25pc_off MLK_DAY

...
Before importing this resource, please create its configuration in the root module. For example:

resource "stripe_coupon" "mlk_day_coupon_25pc_off" {
  # (resource arguments)
}
```

Then after adding these lines to your Terraform file, a plan should result in:

```
$ terraform plan

...
-/+ stripe_coupon.mlk_day_coupon_25pc_off (new resource required)
      id:              "MLK_DAY" => <computed> (forces new resource)
      amount_off:      "4200" => "4200"
      code:            "" => "MLK_DAY"
      created:         "" => <computed>
      currency:        "usd" => "usd"
      duration:        "once" => "once"
      livemode:        "false" => <computed>
      max_redemptions: "1024" => "1024"
      metadata.%:      "2" => "2"
      metadata.mlk:    "<3" => "<3"
      metadata.sales:  "yes" => "yes"
      name:            "King Sales Event" => "King Sales Event"
      redeem_by:       "" => "2019-09-02T12:34:56-08:00" (forces new resource)
      times_redeemed:  "0" => <computed>
      valid:           "true" => <computed>
```

Some updates might require replacing existing resources with new ones.


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

Mozilla Public License Version 2.0 – Franck Verrot – Copyright 2018-2020
