package stripe

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceCustomerPortal() *schema.Resource {
	return &schema.Resource{
		Create: resourceStripeCustomerPortalCreate,
		Read:   resourceStripeCustomerPortalRead,
		Update: resourceStripeCustomerPortalUpdate,
		Delete: resourceStripeCustomerPortalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"business_profile": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"privacy_policy_url": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"terms_of_service_url": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"headline": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
			"features": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_update": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_updates": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"email", "address", "shipping", "phone", "tax_id"}, false),
										},
									},
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
							Optional: true,
						},
						"invoice_history": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
							Optional: true,
						},
						"payment_method_update": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
							Optional: true,
						},
						"subscription_cancel": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cancellation_reason": &schema.Schema{
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": &schema.Schema{
													Type:     schema.TypeBool,
													Required: true,
												},
												"options": &schema.Schema{
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type:         schema.TypeString,
														ValidateFunc: validation.StringInSlice([]string{"too_expensive", "missing_features", "switched_service", "unused", "customer_service", "too_complex", "low_quality", "other"}, false),
													},
												},
											},
										},
										Optional: true,
									},
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Required: true,
									},
									"mode": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"immediately", "at_period_end"}, false),
									},
									"proration_behavior": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"none", "create_prorations"}, false),
									},
								},
							},
							Optional: true,
						},
						"subscription_pause": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
							Optional: true,
						},
						"subscription_update": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_allowed_updates": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"price", "quantity", "promotion_code"}, false),
										},
									},
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Required: true,
									},
									"proration_behavior": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"none", "create_prorations", "always_invoice"}, false),
									},
									"product": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"prices": &schema.Schema{
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
										Optional: true,
									},
								},
							},
							Optional: true,
						},
					},
				},
				Required: true,
			},
			"default_return_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func expandBusinessProfile(businessProfileI []interface{}) *stripe.BillingPortalConfigurationBusinessProfileParams {
	businessProfile := &stripe.BillingPortalConfigurationBusinessProfileParams{}
	for _, v := range businessProfileI {
		businessProfileMap := v.(map[string]interface{})
		if val, ok := businessProfileMap["privacy_policy_url"]; ok {
			businessProfile.PrivacyPolicyURL = stripe.String(val.(string))
		}
		if val, ok := businessProfileMap["terms_of_service_url"]; ok {
			businessProfile.TermsOfServiceURL = stripe.String(val.(string))
		}
		if val, ok := businessProfileMap["headline"]; ok {
			businessProfile.Headline = stripe.String(val.(string))
		}
	}
	return businessProfile
}

