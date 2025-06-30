package mux

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteUpsert(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("/foo", func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte("GET /foo v1"))
	}).Methods("GET")

	router.HandleFunc("/bar", func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte("GET /bar v1"))
	}).Methods("GET")

	// This route is masked because v1 was added to the router fist.
	router.HandleFunc("/foo", func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte("GET /foo v2"))
	}).Methods("GET")

	// This route overrides v1 because it is added using UpsertRoute.
	r := router.NewUnattachedRoute().Path("/bar").Methods("GET").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte("GET /bar v2"))
	})
	router.UpsertRoute(r)

	r = router.NewUnattachedRoute().Path("/zoom").Methods("GET").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte("GET /zoom v1"))
	})
	router.UpsertRoute(r)

	assertResponseBody := func(t *testing.T, s *httptest.Server, method, path, wantBody string) {
		rq, _ := http.NewRequest(method, s.URL+path, nil)
		resp, err := s.Client().Do(rq)
		if err != nil {
			t.Fatalf("unexpected error getting from server: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading body: %v", err)
		}
		if !bytes.Equal(body, []byte(wantBody)) {
			t.Fatalf("response should be hello world, was: %q", string(body))
		}
	}

	t.Run("/foo", func(t *testing.T) {
		s := httptest.NewServer(router)
		defer s.Close()
		assertResponseBody(t, s, "GET", "/foo", "GET /foo v1")
	})
	t.Run("/bar", func(t *testing.T) {
		s := httptest.NewServer(router)
		defer s.Close()
		assertResponseBody(t, s, "GET", "/bar", "GET /bar v2")
	})
	t.Run("/zoom", func(t *testing.T) {
		s := httptest.NewServer(router)
		defer s.Close()
		assertResponseBody(t, s, "GET", "/zoom", "GET /zoom v1")
	})
}
