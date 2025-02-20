package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestShouldHead(t *testing.T) {
	c := New()
	ts := testServer()
	defer ts.Close()

	header, err := c.Head(ts.URL, JSONRequestCallback)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(header) == 0 {
		t.Errorf("No HTTP header: %v", header)
	}

	assertHeader(t, header, "Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	assertHeader(t, header, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func TestShouldGet(t *testing.T) {
	c := New()
	ts := testServer()
	defer ts.Close()

	re, err := c.Get(ts.URL, JSONRequestCallback)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(re.Header) == 0 {
		t.Errorf("No HTTP header: %v", re.Header)
	}

	body := &struct{ SomeProperty string }{}
	DecodeJSON(re.Body, &body)

	if len(body.SomeProperty) == 0 {
		t.Error("body.SomeProperty should not be empty")
	}

	assertHeader(t, re.Header, "Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	assertHeader(t, re.Header, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	assertStatusCode(t, re.StatusCode, http.StatusOK)
}

func TestShouldPost(t *testing.T) {
	c := New()
	ts := testServer()
	defer ts.Close()

	payload := strings.NewReader("{\"someProperty\":\"someValue\"}")
	re, err := c.Post(ts.URL, payload, JSONRequestCallback)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(re.Header) == 0 {
		t.Errorf("No HTTP header: %v", re.Header)
	}

	assertHeader(t, re.Header, "Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	assertHeader(t, re.Header, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	assertStatusCode(t, re.StatusCode, http.StatusCreated)
}

func TestShouldPut(t *testing.T) {
	c := New()
	ts := testServer()
	defer ts.Close()

	payload := EncodeJSON(&struct{ SomeProperty string }{SomeProperty: "struct property value"})
	re, err := c.Put(ts.URL, payload, JSONRequestCallback)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(re.Header) == 0 {
		t.Errorf("No HTTP header: %v", re.Header)
	}

	assertHeader(t, re.Header, "Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	assertHeader(t, re.Header, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	assertStatusCode(t, re.StatusCode, http.StatusOK)
}

func TestShouldPatch(t *testing.T) {
	c := New()
	ts := testServer()
	defer ts.Close()

	payload := strings.NewReader("{\"someProperty\":\"someValue\"}")
	re, err := c.Patch(ts.URL, payload, JSONRequestCallback)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(re.Header) == 0 {
		t.Errorf("No HTTP header: %v", re.Header)
	}

	assertHeader(t, re.Header, "Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	assertHeader(t, re.Header, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	assertStatusCode(t, re.StatusCode, http.StatusOK)
}

func TestShouldOptionsForAllow(t *testing.T) {
	c := New()
	ts := testServer()
	defer ts.Close()

	allows, err := c.OptionsForAllow(ts.URL, JSONRequestCallback)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(allows) == 0 {
		t.Errorf("No HTTP method allowed: %v", allows)
	}

	expected := []string{"POST", "GET", "OPTIONS", "PATCH", "PUT", "DELETE"}
	if reflect.DeepEqual(allows, expected) {
		t.Errorf("Expected allows: [%v] got: [%v]", expected, allows)
	}
}

func TestShouldDelete(t *testing.T) {
	c := New()
	ts := testServer()
	defer ts.Close()

	err := c.Delete(ts.URL, JSONRequestCallback)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(30 * time.Millisecond)

	allowValue := "POST, GET, OPTIONS, PATCH, PUT, DELETE"
	accessControlAllowMethodsValue := "POST, GET, OPTIONS, PATCH, PUT, DELETE"
	accessControlAllowHeadersValue := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Allow", allowValue)
	w.Header().Set("Access-Control-Allow-Methods", accessControlAllowMethodsValue)
	w.Header().Set("Access-Control-Allow-Headers", accessControlAllowHeadersValue)

	switch r.Method {
	case http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"someProperty\":\"someValue\"}"))
	case http.MethodPatch, http.MethodPut:
		defer r.Body.Close()
		rBody, err := io.ReadAll(r.Body)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(rBody)
		}
	case http.MethodPost:
		defer r.Body.Close()
		rBody, err := io.ReadAll(r.Body)
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			w.Write(rBody)
		}
	default:
		return
	}
}

func assertStatusCode(t *testing.T, statusCode, expected int) {
	if statusCode != expected {
		t.Errorf("Expected status code: [%v] got: [%v]", expected, statusCode)
	}
}

func assertHeader(t *testing.T, header http.Header, name, expected string) {
	value := header.Get(name)
	if value != expected {
		t.Errorf("Expected methods: [%v] got: [%v]", expected, value)
	}
}

func testServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(testHandler))
}