func expandFeatures(featuresI []interface{}) *stripe.BillingPortalConfigurationFeaturesParams {
	features := &stripe.BillingPortalConfigurationFeaturesParams{}
	for _, v := range featuresI {
		featuresMap := v.(map[string]interface{})
		if val, ok := featuresMap["customer_update"]; ok {
			customerUpdate := &stripe.BillingPortalConfigurationFeaturesCustomerUpdateParams{}
			cu := val.([]interface{})
			for _, props := range cu {
				p := props.(map[string]interface{})
				if val, ok := p["allowed_updates"]; ok {
					enumsI := val.([]interface{})
					enums := []string{}
					for _, enum := range enumsI {
						enums = append(enums, enum.(string))
					}
					customerUpdate.AllowedUpdates = stripe.StringSlice(enums)
				}
				if val, ok := p["enabled"]; ok {
					customerUpdate.Enabled = stripe.Bool(val.(bool))
				}
			}
			features.CustomerUpdate = customerUpdate
		}

		if val, ok := featuresMap["invoice_history"]; ok {
			invoiceHistory := &stripe.BillingPortalConfigurationFeaturesInvoiceHistoryParams{}
			ih := val.([]interface{})
			for _, props := range ih {
				p := props.(map[string]interface{})
				if val, ok := p["enabled"]; ok {
					invoiceHistory.Enabled = stripe.Bool(val.(bool))
				}
			}
			features.InvoiceHistory = invoiceHistory
		}

		if val, ok := featuresMap["payment_method_update"]; ok {
			paymentMethodUpdate := &stripe.BillingPortalConfigurationFeaturesPaymentMethodUpdateParams{}
			pmu := val.([]interface{})
			for _, props := range pmu {
				p := props.(map[string]interface{})
				if val, ok := p["enabled"]; ok {
					paymentMethodUpdate.Enabled = stripe.Bool(val.(bool))
				}
			}
			features.PaymentMethodUpdate = paymentMethodUpdate
		}

		if val, ok := featuresMap["subscription_cancel"]; ok {
			subscriptionCancel := &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelParams{}
			sc := val.([]interface{})
			for _, props := range sc {
				p := props.(map[string]interface{})
				if val, ok := p["cancellation_reason"]; ok {
					subscriptionCancelReason := &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelCancellationReasonParams{}
					scr := val.([]interface{})
					for _, scrProps := range scr {
						scrP := scrProps.(map[string]interface{})
						if val, ok := scrP["options"]; ok {
							enumsI := val.([]interface{})
							enums := []string{}
							for _, enum := range enumsI {
								enums = append(enums, enum.(string))
							}
							subscriptionCancelReason.Options = stripe.StringSlice(enums)
						}

						if val, ok := scrP["enabled"]; ok {
							subscriptionCancelReason.Enabled = stripe.Bool(val.(bool))
						}
					}
					subscriptionCancel.CancellationReason = subscriptionCancelReason
				}

				if val, ok := p["enabled"]; ok {
					subscriptionCancel.Enabled = stripe.Bool(val.(bool))
				}

				if val, ok := p["mode"]; ok {
					subscriptionCancel.Mode = stripe.String(val.(string))
				}

				if val, ok := p["proration_behavior"]; ok {
					subscriptionCancel.ProrationBehavior = stripe.String(val.(string))
				}
			}
			features.SubscriptionCancel = subscriptionCancel
		}

		if val, ok := featuresMap["subscription_pause"]; ok {
			subscriptionPause := &stripe.BillingPortalConfigurationFeaturesSubscriptionPauseParams{}
			sp := val.([]interface{})
			for _, props := range sp {
				p := props.(map[string]interface{})
				if val, ok := p["enabled"]; ok {
					subscriptionPause.Enabled = stripe.Bool(val.(bool))
				}
			}
			features.SubscriptionPause = subscriptionPause
		}

		if val, ok := featuresMap["subscription_update"]; ok {
			subscriptionUpdate := &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateParams{}
			sp := val.([]interface{})
			for _, props := range sp {
				p := props.(map[string]interface{})
				if val, ok := p["default_allowed_updates"]; ok {
					enumsI := val.([]interface{})
					enums := []string{}
					for _, enum := range enumsI {
						enums = append(enums, enum.(string))
					}
					subscriptionUpdate.DefaultAllowedUpdates = stripe.StringSlice(enums)
				}

				if val, ok := p["enabled"]; ok {
					subscriptionUpdate.Enabled = stripe.Bool(val.(bool))
				}

				if val, ok := p["product"]; ok {
					var productsParams = []*stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams{}
					set := val.(*schema.Set)
					products := set.List()
					for _, i := range products {
						pParams := &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams{}
						finalProduct := i.(map[string]interface{})
						if val, ok := finalProduct["id"]; ok {
							pParams.Product = stripe.String(val.(string))
						}

						if val, ok := finalProduct["prices"]; ok {
							pricesI := val.([]interface{})
							prices := []string{}
							for _, price := range pricesI {
								prices = append(prices, price.(string))
							}
							pParams.Prices = stripe.StringSlice(prices)
						}
						productsParams = append(productsParams, pParams)
					}
					subscriptionUpdate.Products = productsParams
				}

				if val, ok := p["proration_behavior"]; ok {
					subscriptionUpdate.ProrationBehavior = stripe.String(val.(string))
				}
			}
			features.SubscriptionUpdate = subscriptionUpdate
		}
	}
	return features
}

func resourceStripeCustomerPortalCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := &stripe.BillingPortalConfigurationParams{}
	if dru, ok := d.GetOk("default_return_url"); ok {
		params.DefaultReturnURL = stripe.String(dru.(string))
	}

	if bp, ok := d.GetOk("business_profile"); ok {
		params.BusinessProfile = expandBusinessProfile(bp.([]interface{}))
	}

	if bp, ok := d.GetOk("features"); ok {
		params.Features = expandFeatures(bp.([]interface{}))
	}

	params.Metadata = expandMetadata(d)
	portal, err := client.BillingPortalConfigurations.New(params)
	if err == nil {
		log.Printf("[INFO] Customer Portal: %s", portal.ID)
		d.SetId(portal.ID)
	}
	return err
}

func resourceStripeCustomerPortalRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	portal, err := client.BillingPortalConfigurations.Get(d.Id(), nil)

	if err != nil {
		d.SetId("")
	} else {
		d.Set("id", portal.ID)
		d.Set("object", portal.Object)
		d.Set("active", portal.Active)
		d.Set("business_profile", portal.BusinessProfile)
		d.Set("created", portal.Created)
		d.Set("default_return_url", portal.DefaultReturnURL)
		d.Set("features", portal.Features)
		d.Set("is_default", portal.IsDefault)
		d.Set("livemode", portal.Livemode)
		d.Set("metadata", portal.Metadata)
		d.Set("updated", portal.Updated)
	}
	return err
}

func resourceStripeCustomerPortalUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := &stripe.BillingPortalConfigurationParams{}
	if d.HasChange("default_return_url") {
		params.BusinessProfile.Headline = stripe.String(d.Get("default_return_url").(string))
	}

	if d.HasChange("metadata") {
		params.Metadata = expandMetadata(d)
	}

	if d.HasChange("business_profile") {
		_, new := d.GetChange("business_profile")
		params.BusinessProfile = expandBusinessProfile(new.([]interface{}))
	}

	if d.HasChange("features") {
		_, new := d.GetChange("features")
		params.Features = expandFeatures(new.([]interface{}))
	}

	_, err := client.BillingPortalConfigurations.Update(d.Id(), params)
	if err != nil {
		return err
	}
	return resourceStripeCustomerPortalRead(d, m)
}

func resourceStripeCustomerPortalDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("[WARNING] Stripe doesn't allow deleting customer portal via the API. Please remove it manually")
}
