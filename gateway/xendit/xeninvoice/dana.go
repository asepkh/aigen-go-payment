package xeninvoice

import (
	xinvoice "github.com/xendit/xendit-go/invoice"

	"github.com/asepkh/aigen-payment/invoice"
)

func NewDana(inv *invoice.Invoice) (*xinvoice.CreateParams, error) {
	return newBuilder(inv).
		AddPaymentMethod("DANA").
		Build()
}
