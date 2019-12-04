variable "stripe_api_token" {} # populate this by exporting TF_VAR_stripe_api_token

provider "stripe" {
  api_token = "${var.stripe_api_token}"
}

resource "stripe_product" "my_product" {
  name = "My Product"
  type = "service"
}

resource "stripe_product" "my_product_with_id" {
  product_id = "my_product"
  name       = "My Product"
  type       = "service"
}

resource "stripe_plan" "my_product_plan" {
  product  = "${stripe_product.my_product.id}"
  amount   = 12345
  interval = "month" # day week month year
  currency = "usd"
}

resource "stripe_plan" "my_product_metered_plan" {
  product  = "${stripe_product.my_product.id}"
  interval = "month"
  currency = "usd"

  usage_type      = "metered"
  billing_scheme  = "tiered"
  tiers_mode      = "graduated"
  aggregate_usage = "max"

  tier {
    up_to       = 5
    unit_amount = 50
  }

  tier {
    up_to       = 15
    unit_amount = 35
  }

  tier {
    up_to_inf   = true
    unit_amount = 25
  }
}

resource "stripe_plan" "my_product_plan_with_id" {
  plan_id  = "my_plan"

  product  = "${stripe_product.my_product.id}"
  amount   = 3232
  interval = "month" # day week month year
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

output "webhook_secret" {
  sensitive = true
  value     = "${stripe_webhook_endpoint.my_endpoint.secret}"
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
  redeem_by       = "2024-09-02T12:34:56-08:00" # RFC3339
}
