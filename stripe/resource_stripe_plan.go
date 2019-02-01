package stripe

import (
	"fmt"
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
			"nickname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"amount": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"currency": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"interval": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"product": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
	planActive := d.Get("active").(bool)

	// TODO: check interval
	// TODO: check currency

	params := &stripe.PlanParams{
		Amount:    stripe.Int64(planAmount),
		Interval:  stripe.String(planInterval),
		ProductID: stripe.String(planProductID),
		Currency:  stripe.String(planCurrency),
		Active:    stripe.Bool(planActive),
	}

	if planNickname != "" {
		params.Nickname = stripe.String(planNickname)
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
		d.Set("nickname", plan.Nickname)
		d.Set("product", plan.Product)
		d.Set("amount", plan.Amount)
		d.Set("interval", plan.Interval)
		d.Set("currency", plan.Currency)
		d.Set("active", plan.Active)
	}

	return err
}

func resourceStripePlanUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.PlanParams{}

	if d.HasChange("active") {
		params.Active = stripe.Bool(bool(d.Get("active").(bool)))
	}

	d.Partial(true)

	_, err := client.Plans.Update(d.Id(), &params)

	d.SetPartial("active")

	if d.HasChange("amount") {
		return fmt.Errorf("Changing amounts is not possible with the Stripe API")
	}

	d.Partial(false)

	return err
}

func resourceStripePlanDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	_, err := client.Plans.Del(d.Id(), nil)

	if err == nil {
		d.SetId("")
	}

	return err
}
