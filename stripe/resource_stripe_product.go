package stripe

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"

	"fmt"
	"log"
)

func expandAttributes(d *schema.ResourceData) []*string {
	return expandStringList(d, "attributes")
}

func resourceStripeProduct() *schema.Resource {
	return &schema.Resource{
		Create: resourceStripeProductCreate,
		Read:   resourceStripeProductRead,
		Update: resourceStripeProductUpdate,
		Delete: resourceStripeProductDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Unique identifier for the product.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The product’s name, meant to be displayable to the customer. Whenever this product is sold via a subscription, name will show up on associated invoice line item descriptions.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the product is currently available for purchase. Defaults to `true`.",
			},
			"attributes": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"metadata": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.",
			},
			"statement_descriptor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extra information about a product which will appear on your customer’s credit card statement. In the case that multiple products are billed at once, the first statement descriptor will be used.",
			},
			"unit_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A label that represents units of this product in Stripe and on customers’ receipts and invoices. When set, this will be included in associated invoice line item descriptions.",
			},
		},
	}
}

func resourceStripeProductCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	productName := d.Get("name").(string)
	productType := d.Get("type").(string)
	productStatementDescriptor := d.Get("statement_descriptor").(string)
	productUnitLabel := d.Get("unit_label").(string)

	var stripeProductType stripe.ProductType

	switch productType {
	case "good":
		stripeProductType = stripe.ProductTypeGood
	case "service":
		stripeProductType = stripe.ProductTypeService
	default:
		return fmt.Errorf("unknown type: %s", productType)
	}

	params := &stripe.ProductParams{
		Name: stripe.String(productName),
		Type: stripe.String(string(stripeProductType)),
	}

	if productID, ok := d.GetOk("product_id"); ok {
		params.ID = stripe.String(productID.(string))
	}

	if active, ok := d.GetOk("active"); ok {
		params.Active = stripe.Bool(active.(bool))
	}

	params.Attributes = expandAttributes(d)

	params.Metadata = expandMetadata(d)

	if productStatementDescriptor != "" {
		params.StatementDescriptor = stripe.String(productStatementDescriptor)
	}

	if productUnitLabel != "" {
		params.UnitLabel = stripe.String(productUnitLabel)
	}

	product, err := client.Products.New(params)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Created Stripe product: %s", productName)
	d.SetId(product.ID)

	return resourceStripeProductRead(d, m)
}

func resourceStripeProductRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	product, err := client.Products.Get(d.Id(), nil)

	if err != nil {
		return err
	}

	d.Set("product_id", product.ID)
	d.Set("name", product.Name)
	d.Set("type", product.Type)
	d.Set("active", product.Active)
	d.Set("attributes", product.Attributes)
	d.Set("metadata", product.Metadata)
	d.Set("statement_descriptor", product.StatementDescriptor)
	d.Set("unit_label", product.UnitLabel)

	return nil
}

func resourceStripeProductUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.ProductParams{}

	if d.HasChange("name") {
		params.Name = stripe.String(d.Get("name").(string))
	}

	if d.HasChange("type") {
		params.Type = stripe.String(d.Get("type").(string))
	}

	if d.HasChange("active") {
		params.Active = stripe.Bool(d.Get("active").(bool))
	}

	if d.HasChange("attributes") {
		params.Attributes = expandAttributes(d)
	}

	if d.HasChange("metadata") {
		params.Metadata = expandMetadata(d)
	}

	if d.HasChange("statement_descriptor") {
		params.StatementDescriptor = stripe.String(d.Get("statement_descriptor").(string))
	}

	if d.HasChange("unit_label") {
		params.UnitLabel = stripe.String(d.Get("unit_label").(string))
	}

	_, err := client.Products.Update(d.Id(), &params)

	if err != nil {
		return err
	}

	return resourceStripeProductRead(d, m)
}

func resourceStripeProductDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	_, err := client.Products.Del(d.Id(), nil)

	if err == nil {
		d.SetId("")
	}

	return err
}
