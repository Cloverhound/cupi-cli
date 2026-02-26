package client

import (
	"crypto/tls"
	"net/http"
	"time"
)

// NewHTTPClient returns an HTTP client configured for CUC (self-signed certs, 30s timeout)
func NewHTTPClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}
