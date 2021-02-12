package stripe

import (
	"github.com/hashicorp/terraform/helper/schema"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"

	"log"
)

func resourceStripeWebhookEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceStripeWebhookEndpointCreate,
		Read:   resourceStripeWebhookEndpointRead,
		Update: resourceStripeWebhookEndpointUpdate,
		Delete: resourceStripeWebhookEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled_events": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"connect": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"secret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStripeWebhookEndpointCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	url := d.Get("url").(string)

	params := &stripe.WebhookEndpointParams{
		URL:           stripe.String(url),
		EnabledEvents: expandStringList(d, "enabled_events"),
	}

	if connect, ok := d.GetOk("connect"); ok {
		params.Connect = stripe.Bool(connect.(bool))
	}

	webhookEndpoint, err := client.WebhookEndpoints.New(params)

	if err == nil {
		log.Printf("[INFO] Create webbook endpoint: %s", url)
		d.SetId(webhookEndpoint.ID)
		d.Set("secret", webhookEndpoint.Secret)
	}

	return err
}

func resourceStripeWebhookEndpointRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	webhookEndpoint, err := client.WebhookEndpoints.Get(d.Id(), nil)

	if err != nil {
		return err
	}

	d.Set("url", webhookEndpoint.URL)
	d.Set("enabled_events", webhookEndpoint.EnabledEvents)
	d.Set("connect", webhookEndpoint.Application != "")

	return nil
}

func resourceStripeWebhookEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.WebhookEndpointParams{}

	if d.HasChange("url") {
		params.URL = stripe.String(d.Get("url").(string))
	}

	if d.HasChange("enabled_events") {
		params.EnabledEvents = expandStringList(d, "enabled_events")
	}

	if d.HasChange("connect") {
		params.Connect = stripe.Bool(d.Get("connect").(bool))
	}

	_, err := client.WebhookEndpoints.Update(d.Id(), &params)

	if err != nil {
		return err
	}

	return resourceStripeWebhookEndpointRead(d, m)
}

func resourceStripeWebhookEndpointDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	_, err := client.WebhookEndpoints.Del(d.Id(), nil)

	if err == nil {
		d.SetId("")
	}

	return err
}
