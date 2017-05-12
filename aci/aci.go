package main

import "net/url"

type Client struct {
	URL      *url.URL
	Username string
	Password string
}

func NewClient(url, username, password string) (*Client, error) {
	return &Client{
		URL:      &url.URL{Path: url},
		Username: username,
		Password: password,
	}
}
