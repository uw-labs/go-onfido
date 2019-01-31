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

func TestCreateApplicant_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.CreateApplicant(context.Background(), onfido.Applicant{})
	if err == nil {
		t.Fatal()
	}
}

func TestCreateApplicant_ApplicantCreated(t *testing.T) {
	expected := onfido.Applicant{
		ID:        "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Title:     "Mr",
		FirstName: "Foo",
		LastName:  "Bar",
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("POST")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	a, err := client.CreateApplicant(context.Background(), onfido.Applicant{
		Title:     expected.Title,
		FirstName: expected.FirstName,
		LastName:  expected.LastName,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, a.ID)
	assert.Equal(t, expected.Title, a.Title)
	assert.Equal(t, expected.FirstName, a.FirstName)
	assert.Equal(t, expected.LastName, a.LastName)
}

func TestDeleteApplicant_NonOKResponse(t *testing.T) {
	expected := "65643"

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] != expected {
			t.Fatal("expected applicant id was not in the request")
		}
		w.WriteHeader(http.StatusInternalServerError)
	}).Methods("DELETE")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	err := client.DeleteApplicant(context.Background(), expected)
	if err == nil {
		t.Fatal()
	}
}

func TestDeleteApplicant_ValidRequest(t *testing.T) {
	expected := "65643"

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] != expected {
			t.Fatal("expected applicant id was not in the request")
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	err := client.DeleteApplicant(context.Background(), expected)
	if err != nil {
		t.Fatal()
	}
}

func TestGetApplicant_NonOKResponse(t *testing.T) {
	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.GetApplicant(context.Background(), "12432")
	if err == nil {
		t.Fatal()
	}
}

func TestGetApplicant_ValidRequest(t *testing.T) {
	expected := onfido.Applicant{
		ID:        "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Title:     "Mr",
		FirstName: "Foo",
		LastName:  "Bar",
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] != expected.ID {
			t.Fatal("expected applicant id was not in the request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	a, err := client.GetApplicant(context.Background(), expected.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, a.ID)
	assert.Equal(t, expected.Title, a.Title)
	assert.Equal(t, expected.FirstName, a.FirstName)
	assert.Equal(t, expected.LastName, a.LastName)
}

func TestListApplicants_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.ListApplicants()
	if it.Next(context.Background()) == true {
		t.Fatal("expected iterator not to return next item, got next item")
	}
	if it.Err() == nil {
		t.Fatal("expected iterator to return error message, got nil")
	}
}

func TestListApplicants_ApplicantsRetrieved(t *testing.T) {
	expected := onfido.Applicant{
		ID:        "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Title:     "Mr",
		FirstName: "Foo",
		LastName:  "Bar",
	}
	expectedJson, err := json.Marshal(onfido.Applicants{
		Applicants: []*onfido.Applicant{&expected},
	})
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.ListApplicants()
	for it.Next(context.Background()) {
		a := it.Applicant()

		assert.Equal(t, expected.ID, a.ID)
		assert.Equal(t, expected.Title, a.Title)
		assert.Equal(t, expected.FirstName, a.FirstName)
		assert.Equal(t, expected.LastName, a.LastName)
	}
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
}

func TestUpdateApplicant_IDNotSet(t *testing.T) {
	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("PUT")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.UpdateApplicant(context.Background(), onfido.Applicant{})
	if err == nil {
		t.Fatal(err)
	}
}

func TestUpdateApplicant_NonOKResponse(t *testing.T) {
	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
	}).Methods("PUT")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.UpdateApplicant(context.Background(), onfido.Applicant{ID: "3534"})
	if err == nil {
		t.Fatal(err)
	}
}

func TestUpdateApplicant_ValidRequest(t *testing.T) {
	expected := onfido.Applicant{
		ID:        "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Title:     "Mr",
		FirstName: "Foo",
		LastName:  "Bar",
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] != expected.ID {
			t.Fatal("expected applicant id was not in the request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("PUT")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	a, err := client.UpdateApplicant(context.Background(), expected)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, a.ID)
	assert.Equal(t, expected.Title, a.Title)
	assert.Equal(t, expected.FirstName, a.FirstName)
	assert.Equal(t, expected.LastName, a.LastName)
}
