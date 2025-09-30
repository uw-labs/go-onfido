package onfido

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

// Constants
const (
	WebhookSignatureHeader = "X-Signature"
	WebhookTokenEnv        = "ONFIDO_WEBHOOK_TOKEN"
)

// Webhook errors
var (
	ErrInvalidWebhookSignature = errors.New("invalid request, payload hash doesn't match signature")
	ErrMissingWebhookToken     = errors.New("webhook token not found in environmental variable")
)

// Webhook represents a webhook handler
type Webhook struct {
	Token                   string
	SkipSignatureValidation bool
}

// WebhookRequest represents an incoming webhook request from Onfido
type WebhookRequest struct {
	Payload struct {
		ResourceType string `json:"resource_type"`
		Action       string `json:"action"`
		Object       struct {
			ID                 string `json:"id"`
			Status             string `json:"status"`
			CompletedAt        string `json:"completed_at"`         // Deprecated in v3, use CompletedAtISO8601
			CompletedAtISO8601 string `json:"completed_at_iso8601"` // New in v3
			Href               string `json:"href"`
		} `json:"object"`
	} `json:"payload"`
}

// NewWebhookFromEnv creates a new webhook handler using
// configuration from environment variables.
func NewWebhookFromEnv() (*Webhook, error) {
	token := os.Getenv(WebhookTokenEnv)
	if token == "" {
		return nil, ErrMissingWebhookToken
	}
	return NewWebhook(token), nil
}

// NewWebhook creates a new webhook handler
func NewWebhook(token string) *Webhook {
	return &Webhook{
		Token: token,
	}
}

// ValidateSignature validates the request body against the signature header.
func (wh *Webhook) ValidateSignature(body []byte, signature string) error {
	mac := hmac.New(sha1.New, []byte(wh.Token))
	if _, err := mac.Write(body); err != nil {
		return err
	}

	sig, err := hex.DecodeString(signature)
	if err != nil || !hmac.Equal(sig, mac.Sum(nil)) {
		return ErrInvalidWebhookSignature
	}

	return nil
}

// ParseFromRequest parses the webhook request body and returns
// it as WebhookRequest if the request signature is valid.
func (wh *Webhook) ParseFromRequest(req *http.Request) (*WebhookRequest, error) {
	signature := req.Header.Get(WebhookSignatureHeader)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	if !wh.SkipSignatureValidation {
		if err := wh.ValidateSignature(body, signature); err != nil {
			return nil, err
		}
	}

	var wr WebhookRequest
	if err := json.Unmarshal(body, &wr); err != nil {
		return nil, err
	}

	return &wr, nil
}
