package stripe

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	stripe "github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
)

func resourceStripeTaxRate() *schema.Resource {
	return &schema.Resource{
		Create: resourceStripeTaxRateCreate,
		Read:   resourceStripeTaxRateRead,
		Update: resourceStripeTaxRateUpdate,
		Delete: resourceStripeTaxRateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Defaults to `true`. When set to `false`, this tax rate cannot be used with new applications or Checkout Sessions, but will still work for subscriptions and invoices that already have it set.",
			},
			"created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Time at which the object was created. Measured in seconds since the Unix epoch.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An arbitrary string attached to the tax rate for your internal use only. It will not be visible to your customers.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the tax rates as it will appear to your customer on their receipt email, PDF, and the hosted invoice page.",
			},
			"inclusive": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "This specifies if the tax rate is inclusive or exclusive.",
			},
			"jurisdiction": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The jurisdiction for the tax rate. You can use this label field for tax reporting purposes. It also appears on your customer’s invoice.",
			},
			"livemode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Has the value `true` if the object exists in live mode or the value `false` if the object exists in test mode.",
			},
			"metadata": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.",
			},
			"percentage": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "This represents the tax rate percent out of 100.",
			},
		},
	}
}

func resourceStripeTaxRateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	taxRateDisplayName := d.Get("display_name").(string)
	taxRateInclusive := d.Get("inclusive").(bool)
	taxRatePercentage := d.Get("percentage").(float64)
	params := &stripe.TaxRateParams{
		DisplayName: stripe.String(taxRateDisplayName),
		Inclusive:   stripe.Bool(taxRateInclusive),
		Percentage:  stripe.Float64(taxRatePercentage),
	}

	if active, ok := d.GetOk("active"); ok {
		params.Active = stripe.Bool(active.(bool))
	}

	if description, ok := d.GetOk("description"); ok {
		params.Description = stripe.String(description.(string))
	}

	if jurisdiction, ok := d.GetOk("jurisdiction"); ok {
		params.Jurisdiction = stripe.String(jurisdiction.(string))
	}

	params.Metadata = expandMetadata(d)

	Tax, err := client.TaxRates.New(params)

	if err == nil {
		log.Printf("[INFO] Create Tax Rate: %s (%f)", Tax.ID, Tax.Percentage)
		d.SetId(Tax.ID)
		d.Set("display_name", Tax.DisplayName)
		d.Set("inclusive", Tax.Inclusive)
		d.Set("percentage", Tax.Percentage)
		d.Set("created", Tax.Created)
		d.Set("livemode", Tax.Livemode)
	}

	return err
}

func resourceStripeTaxRateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	Tax, err := client.TaxRates.Get(d.Id(), nil)

	if err != nil {
		d.SetId("")
	} else {
		d.Set("active", Tax.Active)
		d.Set("created", Tax.Created)
		d.Set("description", Tax.Description)
		d.Set("display_name", Tax.DisplayName)
		d.Set("inclusive", Tax.Inclusive)
		d.Set("jurisdiction", Tax.Jurisdiction)
		d.Set("livemode", Tax.Livemode)
		d.Set("metadata", Tax.Metadata)
	}

	return err
}

func resourceStripeTaxRateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.TaxRateParams{}

	if d.HasChange("active") {
		params.Active = stripe.Bool(d.Get("active").(bool))
	}

	if d.HasChange("description") {
		params.Description = stripe.String(d.Get("description").(string))
	}

	if d.HasChange("diplay_name") {
		params.DisplayName = stripe.String(d.Get("display_name").(string))
	}

	if d.HasChange("jurisdiction") {
		params.Jurisdiction = stripe.String(d.Get("jurisdiction").(string))
	}

	if d.HasChange("metadata") {
		params.Metadata = expandMetadata(d)
	}

	_, err := client.TaxRates.Update(d.Id(), &params)

	if err != nil {
		return err
	}

	return resourceStripeTaxRateRead(d, m)
}

func resourceStripeTaxRateDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("[WARNING] Stripe doesn't allow deleting tax rates via the API.  Your state file contains at least one (\"%v\") that needs deletion.  Please remove it manually.", d.Get("display_name"))
}
