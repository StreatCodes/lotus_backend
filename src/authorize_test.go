package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestAuthorizeHandler(t *testing.T) {
	payload := `{
		"email": "lotus@example.com",
		"password": "lotus"
	}`

	req, err := http.NewRequest("POST", "/authorize", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authorizeHandler(server))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Expect user_id 0 padded : 40 hex values
	expected := `^\d{10}:[0-9a-f]{80}$`
	result := bytes.TrimSpace(rr.Body.Bytes())

	match, err := regexp.Match(expected, result)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Errorf("handler returned unexpected body: got %#v want %#v",
			rr.Body.String(), expected)
	}
}
