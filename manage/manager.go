package manage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	payment "github.com/asepkh/aigen-go-payment"
	"github.com/asepkh/aigen-go-payment/datastore"
	"github.com/asepkh/aigen-go-payment/gateway/finpay"
	fingateway "github.com/asepkh/aigen-go-payment/gateway/finpay"
	midgateway "github.com/asepkh/aigen-go-payment/gateway/midtrans"
	xengateway "github.com/asepkh/aigen-go-payment/gateway/xendit"
	"github.com/asepkh/aigen-go-payment/invoice"
	"github.com/asepkh/aigen-go-payment/subscription"
	"github.com/asepkh/aigen-go-payment/util/localconfig"
)

// NewManager creates a new payment manager
func NewManager(
	config localconfig.Config,
	secret localconfig.PaymentSecret,
) *Manager {
	return &Manager{
		config:          config,
		xenditGateway:   xengateway.NewGateway(secret.Xendit),
		midtransGateway: midgateway.NewGateway(secret.Midtrans),
		finpayGateway:   fingateway.NewGateway(secret.Finpay),
	}
}

type InvoiceEventFunc func(ctx context.Context, i *invoice.Invoice) error

// Manager handle business logic related to payment gateway
type Manager struct {
	config                      localconfig.Config
	xenditGateway               *xengateway.Gateway
	midtransGateway             *midgateway.Gateway
	finpayGateway               *fingateway.Gateway
	midTransactionRepository    datastore.MidtransTransactionStatusRepository
	invoiceRepository           datastore.InvoiceRepository
	subscriptionRepository      datastore.SubscriptionRepository
	paymentConfigRepository     datastore.PaymentConfigReader
	finpayTransactionRepository datastore.FinpayTransactionStatusRepository

	invoiceCreatedCallback   InvoiceEventFunc
	invoiceProcessedCallback InvoiceEventFunc
	invoiceFailedCallback    InvoiceEventFunc
	invoicePaidCallback      InvoiceEventFunc

	// New callback for when payment gateway callbacks are processed
	paymentCallbackProcessedCallback InvoiceEventFunc
}

// MustPaymentCallbackProcessedEventFunc set event handler for emitting payment callback processed event
func (m *Manager) MustPaymentCallbackProcessedEventFunc(fn InvoiceEventFunc) {
	m.paymentCallbackProcessedCallback = fn
}

// MapMidtransTransactionStatusRepository mapping the midtrans transaction status repository
func (m *Manager) MapMidtransTransactionStatusRepository(repo datastore.MidtransTransactionStatusRepository) error {
	m.midTransactionRepository = repo
	return nil
}

// MustMidtransTransactionStatusRepository mandatory mapping the midtrans transaction status repo interface
func (m *Manager) MustMidtransTransactionStatusRepository(repo datastore.MidtransTransactionStatusRepository) {
	if repo == nil {
		panic(fmt.Errorf("midtrans transaction status repository can't be nil"))
	}
	m.midTransactionRepository = repo
}

// MustInvoiceRepository mandatory mapping the invoice repository
func (m *Manager) MustInvoiceRepository(repo datastore.InvoiceRepository) {
	if repo == nil {
		panic(fmt.Errorf("invoice repository can't be nil"))
	}
	m.invoiceRepository = repo
}

// MustSubscriptionRepository mandatory mapping the subscription repository
func (m *Manager) MustSubscriptionRepository(repo datastore.SubscriptionRepository) {
	if repo == nil {
		panic(fmt.Errorf("invoice repository can't be nil"))
	}
	m.subscriptionRepository = repo
}

// MustPaymentConfigReader mandatory mapping for payment config repository
func (m *Manager) MustPaymentConfigReader(repo datastore.PaymentConfigReader) {
	if repo == nil {
		panic(fmt.Errorf("invoice repository can't be nil"))
	}
	m.paymentConfigRepository = repo
}

// MustInvoiceCreatedEventFunc set event handler for emitting invoice created event
func (m *Manager) MustInvoiceCreatedEventFunc(fn InvoiceEventFunc) {
	m.invoiceCreatedCallback = fn
}

// MustInvoicePaidEventFunc set event handler for emitting invoice processed event
func (m *Manager) MustInvoicePaidEventFunc(fn InvoiceEventFunc) {
	m.invoicePaidCallback = fn
}

// MustInvoiceProcessedEventFunc set event handler for emitting invoice processed event
func (m *Manager) MustInvoiceProcessedEventFunc(fn InvoiceEventFunc) {
	m.invoiceProcessedCallback = fn
}

