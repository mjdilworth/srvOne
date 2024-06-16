package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	srv := New("8080")
	req := httptest.NewRequest(http.MethodGet, "/health/", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Fatalf("expected status OK, got %d", status)
	}

	res := w.Result()
	defer res.Body.Close()

	resp := struct {
		Message string
	}{}

	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("error decoding message: %s", err)
	}

	if resp.Message != "pong" {
		t.Fatalf("expected pong, but got %s", resp.Message)
	}
}

func TestPingErrors(t *testing.T) {
	tbl := []struct {
		name string
		// given
		method  string
		headers map[string]string
		// expected
		status int
	}{
		{
			"should not accept POST",
			http.MethodPost,
			nil,
			http.StatusMethodNotAllowed,
		},
		{
			"should not accept PUT",
			http.MethodPut,
			nil,
			http.StatusMethodNotAllowed,
		},
		{
			"should not accept PATCH",
			http.MethodPatch,
			nil,
			http.StatusMethodNotAllowed,
		},
		{
			"should not accept xml content-type",
			http.MethodGet,
			map[string]string{"content-type": "application/xml"},
			http.StatusNotImplemented,
		},
	}

	for _, tc := range tbl {
		t.Run(tc.name, func(t *testing.T) {
			srv := newServer()
			req := httptest.NewRequest(tc.method, "/ping", nil)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			if status := w.Result().StatusCode; status != tc.status {
				t.Fatalf("expected status %d, got %d", tc.status, status)
			}
		})
	}
}
