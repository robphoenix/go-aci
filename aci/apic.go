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

// Client manages communication with the APIC API
type Client struct {
	Host       *url.URL
	Username   string
	Password   string
	Cookie     string
	httpClient *http.Client
}

// loginRequest is the JSON request for
// authenticating with APIC
type loginRequest struct {
	AAA `json:"aaaUser"`
}

// AAA is part of the authentication process
// that holds authentication attributes
type AAA struct {
	loginAttributes `json:"attributes"`
}

// loginAttributes is the attributes of the APIC
// authentication request and response
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
		Host:     &url.URL{Scheme: "https", Host: host},
		Username: username,
		Password: password,
		httpClient: &http.Client{
			Transport: T,
			Timeout:   15 * time.Second,
		},
	}, nil
}

// newRequest forms an http request for use with an APIC client
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
		return nil, fmt.Errorf("%s request to %s : %v", method, u.String(), err)
	}
	if c.Cookie != "" {
		req.Header.Set("Cookie", c.Cookie)
	}
	return req, nil
}

// do performs APIC client http requests
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"%s response with %s request from %s : %v",
			resp.Status,
			req.Method,
			req.URL.String(),
			err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

// Login authenticates a new APIC session
// adding the apicCookie to the client
func (c *Client) Login() error {
	var lr loginRequest
	lr.Name = c.Username
	lr.Pwd = c.Password
	req, err := c.newRequest("POST", loginPath, lr)
	if err != nil {
		return fmt.Errorf("login for %s: %v", lr.Name, err)
	}

	var la loginAttributes
	resp, err := c.do(req, &la)
	if err != nil {
		return fmt.Errorf("login for %s: %v", lr.Name, err)
	}
	// get auth cookie
	// TODO check cookie name?
	apicCookie := resp.Cookies()[0]
	c.Cookie = apicCookie.String()
	return nil
}