// MustInvoiceFailedEventFunc set event handler for emitting invoice failed event
func (m *Manager) MustInvoiceFailedEventFunc(fn InvoiceEventFunc) {
	m.invoiceFailedCallback = fn
}

// MapFinpayTransactionStatusRepository mapping the finpay transaction status repository
func (m *Manager) MapFinpayTransactionStatusRepository(repo datastore.FinpayTransactionStatusRepository) error {
	m.finpayTransactionRepository = repo
	return nil
}

// MustFinpayTransactionStatusRepository mandatory mapping the finpay transaction status repo interface
func (m *Manager) MustFinpayTransactionStatusRepository(repo datastore.FinpayTransactionStatusRepository) {
	if repo == nil {
		panic(fmt.Errorf("finpay transaction status repository can't be nil"))
	}
	m.finpayTransactionRepository = repo
}

func (m Manager) charger(inv *invoice.Invoice) invoice.PaymentCharger {
	switch payment.NewGateway(inv.Payment.Gateway) {
	case payment.GatewayXendit:
		return &xenditCharger{
			config:        m.config.Xendit,
			XenditGateway: m.xenditGateway,
		}
	case payment.GatewayMidtrans:
		return &midtransCharger{
			MidtransGateway: m.midtransGateway,
		}
	case payment.GatewayFinpay:
		return &finpayCharger{
			FinpayGateway: m.finpayGateway,
		}
	default:
		panic("payment gateway is not found.")
	}
}

// GetPaymentMethods return the payment methods available in payment service
func (m *Manager) GetPaymentMethods(ctx context.Context, opts ...payment.Option) (*PaymentMethodList, error) {
	cfg, err := m.paymentConfigRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	options := payment.Options{}
	for _, o := range opts {
		o(&options)
	}

	return paymentMethodListFromConfig(cfg, options.Price), nil
}

// GetInvoice return invoice given its invoice number
func (m *Manager) GetInvoice(ctx context.Context, number string) (*invoice.Invoice, error) {
	return m.invoiceRepository.FindByNumber(ctx, number)
}

// GenerateInvoice generates new invoice
func (m *Manager) GenerateInvoice(ctx context.Context, gir *GenerateInvoiceRequest) (*invoice.Invoice, error) {

	l := log.Ctx(ctx).With().
		Str("function", "Manager.GenerateInvoice").
		Str("payment_type", string(gir.Payment.PaymentType)).
		Logger()

	var opts []payment.Option
	if gir.Payment.CreditCardDetail != nil {
		opts = append(opts, payment.WithCreditCard(
			gir.Payment.CreditCardDetail.Bank,
			gir.Payment.CreditCardDetail.Installment.Type,
			gir.Payment.CreditCardDetail.Installment.Term,
		))
	}

	l.Debug().Msg("starting to generate the invoice")

	paymentConfig, err := m.paymentConfigRepository.FindByPaymentType(ctx, gir.Payment.PaymentType, opts...)
	if err != nil {
		l.Warn().Err(err).Msg("unable to find the registered payment method")
		return nil, err
	}

	payment, err := invoice.NewPayment(paymentConfig, gir.Payment.PaymentType, gir.Payment.CreditCardDetail)
	if err != nil {
		return nil, err
	}

	dur := 24 * time.Hour
	if gir.Duration.Nanoseconds() != 0 {
		dur = gir.Duration
	} else if paymentConfig.GetPaymentWaitingTime() != nil {
		dur = *paymentConfig.GetPaymentWaitingTime()
	}

	inv := invoice.NewWithDurationLimit(dur)
	l = l.With().
		Str("invoice_number", inv.Number).
		Logger()

	// TODO implement builder pattern to construct invoice
	// TODO check whether it contains the scheme http/https
	if gir.Callback != nil {
		inv.SuccessRedirectURL = gir.Callback.SuccessRedirectURL
		inv.FailureRedirectURL = gir.Callback.FailureRedirectURL
		inv.CallbackURL = gir.Callback.CallbackURL // Set the new callback URL
	}

	var items []invoice.LineItem
	for _, item := range gir.Items {
		i := invoice.NewLineItem(
			item.Name,
			item.Category,
			item.MerchantName,
			item.Description,
			item.Price,
			item.Qty,
			item.Currency,
		)
		items = append(items, *i)
	}
	if err = inv.SetItems(ctx, items); err != nil {
		l.Debug().Err(err).Msg("unable to set items to the invoice")
		return nil, err
	}

	if err = inv.UpsertBillingAddress(gir.Customer.Name, gir.Customer.Email, gir.Customer.PhoneNumber); err != nil {
		l.Debug().Err(err).Msg("unable to set billing address to the invoice")
		return nil, err
	}

	if err = inv.UpdatePaymentMethod(ctx, payment, m.paymentConfigRepository, opts...); err != nil {
		l.Warn().Err(err).Msg("unable to update payment method and all fee")
		return nil, err
	}

	l.Info().Msg("publishing the invoice")
	if err = inv.Publish(ctx); err != nil {
		l.Error().Err(err).Msg("unable to publish the invoice")
		return nil, err
	}

	l.Info().Msg("creating charge request to the payment gateway")
	if err = inv.CreateChargeRequest(ctx, m.charger(inv)); err != nil {
		l.Error().Err(err).Msg("unable to create charge request to payment gateway")
		return nil, err
	}

	if err = m.invoiceRepository.Save(ctx, inv); err != nil {
		return nil, err
	}

	if m.invoiceCreatedCallback != nil {
		go func() {
			err := m.invoiceCreatedCallback(context.Background(), inv)
			if err != nil {
				l.Warn().
					Err(err).
					Msg("failed sending invoice created callback")
			}
		}()
	}

	l.Info().Msg("invoice is created")
	return inv, nil
}

