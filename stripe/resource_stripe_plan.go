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
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"amount": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
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
			},
			"billing_scheme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "per_unit",
			},
			"interval_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
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
			// TODO: tiers
			"tiers_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// TODO: transform_usage
			"trial_period_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"usage_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "licensed",
			},
		},
	}
}

func resourceStripePlanCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	planAmount := int64(d.Get("amount").(int))
	planNickname := d.Get("nickname").(string)
	planInterval := d.Get("interval").(string)
	planCurrency := d.Get("currency").(string)
	planProductID := d.Get("product").(string)

	// TODO: check interval
	// TODO: check currency

	params := &stripe.PlanParams{
		Amount:    stripe.Int64(planAmount),
		Interval:  stripe.String(planInterval),
		ProductID: stripe.String(planProductID),
		Currency:  stripe.String(planCurrency),
	}

	if active, ok := d.GetOk("active"); ok {
		params.Active = stripe.Bool(active.(bool))
	}

	if aggregateUsage, ok := d.GetOk("aggregate_usage"); ok {
		params.AggregateUsage = stripe.String(aggregateUsage.(string))
	}

	if billingScheme, ok := d.GetOk("billing_scheme"); ok {
		params.BillingScheme = stripe.String(billingScheme.(string))
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
		d.Set("active", plan.Active)
		d.Set("aggregate_usage", plan.AggregateUsage)
		d.Set("amount", plan.Amount)
		d.Set("billing_scheme", plan.BillingScheme)
		d.Set("currency", plan.Currency)
		d.Set("interval", plan.Interval)
		d.Set("interval_count", plan.IntervalCount)
		d.Set("metadata", plan.Metadata)
		d.Set("nickname", plan.Nickname)
		d.Set("product", plan.Product)
		d.Set("tiers_mode", plan.TiersMode)
		d.Set("trial_period_days", plan.TrialPeriodDays)
		d.Set("usage_type", plan.UsageType)
	}

	return err
}

func resourceStripePlanUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.PlanParams{}

	if d.HasChange("active") {
		params.Active = stripe.Bool(bool(d.Get("active").(bool)))
	}

	if d.HasChange("aggregate_usage") {
		params.AggregateUsage = stripe.String(d.Get("aggregate_usage").(string))
	}

	if d.HasChange("billing_scheme") {
		params.BillingScheme = stripe.String(d.Get("billing_scheme").(string))
	}

	if d.HasChange("interval_count") {
		params.IntervalCount = stripe.Int64(int64(d.Get("interval_count").(int)))
	}

	if d.HasChange("metadata") {
		params.Metadata = expandMetadata(d)
	}

	if d.HasChange("nickname") {
		params.Nickname = stripe.String(d.Get("nickname").(string))
	}

	if d.HasChange("tiers_mode") {
		params.TiersMode = stripe.String(d.Get("tiers_mode").(string))
	}

	if d.HasChange("trial_period_days") {
		params.TrialPeriodDays = stripe.Int64(int64(d.Get("trial_period_days").(int)))
	}

	if d.HasChange("usage_type") {
		params.UsageType = stripe.String(d.Get("usage_type").(string))
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
