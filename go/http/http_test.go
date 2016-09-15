package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var handlers = []http.HandlerFunc{
	func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("ok"))
		return
	},

	// http: multiple response.WriteHeader calls
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		w.WriteHeader(http.StatusNotFound)
		return
	},
}

func TestHTTP(t *testing.T) {
	for _, f := range handlers {
		testServe(t, f)
	}
}

func testServe(t *testing.T, f http.HandlerFunc) {
	ts := httptest.NewServer(f)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	t.Logf("stats=%v", res.StatusCode)
	if err != nil {
		t.Fatalf("err=%v", err)
	}
}
