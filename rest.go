package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// Client is a exported func() *http.Client that builds a http.Client
var Client = defaultClient

// Timeout is a exported func() time.Duration that returns a time.Duration to setup the http.Client
var Timeout = defaultTimeout

// TransportTimeout is a exported func() time.Duration that returns a time.Duration to setup the http.Client
var TransportTimeout = defaultTransportTimeout

// ResponseEntity struct represents a HTTP response.
type ResponseEntity struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// BodyString resturns a ResponseEntity body as string.
func (re *ResponseEntity) BodyString() string {
	return string(re.Body)
}

func defaultTimeout() time.Duration {
	return 10 * time.Second
}

func defaultTransportTimeout() time.Duration {
	return 5 * time.Second
}

func defaultClient() *http.Client {
	var transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: TransportTimeout(),
		}).Dial,
		TLSHandshakeTimeout: TransportTimeout(),
	}
	return &http.Client{
		Timeout:   Timeout(),
		Transport: transport,
	}
}

func defaultRequestCallback(r *http.Request) {
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Cache-Control", "no-cache")
}

func exchange(client *http.Client, timeout time.Duration, url, method string, body io.Reader, requestCallback func(r *http.Request)) (ResponseEntity, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return ResponseEntity{Header: make(http.Header)}, err
	}

	req = req.WithContext(ctx)

	if requestCallback != nil {
		requestCallback(req)
	}

	res, err := client.Do(req)
	if err != nil {
		return ResponseEntity{Header: make(http.Header)}, err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ResponseEntity{Header: make(http.Header)}, err
	}

	return ResponseEntity{StatusCode: res.StatusCode, Header: res.Header, Body: resBody}, nil
}

// EncodeJSON returns the JSON encoding of v in a reader
func EncodeJSON(v interface{}) io.Reader {
	w := new(bytes.Buffer)
	json.NewEncoder(w).Encode(v)
	return w
}

// DecodeJSON decodes the JSON encoded b into the value pointed to by v.
func DecodeJSON(b []byte, v interface{}) error {
	return json.NewDecoder(bytes.NewReader(b)).Decode(&v)
}

// Exchange generic function that exchanges/requests HTTP operations/verbs
func Exchange(url, method string, body io.Reader, requestCallback func(r *http.Request)) (ResponseEntity, error) {
	return exchange(Client(), Timeout(), url, method, body, requestCallback)
}

// Get gets the content from the given URL
func Get(url string) (ResponseEntity, error) {
	return Exchange(url, http.MethodGet, nil, defaultRequestCallback)
}

// Head returns the headers from the given URL
func Head(url string) (http.Header, error) {
	re, err := Exchange(url, http.MethodHead, nil, defaultRequestCallback)
	return re.Header, err
}

// Post posts body content to the given URL
func Post(url string, body io.Reader) (ResponseEntity, error) {
	return Exchange(url, http.MethodPost, body, defaultRequestCallback)
}

// Put puts the body content to the given URL
func Put(url string, body io.Reader) (ResponseEntity, error) {
	return Exchange(url, http.MethodPut, body, defaultRequestCallback)
}

// Patch patches the body content to the given URL
func Patch(url string, body io.Reader) (ResponseEntity, error) {
	return Exchange(url, http.MethodPatch, body, defaultRequestCallback)
}

// OptionsForAllow returns the allowed HTTP methods
func OptionsForAllow(url string) ([]string, error) {
	re, err := Exchange(url, http.MethodOptions, nil, defaultRequestCallback)
	allowHeader := re.Header.Get("Allow")
	if len(allowHeader) > 0 {
		return strings.Split(allowHeader, ","), err
	}
	return []string{}, err
}

// Delete deletes from the given URL
func Delete(url string) error {
	_, err := Exchange(url, http.MethodDelete, nil, defaultRequestCallback)
	return err
}
