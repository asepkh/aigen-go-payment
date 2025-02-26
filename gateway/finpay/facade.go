package finpay

import (
	"fmt"
	"os"

	payment "github.com/asepkh/aigen-payment"
	"github.com/asepkh/aigen-payment/gateway/finpay/creditcard"
	"github.com/asepkh/aigen-payment/gateway/finpay/qris"
	"github.com/asepkh/aigen-payment/gateway/finpay/retail"
	"github.com/asepkh/aigen-payment/gateway/finpay/va"
	"github.com/asepkh/aigen-payment/invoice"
)

// NewFinpayRequestFromInvoice creates finpay charge request from invoice
func NewFinpayRequestFromInvoice(inv *invoice.Invoice, merchantID string, merchantName string) (interface{}, error) {
	callbackURL := os.Getenv("FINPAY_CALLBACK_URL")
	redirectURL := os.Getenv("FINPAY_SUCCESS_REDIRECT_URL")

	switch inv.Payment.PaymentType {
	case payment.SourceBCAVA:
		req, err := va.NewBCAVA(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	case payment.SourcePermataVA:
		req, err := va.NewPermataVA(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	case payment.SourceMandiriVA:
		req, err := va.NewMandiriVA(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	case payment.SourceBNIVA:
		req, err := va.NewBNIVA(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	case payment.SourceBRIVA:
		req, err := va.NewBRIVA(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	case payment.SourceOtherVA:
		req, err := va.NewOtherVA(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	case payment.SourceAlfamart:
		req, err := retail.NewAlfamart(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	case payment.SourceQRIS:
		req, err := qris.NewQRIS(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	case payment.SourceCreditCard:
		req, err := creditcard.NewCreditCard(inv)
		if err != nil {
			return nil, err
		}
		req.MerchantID = merchantID
		req.MerchantName = merchantName
		req.CallbackURL = callbackURL
		req.RedirectURL = redirectURL
		return req, nil
	default:
		return nil, fmt.Errorf("payment type not supported by Finpay")
	}
}
