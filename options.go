package rvtools

import "net/http"

func WithHttpClient(client http.Client) func(*Client) {
	return func(c *Client) {
		c.HttpClient = client
	}
}
