package ewallet

import (
	"os"

	"github.com/xendit/xendit-go/ewallet"

	"github.com/imrenagi/go-payment/invoice"
)

// NewLinkAja is factory for LinkAja payment with xendit latest charge API
func NewLinkAja(inv *invoice.Invoice) (*ewallet.CreateEWalletChargeParams, error) {
	// Use environment variables as fallbacks
	successRedirectURL := os.Getenv("LINKAJA_SUCCESS_REDIRECT_URL")
	failureRedirectURL := os.Getenv("LINKAJA_FAILURE_REDIRECT_URL")

	// Prioritize per-request callback URLs if available
	if inv.SuccessRedirectURL != "" {
		successRedirectURL = inv.SuccessRedirectURL
	}

	if inv.FailureRedirectURL != "" {
		failureRedirectURL = inv.FailureRedirectURL
	}

	props := map[string]string{
		"success_redirect_url": successRedirectURL,
	}

	// Add failure redirect URL if available
	if failureRedirectURL != "" {
		props["failure_redirect_url"] = failureRedirectURL
	}

	return newBuilder(inv).
		SetPaymentMethod(EWalletIDLinkAja).
		SetChannelProperties(props).
		Build()
}
