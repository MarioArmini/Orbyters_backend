package stripe

import (
	"Orbyters/config"

	"github.com/stripe/stripe-go"
)

func RefreshProducts() error {
	stripe.Key = config.StripeKey

}