// PayInvoice pays an invoice. Invoice can only be paid if it is in right state
func (m *Manager) PayInvoice(ctx context.Context, pir *PayInvoiceRequest) (*invoice.Invoice, error) {

	log := zerolog.Ctx(ctx).With().
		Str("function", "Manager.PayInvoice").
		Str("invoice_number", pir.InvoiceNumber).
		Logger()

	inv, err := m.invoiceRepository.FindByNumber(ctx, pir.InvoiceNumber)
	if errors.Is(err, payment.ErrNotFound) {
		return nil, nil
	}
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		log.Error().Err(err).Msg("unable to find invoice")
		return nil, err
	}

	log.Info().Msg("tying to complete payment")
	err = inv.Pay(ctx, pir.TransactionID)
	if err != nil {
		log.Error().Err(err).Msg("unable to complete the payment")
		return nil, err
	}

	err = m.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}

	if m.invoicePaidCallback != nil {
		go func() {
			err := m.invoicePaidCallback(context.Background(), inv)
			if err != nil {
				log.Warn().
					Err(err).
					Msg("failed sending invoice paid callback")
			}
		}()
	}

	log.Info().Msg("invoice paid")

	return inv, nil
}

// ProcessInvoice used if payment is initiated from user's end. It's either because they are using VA or any payment
// methods that requires payment action from the user after they choose a payment method/see payment instruction
func (m *Manager) ProcessInvoice(ctx context.Context, invoiceNumber string) (*invoice.Invoice, error) {

	log := zerolog.Ctx(ctx).With().
		Str("function", "Manager.ProcessInvoice").
		Str("invoice_number", invoiceNumber).
		Logger()

	inv, err := m.invoiceRepository.FindByNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("processing invoice")

	err = inv.Process(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to process the invoice")
		return nil, err
	}

	err = m.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}

	if m.invoiceProcessedCallback != nil {
		go func() {
			err := m.invoiceProcessedCallback(context.Background(), inv)
			if err != nil {
				log.Warn().
					Err(err).
					Msg("failed sending invoice processed callback")
			}
		}()
	}

	log.Info().Msg("invoice is processed")
	return inv, nil
}

// FailInvoice fails an invoice if the payment is either failed or expired
func (m *Manager) FailInvoice(ctx context.Context, fir *FailInvoiceRequest) (*invoice.Invoice, error) {

	log := zerolog.Ctx(ctx).With().
		Str("function", "Manager.FailInvoice").
		Str("invoice_number", fir.InvoiceNumber).
		Logger()

	inv, err := m.invoiceRepository.FindByNumber(ctx, fir.InvoiceNumber)
	if errors.Is(err, payment.ErrNotFound) {
		return nil, nil
	}
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return nil, err
	}

	log.Info().Msg("trying to fail the invoice")
	err = inv.Fail(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to fail the invoice")
		return nil, err
	}

	err = m.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}

	if m.invoiceFailedCallback != nil {
		go func() {
			err := m.invoiceFailedCallback(context.Background(), inv)
			if err != nil {
				log.Warn().
					Err(err).
					Msg("failed sending invoice failed callback")
			}
		}()
	}

	log.Info().Msg("invoice is failed")
	return inv, nil
}

