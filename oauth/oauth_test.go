package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(Auth))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Error(err)
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("Expect %d, but got %d", http.StatusMovedPermanently, resp.StatusCode)
	}
	if lo := resp.Header.Get("Location"); !strings.Contains(lo, "https://github.com/login/oauth/autorize") {
		t.Errorf("Location should contain %s, but got %s", "https://github.com/login/oauth/autorize", lo)
	}
}
