package onfido_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/utilitywarehouse/go-onfido"
)

func TestNewWebhookFromEnv_MissingToken(t *testing.T) {
	_, err := onfido.NewWebhookFromEnv()
	if err == nil {
		t.Fatal()
	}
	if err != onfido.ErrMissingWebhookToken {
		t.Fatal("expected error to match ErrMissingWebhookToken")
	}
}

func TestNewWebhookFromEnv_TokenSet(t *testing.T) {
	expected := "808yup"
	os.Setenv(onfido.WebhookTokenEnv, expected)
	defer os.Setenv(onfido.WebhookTokenEnv, "")

	wh, err := onfido.NewWebhookFromEnv()
	if err != nil {
		t.Fatal()
	}
	if wh.Token != expected {
		t.Fatalf("expected to see `%s` token but got `%s`", expected, wh.Token)
	}
}

func TestValidateSignature_InvalidSignature(t *testing.T) {
	wh := onfido.Webhook{Token: "abc123"}
	err := wh.ValidateSignature([]byte("hello world"), "invalid")
	if err == nil {
		t.Fatal()
	}
	if err != onfido.ErrInvalidWebhookSignature {
		t.Fatal("expected error to match ErrInvalidWebhookSignature")
	}
}

func TestValidateSignature_ValidSignature(t *testing.T) {
	wh := onfido.Webhook{Token: "abc123"}
	err := wh.ValidateSignature([]byte("hello world"), "fcc98c5b4f306cfe6b5b8fcce03ddb33fc13ae6b")
	if err != nil {
		t.Fatal()
	}
}

func TestParseFromRequest_InvalidSignature(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(onfido.WebhookSignatureHeader, "123")
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("hello world")))

	wh := onfido.Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err == nil {
		t.Fatal()
	}
	if err != onfido.ErrInvalidWebhookSignature {
		t.Fatal("expected error to match ErrInvalidWebhookSignature")
	}
}

func TestParseFromRequest_InvalidJson(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(onfido.WebhookSignatureHeader, "d4163f7af2256fae6ab72cb595d3f9d1dfc6fecc")
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("{\"msg\": \"hello world")))

	wh := onfido.Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err == nil {
		t.Fatal("expected invalid json to raise an error")
	}
}

func TestParseFromRequest_ValidSignature(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(onfido.WebhookSignatureHeader, "d2ef30601350308c1f1c25c5fbf359badb95cbfb")
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("{\"msg\": \"hello world\"}")))

	wh := onfido.Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err != nil {
		t.Fatal()
	}
}
