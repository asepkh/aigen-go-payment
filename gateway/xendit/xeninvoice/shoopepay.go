package xeninvoice

import (
	xinvoice "github.com/xendit/xendit-go/invoice"

	"github.com/asepkh/aigen-payment/invoice"
)

func NewShopeePay(inv *invoice.Invoice) (*xinvoice.CreateParams, error) {
	return newBuilder(inv).
		AddPaymentMethod("SHOPEEPAY").
		Build()
}
