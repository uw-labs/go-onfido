package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// CheckStatus represents a status of a check
type CheckStatus string

// CheckResult represents a result of a check (clear, consider)
type CheckResult string

// Supported check types
const (
	CheckStatusInProgress        CheckStatus = "in_progress"
	CheckStatusAwaitingApplicant CheckStatus = "awaiting_applicant"
	CheckStatusComplete          CheckStatus = "complete"
	CheckStatusWithdrawn         CheckStatus = "withdrawn"
	CheckStatusPaused            CheckStatus = "paused"
	CheckStatusReopened          CheckStatus = "reopened"

	CheckResultClear    CheckResult = "clear"
	CheckResultConsider CheckResult = "consider"
)

// CheckRequest represents a check request to Onfido API
type CheckRequest struct {
	ApplicantID           string                 `json:"applicant_id"`
	ReportNames           []ReportName           `json:"report_names"`
	DocumentIDs           []string               `json:"document_ids,omitempty"`
	ApplicantProvidesData bool                   `json:"applicant_provides_data,omitempty"`
	Asynchronous          *bool                  `json:"asynchronous,omitempty"`
	RedirectURI           string                 `json:"redirect_uri,omitempty"`
	Tags                  []string               `json:"tags,omitempty"`
	SuppressFormEmails    *bool                  `json:"suppress_form_emails,omitempty"`
	WebhookIDs            []string               `json:"webhook_ids,omitempty"`
	USDriversLicence      map[string]interface{} `json:"us_driving_licence,omitempty"`
	ReportConfiguration   map[string]interface{} `json:"report_configuration,omitempty"`
	// Consider is used for Sandbox Testing of multiple report scenarios.
	// see https://documentation.onfido.com/#sandbox-responses
	Consider []string `json:"consider,omitempty"`
}

// Check represents a check in Onfido API
type Check struct {
	ID                    string      `json:"id,omitempty"`
	CreatedAt             *time.Time  `json:"created_at,omitempty"`
	Href                  string      `json:"href,omitempty"`
	ApplicantID           string      `json:"applicant_id,omitempty"`
	ApplicantProvidesData bool        `json:"applicant_provides_data,omitempty"`
	Status                CheckStatus `json:"status,omitempty"`
	Result                CheckResult `json:"result,omitempty"`
	FormURI               string      `json:"form_uri,omitempty"`
	RedirectURI           string      `json:"redirect_uri,omitempty"`
	ResultsURI            string      `json:"results_uri,omitempty"`
	ReportIDs             []string    `json:"report_ids,omitempty"`
	Tags                  []string    `json:"tags,omitempty"`
	WebhookIDs            []string    `json:"webhook_ids,omitempty"`
	Paused                bool        `json:"paused,omitempty"`
	Sandbox               bool        `json:"sandbox,omitempty"`
}

// CheckExpanded represents a check with expanded report objects
type CheckExpanded struct {
	Check
	Reports []*Report `json:"reports,omitempty"`
}

// Checks represents a list of checks in Onfido API
type Checks struct {
	Checks []*Check `json:"checks"`
}

// CreateCheck creates a new check for the provided applicant.
// see https://documentation.onfido.com/?shell#create-check
func (c *Client) CreateCheck(ctx context.Context, cr CheckRequest) (*Check, error) {
	jsonStr, err := json.Marshal(cr)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/checks", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp Check
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// GetCheck retrieves a check by its ID.
// see https://documentation.onfido.com/?shell#retrieve-check
func (c *Client) GetCheck(ctx context.Context, id string) (*Check, error) {
	req, err := c.newRequest("GET", "/checks/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Check
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// GetCheckExpanded retrieves a check by its ID, with
// the Check's Reports expanded within the returned CheckExpanded object.
// see https://documentation.onfido.com/?shell#retrieve-check (Shell) but refer to the JSON
// response object for https://documentation.onfido.com/?php#check-object (PHP) for the expanded contents.
func (c *Client) GetCheckExpanded(ctx context.Context, id string) (*CheckExpanded, error) {
	// Get the Check object. This only includes Report IDs, not the expanded Report objects.
	check, err := c.GetCheck(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create the expanded check with the base check data
	expanded := &CheckExpanded{
		Check:   *check,
		Reports: make([]*Report, len(check.ReportIDs)),
	}

	// For each Report ID in the Check object, fetch (expand) the Report
	for i, reportID := range check.ReportIDs {
		rep, err := c.GetReport(ctx, reportID)
		if err != nil {
			return nil, err
		}
		expanded.Reports[i] = rep
	}

	return expanded, nil
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

// DownloadCheck downloads a PDF summary of a check by its ID.
// see https://documentation.onfido.com/api/latest/#download-check
func (c *Client) DownloadCheck(ctx context.Context, id string) ([]byte, error) {
	req, err := c.newRequest("GET", "/checks/"+id+"/download", nil)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = c.do(ctx, req, &buf)
	return buf.Bytes(), err
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
		nextURL: "/checks?applicant_id=" + applicantID,
		handler: handler,
	}}
}
