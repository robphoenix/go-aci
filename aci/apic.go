package aci

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
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
	AAAUser `json:"aaaUser"`
}

// AAAUser ...
type AAAUser struct {
	Attributes `json:"attributes"`
}

// Attributes ...
type Attributes struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

// NewClient instantiates a new APIC client
func NewClient(host, username, password string) (*Client, error) {
	return &Client{
		Host:     &url.URL{Path: host},
		Username: username,
		Password: password,
		Client:   &http.Client{Transport: T},
	}, nil
}

// Login authenticates a new APIC session
func (c *Client) Login() error {
	a := Attributes{Name: c.Username, Pwd: c.Password}
	l := LoginJSON{AAAUser: AAAUser{Attributes: a}}
	loginURL := c.Host.Path + loginPath
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(l)
	req, err := http.NewRequest("POST", loginURL, b)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// get auth cookie
	// TODO check cookie name?
	apicCookie := resp.Cookies()[0]
	c.Cookie = apicCookie.String()
	fmt.Println("response Status:", resp.Status)
	return nil
}
