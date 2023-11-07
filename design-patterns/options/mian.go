package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Client represents our HTTP client
type Client struct {
	client          *http.Client
	timeout         time.Duration
	userAgent       string
	followRedirects bool
}

// Option is a functional option type that allows us to configure the Client
type Option func(*Client)

// NewClient creates a new HTTP client with default options
func NewClient(options ...Option) *Client {
	client := &Client{
		client:          &http.Client{},
		timeout:         15 * time.Second,      // Default timeout
		userAgent:       "awesome http client", // Default user agent
		followRedirects: false,                 // Default follows redirects
	}

	// Apply all the functional options to configure the client
	for _, opt := range options {
		opt(client)
	}

	return client
}

// WithTimeout is a functional option to set the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(client *Client) {
		client.timeout = timeout
	}
}

// WithUserAgent is a functional option to set the HTTP client user agent
func WithUserAgent(ua string) Option {
	return func(client *Client) {
		client.userAgent = ua
	}
}

// WithoutRedirects is a functional option to disable following redirects
func WithoutRedirects() Option {
	return func(client *Client) {
		client.followRedirects = true
	}
}

// UseInsecureTransport is a functional option to use an insecure HTTP transport
func UseInsecureTransport() Option {
	return func(client *Client) {
		client.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}

// Get performs an HTTP GET request.
func (client *Client) Get(url string) (*http.Response, error) {
	// Use c.client with all the configured options to perform the request.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return client.client.Do(req)
}

func main() {
	// Create a new HTTP client with custom options.
	client := NewClient(
		WithTimeout(10*time.Second),
		WithUserAgent("my custom user agent"),
		UseInsecureTransport(),
	)

	// Use the client to make HTTP requests.
	response, err := client.Get("https://www.googel.com/")
	if err != nil {
		// Handle the error
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Process the response.
	fmt.Println(response.StatusCode)
}
