package stripe

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
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
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"inclusive": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"jurisdiction": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"livemode": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"metadata": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"percentage": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
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
