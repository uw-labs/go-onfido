package onfido_test

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/go-onfido"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLivePhotos_List(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	createdAt := time.Now()

	expected := onfido.LivePhoto{
		ID:           "541d040b-89f8-444b-8921-16b1333bf1c7",
		CreatedAt:    &createdAt,
		Href:         "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		DownloadHref: "https://onfido.com/photo/pdf/1234",
		FileName:     "something.png",
		FileSize:     1234,
		FileType:     "image/png",
	}
	expectedJson, err := json.Marshal(struct {
		LivePhotos []*onfido.LivePhoto `json:"live_photos"`
	}{
		LivePhotos: []*onfido.LivePhoto{&expected},
	})
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/live_photos", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("applicant_id") != applicantID {
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

	it := client.ListLivePhotos(applicantID)
	for it.Next(context.Background()) {
		c := it.LivePhoto()

		assert.Equal(t, expected.ID, c.ID)
		assert.True(t, expected.CreatedAt.Equal(*c.CreatedAt))
		assert.Equal(t, expected.Href, c.Href)
		assert.Equal(t, expected.DownloadHref, c.DownloadHref)
		assert.Equal(t, expected.FileName, c.FileName)
		assert.Equal(t, expected.FileSize, c.FileSize)
		assert.Equal(t, expected.FileType, c.FileType)
	}
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
}
