package stripe

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	stripe "github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
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
			"plan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Unique identifier for the plan.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the plan can be used for new purchases.",
			},
			"amount": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"amount_decimal"},
				Description:   "The unit amount in cents to be charged, represented as a whole integer if possible. Only set if `billing_scheme=per_unit`.",
			},
			"amount_decimal": {
				Type:          schema.TypeFloat,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"amount"},
				Description:   "The unit amount in cents to be charged, represented as a decimal string with at most 12 decimal places. Only set if `billing_scheme=per_unit`.",
			},
			"currency": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Three-letter ISO currency code, in lowercase. Must be a supported currency.",
			},
			"interval": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The frequency at which a subscription is billed. One of `day`, `week`, `month` or `year`.",
			},
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The product whose pricing this plan determines.",
			},
			"aggregate_usage": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies a usage aggregation strategy for plans of `usage_type=metered`. Allowed values are `sum` for summing up all usage during a period, `last_during_period` for using the last usage record reported within a period, `last_ever` for using the last usage record ever (across period bounds) or `max` which uses the usage record with the maximum reported usage during a period. Defaults to `sum`.",
			},
			"billing_scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "per_unit",
				Description: "Describes how to compute the price per period. Either `per_unit` or `tiered`. `per_unit` indicates that the fixed amount (specified in `amount`) will be charged per unit in `quantity` (for plans with `usage_type=licensed`), or per unit of total usage (for plans with `usage_type=metered`). `tiered` indicates that the unit pricing will be computed using a tiering strategy as defined using the `tiers` and `tiers_mode` attributes.",
			},
			"interval_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     1,
				Description: "The number of intervals (specified in the `interval` attribute) between subscription billings. For example, `interval=month` and `interval_count=3` bills every 3 months.",
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
				Description: "A brief description of the plan, hidden from customers.",
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
			"transform_usage": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"divide_by": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "Divide usage by this number.",
						},
						"round": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"down", "up"}, false),
							Description:  "After division, either round the result `up` or `down`.",
						},
					},
				},
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Apply a transformation to the reported usage or set quantity before computing the amount billed. Cannot be combined with `tiers`.",
			},
			"trial_period_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Default number of trial days when subscribing a customer to this plan using `trial_from_plan=true`.",
			},
			"usage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "licensed",
				Description: "Configures how the quantity per period should be determined. Can be either `metered` or `licensed`. `licensed` automatically bills the `quantity` set when adding it to a subscription. `metered` aggregates the total usage based on usage records. Defaults to `licensed`.",
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
		Interval:  stripe.String(planInterval),
		ProductID: stripe.String(planProductID),
		Currency:  stripe.String(planCurrency),
	}

	amount := d.Get("amount").(int)
	amountDecimal := d.Get("amount_decimal").(float64)

	if amountDecimal > 0 {
		params.AmountDecimal = stripe.Float64(float64(amountDecimal))
	} else {
		params.Amount = stripe.Int64(int64(amount))
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
			params.AmountDecimal = nil
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

	if transformUsage, ok := d.GetOk("transform_usage"); ok {
		params.TransformUsage = expandPlanTransformUsage(transformUsage.([]interface{}))
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
		d.Set("transform_usage", flattenPlanTransformUsage(plan.TransformUsage))
		d.Set("trial_period_days", plan.TrialPeriodDays)
		d.Set("usage_type", plan.UsageType)
	}

	return err
}

func flattenPlanTiers(in []*stripe.PlanTier) []map[string]interface{} {
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

func expandPlanTiers(in []interface{}) []*stripe.PlanTierParams {
	out := make([]*stripe.PlanTierParams, len(in))
	for i, v := range in {
		tier := v.(map[string]interface{})
		out[i] = &stripe.PlanTierParams{
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

func flattenPlanTransformUsage(in *stripe.PlanTransformUsage) []map[string]interface{} {
	n := 1
	if in == nil {
		n = 0
	}
	out := make([]map[string]interface{}, n)

	for i := range out {
		out[i] = map[string]interface{}{
			"divide_by": in.DivideBy,
			"round":     in.Round,
		}
	}
	return out
}

func expandPlanTransformUsage(in []interface{}) *stripe.PlanTransformUsageParams {
	if len(in) == 0 {
		return nil
	}

	transformUsage := in[0].(map[string]interface{})
	out := &stripe.PlanTransformUsageParams{
		DivideBy: stripe.Int64(int64(transformUsage["divide_by"].(int))),
		Round:    stripe.String(transformUsage["round"].(string)),
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
