---
page_title: "stripe_product"
subcategory: ""
description: |-
  
---

# stripe_product

## Example Usage

```hcl
resource "stripe_product" "my_product" {
  name = "My Product"
  type = "service"
}
```

## Schema

### Required

- **name** (String) The product’s name, meant to be displayable to the customer. Whenever this product is sold via a subscription, name will show up on associated invoice line item descriptions.
- **type** (String)

### Optional

- **active** (Boolean) Whether the product is currently available for purchase. Defaults to `true`.
- **attributes** (List of String)
- **id** (String) The ID of this resource.
- **metadata** (Map of String) Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.
- **product_id** (String) Unique identifier for the product.
- **statement_descriptor** (String) Extra information about a product which will appear on your customer’s credit card statement. In the case that multiple products are billed at once, the first statement descriptor will be used.
- **unit_label** (String) A label that represents units of this product in Stripe and on customers’ receipts and invoices. When set, this will be included in associated invoice line item descriptions.


