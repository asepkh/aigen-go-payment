package creditcard

import (
	"github.com/asepkh/aigen-payment/invoice"
)

// CreditCardRequest represents a request to create a credit card payment
type CreditCardRequest struct {
	MerchantID    string  `json:"merchant_id"`
	MerchantName  string  `json:"merchant_name"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	CustomerName  string  `json:"customer_name"`
	CustomerEmail string  `json:"customer_email"`
	CustomerPhone string  `json:"customer_phone"`
	Description   string  `json:"description"`
	CallbackURL   string  `json:"callback_url"`
	RedirectURL   string  `json:"redirect_url"`
}

// NewCreditCard creates a new credit card payment request from an invoice
func NewCreditCard(inv *invoice.Invoice) (*CreditCardRequest, error) {
	return &CreditCardRequest{
		MerchantID:    "", // Will be filled by the client
		MerchantName:  "", // Will be filled by the client
		TransactionID: inv.Number,
		Amount:        inv.GetTotal(),
		CustomerName:  inv.BillingAddress.FullName,
		CustomerEmail: inv.BillingAddress.Email,
		CustomerPhone: inv.BillingAddress.PhoneNumber,
		Description:   inv.Title,
		CallbackURL:   inv.CallbackURL,        // Will be overridden if empty
		RedirectURL:   inv.SuccessRedirectURL, // Will be overridden if empty
	}, nil
}
