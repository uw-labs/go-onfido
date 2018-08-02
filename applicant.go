package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"
)

// IDNumberType represents an ID type (ssn, social insurance, etc)
type IDNumberType string

// Supported ID number types
const (
	IDNumberTypeSSN             IDNumberType = "ssn"
	IDNumberTypeSocialInsurance IDNumberType = "social_insurance"
	IDNumberTypeTaxID           IDNumberType = "tax_id"
	IDNumberTypeIdentityCard    IDNumberType = "identity_card"
	IDNumberTypeDrivingLicense  IDNumberType = "driving_license"
)

// IDNumber represents an ID number from the Onfido API
type IDNumber struct {
	Type      IDNumberType `json:"type,omitempty"`
	Value     string       `json:"value,omitempty"`
	StateCode string       `json:"state_code,omitempty"`
}

// Applicants represents a list of applicants from the Onfido API
type Applicants struct {
	Applicants []*Applicant `json:"applicants"`
}

// Applicant represents an applicant from the Onfido API
type Applicant struct {
	ID                string     `json:"id,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	Sandbox           bool       `json:"sandbox,omitempty"`
	Title             string     `json:"title,omitempty"`
	FirstName         string     `json:"first_name,omitempty"`
	LastName          string     `json:"last_name,omitempty"`
	MiddleName        string     `json:"middle_name,omitempty"`
	Email             string     `json:"email,omitempty"`
	Gender            string     `json:"gender,omitempty"`
	DOB               string     `json:"dob,omitempty"`
	Telephone         string     `json:"telephone,omitempty"`
	Mobile            string     `json:"mobile,omitempty"`
	Country           string     `json:"country,omitempty"`
	MothersMaidenName string     `json:"mothers_maiden_name,omitempty"`
	PreviousLastName  string     `json:"previous_last_name,omitempty"`
	Nationality       string     `json:"nationality,omitempty"`
	CountryOfBirth    string     `json:"country_of_birth,omitempty"`
	TownOfBirth       string     `json:"town_of_birth,omitempty"`
	IDNumbers         []IDNumber `json:"id_numbers,omitempty"`
	Addresses         []Address  `json:"addresses,omitempty"`
}

// CreateApplicant creates a new applicant.
// see https://documentation.onfido.com/?shell#create-applicant
func (c *Client) CreateApplicant(ctx context.Context, a Applicant) (*Applicant, error) {
	jsonStr, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/applicants", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp Applicant
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// DeleteApplicant deletes an applicant by its id.
// see https://documentation.onfido.com/?shell#delete-applicant
func (c *Client) DeleteApplicant(ctx context.Context, id string) error {
	req, err := c.newRequest("DELETE", "/applicants/"+id, nil)
	if err != nil {
		return err
	}
	_, err = c.do(ctx, req, nil)
	return err
}

// GetApplicant retrieves an applicant by its id.
// see https://documentation.onfido.com/?shell#retrieve-applicant
func (c *Client) GetApplicant(ctx context.Context, id string) (*Applicant, error) {
	req, err := c.newRequest("GET", "/applicants/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Applicant
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// ApplicantIter represents an applicant iterator
type ApplicantIter struct {
	*iter
}

// Applicant returns the current applicant on the iterator.
func (i *ApplicantIter) Applicant() *Applicant {
	return i.Current().(*Applicant)
}

// ListApplicants retrieves the list of applicants.
// see https://documentation.onfido.com/?shell#list-applicants
func (c *Client) ListApplicants() *ApplicantIter {
	handler := func(body []byte) ([]interface{}, error) {
		var a Applicants
		if err := json.Unmarshal(body, &a); err != nil {
			return nil, err
		}

		values := make([]interface{}, len(a.Applicants))
		for i, v := range a.Applicants {
			values[i] = v
		}
		return values, nil
	}

	return &ApplicantIter{&iter{
		c:       c,
		nextURL: "/applicants",
		handler: handler,
	}}
}

// UpdateApplicant updates an applicant by its id.
// see https://documentation.onfido.com/?shell#update-applicant
func (c *Client) UpdateApplicant(ctx context.Context, a Applicant) (*Applicant, error) {
	if a.ID == "" {
		return nil, errors.New("invalid applicant id")
	}
	jsonStr, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("PUT", "/applicants/"+a.ID, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp Applicant
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}