// CreateSubscription creates new subscription
func (m *Manager) CreateSubscription(ctx context.Context, csr *CreateSubscriptionRequest) (*subscription.Subscription, error) {
	s := csr.ToSubscription()

	if err := s.Start(ctx, m.subscriptionController(payment.GatewayXendit)); err != nil {
		return nil, err
	}

	if err := m.subscriptionRepository.Save(ctx, s); err != nil {
		return nil, err
	}

	return s, nil
}

// PauseSubscription pause active subscription
func (m *Manager) PauseSubscription(ctx context.Context, subsNumber string) (*subscription.Subscription, error) {

	sub, err := m.subscriptionRepository.FindByNumber(ctx, subsNumber)
	if err != nil {
		return nil, err
	}

	if err := sub.Pause(ctx, m.subscriptionController(payment.GatewayXendit)); err != nil {
		return nil, err
	}

	if err := m.subscriptionRepository.Save(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil

}

// ResumeSubscription resume paused subscription
func (m *Manager) ResumeSubscription(ctx context.Context, subsNumber string) (*subscription.Subscription, error) {

	sub, err := m.subscriptionRepository.FindByNumber(ctx, subsNumber)
	if err != nil {
		return nil, err
	}

	if err := sub.Resume(ctx, m.subscriptionController(payment.GatewayXendit)); err != nil {
		return nil, err
	}

	if err := m.subscriptionRepository.Save(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// StopSubscription stop subscription
func (m *Manager) StopSubscription(ctx context.Context, subsNumber string) (*subscription.Subscription, error) {

	sub, err := m.subscriptionRepository.FindByNumber(ctx, subsNumber)
	if err != nil {
		return nil, err
	}

	if err := sub.Stop(ctx, m.subscriptionController(payment.GatewayXendit)); err != nil {
		return nil, err
	}

	if err := m.subscriptionRepository.Save(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (m Manager) subscriptionController(gateway payment.Gateway) subscription.Controller {
	return &xenditSubscriptionController{
		XenditGateway: m.xenditGateway,
	}
}

// ProcessFinpayCallback processes a Finpay callback
func (m *Manager) ProcessFinpayCallback(ctx context.Context, status *finpay.TransactionStatus) error {
	// Store the transaction status
	if err := m.finpayTransactionRepository.Store(ctx, status); err != nil {
		return fmt.Errorf("failed to store finpay transaction status: %w", err)
	}

	// Find the invoice by order ID
	inv, err := m.invoiceRepository.FindByNumber(ctx, status.OrderID)
	if err != nil {
		return fmt.Errorf("failed to find invoice: %w", err)
	}

	// Process the transaction based on the status
	switch status.TransactionState {
	case "settlement", "capture":
		// Payment successful
		if err := inv.MarkAsPaid(); err != nil {
			return fmt.Errorf("failed to mark invoice as paid: %w", err)
		}

		// Update the invoice in the repository
		if err := m.invoiceRepository.Update(ctx, inv); err != nil {
			return fmt.Errorf("failed to update invoice: %w", err)
		}

		// Trigger the invoice paid callback
		if m.invoicePaidCallback != nil {
			if err := m.invoicePaidCallback(ctx, inv); err != nil {
				return fmt.Errorf("failed to trigger invoice paid callback: %w", err)
			}
		}

	case "deny", "cancel", "expire", "failure":
		// Payment failed
		if err := inv.MarkAsFailed(); err != nil {
			return fmt.Errorf("failed to mark invoice as failed: %w", err)
		}

		// Update the invoice in the repository
		if err := m.invoiceRepository.Update(ctx, inv); err != nil {
			return fmt.Errorf("failed to update invoice: %w", err)
		}

		// Trigger the invoice failed callback
		if m.invoiceFailedCallback != nil {
			if err := m.invoiceFailedCallback(ctx, inv); err != nil {
				return fmt.Errorf("failed to trigger invoice failed callback: %w", err)
			}
		}
	}

	// Trigger the payment callback processed callback
	if m.paymentCallbackProcessedCallback != nil {
		if err := m.paymentCallbackProcessedCallback(ctx, inv); err != nil {
			return fmt.Errorf("failed to trigger payment callback processed callback: %w", err)
		}
	}

	return nil
}
