package stripe

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"

	"fmt"
	"log"
)

func extractAttributes(d *schema.ResourceData) []*string {
	attributes := d.Get("attributes").(*schema.Set).List()

	if _, ok := d.GetOk("attributes"); ok {
		convertedAttributes := []*string{}

		for _, attribute := range attributes {
			tmp := attribute.(string)
			convertedAttributes = append(convertedAttributes, &tmp)
		}

		return convertedAttributes
	}

	return nil
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// TODO Restore once we figure out how to properly delete this data with the Go client
			//"attributes": &schema.Schema{
			//	Type: schema.TypeSet,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	Optional: true,
			//},
			//"metadata": &schema.Schema{
			//	Type: schema.TypeMap,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	Optional: true,
			//},
			"statement_descriptor": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"unit_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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

	if active, ok := d.GetOk("active"); ok {
		params.Active = stripe.Bool(active.(bool))
	}

	//params.Attributes = extractAttributes(d)

	//params.Metadata = expandMetadata(d)

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

	d.Set("name", product.Name)
	d.Set("type", product.Type)
	d.Set("active", product.Active)
	//d.Set("attributes", product.Attributes)
	//d.Set("metadata", product.Metadata)
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

	//if d.HasChange("attributes") {
	//	params.Attributes = extractAttributes(d)
	//}
	//
	//if d.HasChange("metadata") {
	//	params.Metadata = expandMetadata(d)
	//}

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
