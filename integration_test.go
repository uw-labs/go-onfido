// +build integration

package onfido

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

)

var (
	applicantID string
	documentID  string
)

var onfidoToken = flag.String("onfidoToken", "", "onfido token used for integration tests")

// ------------------------------------------------------------------
// Applicants
// ------------------------------------------------------------------

func TestIntegrationCreateApplicant_ApplicantCreated(t *testing.T) {
	expected := getDefaultApplicant()
	a, err := getOnfidoClient().CreateApplicant(context.Background(), *getDefaultApplicant())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Title, a.Title)
	assert.Equal(t, expected.Email, a.Email)
	assert.Equal(t, expected.FirstName, a.FirstName)
	assert.Equal(t, expected.LastName, a.LastName)
	assert.Equal(t, expected.Addresses, a.Addresses)
	assert.Equal(t, expected.IDNumbers, a.IDNumbers)

	applicantID = a.ID
}

func TestIntegrationGetApplicant_ApplicantRetrieved(t *testing.T) {
	if applicantID == "" {
		t.Skip("no applicant ID set, check applicant created test. skipping")
	}

	expected := getDefaultApplicant()
	client := getOnfidoClient()

	a, err := client.GetApplicant(context.Background(), applicantID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Title, a.Title)
	assert.Equal(t, expected.Email, a.Email)
	assert.Equal(t, expected.FirstName, a.FirstName)
	assert.Equal(t, expected.LastName, a.LastName)
	assert.Equal(t, expected.Addresses, a.Addresses)
	assert.Equal(t, expected.IDNumbers, a.IDNumbers)
}

func TestIntegrationUpdateApplicant_ApplicantUpdated(t *testing.T) {
	if applicantID == "" {
		t.Skip("no applicant ID set, check applicant created test. skipping")
	}

	expected := getDefaultApplicant()
	expected.ID = applicantID
	expected.FirstName = "John"
	expected.LastName = "Doe"

	a, err := getOnfidoClient().UpdateApplicant(context.Background(), *expected)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.FirstName, a.FirstName)
	assert.Equal(t, expected.LastName, a.LastName)
}

func TestIntegrationListApplicants_ApplicantsListed(t *testing.T) {
	client := getOnfidoClient()
	iterated := false

	it := client.ListApplicants()
	for it.Next(context.Background()) {
		a := it.Applicant()
		assert.NotEmpty(t, a.FirstName)
		assert.NotEmpty(t, a.LastName)
		iterated = true
	}
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
	if !iterated {
		t.Fatal("no applicant returned by iterator")
	}
}

// ------------------------------------------------------------------
// Documents
// ------------------------------------------------------------------

func TestIntegrationUploadDocument_DocumentUploaded(t *testing.T) {
	if applicantID == "" {
		t.Skip("no applicant ID set, check applicant created test. skipping")
	}

	file, err := os.Open("./examples/id-card.jpg")
	if err != nil {
		t.Fatal(err)
	}

	expected := getDefaultDocument()
	d, err := getOnfidoClient().UploadDocument(context.Background(), applicantID, *expected)
	if err != nil {
		t.Fatal(err)
	}

	if d.FileName != filepath.Base(file.Name()) || d.Type != expected.Type || d.Side != expected.Side {
		t.Fatalf(
			"document uploaded did not match expected\nexpected filename: %s, type: %s, side: %s\ngot filename: %s, type: %s, side: %s",
			filepath.Base(file.Name()),
			expected.Type,
			expected.Side,
			d.FileName,
			d.Type,
			d.Side,
		)
	}

	documentID = d.ID
}

func TestIntegrationGetDocument_DocumentRetrieved(t *testing.T) {
	if documentID == "" {
		t.Skip("no document ID set, check document upload test. skipping")
	}

	expected := getDefaultDocument()
	file := expected.File.(*os.File)
	d, err := getOnfidoClient().GetDocument(context.Background(), applicantID, documentID)
	if err != nil {
		t.Fatal(err)
	}

	if d.FileName != filepath.Base(file.Name()) || d.Type != expected.Type || d.Side != expected.Side {
		t.Fatalf(
			"document uploaded did not match expected\nexpected filename: %s, type: %s, side: %s\ngot filename: %s, type: %s, side: %s",
			filepath.Base(file.Name()),
			expected.Type,
			expected.Side,
			d.FileName,
			d.Type,
			d.Side,
		)
	}
}

// ------------------------------------------------------------------
// Delete Applicant
// ------------------------------------------------------------------

func TestIntegrationDeleteApplicant_ApplicantDeleted(t *testing.T) {
	if applicantID == "" {
		t.Skip("no applicant ID set, check applicant created test. skipping")
	}

	if err := getOnfidoClient().DeleteApplicant(context.Background(), applicantID); err != nil {
		t.Fatal(err)
	}
	applicantID = ""
}

// ------------------------------------------------------------------
// Helpers
// ------------------------------------------------------------------

func getDefaultApplicant() *onfido.Applicant {
	return &onfido.Applicant{
		Title:     "Mr",
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "foo@bar.com",
		IDNumbers: []onfido.IDNumber{
			{
				Type:  onfido.IDNumberTypeIdentityCard,
				Value: "1234567",
			},
		},
		Addresses: []onfido.Address{
			{
				FlatNumber: "10",
				Street:     "Baker Street",
				Town:       "London",
				Postcode:   "W1U 8ED",
				Country:    "GBR",
				StartDate:  "2017-12-05",
			},
		},
	}
}

func getDefaultDocument() *onfido.DocumentRequest {
	file, err := os.Open("./examples/id-card.jpg")
	if err != nil {
		panic(err)
	}

	return &onfido.DocumentRequest{
		File: file,
		Type: onfido.DocumentTypeIDCard,
		Side: onfido.DocumentSideFront,
	}
}

func getOnfidoClient() *onfido.Client {
	if *onfidoToken == "" {
		panic("onfido token not set")
	}
	client := onfido.NewClient(*onfidoToken)
	if client.Token.Prod() {
		panic("do not use a production token for integration tests")
	}
	return client
}
