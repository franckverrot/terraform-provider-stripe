package stripe

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("STRIPE_API_TOKEN", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"stripe_product": resourceStripeProduct(),
			"stripe_plan":    resourceStripePlan(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIToken: d.Get("api_token").(string),
	}

	log.Println("[INFO] Initializing Stripe client")
	return config.Client()
}
