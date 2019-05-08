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

func TestCreateCheck_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
		assert.NoError(t, wErr)
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
		Result:      onfido.CheckResultClear,
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
	expectedJSON, err := json.Marshal(expected)
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
		_, wErr := w.Write(expectedJSON)
		assert.NoError(t, wErr)
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
		_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
		assert.NoError(t, wErr)
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
	expected := onfido.CheckRetrieved{
		ID:          "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        onfido.CheckTypeExpress,
		Status:      "complete",
		Result:      onfido.CheckResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Reports:     []string{"7410a943-8f00-43d8-98de-36a774196d86"},
		Tags:        []string{"my-tag"},
	}
	expectedJSON, err := json.Marshal(expected)
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
		_, wErr := w.Write(expectedJSON)
		assert.NoError(t, wErr)
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
	assert.EqualValues(t, expected.Reports, c.Reports)
}

func TestGetCheckExpanded_NoReports(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := onfido.CheckRetrieved{
		ID:          "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        onfido.CheckTypeExpress,
		Status:      "complete",
		Result:      onfido.CheckResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Reports:     []string{},
		Tags:        []string{"my-tag"},
	}
	expectedJSON, err := json.Marshal(expected)
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
		_, wErr := w.Write(expectedJSON)
		assert.NoError(t, wErr)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	c, err := client.GetCheckExpanded(context.Background(), applicantID, expected.ID)
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
	assert.Len(t, c.Reports, 0)
}

func TestGetCheckExpanded_NonOkResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
		assert.NoError(t, wErr)
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.GetCheckExpanded(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestGetCheckExpanded_HasReports(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	checkID := "ce62d838-56f8-4ea5-98be-e7166d1dc33d"
	report1ID := "1fd6fec0-456f-443a-b75d-b048f47c34f7"
	report2ID := "6ec6c029-469e-4c9e-91f3-beeb3fbc175e"

	expected := onfido.CheckRetrieved{
		ID:          checkID,
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        onfido.CheckTypeExpress,
		Status:      "complete",
		Result:      onfido.CheckResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Reports:     []string{report1ID, report2ID},
		Tags:        []string{"my-tag"},
	}
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	// Expected Report 1
	expectedReport1 := onfido.Report{
		ID:        report1ID,
		Name:      onfido.ReportNameDocument,
		Status:    "complete",
		Result:    onfido.ReportResultClear,
		SubResult: onfido.ReportSubResultClear,
		Variant:   onfido.ReportVariantStandard,
		Href:      "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
	}
	expectedReport1Json, err := json.Marshal(expectedReport1)
	if err != nil {
		t.Fatal(err)
	}

	// Expected Report 2
	expectedReport2 := onfido.Report{
		ID:        report2ID,
		Name:      onfido.ReportNameDocument,
		Status:    "complete",
		Result:    onfido.ReportResultClear,
		SubResult: onfido.ReportSubResultClear,
		Variant:   onfido.ReportVariantStandard,
		Href:      "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
	}
	expectedReport2Json, err := json.Marshal(expectedReport2)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	// Return the requested Report
	m.HandleFunc("/checks/{checkId}/reports/{reportId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, checkID, vars["checkId"])
		assert.Contains(t, expected.Reports, vars["reportId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch vars["reportId"] {
		case report1ID:
			_, wErr := w.Write(expectedReport1Json)
			assert.NoError(t, wErr)
		case report2ID:
			_, wErr := w.Write(expectedReport2Json)
			assert.NoError(t, wErr)
		}
	}).Methods("GET")

	// Return the requested Check
	m.HandleFunc("/applicants/{applicantId}/checks/{checkId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])
		assert.Equal(t, expected.ID, vars["checkId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, wErr := w.Write(expectedJSON)
		assert.NoError(t, wErr)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	c, err := client.GetCheckExpanded(context.Background(), applicantID, expected.ID)
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
	assert.Len(t, c.Reports, 2)
	assert.ElementsMatch(t, c.Reports, []*onfido.Report{&expectedReport1, &expectedReport2})
}

func TestGetCheckExpanded_HasReports_NonOkResponse(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	checkID := "ce62d838-56f8-4ea5-98be-e7166d1dc33d"
	report1ID := "1fd6fec0-456f-443a-b75d-b048f47c34f7"
	report2ID := "returns-error-status"

	expected := onfido.CheckRetrieved{
		ID:          checkID,
		Href:        "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		Type:        onfido.CheckTypeExpress,
		Status:      "complete",
		Result:      onfido.CheckResultClear,
		DownloadURI: "https://onfido.com/dashboard/pdf/1234",
		FormURI:     "https://onfido.com/information/1234",
		RedirectURI: "https://somewhere.else",
		ResultsURI:  "https://onfido.com/dashboard/information_requests/1234",
		Reports:     []string{report1ID, report2ID},
		Tags:        []string{"my-tag"},
	}
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	// Expected Report 1
	expectedReport1 := onfido.Report{
		ID:        report1ID,
		Name:      onfido.ReportNameDocument,
		Status:    "complete",
		Result:    onfido.ReportResultClear,
		SubResult: onfido.ReportSubResultClear,
		Variant:   onfido.ReportVariantStandard,
		Href:      "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
	}
	expectedReport1Json, err := json.Marshal(expectedReport1)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	// Return the requested Report
	m.HandleFunc("/checks/{checkId}/reports/{reportId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, checkID, vars["checkId"])
		assert.Contains(t, expected.Reports, vars["reportId"])

		w.Header().Set("Content-Type", "application/json")

		switch vars["reportId"] {
		case report1ID:
			w.WriteHeader(http.StatusOK)
			_, wErr := w.Write(expectedReport1Json)
			assert.NoError(t, wErr)
		case report2ID:
			w.WriteHeader(http.StatusForbidden)
			_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
			assert.NoError(t, wErr)
		}
	}).Methods("GET")

	// Return the requested Check
	m.HandleFunc("/applicants/{applicantId}/checks/{checkId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])
		assert.Equal(t, expected.ID, vars["checkId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, wErr := w.Write(expectedJSON)
		assert.NoError(t, wErr)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err = client.GetCheckExpanded(context.Background(), applicantID, expected.ID)
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestResumeCheck_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
		assert.NoError(t, wErr)
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
	expectedJSON, err := json.Marshal(expected)
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
		_, wErr := w.Write(expectedJSON)
		assert.NoError(t, wErr)
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
		_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
		assert.NoError(t, wErr)
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
		Result:      onfido.CheckResultClear,
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
	expectedJSON, err := json.Marshal(onfido.Checks{
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
		_, wErr := w.Write(expectedJSON)
		assert.NoError(t, wErr)
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
