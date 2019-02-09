package onfido

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestNewWebhookFromEnv_MissingToken(t *testing.T) {
	_, err := NewWebhookFromEnv()
	if err == nil {
		t.Fatal()
	}
	if err != ErrMissingWebhookToken {
		t.Fatal("expected error to match ErrMissingWebhookToken")
	}
}

func TestNewWebhookFromEnv_TokenSet(t *testing.T) {
	expected := "808yup"
	os.Setenv(WebhookTokenEnv, expected)
	defer os.Setenv(WebhookTokenEnv, "")

	wh, err := NewWebhookFromEnv()
	if err != nil {
		t.Fatal()
	}
	if wh.Token != expected {
		t.Fatalf("expected to see `%s` token but got `%s`", expected, wh.Token)
	}
}

func TestValidateSignature_InvalidSignature(t *testing.T) {
	wh := Webhook{Token: "abc123"}
	err := wh.ValidateSignature([]byte("hello world"), "invalid")
	if err == nil {
		t.Fatal()
	}
	if err != ErrInvalidWebhookSignature {
		t.Fatal("expected error to match ErrInvalidWebhookSignature")
	}
}

func TestValidateSignature_ValidSignature(t *testing.T) {
	wh := Webhook{Token: "abc123"}
	err := wh.ValidateSignature([]byte("hello world"), "fcc98c5b4f306cfe6b5b8fcce03ddb33fc13ae6b")
	if err != nil {
		t.Fatal()
	}
}

func TestParseFromRequest_InvalidSignature(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(WebhookSignatureHeader, "123")
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("hello world")))

	wh := Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err == nil {
		t.Fatal()
	}
	if err != ErrInvalidWebhookSignature {
		t.Fatal("expected error to match ErrInvalidWebhookSignature")
	}
}

func TestParseFromRequest_InvalidJson(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(WebhookSignatureHeader, "d4163f7af2256fae6ab72cb595d3f9d1dfc6fecc")
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("{\"msg\": \"hello world")))

	wh := Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err == nil {
		t.Fatal("expected invalid json to raise an error")
	}
}

func TestParseFromRequest_ValidSignature(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(WebhookSignatureHeader, "d2ef30601350308c1f1c25c5fbf359badb95cbfb")
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("{\"msg\": \"hello world\"}")))

	wh := Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err != nil {
		t.Fatal()
	}
}

func ExampleWebhook_ParseFromRequest() {
	wh, err := NewWebhookFromEnv()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/webhook/onfido", func(w http.ResponseWriter, req *http.Request) {
		whReq, err := wh.ParseFromRequest(req)
		if err != nil {
			if err == ErrInvalidWebhookSignature {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid signature"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error occurred"))
			return
		}

		fmt.Fprintf(w, "Webhook: %+v\n", whReq)
	})

	http.ListenAndServe(":8080", nil)

}
