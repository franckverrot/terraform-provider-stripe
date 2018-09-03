package stripe

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"

	"fmt"
	"log"
)

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
		},
	}
}

func resourceStripeProductCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	productName := d.Get("name").(string)
	productType := d.Get("type").(string)

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
	product, err := client.Products.New(params)

	if err != nil {
		return err
	} else {
		log.Printf("[INFO] Create product: %s", productName)
		d.SetId(product.ID)
	}
	return nil
}

func resourceStripeProductRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	product, err := client.Products.Get(d.Id(), nil)

	if err != nil {
		return err
	} else {
		d.Set("name", product.Name)
		d.Set("type", product.Type)
	}
	return nil
}

func resourceStripeProductUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.ProductParams{}

	if d.HasChange("name") {
		params.Name = stripe.String(d.Get("name").(string))
	}
	if d.HasChange("type") {
		params.Type = stripe.String(d.Get("Type").(string))
	}

	_, err := client.Products.Update(d.Id(), &params)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func resourceStripeProductDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	_, err := client.Products.Del(d.Id(), nil)

	if err != nil {
		return err
	} else {
		d.SetId("")
		return nil
	}
}
