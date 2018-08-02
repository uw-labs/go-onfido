package onfido

import (
	"encoding/json"
	"errors"
	"net/url"
)

var (
	// ErrEmptyPostcode means that an empty postcode param was passed
	ErrEmptyPostcode = errors.New("empty postcode")
)

// Addresses represents a list of addresses from the Onfido API
type Addresses struct {
	Addresses []*Address `json:"addresses"`
}

// Address represents an address from the Onfido API
type Address struct {
	FlatNumber     string `json:"flat_number"`
	BuildingNumber string `json:"building_number"`
	BuildingName   string `json:"building_name"`
	Street         string `json:"street"`
	SubStreet      string `json:"sub_street"`
	Town           string `json:"town"`
	State          string `json:"state"`
	Postcode       string `json:"postcode"`
	Country        string `json:"country"`

	// Applicant specific
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// PickerIter represents an address picker iterator
type PickerIter struct {
	*iter
}

// Address returns the current address on the iterator.
func (i *PickerIter) Address() *Address {
	return i.Current().(*Address)
}

// PickAddresses retrieves the list of addresses matched against the provided postcode.
// see https://documentation.onfido.com/?shell#address-picker
func (c *Client) PickAddresses(postcode string) *PickerIter {
	if postcode == "" {
		return &PickerIter{&iter{
			err: ErrEmptyPostcode,
		}}
	}
	handler := func(body []byte) ([]interface{}, error) {
		var a Addresses
		if err := json.Unmarshal(body, &a); err != nil {
			return nil, err
		}

		values := make([]interface{}, len(a.Addresses))
		for i, v := range a.Addresses {
			values[i] = v
		}
		return values, nil
	}

	params := make(url.Values)
	params.Set("postcode", postcode)

	return &PickerIter{&iter{
		c:       c,
		nextURL: "addresses/pick?" + params.Encode(),
		handler: handler,
	}}
}
