---
page_title: "stripe_price"
subcategory: ""
description: |-
  
---

# stripe_price

## Schema

### Required

- **currency** (String) Three-letter ISO currency code, in lowercase. Must be a supported currency.

### Optional

- **active** (Boolean) Whether the price can be used for new purchases.
- **billing_scheme** (String) Describes how to compute the price per period. Either `per_unit` or `tiered`. `per_unit` indicates that the fixed amount (specified in `unit_amount` or `unit_amount_decimal`) will be charged per unit in quantity (for prices with `usage_type=licensed`), or per unit of total usage (for prices with `usage_type=metered`). `tiered` indicates that the unit pricing will be computed using a tiering strategy as defined using the `tiers` and `tiers_mode` attributes.
- **id** (String) The ID of this resource.
- **metadata** (Map of String) Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.
- **nickname** (String) A brief description of the price, hidden from customers.
- **price_id** (String) Unique identifier for the price.
- **product** (String) The ID of the product this price is associated with.
- **recurring** (Map of String) The recurring components of a price such as `interval` and `usage_type`.
- **tier** (Block List) (see [below for nested schema](#nestedblock--tier))
- **tiers_mode** (String) Defines if the tiering price should be `graduated` or `volume` based. In `volume`-based tiering, the maximum quantity within a period determines the per unit price. In `graduated` tiering, pricing can change as the quantity grows.
- **unit_amount** (Number) The unit amount in cents to be charged, represented as a whole integer if possible. Only set if `billing_scheme=per_unit`.
- **unit_amount_decimal** (Number) The unit amount in cents to be charged, represented as a decimal string with at most 12 decimal places. Only set if `billing_scheme=per_unit`.

### Read-Only

- **created** (Number) Time at which the object was created. Measured in seconds since the Unix epoch.
- **livemode** (Boolean) Has the value `true` if the object exists in live mode or the value `false` if the object exists in test mode.

<a id="nestedblock--tier"></a>
### Nested Schema for `tier`

Optional:

- **flat_amount** (Number) Price for the entire tier.
- **flat_amount_decimal** (Number) Same as `flat_amount`, but contains a decimal value with at most 12 decimal places.
- **unit_amount** (Number) Per unit price for units relevant to the tier.
- **unit_amount_decimal** (Number) Same as `unit_amount`, but contains a decimal value with at most 12 decimal places.
- **up_to** (Number) Up to and including to this quantity will be contained in the tier.
- **up_to_inf** (Boolean)


