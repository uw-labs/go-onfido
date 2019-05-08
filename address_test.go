package onfido_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	onfido "github.com/uw-labs/go-onfido"
)

func TestPickAddresses_EmptyPostcode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
		assert.NoError(t, wErr)
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.PickAddresses("")
	if it.Next(context.Background()) == true {
		t.Fatal("expected iterator not to return next item, got next item")
	}
	if it.Err() != onfido.ErrEmptyPostcode {
		t.Fatal("expected iterator to error with empty postcode")
	}
}

func TestPickAddresses_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
		assert.NoError(t, wErr)
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.PickAddresses("SW1 XAM")
	if it.Next(context.Background()) == true {
		t.Fatal("expected iterator not to return next item, got next item")
	}
	if it.Err() == nil {
		t.Fatalf("expected iterator to return error message, got nil: %s", it.Err())
	}
}

func TestPickAddresses_ApplicantsRetrieved(t *testing.T) {
	expected := onfido.Address{
		BuildingNumber: "20",
		Street:         "Sandbanks Way",
		Town:           "Aldershot",
		Postcode:       "SAP POP",
		Country:        "GBR",
	}
	expectedJSON, err := json.Marshal(onfido.Addresses{
		Addresses: []*onfido.Address{&expected},
	})
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/addresses/pick", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expected.Postcode, r.URL.Query().Get("postcode"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, wErr := w.Write(expectedJSON)
		assert.NoError(t, wErr)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.PickAddresses(expected.Postcode)
	for it.Next(context.Background()) {
		a := it.Address()

		assert.Equal(t, expected.BuildingNumber, a.BuildingNumber)
		assert.Equal(t, expected.Street, a.Street)
		assert.Equal(t, expected.Town, a.Town)
		assert.Equal(t, expected.Postcode, a.Postcode)
		assert.Equal(t, expected.Country, a.Country)
	}
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
}
