package stripe

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	stripe "github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
)

func resourceStripePrice() *schema.Resource {
	return &schema.Resource{
		Create: resourceStripePriceCreate,
		Read:   resourceStripePriceRead,
		Update: resourceStripePriceUpdate,
		Delete: resourceStripePriceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"price_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Unique identifier for the price.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the price can be used for new purchases.",
			},
			"currency": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Three-letter ISO currency code, in lowercase. Must be a supported currency.",
			},
			"metadata": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.",
			},
			"nickname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A brief description of the price, hidden from customers.",
			},
			"product": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the product this price is associated with.",
			},
			"recurring": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The recurring components of a price such as `interval` and `usage_type`.",
			},
			"unit_amount": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The unit amount in cents to be charged, represented as a whole integer if possible. Only set if `billing_scheme=per_unit`.",
			},
			"unit_amount_decimal": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The unit amount in cents to be charged, represented as a decimal string with at most 12 decimal places. Only set if `billing_scheme=per_unit`.",
			},
			"billing_scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Describes how to compute the price per period. Either `per_unit` or `tiered`. `per_unit` indicates that the fixed amount (specified in `unit_amount` or `unit_amount_decimal`) will be charged per unit in quantity (for prices with `usage_type=licensed`), or per unit of total usage (for prices with `usage_type=metered`). `tiered` indicates that the unit pricing will be computed using a tiering strategy as defined using the `tiers` and `tiers_mode` attributes.",
			},
			"created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Time at which the object was created. Measured in seconds since the Unix epoch.",
			},
			"livemode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Has the value `true` if the object exists in live mode or the value `false` if the object exists in test mode.",
			},
			"tier": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"up_to": {
							Type:          schema.TypeInt,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"tier.up_to_inf"},
							Description:   "Up to and including to this quantity will be contained in the tier.",
						},
						"up_to_inf": {
							Type:          schema.TypeBool,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"tier.up_to"},
						},
						"flat_amount": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Price for the entire tier.",
						},
						"flat_amount_decimal": {
							Type:        schema.TypeFloat,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Same as `flat_amount`, but contains a decimal value with at most 12 decimal places.",
						},
						"unit_amount": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Per unit price for units relevant to the tier.",
						},
						"unit_amount_decimal": {
							Type:        schema.TypeFloat,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Same as `unit_amount`, but contains a decimal value with at most 12 decimal places.",
						},
					},
				},
				Optional:    true,
				ForceNew:    true,
				Description: "Each element represents a pricing tier. This parameter requires `billing_scheme` to be set to `tiered`. See also the documentation for `billing_scheme`.",
			},
			"tiers_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Defines if the tiering price should be `graduated` or `volume` based. In `volume`-based tiering, the maximum quantity within a period determines the per unit price. In `graduated` tiering, pricing can change as the quantity grows.",
			},
		},
	}
}

func expandPriceRecurring(recurring map[string]interface{}) (*stripe.PriceRecurringParams, error) {
	params := &stripe.PriceRecurringParams{}
	parsed := expandStringMap(recurring)

	if aggregateUsage, ok := parsed["aggregate_usage"]; ok {
		params.AggregateUsage = stripe.String(aggregateUsage)
	}

	if interval, ok := parsed["interval"]; ok {
		params.Interval = stripe.String(interval)
	}

	if intervalCount, ok := parsed["interval_count"]; ok {
		intervalCountInt, err := strconv.ParseInt(intervalCount, 10, 64)
		if err != nil {
			return nil, errors.New("interval_count must be a string, representing an int (e.g. \"52\")")
		}
		params.IntervalCount = stripe.Int64(intervalCountInt)
	}

	if usageType, ok := parsed["usage_type"]; ok {
		params.UsageType = stripe.String(usageType)
	}

	return params, nil
}

func resourceStripePriceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	nickname := d.Get("nickname").(string)
	currency := d.Get("currency").(string)

	params := &stripe.PriceParams{
		Currency: stripe.String(currency),
	}

	if active, ok := d.GetOk("active"); ok {
		params.Active = stripe.Bool(active.(bool))
	}

	params.Metadata = expandMetadata(d)

	if _, ok := d.GetOk("nickname"); ok {
		params.Nickname = stripe.String(nickname)
	}

	if tiersMode, ok := d.GetOk("tiers_mode"); ok {
		params.TiersMode = stripe.String(tiersMode.(string))
	}

	if tiers, ok := d.GetOk("tier"); ok {
		params.Tiers = expandPriceTiers(tiers.([]interface{}))
	}

	if product, ok := d.GetOk("product"); ok {
		params.Product = stripe.String(product.(string))
	}

	if recurring, ok := d.GetOk("recurring"); ok {
		recurringParams, err := expandPriceRecurring(recurring.(map[string]interface{}))
		if err != nil {
			return err
		}
		params.Recurring = recurringParams
	}

	if unitAmount, ok := d.GetOk("unit_amount"); ok {
		params.UnitAmount = stripe.Int64(int64(unitAmount.(int)))
	}

	if unitAmountDecimal, ok := d.GetOk("unit_amount_decimal"); ok {
		params.UnitAmountDecimal = stripe.Float64(unitAmountDecimal.(float64))
	}

	if billingScheme, ok := d.GetOk("billing_scheme"); ok {
		params.BillingScheme = stripe.String(billingScheme.(string))
	}

	price, err := client.Prices.New(params)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created Stripe price: %s", nickname)
	d.SetId(price.ID)

	return resourceStripePriceRead(d, m)
}

func resourceStripePriceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	price, err := client.Prices.Get(d.Id(), nil)

	if err != nil {
		d.SetId("")
	} else {
		d.Set("price_id", price.ID)
		d.Set("active", price.Active)
		d.Set("created", price.Created)
		d.Set("currency", price.Currency)
		d.Set("livemode", price.Livemode)
		d.Set("metadata", price.Metadata)
		d.Set("nickname", price.Nickname)
		if price.Product != nil {
			d.Set("product", price.Product.ID)
		}
		d.Set("recurring", price.Active)
		d.Set("unit_amount", price.UnitAmount)
		d.Set("unit_amount_decimal", price.UnitAmountDecimal)
		d.Set("tiers_mode", price.TiersMode)
		// Stripe's API doesn't return tiers.
		// d.Set("tier", flattenPriceTiers(price.Tiers))
		d.Set("billing_scheme", price.BillingScheme)
	}

	return err
}

func flattenPriceTiers(in []*stripe.PriceTier) []map[string]interface{} {
	out := make([]map[string]interface{}, len(in))
	for i, tier := range in {
		out[i] = map[string]interface{}{
			"up_to":               tier.UpTo,
			"up_to_inf":           tier.UpTo == 0,
			"flat_amount":         tier.FlatAmount,
			"flat_amount_decimal": tier.FlatAmountDecimal,
			"unit_amount":         tier.UnitAmount,
			"unit_amount_decimal": tier.UnitAmountDecimal,
		}
	}
	return out
}

func expandPriceTiers(in []interface{}) []*stripe.PriceTierParams {
	out := make([]*stripe.PriceTierParams, len(in))
	for i, v := range in {
		tier := v.(map[string]interface{})
		out[i] = &stripe.PriceTierParams{
			UpTo:    stripe.Int64(int64(tier["up_to"].(int))),
			UpToInf: stripe.Bool(tier["up_to_inf"].(bool)),
		}
		if tier["flat_amount"] != nil {
			out[i].FlatAmount = stripe.Int64(int64(tier["flat_amount"].(int)))
		} else if tier["flat_amount_decimal"] != nil {
			out[i].FlatAmountDecimal = stripe.Float64(tier["flat_amount_decimal"].(float64))
		}
		if tier["unit_amount"] != nil {
			out[i].UnitAmount = stripe.Int64(int64(tier["unit_amount"].(int)))
		} else if tier["unit_amount_decimal"] != nil {
			out[i].UnitAmountDecimal = stripe.Float64(tier["unit_amount_decimal"].(float64))
		}
	}
	return out
}

func resourceStripePriceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.PriceParams{}

	if d.HasChange("active") {
		params.Active = stripe.Bool(d.Get("active").(bool))
	}

	if d.HasChange("metadata") {
		params.Metadata = expandMetadata(d)
	}

	if d.HasChange("nickname") {
		params.Nickname = stripe.String(d.Get("nickname").(string))
	}

	_, err := client.Prices.Update(d.Id(), &params)
	if err != nil {
		return err
	}

	return resourceStripePriceRead(d, m)
}

func resourceStripePriceDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("[WARNING] Stripe doesn't allow deleting prices via the API. Your state file contains at least one (\"%v\") that needs deletion. Please remove it manually.", d.Id())
}
