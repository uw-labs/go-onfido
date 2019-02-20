package onfido

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCheckUnmarshal(t *testing.T) {
	t.Parallel()
	testWithStringReports := `
{
    "id": "4b6ff41d-4562-f66c-6sadf-24857b1a380f",
    "created_at": "2019-12-20T14:47:31Z",
    "status": "complete",
    "redirect_uri": null,
    "type": "express",
    "result": "clear",
    "sandbox": true,
    "report_type_groups": [
        "156423"
    ],
    "tags": [],
    "results_uri": "https://onfido.com/dashboard/information_requests/4b6ff41d-4562-f66c-6sadf-24857b1a380f",
    "download_uri": "https://onfido.com/dashboard/pdf/information_requests/4b6ff41d-4562-f66c-6sadf-24857b1a380f",
    "form_uri": null,
    "href": "/v2/applicants/4b6ff41d-4562-f66c-6sadf-24857b1a380f/checks/4b6ff41d-4562-f66c-6sadf-24857b1a380f",
    "reports": [
        "4b6ff41d-4562-f66c-6sadf-24857b1a380f"
    ],
    "paused": false,
    "version": "2.0"
}
`
	testWithFullReports := `
{
    "id": "4b6ff41d-4562-f66c-6sadf-24857b1a380f",
    "created_at": "2019-12-20T14:47:31Z",
    "status": "complete",
    "redirect_uri": null,
    "type": "express",
    "result": "clear",
    "sandbox": true,
    "report_type_groups": [
        "156423"
    ],
    "tags": [],
    "results_uri": "https://onfido.com/dashboard/information_requests/4b6ff41d-4562-f66c-6sadf-24857b1a380f",
    "download_uri": "https://onfido.com/dashboard/pdf/information_requests/4b6ff41d-4562-f66c-6sadf-24857b1a380f",
    "form_uri": null,
    "href": "/v2/applicants/4b6ff41d-4562-f66c-6sadf-24857b1a380f/checks/4b6ff41d-4562-f66c-6sadf-24857b1a380f",
    "reports": [
        {
            "created_at": "2019-02-20T14:57:26Z",
            "href": "/v2/checks/4b6ff41d-4562-f66c-6sadf-24857b1a380f/reports/4b6ff41d-4562-f66c-6sadf-24857b1a380f",
            "id": "4b6ff41d-4562-f66c-6sadf-24857b1a380f",
            "name": "watchlist",
            "properties": {
                "records": []
            },
            "result": "clear",
            "status": "complete",
            "sub_result": null,
            "variant": "kyc",
            "breakdown": {
                "adverse_media": {
                    "result": "clear"
                },
                "sanction": {
                    "result": "clear"
                },
                "politically_exposed_person": {
                    "result": "clear"
                }
            }
        }
    ],
    "paused": false,
    "version": "2.0"
}
`
	tts := []struct {
		name string
		in   string
	}{
		{
			name: "check with report ids",
			in:   testWithStringReports,
		},
		{
			name: "check with full reports",
			in:   testWithFullReports,
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			var c Check
			err := json.Unmarshal([]byte(tt.in), &c)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCreateCheck_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.CreateCheck(context.Background(), "", CheckRequest{})
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestCreateCheck_CheckCreated(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := Check{
		ID:          "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        CheckTypeExpress,
		Status:      "complete",
		Result:      CheckResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Tags:        []string{"my-tag"},
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

	client := NewClient("123")
	client.Endpoint = srv.URL

	reports := []*Report{
		{
			ID:     "7410a943-8f00-43d8-98de-36a774196d86",
			Name:   ReportNameDocument,
			Result: ReportResultClear,
		},
	}

	c, err := client.CreateCheck(context.Background(), applicantID, CheckRequest{
		Type:              expected.Type,
		RedirectURI:       expected.RedirectURI,
		Reports:           reports,
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

	client := NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.GetCheck(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestGetCheck_CheckRetrieved(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := Check{
		ID:          "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        CheckTypeExpress,
		Status:      "complete",
		Result:      CheckResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Tags:        []string{"my-tag"},
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

	client := NewClient("123")
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

	client := NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.ResumeCheck(context.Background(), "")
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestResumeCheck_CheckCreated(t *testing.T) {
	expected := Check{
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

	client := NewClient("123")
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

	client := NewClient("123")
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
	expected := Check{
		ID:          "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        CheckTypeExpress,
		Status:      "complete",
		Result:      CheckResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Tags:        []string{"my-tag"},
	}
	expectedJson, err := json.Marshal(Checks{
		Checks: []*Check{&expected},
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

	client := NewClient("123")
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

func ExampleClient_CreateCheck() {
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

	check, err := client.CreateCheck(ctx, applicant.ID, CheckRequest{
		Type: CheckTypeStandard,
		Reports: []*Report{
			{
				Name: ReportNameDocument,
			},
			{
				Name:    ReportNameIdentity,
				Variant: ReportVariantKYC,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Form: %+v\n", check.FormURI)

}
