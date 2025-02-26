package finpay

import (
	"testing"

	"github.com/asepkh/aigen-payment/util/localconfig"
)

func newTestGateway(t *testing.T) *Gateway {
	merchantID := "test-merchant-id"
	creds := localconfig.APICredential{
		SecretKey:  "test-secret-key",
		ClientKey:  "test-client-key",
		MerchantID: &merchantID,
	}
	return NewGateway(creds)
}
