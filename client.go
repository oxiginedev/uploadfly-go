package uploadfly

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const baseEndpoint = "https://api.uploadfly.cloud"

type Client struct {
	c      *http.Client
	apiKey string
}

type apiKeyAuthTransport struct {
	transport http.RoundTripper
	key       string
}

func (a *apiKeyAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.key))
	return a.transport.RoundTrip(r)
}

func New(opts ...Option) (*Client, error) {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	if len(strings.TrimSpace(c.apiKey)) == 0 {
		return nil, errors.New("please provide a vaild api key")
	}

	if c.c == nil {
		c.c = &http.Client{
			Transport: &apiKeyAuthTransport{
				transport: http.DefaultTransport,
				key:       c.apiKey,
			},
			Timeout: time.Second * 5,
		}
	}

	c.apiKey = ""

	return c, nil
}

type Option func(c *Client)

func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		c.c = h
	}
}

type UploadFileOption struct {
	Files []*multipart.File
}

type UploadedFile struct {
	URL  string `json:"url"`
	Path string `json:"path"`
	Type string `json:"type"`
	Size string `json:"size"`
	Name string `json:"name"`
}

func (c *Client) Upload(opt *UploadFileOption) (*UploadedFile, error) {
	return nil, nil
}

// This deletes a single file
func (c *Client) Delete(fileURL string) error {
	return nil
}
