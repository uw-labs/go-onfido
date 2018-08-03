package onfido_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/go-onfido"
)

func TestCreateCheck_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.CreateCheck(context.Background(), "", onfido.CheckRequest{})
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestCreateCheck_CheckCreated(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := onfido.Check{
		ID:          "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        onfido.CheckTypeExpress,
		Status:      "complete",
		Result:      onfido.ReportResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Reports: []*onfido.Report{
			{
				ID:     "7410a943-8f00-43d8-98de-36a774196d86",
				Name:   onfido.ReportNameDocument,
				Result: onfido.ReportResultClear,
			},
		},
		Tags: []string{"my-tag"},
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}/checks", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] != applicantID {
			t.Fatal("expected applicant id was not in the request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("POST")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	c, err := client.CreateCheck(context.Background(), applicantID, onfido.CheckRequest{
		Type:              expected.Type,
		RedirectURI:       expected.RedirectURI,
		Reports:           expected.Reports,
		Tags:              expected.Tags,
		SupressFormEmails: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, c.ID)
	assert.Equal(t, expected.Href, c.Href)
	assert.Equal(t, expected.Type, c.Type)
	assert.Equal(t, expected.Status, c.Status)
	assert.Equal(t, expected.Result, c.Result)
	assert.Equal(t, expected.DownloadURI, c.DownloadURI)
	assert.Equal(t, expected.FormURI, c.FormURI)
	assert.Equal(t, expected.RedirectURI, c.RedirectURI)
	assert.Equal(t, expected.ResultsURI, c.ResultsURI)
}

func TestGetCheck_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.GetCheck(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestGetCheck_CheckRetrieved(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := onfido.Check{
		ID:          "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        onfido.CheckTypeExpress,
		Status:      "complete",
		Result:      onfido.ReportResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Reports: []*onfido.Report{
			{
				ID:     "7410a943-8f00-43d8-98de-36a774196d86",
				Name:   onfido.ReportNameDocument,
				Result: onfido.ReportResultClear,
			},
		},
		Tags: []string{"my-tag"},
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{applicantId}/checks/{checkId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])
		assert.Equal(t, expected.ID, vars["checkId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	c, err := client.GetCheck(context.Background(), applicantID, expected.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, c.ID)
	assert.Equal(t, expected.Href, c.Href)
	assert.Equal(t, expected.Type, c.Type)
	assert.Equal(t, expected.Status, c.Status)
	assert.Equal(t, expected.Result, c.Result)
	assert.Equal(t, expected.DownloadURI, c.DownloadURI)
	assert.Equal(t, expected.FormURI, c.FormURI)
	assert.Equal(t, expected.RedirectURI, c.RedirectURI)
	assert.Equal(t, expected.ResultsURI, c.ResultsURI)
}

func TestResumeCheck_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.ResumeCheck(context.Background(), "")
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestResumeCheck_CheckCreated(t *testing.T) {
	expected := onfido.Check{
		ID:     "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Status: "in_progress",
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/checks/{id}/resume", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] != expected.ID {
			t.Fatal("expected check id was not in the request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("POST")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	c, err := client.ResumeCheck(context.Background(), expected.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, c.ID)
	assert.Equal(t, expected.Status, c.Status)
}

func TestListChecks_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.ListChecks("")
	if it.Next(context.Background()) == true {
		t.Fatal("expected iterator not to return next item, got next item")
	}
	if it.Err() == nil {
		t.Fatal("expected iterator to return error message, got nil")
	}
}

func TestListChecks_ChecksRetrieved(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := onfido.Check{
		ID:          "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        onfido.CheckTypeExpress,
		Status:      "complete",
		Result:      onfido.ReportResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Reports: []*onfido.Report{
			{
				ID:     "7410a943-8f00-43d8-98de-36a774196d86",
				Name:   onfido.ReportNameDocument,
				Result: onfido.ReportResultClear,
			},
		},
		Tags: []string{"my-tag"},
	}
	expectedJson, err := json.Marshal(onfido.Checks{
		Checks: []*onfido.Check{&expected},
	})
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{id}/checks", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] != applicantID {
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

	it := client.ListChecks(applicantID)
	for it.Next(context.Background()) {
		c := it.Check()

		assert.Equal(t, expected.ID, c.ID)
		assert.Equal(t, expected.Href, c.Href)
		assert.Equal(t, expected.Type, c.Type)
		assert.Equal(t, expected.Status, c.Status)
		assert.Equal(t, expected.Result, c.Result)
		assert.Equal(t, expected.DownloadURI, c.DownloadURI)
		assert.Equal(t, expected.FormURI, c.FormURI)
		assert.Equal(t, expected.RedirectURI, c.RedirectURI)
		assert.Equal(t, expected.ResultsURI, c.ResultsURI)
	}
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
}
