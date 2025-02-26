package finpay

import (
	"encoding/json"
	"time"

	payment "github.com/asepkh/aigen-go-payment"
)

// TransactionStatus represents finpay transaction status
type TransactionStatus struct {
	payment.Model
	TransactionID    string    `json:"transaction_id" gorm:"index:finpay_transaction_id_idx"`
	OrderID          string    `json:"order_id" gorm:"index:finpay_order_id_idx"`
	PaymentType      string    `json:"payment_type"`
	TransactionTime  time.Time `json:"transaction_time"`
	TransactionState string    `json:"transaction_state"`
	StatusCode       string    `json:"status_code"`
	StatusMessage    string    `json:"status_message"`
	GrossAmount      string    `json:"gross_amount"`
	RawJSON          string    `json:"raw_json" gorm:"type:text"`
}

// TableName returns the table name for the model
func (TransactionStatus) TableName() string {
	return "finpay_transaction_status"
}

// NewTransactionStatusFromJSON creates a new transaction status from JSON
func NewTransactionStatusFromJSON(jsonStr []byte) (*TransactionStatus, error) {
	var ts TransactionStatus
	if err := json.Unmarshal(jsonStr, &ts); err != nil {
		return nil, err
	}
	ts.RawJSON = string(jsonStr)
	return &ts, nil
}
