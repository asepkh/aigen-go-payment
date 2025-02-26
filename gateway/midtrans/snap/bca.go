package snap

import (
	"github.com/midtrans/midtrans-go/snap"

	"github.com/asepkh/aigen-payment/invoice"
)

func NewBCAVA(inv *invoice.Invoice) (*snap.Request, error) {
	return newBuilder(inv).
		AddPaymentMethods(snap.PaymentTypeBCAVA).
		Build()
}
