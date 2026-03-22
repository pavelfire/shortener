package save

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"shortener/internal/storage"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
)

type mockURLSaver struct {
	saveURLFunc func(urlToSave string, alias string) (int64, error)
	getURLFunc  func(alias string) (string, error)
}

func (m *mockURLSaver) SaveURL(urlToSave string, alias string) (int64, error) {
	if m.saveURLFunc != nil {
		return m.saveURLFunc(urlToSave, alias)
	}
	return 0, nil
}

func (m *mockURLSaver) GetURL(alias string) (string, error) {
	if m.getURLFunc != nil {
		return m.getURLFunc(alias)
	}
	return "", storage.ErrURLNotFound
}

func newRequest(method, path string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, path, body)
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, "test-request-id"))
	return r
}

func TestSave_Success(t *testing.T) {
	const alias = "custom1"
	saver := &mockURLSaver{
		saveURLFunc: func(urlToSave string, a string) (int64, error) {
			if urlToSave != "https://example.com" || a != alias {
				t.Errorf("unexpected args: url=%q, alias=%q", urlToSave, a)
			}
			return 1, nil
		},
		getURLFunc: func(a string) (string, error) {
			return "", storage.ErrURLNotFound
		},
	}

	handler := New(slog.Default(), saver)
	reqBody := map[string]string{"url": "https://example.com", "alias": alias}
	body, _ := json.Marshal(reqBody)

	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "OK" {
		t.Errorf("expected status OK, got %q", resp.Status)
	}
	if resp.Alias != alias {
		t.Errorf("expected alias %q, got %q", alias, resp.Alias)
	}
}

func TestSave_SuccessAutoAlias(t *testing.T) {
	var capturedAlias string
	saver := &mockURLSaver{
		saveURLFunc: func(urlToSave string, alias string) (int64, error) {
			capturedAlias = alias
			return 1, nil
		},
		getURLFunc: func(alias string) (string, error) {
			return "", storage.ErrURLNotFound
		},
	}

	handler := New(slog.Default(), saver)
	reqBody := map[string]string{"url": "https://example.com"}
	body, _ := json.Marshal(reqBody)

	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "OK" {
		t.Errorf("expected status OK, got %q", resp.Status)
	}
	if resp.Alias == "" {
		t.Error("expected non-empty auto-generated alias")
	}
	if capturedAlias != resp.Alias {
		t.Errorf("alias mismatch: saved %q, returned %q", capturedAlias, resp.Alias)
	}
}

func TestSave_InvalidJSON(t *testing.T) {
	saver := &mockURLSaver{}

	handler := New(slog.Default(), saver)
	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader([]byte(`{invalid json`)))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 (render default), got %d", rec.Code)
	}

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "Error" {
		t.Errorf("expected status Error, got %q", resp.Status)
	}
	if resp.Error != "failed to decode request" {
		t.Errorf("unexpected error: %q", resp.Error)
	}
}

func TestSave_ValidationMissingURL(t *testing.T) {
	saver := &mockURLSaver{}

	handler := New(slog.Default(), saver)
	reqBody := map[string]string{}
	body, _ := json.Marshal(reqBody)

	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "Error" {
		t.Errorf("expected status Error, got %q", resp.Status)
	}
	if resp.Error == "" {
		t.Error("expected validation error message")
	}
}

func TestSave_ValidationInvalidURL(t *testing.T) {
	saver := &mockURLSaver{}

	handler := New(slog.Default(), saver)
	reqBody := map[string]string{"url": "not-a-valid-url"}
	body, _ := json.Marshal(reqBody)

	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "Error" {
		t.Errorf("expected status Error, got %q", resp.Status)
	}
}

func TestSave_URLExists(t *testing.T) {
	saver := &mockURLSaver{
		saveURLFunc: func(urlToSave string, alias string) (int64, error) {
			return 0, storage.ErrURLExists
		},
		getURLFunc: func(alias string) (string, error) {
			return "", storage.ErrURLNotFound
		},
	}

	handler := New(slog.Default(), saver)
	reqBody := map[string]string{"url": "https://example.com", "alias": "exist1"}
	body, _ := json.Marshal(reqBody)

	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "Error" {
		t.Errorf("expected status Error, got %q", resp.Status)
	}
	if resp.Error != "url already exists" {
		t.Errorf("expected 'url already exists', got %q", resp.Error)
	}
}

func TestSave_SaveURLError(t *testing.T) {
	saver := &mockURLSaver{
		saveURLFunc: func(urlToSave string, alias string) (int64, error) {
			return 0, errors.New("database error")
		},
		getURLFunc: func(alias string) (string, error) {
			return "", storage.ErrURLNotFound
		},
	}

	handler := New(slog.Default(), saver)
	reqBody := map[string]string{"url": "https://example.com", "alias": "test1"}
	body, _ := json.Marshal(reqBody)

	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "Error" {
		t.Errorf("expected status Error, got %q", resp.Status)
	}
	if resp.Error != "failed to save url" {
		t.Errorf("expected 'failed to save url', got %q", resp.Error)
	}
}

func TestSave_GetURLErrorDuringAliasCheck(t *testing.T) {
	saver := &mockURLSaver{
		saveURLFunc: func(urlToSave string, alias string) (int64, error) {
			return 1, nil
		},
		getURLFunc: func(alias string) (string, error) {
			return "", errors.New("storage unavailable")
		},
	}

	handler := New(slog.Default(), saver)
	reqBody := map[string]string{"url": "https://example.com"}
	body, _ := json.Marshal(reqBody)

	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "Error" {
		t.Errorf("expected status Error, got %q", resp.Status)
	}
	if resp.Error != "failed to check alias" {
		t.Errorf("expected 'failed to check alias', got %q", resp.Error)
	}
}

func TestSave_FailedToGenerateUniqueAlias(t *testing.T) {
	callCount := 0
	saver := &mockURLSaver{
		getURLFunc: func(alias string) (string, error) {
			callCount++
			return "https://existing.com", nil
		},
	}

	handler := New(slog.Default(), saver)
	reqBody := map[string]string{"url": "https://example.com"}
	body, _ := json.Marshal(reqBody)

	req := newRequest(http.MethodPost, "/api/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "Error" {
		t.Errorf("expected status Error, got %q", resp.Status)
	}
	if resp.Error != "failed to generate unique alias" {
		t.Errorf("expected 'failed to generate unique alias', got %q", resp.Error)
	}
}
