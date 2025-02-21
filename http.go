package rest

import (
	"net"
	"net/http"
	"time"
)

func timeout() time.Duration {
	return 10 * time.Second
}

func transportTimeout() time.Duration {
	return 5 * time.Second
}

func buildHTTPClient() *http.Client {
	var transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: transportTimeout(),
		}).Dial,
		TLSHandshakeTimeout: transportTimeout(),
	}
	return &http.Client{
		Timeout:   timeout(),
		Transport: transport,
	}
}
