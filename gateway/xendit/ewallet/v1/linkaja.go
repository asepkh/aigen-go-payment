package ewallet

import (
	"os"

	goxendit "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"

	"github.com/asepkh/aigen-go-payment/invoice"
)

// NewLinkAja create xendit payment request for LinkAja
func NewLinkAja(inv *invoice.Invoice) (*ewallet.CreatePaymentParams, error) {
	// Use environment variables as fallbacks
	callbackURL := os.Getenv("LINKAJA_LEGACY_CALLBACK_URL")
	redirectURL := os.Getenv("LINKAJA_LEGACY_REDIRECT_URL")

	// Prioritize per-request callback URLs if available
	if inv.SuccessRedirectURL != "" {
		redirectURL = inv.SuccessRedirectURL
	}

	return newBuilder(inv).
		SetPaymentMethod(goxendit.EWalletTypeLINKAJA).
		SetCallback(callbackURL).
		SetRedirect(redirectURL).
		Build()
}
