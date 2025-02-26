package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/asepkh/aigen-payment/gateway/finpay"
	"github.com/rs/zerolog/log"
)

// FinpayCallbackHandler handles callbacks from Finpay
func (s *Server) FinpayCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read Finpay callback body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validate the signature
		signature := r.Header.Get("X-Finpay-Signature")
		if !s.validateFinpaySignature(signature, body) {
			log.Error().Msg("Invalid Finpay signature")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Parse the transaction status
		var status finpay.TransactionStatus
		if err := json.Unmarshal(body, &status); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal Finpay callback body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Store the transaction status
		status.RawJSON = string(body)
		if err := s.storeFinpayTransactionStatus(&status); err != nil {
			log.Error().Err(err).Msg("Failed to store Finpay transaction status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Process the invoice based on the transaction status
		if err := s.Manager.ProcessFinpayCallback(r.Context(), &status); err != nil {
			log.Error().Err(err).Msg("Failed to process Finpay callback")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}
}

// validateFinpaySignature validates the signature from Finpay
func (s *Server) validateFinpaySignature(signature string, body []byte) bool {
	// Implementation would depend on Finpay's signature validation mechanism
	// This is a placeholder for the actual implementation
	return true
}

// storeFinpayTransactionStatus stores the transaction status in the database
func (s *Server) storeFinpayTransactionStatus(status *finpay.TransactionStatus) error {
	// Since the Manager.ProcessFinpayCallback method already stores the transaction status,
	// we can simply return nil here. The actual storage will happen in the ProcessFinpayCallback call.
	return nil
}
