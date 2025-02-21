package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

// ResponseEntity struct represents a HTTP response.
type ResponseEntity struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

type Client struct {
	httpClient *http.Client
	timeout    time.Duration
}

func New() *Client {
	return &Client{
		httpClient: buildHTTPClient(),
		timeout:    timeout(),
	}
}

// BodyReader resturns a ResponseEntity body as a Reader.
func (re *ResponseEntity) BodyReader() *bytes.Reader {
	return bytes.NewReader(re.Body)
}

// BodyString resturns a ResponseEntity body as a string.
func (re *ResponseEntity) BodyString() string {
	return string(re.Body)
}

func JSONRequestCallback(r *http.Request) {
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
	resBody, err := io.ReadAll(res.Body)
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
func (c *Client) Exchange(url, method string, body io.Reader, requestCallback func(r *http.Request)) (ResponseEntity, error) {
	return exchange(c.httpClient, c.timeout, url, method, body, requestCallback)
}

// Get gets the content from the given URL
func (c *Client) Get(url string, requestCallback func(r *http.Request)) (ResponseEntity, error) {
	return c.Exchange(url, http.MethodGet, nil, requestCallback)
}

// Head returns the headers from the given URL
func (c *Client) Head(url string, requestCallback func(r *http.Request)) (http.Header, error) {
	re, err := c.Exchange(url, http.MethodHead, nil, requestCallback)
	return re.Header, err
}

// Post posts body content to the given URL
func (c *Client) Post(url string, body io.Reader, requestCallback func(r *http.Request)) (ResponseEntity, error) {
	return c.Exchange(url, http.MethodPost, body, requestCallback)
}

// Put puts the body content to the given URL
func (c *Client) Put(url string, body io.Reader, requestCallback func(r *http.Request)) (ResponseEntity, error) {
	return c.Exchange(url, http.MethodPut, body, requestCallback)
}

// Patch patches the body content to the given URL
func (c *Client) Patch(url string, body io.Reader, requestCallback func(r *http.Request)) (ResponseEntity, error) {
	return c.Exchange(url, http.MethodPatch, body, requestCallback)
}

// OptionsForAllow returns the allowed HTTP methods
func (c *Client) OptionsForAllow(url string, requestCallback func(r *http.Request)) ([]string, error) {
	re, err := c.Exchange(url, http.MethodOptions, nil, requestCallback)
	allowHeader := re.Header.Get("Allow")
	if len(allowHeader) > 0 {
		return strings.Split(allowHeader, ","), err
	}
	return []string{}, err
}

// Delete deletes from the given URL
func (c *Client) Delete(url string, requestCallback func(r *http.Request)) error {
	_, err := c.Exchange(url, http.MethodDelete, nil, requestCallback)
	return err
}
