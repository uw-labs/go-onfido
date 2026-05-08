package onfido_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/uw-labs/go-onfido"
)

func TestNewWebhookFromEnv_MissingToken(t *testing.T) {
	t.Parallel()
	_, err := onfido.NewWebhookFromEnv()
	if err == nil {
		t.Fatal()
	}
	if !errors.Is(err, onfido.ErrMissingWebhookToken) {
		t.Fatal("expected error to match ErrMissingWebhookToken")
	}
}

func TestNewWebhookFromEnv_TokenSet(t *testing.T) {
	expected := "808yup"
	t.Setenv(onfido.WebhookTokenEnv, expected)

	wh, err := onfido.NewWebhookFromEnv()
	if err != nil {
		t.Fatal()
	}
	if wh.Token != expected {
		t.Fatalf("expected to see `%s` token but got `%s`", expected, wh.Token)
	}
}

func TestValidateSignature_InvalidSignature(t *testing.T) {
	t.Parallel()
	wh := onfido.Webhook{Token: "abc123"}
	err := wh.ValidateSignature([]byte("hello world"), "invalid")
	if err == nil {
		t.Fatal()
	}
	if !errors.Is(err, onfido.ErrInvalidWebhookSignature) {
		t.Fatal("expected error to match ErrInvalidWebhookSignature")
	}
}

func TestValidateSignature_ValidSignature(t *testing.T) {
	t.Parallel()
	wh := onfido.Webhook{Token: "abc123"}
	err := wh.ValidateSignature([]byte("hello world"), "fcc98c5b4f306cfe6b5b8fcce03ddb33fc13ae6b")
	if err != nil {
		t.Fatal()
	}
}

func TestParseFromRequest_InvalidSignature(t *testing.T) {
	t.Parallel()
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(onfido.WebhookSignatureHeader, "123")
	req.Body = io.NopCloser(bytes.NewBufferString("{\"msg\": \"hello world\"}"))

	wh := onfido.Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err == nil {
		t.Fatal()
	}
	if !errors.Is(err, onfido.ErrInvalidWebhookSignature) {
		t.Fatal("expected error to match ErrInvalidWebhookSignature")
	}
}

func TestParseFromRequest_SkipSignatureValidation(t *testing.T) {
	t.Parallel()
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Body = io.NopCloser(bytes.NewBufferString("{\"msg\": \"hello world\"}"))

	wh := onfido.Webhook{Token: "abc123", SkipSignatureValidation: true}
	_, err := wh.ParseFromRequest(req)
	if err != nil {
		t.Errorf("expected no error as signature validation should have been skipped: %s", err.Error())
	}
}

func TestParseFromRequest_InvalidJson(t *testing.T) {
	t.Parallel()
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(onfido.WebhookSignatureHeader, "d4163f7af2256fae6ab72cb595d3f9d1dfc6fecc")
	req.Body = io.NopCloser(bytes.NewBufferString("{\"msg\": \"hello world"))

	wh := onfido.Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err == nil {
		t.Fatal("expected invalid json to raise an error")
	}
}

func TestParseFromRequest_ValidSignature(t *testing.T) {
	t.Parallel()
	req := &http.Request{
		Header: make(map[string][]string),
	}
	req.Header.Add(onfido.WebhookSignatureHeader, "d2ef30601350308c1f1c25c5fbf359badb95cbfb")
	req.Body = io.NopCloser(bytes.NewBufferString("{\"msg\": \"hello world\"}"))

	wh := onfido.Webhook{Token: "abc123"}
	_, err := wh.ParseFromRequest(req)
	if err != nil {
		t.Fatal()
	}
}
