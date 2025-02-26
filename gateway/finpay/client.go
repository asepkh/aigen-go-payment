package finpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Client represents a Finpay API client
type Client struct {
	BaseURL    string
	SecretKey  string
	ClientKey  string
	MerchantID string
	HTTPClient *http.Client
}

// NewClient creates a new Finpay API client
func NewClient(gateway *Gateway) *Client {
	return &Client{
		BaseURL:    gateway.GetBaseURL(),
		SecretKey:  gateway.secretKey,
		ClientKey:  gateway.clientKey,
		MerchantID: gateway.merchantID,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Response represents a generic Finpay API response
type Response struct {
	Status        string      `json:"status"`
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	TransactionID string      `json:"transaction_id"`
	OrderID       string      `json:"order_id"`
	RedirectURL   string      `json:"redirect_url"`
	PaymentURL    string      `json:"payment_url"`
	Token         string      `json:"token"`
	Data          interface{} `json:"data"`
}

// generateSignature generates a signature for the request
func (c *Client) generateSignature(payload []byte, timestamp string) string {
	stringToSign := c.MerchantID + ":" + timestamp + ":" + string(payload)
	h := hmac.New(sha256.New, []byte(c.SecretKey))
	h.Write([]byte(stringToSign))
	return hex.EncodeToString(h.Sum(nil))
}

// Post makes a POST request to the Finpay API
func (c *Client) Post(endpoint string, payload interface{}) (*Response, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Generate timestamp for signature
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := c.generateSignature(jsonPayload, timestamp)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TIMESTAMP", timestamp)
	req.Header.Set("X-CLIENT-KEY", c.ClientKey)
	req.Header.Set("X-MERCHANT-ID", c.MerchantID)
	req.Header.Set("X-SIGNATURE", signature)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &response, nil
}
