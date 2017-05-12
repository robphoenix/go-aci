package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	loginPath = "api/aaaLogin.json"
)

var (
	// T is the TLS config
	// https://stackoverflow.com/questions/41250665/go-https-client-issue-remote-error-tls-handshake-failure#
	T = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS11,
		},
	}
)

// Client is
type Client struct {
	Host     *url.URL
	Username string
	Password string
	Cookie   string // APIC login token (apicCookie)
	Client   *http.Client
}

// LoginJSON represents the JSON needed for authentication
type LoginJSON struct {
	AAAUser struct {
		Attributes struct {
			Name string `json:"name"`
			Pwd  string `json:"pwd"`
		} `json:"attributes"`
	} `json:"aaaUser"`
}

// NewClient instantiates a new APIC client
func NewClient(url, username, password string) (*Client, error) {
	return &Client{
		Host:     &url.URL{Path: url},
		Username: username,
		Password: password,
		Client:   &http.Client{Transport: T},
	}
}

// Login authenticates a new APIC session
func (c *Client) Login() error {
	l := LoginJSON{AAAUser{Attributes{
		Name: c.Username,
		Pwd:  c.Password,
	}}}
	loginURL := c.Host.Path + loginPath
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(l)
	resp, err := http.Post(loginURL, "application/json; charset=utf-8", b)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// get auth cookie
	// TODO check cookie name?
	apicCookie := resp.Cookies()[0]
	c.Cookie = apicCookie.String()
}
