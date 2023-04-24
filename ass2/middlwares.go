package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuthenticatedUser(t *testing.T) {
	// Create a new instance of the application.
	app := newTestApplication(t)

	// Create a new HTTP test request.
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Call the requireAuthenticatedUser function with a dummy handler.
	handler := app.requireAuthenticatedUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We don't need to do anything here.
	}))

	// Create a new HTTP response recorder.
	rr := httptest.NewRecorder()

	// Call the handler function.
	handler.ServeHTTP(rr, req)

	// Check that the response status code is 401 (unauthorized).
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want %d; got %d", http.StatusUnauthorized, rr.Code)
	}

	// Set an Authorization header on the request.
	req.Header.Set("Authorization", "Bearer abc123")

	// Call the handler function again.
	handler.ServeHTTP(rr, req)

	// Check that the response status code is now 200 (OK).
	if rr.Code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rr.Code)
	}
}
