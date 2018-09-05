package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestCreatePageHandler(t *testing.T) {
	payload := `{
		"title": "About Us",
		"slug": "about-us",
		"parent": null,
		"sort": null
	}`

	req, err := http.NewRequest("POST", "/api/pages", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePageHandler(server))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Check the response body is what we expect.
	//Expect newly created page_id
	expected := `^\d+$`
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
