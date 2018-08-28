package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	//Drop all rows and reset the id sequence
	fmt.Println("Running user init")
	server := GetConfig()
	_, err := server.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY;")
	if err != nil {
		log.Fatal("Failed to delete all from users table: " + err.Error())
		return
	}
}

func TestCreateUserHandler(t *testing.T) {
	payload := `{
		"name": "lotus",
		"email": "lotus@example.com",
		"password": "secur3rBo!z"
	}`

	req, err := http.NewRequest("POST", "/users", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUserHandler(GetConfig()))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `1`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %#v want %#v",
			rr.Body.String(), expected)
	}
}
