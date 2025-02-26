package ewallet

import (
	"os"

	"github.com/xendit/xendit-go/ewallet"

	"github.com/asepkh/aigen-go-payment/invoice"
)

// NewDana is factory for Dana payment with xendit latest charge API
func NewDana(inv *invoice.Invoice) (*ewallet.CreateEWalletChargeParams, error) {
	// Use environment variables as fallbacks
	successRedirectURL := os.Getenv("DANA_SUCCESS_REDIRECT_URL")
	failureRedirectURL := os.Getenv("DANA_FAILURE_REDIRECT_URL")

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
		SetPaymentMethod(EWalletIDDana).
		SetChannelProperties(props).
		Build()
}
