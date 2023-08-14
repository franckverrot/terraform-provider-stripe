terraform {
  required_providers {
    stripe = {
      source = "franckverrot/stripe"
      version = "1.9.0"
    }
  }
}

variable "stripe_api_token" {} # populate this by exporting TF_VAR_stripe_api_token

provider "stripe" {
  api_token = var.stripe_api_token
}

resource "stripe_product" "my_product" {
  name = "My Product"
  type = "service"
}

resource "stripe_product" "free_product" {
  name = "My Free Product"
  type = "service"
}

resource "stripe_product" "my_product_with_id" {
  product_id = "my_product"
  name       = "My Product"
  type       = "service"
}

resource "stripe_plan" "my_product_plan" {
  product  = stripe_product.my_product.id
  amount   = 12345
  interval = "month" # day week month year
  currency = "usd"
}

resource "stripe_plan" "free_product_plan" {
  product  = stripe_product.free_product.id
  amount   = 0
  interval = "month"
  currency = "usd"
}

resource "stripe_plan" "my_product_metered_plan" {
  product  = stripe_product.my_product.id
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
  plan_id = "my_plan"

  product  = stripe_product.my_product.id
  amount   = 3232
  interval = "month" # day week month year
  currency = "usd"
}

resource "stripe_plan" "my_decimal_product_plan" {
  product        = stripe_product.my_product.id
  amount_decimal = 123.45
  interval       = "month" # day week month year
  currency       = "usd"
}

resource "stripe_plan" "my_transformed_product_plan" {
  product  = stripe_product.my_product.id
  amount   = 2401
  interval = "month" # day week month year
  currency = "usd"
  transform_usage {
    divide_by = 7
    round     = "down"
  }
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
  value     = stripe_webhook_endpoint.my_endpoint.secret
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

resource "stripe_tax_rate" "my_tax_rate" {
  active       = true
  percentage   = 21
  display_name = "Twenty-one percent tax rate"
  inclusive    = true
}

resource "stripe_price" "my_price" {
  active   = true
  currency = "usd" # lowercase
  metadata = {
    blm = "always"
  }
  nickname    = "my price"
  product     = stripe_product.my_product.id
  unit_amount = 1337
  recurring = {
    interval       = "month"
    interval_count = 1
    usage_type     = "licensed"
  }
  billing_scheme = "per_unit"
}

resource "stripe_price" "my_graduated_price" {
  active   = true
  currency = "usd" # lowercase
  metadata = {
    blm = "always"
  }
  nickname = "my graduated price"
  product  = stripe_product.my_product.id

  recurring = {
    interval       = "month"
    interval_count = 1
    usage_type     = "licensed"
  }

  billing_scheme = "tiered"
  tiers_mode     = "graduated"

  tier {
    up_to       = 10
    unit_amount = 10
  }

  tier {
    up_to       = 20
    unit_amount = 20
  }

  tier {
    up_to_inf   = true
    unit_amount = 100
  }
}

resource "stripe_customer_portal" "customer_portal" {

  business_profile {
    headline             = "Headline"
    terms_of_service_url = "https://terms-of-service-url.example"
    privacy_policy_url   = "https://privacy-policy-url.example"
  }

  features {

    customer_update {
      allowed_updates = ["email", "address", "shipping", "phone", "tax_id"]
      enabled         = true
    }

    invoice_history {
      enabled = true
    }

    payment_method_update {
      enabled = true
    }

    subscription_cancel {
      cancellation_reason {
        enabled = true
        options = ["too_expensive", "missing_features", "switched_service", "unused", "customer_service", "too_complex", "low_quality", "other"]
      }
      enabled            = true
      mode               = "at_period_end"
      proration_behavior = "none"

    }

    subscription_pause {
      enabled = true
    }

    subscription_update {
      default_allowed_updates = ["price", "quantity", "promotion_code"]
      enabled                 = true
      proration_behavior      = "none"

      product {
        id  = stripe_product.my_product.id
        prices = [stripe_price.my_price.id]
      }

    }

  }

  metadata = {
    key = "val"
  }

  default_return_url = "https://return.example"

}
