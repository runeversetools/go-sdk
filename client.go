package rvtools

import (
	"net/http"
	"time"
)

type Client struct {
	Host       string
	ApiKey     string
	HttpClient http.Client
}

func NewClient(host string, apiKey string, options ...func(*Client)) *Client {
	client := &Client{
		Host:       host,
		ApiKey:     apiKey,
		HttpClient: http.Client{Timeout: 5 * time.Second},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func NewLocalClient(apiKey string) *Client {
	return NewClient("https://api.runeverse.local", apiKey)
}

func NewRemoteClient(apiKey string) *Client {
	return NewClient("https://api.runeverse.tools", apiKey)
}
