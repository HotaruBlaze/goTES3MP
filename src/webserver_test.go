// webserver_test.go

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_webserver_statusHandler(t *testing.T) {
	// Create a request to the /status endpoint
	req := httptest.NewRequest("GET", "/status", nil)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the status handler function with the ResponseRecorder and Request
	status(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, rr.Code)
	}
}
