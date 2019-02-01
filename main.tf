variable "stripe_api_token" {} # populate this by exporting TF_VAR_api_token

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
