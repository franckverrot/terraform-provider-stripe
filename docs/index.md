---
page_title: "Stripe Provider"
subcategory: ""
description: |-
  
---

# Stripe Provider

This provider enables Stripe merchants to manage certain parts of their
Stripe infrastructure—products, plans, webhook endpoints—via Terraform.

Example use cases:

* Create and update resources in a repeatable manner
* Clone resources across multiple Stripe accounts (e.g. different locales or brands)

Use the navigation to the left to read about the available resources.

## Example usage

```hcl
provider "stripe" {
  api_token = "sk_live_..."
}

resource "stripe_product" "my_product" {
# ...
}
```

### Environment variables

You can provide your API key via `STRIPE_API_TOKEN` environment variable, representing your Stripe API key. When using this method, you may omit the Stripe `provider` block entirely:

```hcl
resource "stripe_product" "my_product" {
# ...
}
```

Usage:

```bash
$ export STRIPE_API_TOKEN="sk_live_..."
$ terraform plan
```

## Schema

### Optional

- **api_token** (String) – Stripe API key from https://dashboard.stripe.com/apikeys
