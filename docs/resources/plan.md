---
page_title: "stripe_plan"
subcategory: ""
description: |-
  
---

# stripe_plan

## Example Usage

```hcl
resource "stripe_plan" "my_product_plan" {
  product  = "prod_abc"
  amount   = 12345
  interval = "month"
  currency = "usd"
}
```

## Import

This is an example of the import command being applied to the resource named `stripe_plan.my_product_plan`. The resource ID is a product ID.

```bash
$ terraform import stripe_plan.my_product_plan prod_abc
```

## Schema

### Required

- **currency** (String) Three-letter ISO currency code, in lowercase. Must be a supported currency.
- **interval** (String) The frequency at which a subscription is billed. One of `day`, `week`, `month` or `year`.
- **product** (String) The product whose pricing this plan determines.

### Optional

- **active** (Boolean) Whether the plan can be used for new purchases.
- **aggregate_usage** (String) Specifies a usage aggregation strategy for plans of `usage_type=metered`. Allowed values are `sum` for summing up all usage during a period, `last_during_period` for using the last usage record reported within a period, `last_ever` for using the last usage record ever (across period bounds) or `max` which uses the usage record with the maximum reported usage during a period. Defaults to `sum`.
- **amount** (Number) The unit amount in cents to be charged, represented as a whole integer if possible. Only set if `billing_scheme=per_unit`.
- **amount_decimal** (Number) The unit amount in cents to be charged, represented as a decimal string with at most 12 decimal places. Only set if `billing_scheme=per_unit`.
- **billing_scheme** (String) Describes how to compute the price per period. Either `per_unit` or `tiered`. `per_unit` indicates that the fixed amount (specified in `amount`) will be charged per unit in `quantity` (for plans with `usage_type=licensed`), or per unit of total usage (for plans with `usage_type=metered`). `tiered` indicates that the unit pricing will be computed using a tiering strategy as defined using the `tiers` and `tiers_mode` attributes.
- **id** (String) The ID of this resource.
- **interval_count** (Number) The number of intervals (specified in the `interval` attribute) between subscription billings. For example, `interval=month` and `interval_count=3` bills every 3 months.
- **metadata** (Map of String) Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.
- **nickname** (String) A brief description of the plan, hidden from customers.
- **plan_id** (String) Unique identifier for the plan.
- **tier** (Block List) (see [below for nested schema](#nestedblock--tier))
- **tiers_mode** (String) Defines if the tiering price should be `graduated` or `volume` based. In `volume`-based tiering, the maximum quantity within a period determines the per unit price. In `graduated` tiering, pricing can change as the quantity grows.
- **transform_usage** (Block List, Max: 1) (see [below for nested schema](#nestedblock--transform_usage))
- **trial_period_days** (Number) Default number of trial days when subscribing a customer to this plan using `trial_from_plan=true`.
- **usage_type** (String) Configures how the quantity per period should be determined. Can be either `metered` or `licensed`. `licensed` automatically bills the `quantity` set when adding it to a subscription. `metered` aggregates the total usage based on usage records. Defaults to `licensed`.

<a id="nestedblock--tier"></a>
### Nested Schema for `tier`

Optional:

- **flat_amount** (Number) Price for the entire tier.
- **flat_amount_decimal** (Number) Same as `flat_amount`, but contains a decimal value with at most 12 decimal places.
- **unit_amount** (Number) Per unit price for units relevant to the tier.
- **unit_amount_decimal** (Number) Same as `unit_amount`, but contains a decimal value with at most 12 decimal places.
- **up_to** (Number) Up to and including to this quantity will be contained in the tier.
- **up_to_inf** (Boolean)


<a id="nestedblock--transform_usage"></a>
### Nested Schema for `transform_usage`

Required:

- **divide_by** (Number) Divide usage by this number.
- **round** (String) After division, either round the result `up` or `down`.


