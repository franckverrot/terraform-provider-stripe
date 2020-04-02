package stripe

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

func resourceStripePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceStripePlanCreate,
		Read:   resourceStripePlanRead,
		Update: resourceStripePlanUpdate,
		Delete: resourceStripePlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"plan_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"amount": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"amount_decimal": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"currency": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interval": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aggregate_usage": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"billing_scheme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "per_unit",
			},
			"interval_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
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
			// TODO: transform_usage
			"trial_period_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"usage_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "licensed",
			},
		},
	}
}

func resourceStripePlanCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	planNickname := d.Get("nickname").(string)
	planInterval := d.Get("interval").(string)
	planCurrency := d.Get("currency").(string)
	planProductID := d.Get("product").(string)

	// TODO: check interval
	// TODO: check currency

	params := &stripe.PlanParams{
		Interval:	stripe.String(planInterval),
		ProductID:	stripe.String(planProductID),
		Currency:	stripe.String(planCurrency),
	}

	if amount, ok := d.GetOk("amount"); ok {
		params.Amount = stripe.Int64(int64(amount.(int)))
	} else if amountDecimal, ok := d.GetOk("amount_decimal"); ok {
		params.AmountDecimal = stripe.Float64(amountDecimal.(float64))
	}

	if id, ok := d.GetOk("plan_id"); ok {
		params.ID = stripe.String(id.(string))
	}

	if active, ok := d.GetOk("active"); ok {
		params.Active = stripe.Bool(active.(bool))
	}

	if aggregateUsage, ok := d.GetOk("aggregate_usage"); ok {
		params.AggregateUsage = stripe.String(aggregateUsage.(string))
	}

	if billingScheme, ok := d.GetOk("billing_scheme"); ok {
		params.BillingScheme = stripe.String(billingScheme.(string))
		if billingScheme == "tiered" {
			params.Amount = nil
		}
	}

	if intervalCount, ok := d.GetOk("interval_count"); ok {
		params.IntervalCount = stripe.Int64(int64(intervalCount.(int)))
	}

	params.Metadata = expandMetadata(d)

	if _, ok := d.GetOk("nickname"); ok {
		params.Nickname = stripe.String(planNickname)
	}

	if tiersMode, ok := d.GetOk("tiers_mode"); ok {
		params.TiersMode = stripe.String(tiersMode.(string))
	}

	if tiers, ok := d.GetOk("tier"); ok {
		params.Tiers = expandPlanTiers(tiers.([]interface{}))
	}

	if trialPeriodDays, ok := d.GetOk("trial_period_days"); ok {
		params.TrialPeriodDays = stripe.Int64(int64(trialPeriodDays.(int)))
	}

	if usageType, ok := d.GetOk("usage_type"); ok {
		params.UsageType = stripe.String(usageType.(string))
	}

	plan, err := client.Plans.New(params)

	if err == nil {
		if planNickname != "" {
			log.Printf("[INFO] Create plan: %s (%s)", plan.Nickname, plan.ID)
		} else {
			log.Printf("[INFO] Create anonymous plan: %s", plan.ID)
		}
		d.SetId(plan.ID)
	}

	return err
}

func resourceStripePlanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	plan, err := client.Plans.Get(d.Id(), nil)

	if err != nil {
		d.SetId("")
	} else {
		d.Set("plan_id", plan.ID)
		d.Set("active", plan.Active)
		d.Set("aggregate_usage", plan.AggregateUsage)
		d.Set("amount", plan.Amount)
		d.Set("amount_decimal", plan.AmountDecimal)
		d.Set("billing_scheme", plan.BillingScheme)
		d.Set("currency", plan.Currency)
		d.Set("interval", plan.Interval)
		d.Set("interval_count", plan.IntervalCount)
		d.Set("metadata", plan.Metadata)
		d.Set("nickname", plan.Nickname)
		d.Set("product", plan.Product)
		d.Set("tiers_mode", plan.TiersMode)
		d.Set("tier", flattenPlanTiers(plan.Tiers))
		d.Set("trial_period_days", plan.TrialPeriodDays)
		d.Set("usage_type", plan.UsageType)
	}

	return err
}

func flattenPlanTiers(in []*stripe.PlanTier) []map[string]interface{} {
	out := make([]map[string]interface{}, len(in))
	for i, tier := range in {
		out[i] = map[string]interface{}{
			"up_to":		tier.UpTo,
			"up_to_inf":		tier.UpTo == 0,
			"flat_amount":		tier.FlatAmount,
			"flat_amount_decimal":	tier.FlatAmountDecimal,
			"unit_amount":		tier.UnitAmount,
			"unit_amount_decimal":	tier.UnitAmountDecimal,
		}
	}
	return out
}

func expandPlanTiers(in []interface{}) []*stripe.PlanTierParams {
	out := make([]*stripe.PlanTierParams, len(in))
	for i, v := range in {
		tier := v.(map[string]interface{})
		out[i] = &stripe.PlanTierParams{
			UpTo:		stripe.Int64(int64(tier["up_to"].(int))),
			UpToInf:	stripe.Bool(tier["up_to_inf"].(bool)),
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

func resourceStripePlanUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.PlanParams{}

	if d.HasChange("plan_id") {
		params.ID = stripe.String(d.Get("plan_id").(string))
	}

	if d.HasChange("active") {
		params.Active = stripe.Bool(bool(d.Get("active").(bool)))
	}

	if d.HasChange("metadata") {
		params.Metadata = expandMetadata(d)
	}

	if d.HasChange("nickname") {
		params.Nickname = stripe.String(d.Get("nickname").(string))
	}

	if d.HasChange("trial_period_days") {
		params.TrialPeriodDays = stripe.Int64(int64(d.Get("trial_period_days").(int)))
	}

	_, err := client.Plans.Update(d.Id(), &params)

	if err != nil {
		return err
	}

	return resourceStripePlanRead(d, m)
}

func resourceStripePlanDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	_, err := client.Plans.Del(d.Id(), nil)

	if err == nil {
		d.SetId("")
	}

	return err
}
