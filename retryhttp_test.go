package retryablehttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_defaults(t *testing.T) {
	c := NewClient()
	if c.RetryMax != DefaultRetryMax {
		t.Errorf("expected RetryMax=%d, got %d", DefaultRetryMax, c.RetryMax)
	}
	if c.RetryWaitMin != DefaultRetryWaitMin {
		t.Errorf("expected RetryWaitMin=%s, got %s", DefaultRetryWaitMin, c.RetryWaitMin)
	}
	if c.RetryWaitMax != DefaultRetryWaitMax {
		t.Errorf("expected RetryWaitMax=%s, got %s", DefaultRetryWaitMax, c.RetryWaitMax)
	}
	if c.HTTPClient == nil {
		t.Error("expected non-nil HTTPClient")
	}
}

func TestGet_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient()
	resp, err := c.Get(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGet_retries_on_500(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient()
	c.RetryWaitMin = 0
	resp, err := c.Get(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 after retries, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestGet_exhausts_retries(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := NewClient()
	c.RetryMax = 2
	c.RetryWaitMin = 0
	_, err := c.Get(srv.URL)
	if err == nil {
		t.Error("expected error after exhausting retries")
	}
}
