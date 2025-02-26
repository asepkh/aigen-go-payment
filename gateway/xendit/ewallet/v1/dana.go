package ewallet

import (
	"os"

	goxendit "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"

	"github.com/asepkh/aigen-go-payment/invoice"
)

// NewDana create xendit payment request for Dana
func NewDana(inv *invoice.Invoice) (*ewallet.CreatePaymentParams, error) {
	// Use environment variables as fallbacks
	callbackURL := os.Getenv("DANA_LEGACY_CALLBACK_URL")
	redirectURL := os.Getenv("DANA_LEGACY_REDIRECT_URL")

	// Prioritize per-request callback URLs if available
	if inv.SuccessRedirectURL != "" {
		redirectURL = inv.SuccessRedirectURL
	}

	return newBuilder(inv).
		SetPaymentMethod(goxendit.EWalletTypeDANA).
		SetCallback(callbackURL).
		SetRedirect(redirectURL).
		Build()
}
