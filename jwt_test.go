package onfido

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

func TestNewSdkToken_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	token, err := client.NewSdkToken(context.Background(), "123", "https://*.onfido.com/documentation/*")
	if err == nil {
		t.Fatal("expected to see an error")
	}
	if token != nil {
		t.Fatal("token returned")
	}
}

func TestNewSdkToken_ApplicantsRetrieved(t *testing.T) {
	expected := SdkToken{
		ApplicantID: "klj25h2jk5j4k5jk35",
		Referrer:    "https://*.uw-labs.co.uk/documentation/*",
		Token:       "423423m4n234czxKJKDLF",
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/sdk_token", func(w http.ResponseWriter, r *http.Request) {
		var tk SdkToken
		json.NewDecoder(r.Body).Decode(&tk)
		assert.Equal(t, expected.ApplicantID, tk.ApplicantID)
		assert.Equal(t, expected.Referrer, tk.Referrer)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("POST")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	token, err := client.NewSdkToken(context.Background(), expected.ApplicantID, expected.Referrer)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ApplicantID, token.ApplicantID)
	assert.Equal(t, expected.Referrer, token.Referrer)
	assert.Equal(t, expected.Token, token.Token)
}

func ExampleClient_NewSdkToken() {
	ctx := context.Background()

	client, err := NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	if client.Token.Prod() {
		panic("onfido token is only for production use")
	}

	applicant, err := client.CreateApplicant(ctx, Applicant{
		Email:     "rcrowe@example.co.uk",
		FirstName: "Rob",
		LastName:  "Crowe",
		Addresses: []Address{
			{
				BuildingNumber: "18",
				Street:         "Wind Corner",
				Town:           "Crawley",
				State:          "West Sussex",
				Postcode:       "NW9 5AB",
				Country:        "GBR",
				StartDate:      "2018-02-10",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	t, err := client.NewSdkToken(ctx, applicant.ID, "https://*.onfido.com/documentation/*")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %v\n", t.Token)

	if err := client.DeleteApplicant(ctx, applicant.ID); err != nil {
		panic(err)
	}
}
