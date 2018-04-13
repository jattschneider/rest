package rest

import (
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "3")
}

const allowValue = "POST, GET, OPTIONS, PATCH, PUT, DELETE"

const accessControlAllowMethodsValue = "POST, GET, OPTIONS, PATCH, PUT, DELETE"

const accessControlAllowHeadersValue = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

func testHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Allow", allowValue)
	w.Header().Set("Access-Control-Allow-Methods", accessControlAllowMethodsValue)
	w.Header().Set("Access-Control-Allow-Headers", accessControlAllowHeadersValue)

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"someProperty\":\"someValue\"}"))
	case http.MethodPatch, http.MethodPost, http.MethodPut:
		rBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			w.WriteHeader(http.StatusOK)
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

func TestShouldHead(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	header, err := Head(ts.URL)
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
	ts := testServer()
	defer ts.Close()

	re, err := Get(ts.URL)
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
}

func TestShouldPost(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	payload := strings.NewReader("{\"someProperty\":\"someValue\"}")
	re, err := Post(ts.URL, payload)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(re.Header) == 0 {
		t.Errorf("No HTTP header: %v", re.Header)
	}

	assertHeader(t, re.Header, "Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	assertHeader(t, re.Header, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func TestShouldPut(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	payload := EncodeJSON(&struct{ SomeProperty string }{SomeProperty: "struct property value"})
	re, err := Put(ts.URL, payload)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(re.Header) == 0 {
		t.Errorf("No HTTP header: %v", re.Header)
	}

	assertHeader(t, re.Header, "Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	assertHeader(t, re.Header, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func TestShouldPatch(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	payload := strings.NewReader("{\"someProperty\":\"someValue\"}")
	re, err := Patch(ts.URL, payload)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(re.Header) == 0 {
		t.Errorf("No HTTP header: %v", re.Header)
	}

	assertHeader(t, re.Header, "Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	assertHeader(t, re.Header, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func TestShouldOptionsForAllow(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	allows, err := OptionsForAllow(ts.URL)
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
