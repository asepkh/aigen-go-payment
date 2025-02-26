package sql

import (
	"context"
	"fmt"

	payment "github.com/asepkh/aigen-go-payment"
	"github.com/asepkh/aigen-go-payment/gateway/finpay"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// NewFinpayTransactionRepository creates a new repository for Finpay transaction status
func NewFinpayTransactionRepository(db *gorm.DB) *FinpayTransactionRepository {
	return &FinpayTransactionRepository{
		DB: db,
	}
}

// FinpayTransactionRepository is a repository for Finpay transaction status
type FinpayTransactionRepository struct {
	DB *gorm.DB
}

// Store saves a Finpay transaction status to the database
func (r *FinpayTransactionRepository) Store(ctx context.Context, status *finpay.TransactionStatus) error {
	log := zerolog.Ctx(ctx).With().Str("function", "FinpayTransactionRepository.Store").Logger()

	if err := r.DB.Save(status).Find(&status).Error; err != nil {
		log.Error().Err(err).Msg("can't save Finpay transaction status")
		return payment.ErrDatabase
	}
	return nil
}

// FindByTransactionID finds a Finpay transaction status by transaction ID
func (r *FinpayTransactionRepository) FindByTransactionID(ctx context.Context, transactionID string) (*finpay.TransactionStatus, error) {
	log := zerolog.Ctx(ctx).With().Str("function", "FinpayTransactionRepository.FindByTransactionID").Logger()

	var status finpay.TransactionStatus
	req := r.DB.
		Where("transaction_id = ?", transactionID).
		First(&status)

	if req.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("payment status for transaction %s %w", transactionID, payment.ErrNotFound)
	}
	if req.Error != nil {
		log.Error().Err(req.Error).Msg("can't find Finpay transaction status")
		return nil, payment.ErrDatabase
	}

	return &status, nil
}

// FindByOrderID finds a Finpay transaction status by order ID
func (r *FinpayTransactionRepository) FindByOrderID(ctx context.Context, orderID string) (*finpay.TransactionStatus, error) {
	log := zerolog.Ctx(ctx).With().Str("function", "FinpayTransactionRepository.FindByOrderID").Logger()

	var status finpay.TransactionStatus
	req := r.DB.
		Where("order_id = ?", orderID).
		First(&status)

	if req.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("payment status for order %s %w", orderID, payment.ErrNotFound)
	}
	if req.Error != nil {
		log.Error().Err(req.Error).Msg("can't find Finpay transaction status")
		return nil, payment.ErrDatabase
	}

	return &status, nil
}
