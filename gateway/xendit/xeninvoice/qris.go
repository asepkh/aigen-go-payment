package xeninvoice

import (
	xinvoice "github.com/xendit/xendit-go/invoice"

	"github.com/asepkh/aigen-go-payment/invoice"
)

func NewQRIS(inv *invoice.Invoice) (*xinvoice.CreateParams, error) {
	return newBuilder(inv).
		AddPaymentMethod("LINKAJA").
		AddPaymentMethod("DANA").
		AddPaymentMethod("OVO").
		Build()
}
