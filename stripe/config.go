package stripe

import (
	"log"

	"github.com/stripe/stripe-go/client"
)

// Config stores Stripe's API configuration
type Config struct {
	APIToken string
}

// Client returns a new Client for accessing Stripe.
func (c *Config) Client() (*client.API, error) {
	client := &client.API{}
	client.Init(c.APIToken, nil)
	log.Printf("[INFO] Stripe Client configured.")

	return client, nil
}
