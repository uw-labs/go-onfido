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

func TestCreateWebhook_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.CreateWebhook(context.Background(), onfido.WebhookRefRequest{})
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestCreateWebhook_WebhookCreated(t *testing.T) {
	expected := onfido.WebhookRef{
		ID:           "fcb73186-0733-4f6f-9c57-d9d5ef979443",
		URL:          "https://webhookendpoint.url",
		Enabled:      true,
		Href:         "/v2/webhooks/fcb73186-0733-4f6f-9c57-d9d5ef979443",
		Token:        "ExampleToken",
		Environments: []onfido.WebhookEnvironment{onfido.WebhookEnvironmentSandbox, onfido.WebhookEnvironmentLive},
		Events:       []onfido.WebhookEvent{onfido.WebhookEventCheckStarted, onfido.WebhookEventCheckCompleted},
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("POST")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	wh, err := client.CreateWebhook(context.Background(), onfido.WebhookRefRequest{
		URL:          "https://webhookendpoint.url",
		Enabled:      true,
		Environments: []onfido.WebhookEnvironment{onfido.WebhookEnvironmentSandbox, onfido.WebhookEnvironmentLive},
		Events:       []onfido.WebhookEvent{onfido.WebhookEventCheckStarted, onfido.WebhookEventCheckCompleted},
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, wh.ID)
	assert.Equal(t, expected.URL, wh.URL)
	assert.Equal(t, expected.Href, wh.Href)
	assert.Equal(t, expected.Token, wh.Token)
	assert.Equal(t, expected.Enabled, wh.Enabled)
	assert.Equal(t, expected.Environments, wh.Environments)
	assert.Equal(t, expected.Events, wh.Events)
}

func TestListWebhooks_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.ListWebhooks()
	if it.Next(context.Background()) == true {
		t.Fatal("expected iterator not to return next item, got next item")
	}
	if it.Err() == nil {
		t.Fatal("expected iterator to return error message, got nil")
	}
}

func TestListWebhooks_WebhooksRetrieved(t *testing.T) {
	expected := onfido.WebhookRef{
		ID:           "fcb73186-0733-4f6f-9c57-d9d5ef979443",
		URL:          "https://webhookendpoint.url",
		Enabled:      true,
		Href:         "/v2/webhooks/fcb73186-0733-4f6f-9c57-d9d5ef979443",
		Token:        "ExampleToken",
		Environments: []onfido.WebhookEnvironment{onfido.WebhookEnvironmentSandbox, onfido.WebhookEnvironmentLive},
		Events:       []onfido.WebhookEvent{onfido.WebhookEventCheckStarted, onfido.WebhookEventCheckCompleted},
	}
	expectedJson, err := json.Marshal(onfido.WebhookRefs{
		WebhookRefs: []*onfido.WebhookRef{&expected},
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

	it := client.ListWebhooks()
	for it.Next(context.Background()) {
		wh := it.WebhookRef()

		assert.Equal(t, expected.ID, wh.ID)
		assert.Equal(t, expected.URL, wh.URL)
		assert.Equal(t, expected.Href, wh.Href)
		assert.Equal(t, expected.Token, wh.Token)
		assert.Equal(t, expected.Enabled, wh.Enabled)
		assert.Equal(t, expected.Environments, wh.Environments)
		assert.Equal(t, expected.Events, wh.Events)
	}
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
}
