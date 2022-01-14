---
page_title: "stripe_webhook_endpoint"
subcategory: ""
description: |-
  
---

# stripe_webhook_endpoint

## Example Usage

```hcl
resource "stripe_webhook_endpoint" "my_endpoint" {
  url = "https://mydomain.example.com/webhook"

  enabled_events = [
    "charge.succeeded",
    "charge.failed",
    "source.chargeable",
  ]
}
```

## Schema

### Required

- **enabled_events** (List of String) The list of events to enable for this endpoint. `['*']` indicates that all events are enabled, except those that require explicit selection.
- **url** (String) The URL of the webhook endpoint.

### Optional

- **connect** (Boolean)
- **id** (String) The ID of this resource.

### Read-Only

- **secret** (String) The endpoint’s secret, used to generate webhook signatures. Only returned at creation.


