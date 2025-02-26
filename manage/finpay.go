package manage

import (
	"context"
	"fmt"
	"os"

	payment "github.com/asepkh/aigen-payment"
	"github.com/asepkh/aigen-payment/gateway/finpay"
	"github.com/asepkh/aigen-payment/invoice"
)

type finpayCharger struct {
	FinpayGateway *finpay.Gateway
	MerchantName  string
}

func (c finpayCharger) Create(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {
	// Get merchant name from environment or use a default
	merchantName := os.Getenv("FINPAY_MERCHANT_NAME")
	if merchantName == "" {
		merchantName = "Merchant"
	}

	req, err := finpay.NewFinpayRequestFromInvoice(inv, c.FinpayGateway.MerchantID(), merchantName)
	if err != nil {
		return nil, fmt.Errorf("failed to create finpay request: %w", err)
	}

	client := finpay.NewClient(c.FinpayGateway)

	// Determine the endpoint based on the payment type
	var endpoint string
	switch inv.Payment.PaymentType {
	case payment.SourceBCAVA:
		endpoint = "/v1/va/bca/create"
	case payment.SourceBNIVA:
		endpoint = "/v1/va/bni/create"
	case payment.SourceBRIVA:
		endpoint = "/v1/va/bri/create"
	case payment.SourceMandiriVA:
		endpoint = "/v1/va/mandiri/create"
	case payment.SourcePermataVA:
		endpoint = "/v1/va/permata/create"
	case payment.SourceOtherVA:
		endpoint = "/v1/va/other/create"
	case payment.SourceAlfamart:
		endpoint = "/v1/retail/alfamart/create"
	case payment.SourceQRIS:
		endpoint = "/v1/qris/create"
	case payment.SourceCreditCard:
		endpoint = "/v1/cc/create"
	default:
		return nil, fmt.Errorf("payment type not supported by Finpay")
	}

	response, err := client.Post(endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API call to Finpay: %w", err)
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("Finpay API error: %s - %s", response.StatusCode, response.StatusMessage)
	}

	return &invoice.ChargeResponse{
		TransactionID: response.TransactionID,
		PaymentToken:  response.Token,
		PaymentURL:    response.PaymentURL,
	}, nil
}

func (c finpayCharger) Gateway() payment.Gateway {
	return payment.GatewayFinpay
}
