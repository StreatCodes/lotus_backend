package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestCreateComponentHandler(t *testing.T) {
	payload := `{
		"page_id": 1,
		"sort": 1,
		"data": "{\"width\":3, \"type\": \"text-block\", \"text\": \"Some Sample Text\"}"
	}`

	req, err := http.NewRequest("POST", "/api/components", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateComponentHandler(server))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Check the response body is what we expect.
	//Expect newly created component_id
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
