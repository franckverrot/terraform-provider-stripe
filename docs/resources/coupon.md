---
page_title: "stripe_coupon"
subcategory: ""
description: |-
  
---

# stripe_coupon

## Example Usage

```hcl
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

## Import

This is an example of the import command being applied to the resource named `stripe_coupon.mlk_day_coupon_25pc_off`. The resource ID is a coupon code.

```bash
$ terraform import stripe_coupon.mlk_day_coupon_25pc_off MLK_DAY
```

## Schema

### Required

- **code** (String) The unique identifier of the coupon.
- **duration** (String) One of `forever`, `once`, and `repeating`. Describes how long a customer who applies this coupon will get the discount.

### Optional

- **amount_off** (Number) Amount (in the `currency` specified) that will be taken off the subtotal of any invoices for this customer.
- **currency** (String) If `amount_off` has been set, the three-letter ISO code for the currency of the amount to take off.
- **duration_in_months** (Number) If `duration` is `repeating`, the number of months the coupon applies. Null if coupon `duration` is `forever` or `once`.
- **id** (String) The ID of this resource.
- **max_redemptions** (Number) Maximum number of times this coupon can be redeemed, in total, across all customers, before it is no longer valid.
- **metadata** (Map of String) Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.
- **name** (String) Name of the coupon displayed to customers on for instance invoices or receipts.
- **percent_off** (Number) Percent that will be taken off the subtotal of any invoices for this customer for the duration of the coupon. For example, a coupon with percent_off of 50 will make a $100 invoice $50 instead.
- **redeem_by** (String) Date after which the coupon can no longer be redeemed.

### Read-Only

- **created** (Number) Time at which the object was created. Measured in seconds since the Unix epoch.
- **livemode** (Boolean) Has the value `true` if the object exists in live mode or the value `false` if the object exists in test mode.
- **times_redeemed** (Number) Number of times this coupon has been applied to a customer.
- **valid** (Boolean) Taking account of the above properties, whether this coupon can still be applied to a customer.


