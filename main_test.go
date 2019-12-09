package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sterks/ftpReader/controller"
)

func TestMain(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.HomeHandler)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `Test1`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}

}
