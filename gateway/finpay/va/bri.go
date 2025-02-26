package va

import (
	"github.com/asepkh/aigen-payment/invoice"
)

// BRIVARequest represents a request to create a BRI virtual account
type BRIVARequest struct {
	MerchantID    string  `json:"merchant_id"`
	MerchantName  string  `json:"merchant_name"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	CustomerName  string  `json:"customer_name"`
	CustomerEmail string  `json:"customer_email"`
	CustomerPhone string  `json:"customer_phone"`
	Description   string  `json:"description"`
	ExpiryTime    int64   `json:"expiry_time"` // in seconds
	CallbackURL   string  `json:"callback_url"`
	RedirectURL   string  `json:"redirect_url"`
}

// NewBRIVA creates a new BRI virtual account request from an invoice
func NewBRIVA(inv *invoice.Invoice) (*BRIVARequest, error) {
	// Calculate expiry time based on invoice waiting time
	var expiryTime int64 = 86400 // Default 1 day in seconds
	if inv.Payment.WaitingDuration() != nil {
		expiryTime = int64(inv.Payment.WaitingDuration().Seconds())
	}

	return &BRIVARequest{
		MerchantID:    "", // Will be filled by the client
		MerchantName:  "", // Will be filled by the client
		TransactionID: inv.Number,
		Amount:        inv.GetTotal(),
		CustomerName:  inv.BillingAddress.FullName,
		CustomerEmail: inv.BillingAddress.Email,
		CustomerPhone: inv.BillingAddress.PhoneNumber,
		Description:   inv.Title,
		ExpiryTime:    expiryTime,
		CallbackURL:   inv.CallbackURL,        // Will be overridden if empty
		RedirectURL:   inv.SuccessRedirectURL, // Will be overridden if empty
	}, nil
}
