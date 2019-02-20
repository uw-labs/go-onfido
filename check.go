package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// CheckType represents a check type (express, standard)
type CheckType string

// CheckStatus represents a status of a check
type CheckStatus string

// CheckResult represents a result of a check (clear, consider)
type CheckResult string

// Supported check types
const (
	CheckTypeExpress  CheckType = "express"
	CheckTypeStandard CheckType = "standard"

	CheckStatusInProgress        CheckStatus = "in_progress"
	CheckStatusAwaitingApplicant CheckStatus = "awaiting_applicant"
	CheckStatusComplete          CheckStatus = "complete"
	CheckStatusWithdrawn         CheckStatus = "withdrawn"
	CheckStatusPaused            CheckStatus = "paused"
	CheckStatusReopened          CheckStatus = "reopened"

	CheckResultClear    CheckResult = "clear"
	CheckResultConsider CheckResult = "consider"
)

// CheckRequest represents a check request to Onfido API. It contains the consider field which is used for Sandbox
// Testing of multiple report scenarios. See https://documentation.onfido.com/#sandbox-responses
type CheckRequest struct {
	Type                    CheckType `json:"type"`
	RedirectURI             string    `json:"redirect_uri,omitempty"`
	Reports                 []*Report `json:"reports"`
	Tags                    []string  `json:"tags,omitempty"`
	SupressFormEmails       bool      `json:"suppress_form_emails,omitempty"`
	Async                   bool      `json:"async,omitempty"`
	ChargeApplicantForCheck bool      `json:"charge_applicant_for_check,omitempty"`
	// Used for sandbox testing
	Consider []ReportName `json:"consider,omitempty"`
}

// Check represents a check in Onfido API
type Check struct {
	ID          string      `json:"id,omitempty"`
	CreatedAt   *time.Time  `json:"created_at,omitempty"`
	Href        string      `json:"href,omitempty"`
	Type        CheckType   `json:"type,omitempty"`
	Status      CheckStatus `json:"status,omitempty"`
	Result      CheckResult `json:"result,omitempty"`
	DownloadURI string      `json:"download_uri,omitempty"`
	FormURI     string      `json:"form_uri,omitempty"`
	RedirectURI string      `json:"redirect_uri,omitempty"`
	ResultsURI  string      `json:"results_uri,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
}

// Checks represents a list of checks in Onfido API
type Checks struct {
	Checks []*Check `json:"checks"`
}

// CreateCheck creates a new check for the provided applicant.
// see https://documentation.onfido.com/?shell#create-check
func (c *Client) CreateCheck(ctx context.Context, applicantID string, cr CheckRequest) (*Check, error) {
	jsonStr, err := json.Marshal(cr)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/applicants/"+applicantID+"/checks", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp Check
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// GetCheck retrieves a check for the provided applicant by its ID.
// see https://documentation.onfido.com/?shell#retrieve-check
func (c *Client) GetCheck(ctx context.Context, applicantID, id string) (*Check, error) {
	req, err := c.newRequest("GET", "/applicants/"+applicantID+"/checks/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Check
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// ResumeCheck resumes a paused check by its ID.
// see https://documentation.onfido.com/?shell#resume-check
func (c *Client) ResumeCheck(ctx context.Context, id string) (*Check, error) {
	req, err := c.newRequest("POST", "/checks/"+id+"/resume", nil)
	if err != nil {
		return nil, err
	}

	var resp Check
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// CheckIter represents a check iterator
type CheckIter struct {
	*iter
}

// Check returns the current item in the iterator as a Check.
func (i *CheckIter) Check() *Check {
	return i.Current().(*Check)
}

// ListChecks retrieves the list of checks for the provided applicant.
// see https://documentation.onfido.com/?shell#list-checks
func (c *Client) ListChecks(applicantID string) *CheckIter {
	handler := func(body []byte) ([]interface{}, error) {
		var r Checks
		if err := json.Unmarshal(body, &r); err != nil {
			return nil, err
		}

		values := make([]interface{}, len(r.Checks))
		for i, v := range r.Checks {
			values[i] = v
		}
		return values, nil
	}

	return &CheckIter{&iter{
		c:       c,
		nextURL: "/applicants/" + applicantID + "/checks",
		handler: handler,
	}}
}
