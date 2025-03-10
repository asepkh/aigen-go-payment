//go:generate mockery -dir . -name PaymentConfigReader -output ./mocks -filename config_reader.go

package datastore

import (
	"context"

	"github.com/asepkh/aigen-go-payment/subscription"

	payment "github.com/asepkh/aigen-go-payment"
	"github.com/asepkh/aigen-go-payment/config"
	"github.com/asepkh/aigen-go-payment/gateway/finpay"
	"github.com/asepkh/aigen-go-payment/gateway/midtrans"
	"github.com/asepkh/aigen-go-payment/invoice"
)

// MidtransTransactionStatusRepository is an interface for
// the storage of midtrans transaction status.
type MidtransTransactionStatusRepository interface {
	Save(ctx context.Context, status *midtrans.TransactionStatus) error
	FindByOrderID(ctx context.Context, orderID string) (*midtrans.TransactionStatus, error)
}

// InvoiceRepository is an interface for invoice storage
type InvoiceRepository interface {
	FindByNumber(ctx context.Context, number string) (*invoice.Invoice, error)
	Save(ctx context.Context, invoice *invoice.Invoice) error
	Update(ctx context.Context, invoice *invoice.Invoice) error
}

// PaymentConfigReader is interface for reading payment configuration
type PaymentConfigReader interface {
	FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.Option) (config.FeeConfigReader, error)
	FindAll(ctx context.Context) (*config.PaymentConfig, error)
}

// SubscriptionRepository is an interface for subscription store
type SubscriptionRepository interface {
	Save(ctx context.Context, subs *subscription.Subscription) error
	FindByNumber(ctx context.Context, number string) (*subscription.Subscription, error)
}

// FinpayTransactionStatusRepository is a repository for Finpay transaction status
type FinpayTransactionStatusRepository interface {
	Store(context.Context, *finpay.TransactionStatus) error
	FindByTransactionID(context.Context, string) (*finpay.TransactionStatus, error)
	FindByOrderID(context.Context, string) (*finpay.TransactionStatus, error)
}
