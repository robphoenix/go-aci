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

// LoginJSON represents the JSON needed for authentication
type LoginJSON struct {
	AAAUser `json:"aaaUser"`
}

// AAAUser ...
type AAAUser struct {
	loginAttributes `json:"attributes"`
}

// Attributes ...
type loginAttributes struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
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

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}

// Login authenticates a new APIC session
func (c *Client) Login() error {
	l := LoginJSON{AAAUser: AAAUser{loginAttributes: loginAttributes{
		Name: c.Username,
		Pwd:  c.Password,
	}}}

	req, err := c.newRequest("POST", loginPath, l)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
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
