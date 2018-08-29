package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	payload := `{
		"name": "john",
		"email": "john_doe@example.com",
		"password": "J0hn_1s_c00l!"
	}`

	req, err := http.NewRequest("POST", "/api/users", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUserHandler(server))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.

	//Expect user_id
	expected := `^\d+$`
	result := bytes.TrimSpace(rr.Body.Bytes())

	match, err := regexp.Match(expected, result)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Errorf("handler returned unexpected body: got %#v want number",
			rr.Body.String())
	}
}
