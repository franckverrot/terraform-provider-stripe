variable "stripe_api_token" {} # populate this by exporting TF_VAR_stripe_api_token

provider "stripe" {
  api_token = "${var.stripe_api_token}"
}

resource "stripe_product" "my_product" {
  name = "My Product"
  type = "service"
}

resource "stripe_plan" "my_product_plan" {
  product  = "${stripe_product.my_product.id}"
  amount   = 12345
  interval = "month"                           # day week month year
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
  redeem_by       = "2019-09-02T12:34:56-08:00" # RFC3339
}
