package onfido

import (
	"bytes"
	"context"
	"encoding/json"
)

// SdkToken represents the response for a request for a JWT token
type SdkToken struct {
	ApplicantID   string `json:"applicant_id,omitempty"`
	Referrer      string `json:"referrer,omitempty"`
	Token         string `json:"token,omitempty"`
	ApplicationID string `json:"application_id,omitempty"`
}

// NewSdkToken returns a JWT token to used by the Javascript SDK
func (c *Client) NewSdkToken(ctx context.Context, id, referrer string) (*SdkToken, error) {
	t := &SdkToken{
		ApplicantID: id,
		Referrer:    referrer,
	}
	return c.sdkTokenRequest(ctx, t)
}

// NewSDKTokenForMobileApp returns a JWT token to used by Onfido's Mobile Application components.
// These require an Application ID rather than a Referrer URL.
// See https://github.com/onfido/onfido-ios-sdk#31-sdk-tokens
func (c *Client) NewSDKTokenForMobileApp(ctx context.Context, id, applicationID string) (*SdkToken, error) {
	t := &SdkToken{
		ApplicantID:   id,
		ApplicationID: applicationID,
	}
	return c.sdkTokenRequest(ctx, t)
}

func (c *Client) sdkTokenRequest(ctx context.Context, t *SdkToken) (*SdkToken, error) {
	jsonStr, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/sdk_token", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp SdkToken
	if _, err := c.do(ctx, req, &resp); err != nil {
		return nil, err
	}
	t.Token = resp.Token
	return t, err
}
