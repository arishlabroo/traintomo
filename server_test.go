package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandlesRequests(t *testing.T) {
	s := NewServer()
	req := httptest.NewRequest("GET", "/postschedule", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Error("failed request")
	}
}
