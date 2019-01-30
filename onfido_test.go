package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_MsgSet(t *testing.T) {
	err := Error{}
	err.Err.Msg = "some error message"
	if err.Error() != err.Err.Msg {
		t.Fatal()
	}
}

func TestError_UseHttpResp(t *testing.T) {
	err := Error{
		Resp: &http.Response{
			StatusCode: http.StatusTeapot,
		},
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("status code %d", http.StatusTeapot)) {
		t.Fatal()
	}
}

func TestError_FallbackMsg(t *testing.T) {
	err := Error{}
	if err.Error() != "an unknown error occurred" {
		t.Fatal()
	}
}

func TestToken_IsProd(t *testing.T) {
	tokens := []struct {
		token  string
		isProd bool
	}{
		{"prod_122333", true},
		{"122333", true},
		{"test_122333", false},
	}

	for _, expected := range tokens {
		token := Token(expected.token)
		if token.Prod() != expected.isProd {
			t.Fatal()
		}
	}
}

func TestNewClientFromEnv_NoToken(t *testing.T) {
	os.Setenv(TokenEnv, "")
	if _, err := NewClientFromEnv(); err == nil {
		t.Fatal()
	}
}

func TestNewClientFromEnv_EnvSet(t *testing.T) {
	expectedToken := "lk3j6323j442"
	os.Setenv(TokenEnv, expectedToken)
	defer os.Setenv(TokenEnv, "")

	client, err := NewClientFromEnv()
	if err != nil {
		t.Fatal()
	}
	if client.Token.String() != expectedToken {
		t.Fatalf("expected token to be `%s` but got `%s`", expectedToken, client.Token)
	}
}

func TestNewRequest_WithFullURL(t *testing.T) {
	client := NewClient("123")
	req, err := client.newRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != "GET" {
		t.Fatalf("expected method of `GET` but got `%s`", req.Method)
	}
	if req.URL.String() != "https://example.com" {
		t.Fatalf("exptected URL of `https://example.com` but got `%s`", req.URL)
	}
}

func TestNewRequest_WithPathUri(t *testing.T) {
	expectedURL := "https://api.onfido.com/v2/applicants"
	client := NewClient("123")
	uris := []string{"/applicants", "applicants"}

	for _, uri := range uris {
		req, err := client.newRequest("GET", uri, nil)
		if err != nil {
			t.Fatal(err)
		}
		if req.Method != "GET" {
			t.Fatalf("expected method of `GET` but got `%s`", req.Method)
		}
		if req.URL.String() != expectedURL {
			t.Fatalf("exptected URL of `%s` but got `%s`", expectedURL, req.URL)
		}
	}
}

func TestNewRequest_TokenSet(t *testing.T) {
	expectedToken := "io2h54k2j3h52jk"
	client := NewClient(expectedToken)
	req, err := client.newRequest("get", "/foo", nil)
	if err != nil {
		t.Fatal()
	}
	if req.Header.Get("Authorization") != fmt.Sprintf("Token token=%s", expectedToken) {
		t.Fatalf("expected to see Authorization header of `%s` but got `%s`",
			fmt.Sprintf("Token token=%s", expectedToken),
			req.Header.Get("Authorization"))
	}
}

func TestDo_RequestErrors(t *testing.T) {
	expected := errors.New("TestJson_RequestErrors")

	client := NewClient("123")
	client.HTTPClient = &stubbedHTTPClient{err: expected}

	_, err := client.do(context.Background(), &http.Request{}, nil)
	if err == nil {
		t.Fatal()
	}
	if err != expected {
		t.Fatalf("expected to see error `%s` but got `%s`", expected, err)
	}
}

func TestDo_InvalidStatusCode(t *testing.T) {
	client := NewClient("123")
	client.HTTPClient = &stubbedHTTPClient{resp: &http.Response{StatusCode: http.StatusForbidden}}

	_, err := client.do(context.Background(), &http.Request{}, nil)
	if err == nil {
		t.Fatal()
	}
	onfidoErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected to see `onfido.OnfidoError` but got %T", err)
	}
	if onfidoErr.Resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected to see status code `%d` but got `%d`", http.StatusForbidden, onfidoErr.Resp.StatusCode)
	}
}

func TestDo_InvalidStatusCode_InvalidJsonParsed(t *testing.T) {
	resp := &http.Response{
		Header:     make(map[string][]string),
		StatusCode: http.StatusBadGateway,
	}
	resp.Header.Add("Content-Type", "application/json")
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("hello")))

	client := NewClient("123")
	client.HTTPClient = &stubbedHTTPClient{resp: resp}

	_, err := client.do(context.Background(), &http.Request{}, &Applicant{})
	if err == nil {
		t.Fatal("expected to see an error after the body was unable to be parsed as JSON")
	}
	if _, ok := err.(*Error); ok {
		t.Fatal("json should not have been parsed")
	}
}

func TestDo_InvalidStatusCode_JsonParsed(t *testing.T) {
	expected := Error{
		Err: struct {
			ID     string      `json:"id"`
			Type   string      `json:"type"`
			Msg    string      `json:"message"`
			Fields ErrorFields `json:"fields"`
		}{
			ID:   "123",
			Type: "foo",
			Msg:  "some msg",
			Fields: map[string][]string{
				"first_name": []string{"can't be blank"},
				"last_name":  []string{"can't be blank", "is too short (minimum is 2 characters)"},
			},
		},
	}
	encodedErr, err := json.Marshal(expected)
	if err != nil {
		panic(err)
	}

	resp := &http.Response{
		Header:     make(map[string][]string),
		StatusCode: http.StatusBadGateway,
	}
	resp.Header.Add("Content-Type", "application/json")
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(encodedErr))

	client := NewClient("123")
	client.HTTPClient = &stubbedHTTPClient{resp: resp}

	_, err = client.do(context.Background(), &http.Request{}, &Applicant{})
	if err == nil {
		t.Fatal("expected to see an error after the body was unable to be parsed as JSON")
	}
	onfidoErr, ok := err.(*Error)
	if !ok {
		t.Fatal("failed to parse error json")
	}
	assert.Equal(t, expected.Err.ID, onfidoErr.Err.ID)
	assert.Equal(t, expected.Err.Type, onfidoErr.Err.Type)
	assert.Equal(t, expected.Err.Msg, onfidoErr.Err.Msg)
	for name, value := range expected.Err.Fields {
		assert.Contains(t, onfidoErr.Err.Fields, name)
		assert.ElementsMatch(t, onfidoErr.Err.Fields[name], value)
	}
}

func TestDo_InvalidJsonResponse(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusOK}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("hello")))

	client := NewClient("123")
	client.HTTPClient = &stubbedHTTPClient{resp: resp}

	_, err := client.do(context.Background(), &http.Request{}, &Applicant{})
	if err == nil {
		t.Fatal("expected to see an error after the body was unable to be parsed as JSON")
	}
}

type stubbedHTTPClient struct {
	resp *http.Response
	err  error
}

func (c *stubbedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.resp, nil
}
