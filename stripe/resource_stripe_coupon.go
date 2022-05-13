package stripe

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceStripeCoupon() *schema.Resource {
	return &schema.Resource{
		Create: resourceStripeCouponCreate,
		Read:   resourceStripeCouponRead,
		Update: resourceStripeCouponUpdate,
		Delete: resourceStripeCouponDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true, // require it as the default one is more trouble than it's worth
			},
			"amount_off": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"currency": &schema.Schema{
				Type:     schema.TypeString, // <- check values
				Optional: true,
				ForceNew: true,
			},
			"duration": &schema.Schema{
				Type:     schema.TypeString,
				Required: true, // forever | once | repeating
				ForceNew: true,
			},
			"duration_in_months": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"max_redemptions": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
				ForceNew: true,
			},
			"metadata": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"percent_off": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"redeem_by": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed
			"valid": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"livemode": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"times_redeemed": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceStripeCouponCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	couponID := d.Get("code").(string)
	params := &stripe.CouponParams{
		ID: stripe.String(couponID),
	}

	couponDuration := d.Get("duration").(string)
	validDurations := map[string]bool{
		"repeating": true,
		"once":      true,
		"forever":   true,
	}
	if !(validDurations)[couponDuration] {
		formattedKeys := "( " + strings.Join(getMapKeys(validDurations), " | ") + " )"
		return fmt.Errorf("\"%s\" is not a valid value for \"duration\", expected one of %s", couponDuration, formattedKeys)
	}

	if name, ok := d.GetOk("name"); ok {
		params.Name = stripe.String(name.(string))
	}

	if durationInMonths, ok := d.GetOk("duration_in_months"); ok {
		if couponDuration != "repeating" {
			return fmt.Errorf("can't set duration in months if event is not repeating")
		}
		params.DurationInMonths = stripe.Int64(int64(durationInMonths.(int)))
	}

	if couponDuration != "" {
		params.Duration = stripe.String(couponDuration)
	}

	if percentOff, ok := d.GetOk("percent_off"); ok {
		params.PercentOff = stripe.Float64(percentOff.(float64))
	}

	if amountOff, ok := d.GetOk("amount_off"); ok {
		params.AmountOff = stripe.Int64(int64(amountOff.(int)))
	}

	if maxRedemptions, ok := d.GetOk("max_redemptions"); ok {
		params.MaxRedemptions = stripe.Int64(int64(maxRedemptions.(int)))
	}

	if currency, ok := d.GetOk("currency"); ok {
		if &params.AmountOff == nil {
			return fmt.Errorf("can't set currency when using percent off")
		}
		params.Currency = stripe.String(currency.(string))
	}

	if redeemByStr, ok := d.GetOk("redeem_by"); ok {
		redeemByTime, err := time.Parse(time.RFC3339, redeemByStr.(string))

		if err != nil {
			return fmt.Errorf("can't convert time \"%s\" to time.  Please check if it's RFC3339-compliant", redeemByStr)
		}

		params.RedeemBy = stripe.Int64(redeemByTime.Unix())
	}

	params.Metadata = expandMetadata(d)

	coupon, err := client.Coupons.New(params)

	if err == nil {
		log.Printf("[INFO] Create coupon: %s (%s)", coupon.Name, coupon.ID)
		d.SetId(coupon.ID)
		d.Set("valid", coupon.Valid)
		d.Set("created", coupon.Created)
		d.Set("times_redeemed", coupon.TimesRedeemed)
		d.Set("livemode", coupon.Livemode)
	}

	return err
}

func resourceStripeCouponRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	coupon, err := client.Coupons.Get(d.Id(), nil)

	if err != nil {
		d.SetId("")
	} else {
		d.Set("code", d.Id())
		d.Set("amount_off", coupon.AmountOff)
		d.Set("currency", coupon.Currency)
		d.Set("duration", coupon.Duration)
		d.Set("duration_in_months", coupon.DurationInMonths)
		d.Set("livemode", coupon.Livemode)
		d.Set("max_redemptions", coupon.MaxRedemptions)
		d.Set("metadata", coupon.Metadata)
		d.Set("name", coupon.Name)
		d.Set("percent_off", coupon.PercentOff)
		d.Set("redeem_by", coupon.RedeemBy)
		d.Set("times_redeemed", coupon.TimesRedeemed)
		d.Set("valid", coupon.Valid)
		d.Set("created", coupon.Valid)
	}

	return err
}

func resourceStripeCouponUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.CouponParams{}

	if d.HasChange("metadata") {
		params.Metadata = expandMetadata(d)
	}

	if d.HasChange("name") {
		params.Name = stripe.String(d.Get("name").(string))
	}

	_, err := client.Coupons.Update(d.Id(), &params)

	if err != nil {
		return err
	}

	return resourceStripeCouponRead(d, m)
}

func resourceStripeCouponDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	_, err := client.Coupons.Del(d.Id(), nil)

	if err == nil {
		d.SetId("")
	}

	return err
}
