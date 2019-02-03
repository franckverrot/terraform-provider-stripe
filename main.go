package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/stripe/terraform-provider-stripe/stripe"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: stripe.Provider})
}
