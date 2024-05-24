package uploadfly

import (
	"bytes"
	"encoding/json"
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

type DeleteFileOption struct {
	FileURL string `json:"file_url"`
}

// This deletes a single file
func (c *Client) Delete(opts *DeleteFileOption) error {
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(&opts); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/delete", baseEndpoint), buf)
	if err != nil {
		return err
	}

	res, err := c.c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode > http.StatusCreated {
		var w struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(res.Body).Decode(&w); err != nil {
			return err
		}

		return errors.New(w.Message)
	}

	return nil
}
