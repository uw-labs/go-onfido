package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

)

func TestUploadDocument_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	docReq := DocumentRequest{
		File: bytes.NewReader([]byte("test")),
		Type: DocumentTypeIDCard,
		Side: DocumentSideFront,
	}

	_, err := client.UploadDocument(context.Background(), "", docReq)
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestUploadDocument_DocumentUploaded(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := Document{
		ID:           "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:         "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		DownloadHref: "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86/download",
		FileName:     "localfile.png",
		FileType:     "png",
		FileSize:     282123,
		Type:         DocumentTypePassport,
		Side:         DocumentSideBack,
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{applicantId}/documents", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("POST")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	d, err := client.UploadDocument(context.Background(), applicantID, DocumentRequest{
		File: bytes.NewReader([]byte("test")),
		Type: expected.Type,
		Side: expected.Side,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, d.ID)
	assert.Equal(t, expected.Href, d.Href)
	assert.Equal(t, expected.DownloadHref, d.DownloadHref)
	assert.Equal(t, expected.FileName, d.FileName)
	assert.Equal(t, expected.FileType, d.FileType)
	assert.Equal(t, expected.FileSize, d.FileSize)
	assert.Equal(t, expected.Type, d.Type)
	assert.Equal(t, expected.Side, d.Side)
}

func TestGetDocument_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.GetDocument(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestGetDocument_DocumentRetrieved(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := Document{
		ID:           "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:         "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		DownloadHref: "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86/download",
		FileName:     "localfile.png",
		FileType:     "png",
		FileSize:     282123,
		Type:         DocumentTypePassport,
		Side:         DocumentSideBack,
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{applicantId}/documents/{documentId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])
		assert.Equal(t, expected.ID, vars["documentId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	d, err := client.GetDocument(context.Background(), applicantID, expected.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, d.ID)
	assert.Equal(t, expected.Href, d.Href)
	assert.Equal(t, expected.DownloadHref, d.DownloadHref)
	assert.Equal(t, expected.FileName, d.FileName)
	assert.Equal(t, expected.FileType, d.FileType)
	assert.Equal(t, expected.FileSize, d.FileSize)
	assert.Equal(t, expected.Type, d.Type)
	assert.Equal(t, expected.Side, d.Side)
}

func TestListDocuments_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	it := client.ListDocuments("")
	if it.Next(context.Background()) == true {
		t.Fatal("expected iterator not to return next item, got next item")
	}
	if it.Err() == nil {
		t.Fatal("expected iterator to return error message, got nil")
	}
}

func TestListDocuments_DocumentsRetrieved(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := Document{
		ID:           "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:         "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		DownloadHref: "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86/download",
		FileName:     "localfile.png",
		FileType:     "png",
		FileSize:     282123,
		Type:         DocumentTypePassport,
		Side:         DocumentSideBack,
	}
	expectedJson, err := json.Marshal(Documents{
		Documents: []*Document{&expected},
	})
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{applicantId}/documents", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := NewClient("123")
	client.Endpoint = srv.URL

	it := client.ListDocuments(applicantID)
	for it.Next(context.Background()) {
		d := it.Document()

		assert.Equal(t, expected.ID, d.ID)
		assert.Equal(t, expected.Href, d.Href)
		assert.Equal(t, expected.DownloadHref, d.DownloadHref)
		assert.Equal(t, expected.FileName, d.FileName)
		assert.Equal(t, expected.FileType, d.FileType)
		assert.Equal(t, expected.FileSize, d.FileSize)
		assert.Equal(t, expected.Type, d.Type)
		assert.Equal(t, expected.Side, d.Side)
	}
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
}
