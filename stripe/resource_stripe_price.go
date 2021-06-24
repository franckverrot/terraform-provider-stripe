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
			"price_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"currency": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metadata": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"nickname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"product": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"recurring": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"unit_amount": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"unit_amount_decimal": &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"billing_scheme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"livemode": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tier": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"up_to": &schema.Schema{
							Type:          schema.TypeInt,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"tier.up_to_inf"},
						},
						"up_to_inf": &schema.Schema{
							Type:          schema.TypeBool,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"tier.up_to"},
						},
						"flat_amount": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"flat_amount_decimal": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"unit_amount": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"unit_amount_decimal": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
				Optional: true,
				ForceNew: true,
			},
			"tiers_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

	if unitAmount, ok := d.GetOkExists("unit_amount"); ok {
		params.UnitAmount = stripe.Int64(int64(unitAmount.(int)))
	}

	if unitAmountDecimal, ok := d.GetOkExists("unit_amount_decimal"); ok {
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
