package httpclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	Http *http.Client
}

func New() *Client {
	return &Client{
		Http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Request(method, url string, payload interface{}) (*http.Response, error) {
	var body []byte

	if payload != nil {
		body, _ = json.Marshal(payload)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.Http.Do(req)
}
