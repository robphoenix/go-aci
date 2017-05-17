package aci

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	create       = "created"
	modify       = "modified"
	createModify = "created,modified"
	delete       = "deleted"
	loginPath    = "api/aaaLogin.json"
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
	Host       *url.URL
	Username   string
	Password   string
	Cookie     string // APIC login token (apicCookie)
	httpClient *http.Client
}

// LoginRequest represents the JSON needed for authentication
type LoginRequest struct {
	AAA `json:"aaaUser"`
}

// LoginResponse ...
type LoginResponse struct {
	Imdata []struct {
		AAA `json:"aaaLogin"`
	} `json:"imdata"`
}

// AAA ...
type AAA struct {
	loginAttributes `json:"attributes"`
}

// Attributes ...
type loginAttributes struct {
	Name                   string `json:"name,omitempty"`
	Pwd                    string `json:"pwd,omitempty"`
	FirstLoginTime         string `json:"firstLoginTime,omitempty"`
	FirstName              string `json:"firstName,omitempty"`
	LastName               string `json:"lastName,omitempty"`
	MaximumLifetimeSeconds string `json:"maximumLifetimeSeconds,omitempty"`
	Node                   string `json:"node,omitempty"`
	RefreshTimeoutSeconds  string `json:"refreshTimeoutSeconds,omitempty"`
	RestTimeoutSeconds     string `json:"restTimeoutSeconds,omitempty"`
	SessionID              string `json:"sessionId,omitempty"`
	SiteFingerprint        string `json:"siteFingerprint,omitempty"`
	Token                  string `json:"token,omitempty"`
	UserName               string `json:"userName,omitempty"`
	Version                string `json:"version,omitempty"`
}

// NewClient instantiates a new APIC client
func NewClient(host, username, password string) (*Client, error) {
	return &Client{
		Host:       &url.URL{Scheme: "https", Host: host},
		Username:   username,
		Password:   password,
		httpClient: &http.Client{Transport: T},
	}, nil
}

func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	u := url.URL{Scheme: c.Host.Scheme, Host: c.Host.Host, Path: path}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if c.Cookie != "" {
		req.Header.Set("Cookie", c.Cookie)
	}
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

// Login authenticates a new APIC session
func (c *Client) Login() error {
	a := loginAttributes{
		Name: c.Username,
		Pwd:  c.Password,
	}
	l := LoginRequest{AAA: AAA{
		loginAttributes: a,
	}}
	req, err := c.newRequest("POST", loginPath, l)
	if err != nil {
		return err
	}

	var lr LoginResponse
	resp, err := c.do(req, &lr)
	if err != nil {
		return err
	}
	// get auth cookie
	// TODO check cookie name?
	apicCookie := resp.Cookies()[0]
	c.Cookie = apicCookie.String()
	fmt.Println("response Status:", resp.Status)
	return nil
}
