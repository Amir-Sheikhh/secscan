package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	server := NewServer(Config{
		AllowedOrigins:      []string{"http://localhost:3000"},
		ScanTimeout:         2 * time.Second,
		EnableActiveProbes:  false,
		EnableExternalIntel: false,
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"status":"ok"`)) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestCreateScanRejectsBlockedTarget(t *testing.T) {
	t.Parallel()

	server := NewServer(Config{
		AllowedOrigins:      []string{"http://localhost:3000"},
		ScanTimeout:         2 * time.Second,
		EnableActiveProbes:  false,
		EnableExternalIntel: false,
	})

	body, _ := json.Marshal(map[string]any{
		"url": "http://127.0.0.1",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/scan", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d with body %s", rec.Code, rec.Body.String())
	}
}
