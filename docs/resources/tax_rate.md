---
page_title: "stripe_tax_rate"
subcategory: ""
description: |-
  
---

# stripe_tax_rate

## Schema

### Required

- **active** (Boolean) Defaults to `true`. When set to `false`, this tax rate cannot be used with new applications or Checkout Sessions, but will still work for subscriptions and invoices that already have it set.
- **display_name** (String) The display name of the tax rates as it will appear to your customer on their receipt email, PDF, and the hosted invoice page.
- **inclusive** (Boolean) This specifies if the tax rate is inclusive or exclusive.
- **percentage** (Number) This represents the tax rate percent out of 100.

### Optional

- **description** (String) An arbitrary string attached to the tax rate for your internal use only. It will not be visible to your customers.
- **id** (String) The ID of this resource.
- **jurisdiction** (String) The jurisdiction for the tax rate. You can use this label field for tax reporting purposes. It also appears on your customer’s invoice.
- **metadata** (Map of String) Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.

### Read-Only

- **created** (Number) Time at which the object was created. Measured in seconds since the Unix epoch.
- **livemode** (Boolean) Has the value `true` if the object exists in live mode or the value `false` if the object exists in test mode.


