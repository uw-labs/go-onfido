package onfido_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
	onfido "github.com/uw-labs/go-onfido"
)

func TestNewSdkToken_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, wErr := w.Write([]byte("{\"error\": \"things went bad\"}"))
		assert.NoError(t, wErr)
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
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
	testCases := []onfido.SdkToken{
		{
			ApplicantID: "klj25h2jk5j4k5jk35",
			Referrer:    "https://*.uw-labs.co.uk/documentation/*",
			Token:       "423423m4n234czxKJKDLF",
		},
		{
			ApplicantID:   "maf92h1qa5j4g3si34",
			ApplicationID: "com.ios.application",
			Token:         "534534m4n234czxQIKKLF",
		},
	}

	for i := range testCases {
		expected := testCases[i] // pinning demanded by golint!
		expectedJSON, err := json.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}

		m := mux.NewRouter()
		m.HandleFunc("/sdk_token", func(w http.ResponseWriter, r *http.Request) {
			var tk onfido.SdkToken
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&tk))
			assert.Equal(t, expected.ApplicantID, tk.ApplicantID)
			assert.Equal(t, expected.Referrer, tk.Referrer)
			assert.Equal(t, expected.ApplicationID, tk.ApplicationID)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, wErr := w.Write(expectedJSON)
			assert.NoError(t, wErr)
		}).Methods("POST")
		srv := httptest.NewServer(m)
		defer srv.Close()

		client := onfido.NewClient("123")
		client.Endpoint = srv.URL

		var sdkToken *onfido.SdkToken
		switch {
		case expected.Referrer != "":
			sdkToken, err = client.NewSdkToken(context.Background(), expected.ApplicantID, expected.Referrer)
			if err != nil {
				t.Fatal(err)
			}
		case expected.ApplicationID != "":
			sdkToken, err = client.NewSDKTokenForMobileApp(context.Background(), expected.ApplicantID, expected.ApplicationID)
			if err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatal("neither a Referrer or Application ID was specified in the Token request")
		}

		assert.Equal(t, expected.ApplicantID, sdkToken.ApplicantID)
		assert.Equal(t, expected.Referrer, sdkToken.Referrer)
		assert.Equal(t, expected.ApplicationID, sdkToken.ApplicationID)
		assert.Equal(t, expected.Token, sdkToken.Token)
	}
}
