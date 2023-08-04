# Changelog

## UNRELEASED

  * /

## June 20th 2022 (v1.9.0)

  * Add support for customer portal
  * Update Stripe SDK's version requirement


## January 30th 2021 (v1.8.0)

  * Make sure product prices are read properly


## November 15th 2020 (v1.7.0)

  * Add support tiered prices
  * Add support for free plans


## September 5th 2020 (v1.6.1)

  * Add support for `unit_amount_decimal` for prices


## August 8th 2020 (v1.6.0)

  * Add support for prices
  * Update dependencies


## June 28th 2020 (v1.5.0)

  * Add more consistency to plan creations
    * force a new plan when changing/setting plan_id
    * enforce using either `amount` or `amount_decimal`, not both
  * Add support for decimal pricing and transform_usage
  * Update dependencies
  * Adapt connect account handling in webhooks
  * Check existence of application in order to get connect-status
  * Force new resource if connect status changes
  * Improves docs


## February 1st 2020 (v1.4.0)

  * Add tax rate


## January 5th 2020

  * Set coupon ID as its code
  * Update `redeem_by` as it was set in the past
  * Add `product_id` to `stripe_product`
  * Fix the update webhook url
  * Update dependencies


## August 30th 2019

  * Support for `plan_id` in plan resources
  * Add tiers to plan resources


## August 1st 2019

  * Support for Terraform 0.12


## Jul 21st 2019

  * A change to a plan's amount, currency or interval will force a new resource


## Apr 17th 2019

  * Add example on how to import existing resources


## Feb 24th 2019

  * Add support for coupons!


## Feb 14th 2019

  * Add support for more attributes in plan, product and webhooks
  * Add support for webhook secrets


## Feb 3rd 2019

  * Set app info for user agent


## Jan 29th 2019

  * Add `stripe_webhook_endpoint` resource


## Feb 1st 2019

  * Ensure right params are tracked
  * Make products active by default
  * Make plans active by default


## Sep 3rd 2018

  * Add more documentation on what is supported
  * Add Statement Descriptor, UnitLabel & Active to Product
  * Make Product and Plan importable
  * Initial commit
