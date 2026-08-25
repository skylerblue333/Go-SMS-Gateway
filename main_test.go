package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthAndReadiness(t *testing.T) {
	for _, path := range []string{"/healthz", "/readyz"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		routes().ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d", path, rr.Code)
		}
	}
}

func TestValidateMessage(t *testing.T) {
	response, err := validateMessage(messageRequest{To: "+14155550123", Body: "hello", SenderID: "Sky"})
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if !response.Valid || response.ProviderDispatch || response.BodyRunes != 5 || response.To != "+14155550123" {
		t.Fatalf("unexpected validation response: %+v", response)
	}
}

func TestValidateEndpointRejectsInvalidInputs(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"invalid phone", `{"to":"4155550123","body":"hello"}`},
		{"empty body", `{"to":"+14155550123","body":"   "}`},
		{"unknown field", `{"to":"+14155550123","body":"hello","dispatch":true}`},
		{"second object", `{"to":"+14155550123","body":"hello"}{}`},
		{"trailing garbage", `{"to":"+14155550123","body":"hello"}x`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/messages/validate", strings.NewReader(tc.body))
			rr := httptest.NewRecorder()
			routes().ServeHTTP(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", rr.Code, rr.Body.String())
			}
		})
	}
}

func TestValidateEndpointSuccessDoesNotClaimDispatch(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/messages/validate", strings.NewReader(`{"to":"+14155550123","body":"Hello from Sky"}`))
	rr := httptest.NewRecorder()
	routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var response validationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.ProviderDispatch {
		t.Fatal("validator must not claim provider dispatch")
	}
}

func TestValidateEndpointMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/messages/validate", nil)
	rr := httptest.NewRecorder()
	routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
