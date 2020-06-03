package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// WebhookEnvironment represents an environment type (see `WebhookEnvironment*` constants for possible values)
type WebhookEnvironment string

// WebhookEvent represents an event type (see `WebhookEvent*` constants for possible values)
type WebhookEvent string

// Constants
const (
	WebhookEnvironmentSandbox WebhookEnvironment = "sandbox"
	WebhookEnvironmentLive    WebhookEnvironment = "live"

	WebhookEventReportWithdrawn        WebhookEvent = "report.withdrawn"
	WebhookEventReportResumed          WebhookEvent = "report.resumed"
	WebhookEventReportCancelled        WebhookEvent = "report.cancelled"
	WebhookEventReportAwaitingApproval WebhookEvent = "report.awaiting_approval"
	WebhookEventReportInitiated        WebhookEvent = "report.initiated"
	WebhookEventReportCompleted        WebhookEvent = "report.completed"
	WebhookEventCheckStarted           WebhookEvent = "check.started"
	WebhookEventCheckReopened          WebhookEvent = "check.reopened"
	WebhookEventCheckWithdrawn         WebhookEvent = "check.withdrawn"
	WebhookEventCheckCompleted         WebhookEvent = "check.completed"
	WebhookEventCheckFormOpened        WebhookEvent = "check.form_opened"
	WebhookEventCheckFormCompleted     WebhookEvent = "check.form_completed"
)

// WebhookRefRequest represents a webhook request to Onfido API
type WebhookRefRequest struct {
	URL          string               `json:"url"` // Onfido requires that this must be HTTPS
	Enabled      bool                 `json:"enabled"`
	Environments []WebhookEnvironment `json:"environments,omitempty"` // If omitted then Onfido will default to both
	Events       []WebhookEvent       `json:"events,omitempty"`       // If omitted then Onfido will default to all
}

// WebhookRef represents a webhook in Onfido API
type WebhookRef struct {
	ID           string               `json:"id,omitempty"`
	URL          string               `json:"url,omitempty"`
	Enabled      bool                 `json:"enabled"`
	Href         string               `json:"href,omitempty"`
	Token        string               `json:"token,omitempty"`
	Environments []WebhookEnvironment `json:"environments,omitempty"`
	Events       []WebhookEvent       `json:"events,omitempty"`
}

// WebhookRefs represents a list of webhooks in Onfido API
type WebhookRefs struct {
	WebhookRefs []*WebhookRef `json:"webhooks"`
}

// CreateWebhook register a new webhook.
// see https://documentation.onfido.com/#register-webhook
func (c *Client) CreateWebhook(ctx context.Context, wr WebhookRefRequest) (*WebhookRef, error) {
	jsonStr, err := json.Marshal(wr)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/webhooks", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp WebhookRef
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// UpdateWebhook updates a previously created webhook.
// https://documentation.onfido.com/v2/#edit-webhook
func (c *client) UpdateWebhook(ctx context.Context, id string, wr WebhookRefRequest) (*WebhookRef, error) {
	jsonStr, err := json.Marshal(wr)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(http.MethodPut, "/webhooks/"+id, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp WebhookRef
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// WebhookRefIter represents a webhook iterator
type WebhookRefIter struct {
	*iter
}

// WebhookRef returns the current item in the iterator as a WebhookRef.
func (i *WebhookRefIter) WebhookRef() *WebhookRef {
	return i.Current().(*WebhookRef)
}

// ListWebhooks retrieves the list of webhooks.
// see https://documentation.onfido.com/#list-webhooks
func (c *Client) ListWebhooks() *WebhookRefIter {
	handler := func(body []byte) ([]interface{}, error) {
		var r WebhookRefs
		if err := json.Unmarshal(body, &r); err != nil {
			return nil, err
		}

		values := make([]interface{}, len(r.WebhookRefs))
		for i, v := range r.WebhookRefs {
			values[i] = v
		}
		return values, nil
	}

	return &WebhookRefIter{&iter{
		c:       c,
		nextURL: "/webhooks/",
		handler: handler,
	}}
}
