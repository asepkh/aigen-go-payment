package finpay

import (
	"os"

	"github.com/asepkh/aigen-payment/util/localconfig"
)

// NewGateway creates new finpay payment gateway
func NewGateway(creds localconfig.APICredential) *Gateway {
	gateway := Gateway{
		secretKey: creds.SecretKey,
		clientKey: creds.ClientKey,
	}

	// Handle merchantID which is a pointer
	if creds.MerchantID != nil {
		gateway.merchantID = *creds.MerchantID
	}

	// Set environment (production or sandbox)
	switch os.Getenv("ENVIRONMENT") {
	case "prod":
		gateway.isProduction = true
	default:
		gateway.isProduction = false
	}

	return &gateway
}

// Gateway stores finpay gateway and client
type Gateway struct {
	secretKey    string
	clientKey    string
	merchantID   string
	isProduction bool
}

// NotificationValidationKey returns finpay server key used for validating
// finpay transaction status
func (g Gateway) NotificationValidationKey() string {
	return g.secretKey
}

// IsProduction returns whether the gateway is in production mode
func (g Gateway) IsProduction() bool {
	return g.isProduction
}

// GetBaseURL returns the base URL for API calls
func (g Gateway) GetBaseURL() string {
	if g.isProduction {
		return "https://api.finpay.co.id"
	}
	return "https://sandbox-api.finpay.co.id"
}

// MerchantID returns the merchant ID
func (g Gateway) MerchantID() string {
	return g.merchantID
}

// ClientKey returns the client key
func (g Gateway) ClientKey() string {
	return g.clientKey
}

// SecretKey returns the secret key
func (g Gateway) SecretKey() string {
	return g.secretKey
}
