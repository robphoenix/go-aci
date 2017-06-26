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
	// ErrAlreadyDiscovered - Can't remove node identity policy - Node TEP-1-102 is already discovered. Please decommission first.
	ErrAlreadyDiscovered = "107"
)

var (
	httpTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
	}
	// tlsConfig is the TLS config
	// https://stackoverflow.com/questions/41250665/go-https-client-issue-remote-error-tls-handshake-failure#
	tlsConfig = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true,
		MinVersion:               tls.VersionTLS11,
		MaxVersion:               tls.VersionTLS11,
	}
	clientTimeout = 15 * time.Second
)

// ClientOptions specifies the client connection options
type ClientOptions struct {
	Host     string
	Username string
	Password string
}

// Client manages communication with the APIC API
type Client struct {
	Host       *url.URL
	Username   string
	Password   string
	cookie     string
	httpClient *http.Client
}

// NewClient instantiates a new APIC client
func NewClient(o ClientOptions) (*Client, error) {
	return &Client{
		Host:     &url.URL{Scheme: "https", Host: o.Host},
		Username: o.Username,
		Password: o.Password,
		httpClient: &http.Client{
			Transport: httpTransport,
			Timeout:   clientTimeout,
		},
	}, nil
}

// NewRequest forms an http request for use with an APIC client
func (c *Client) NewRequest(method string, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := c.Host.ResolveReference(rel)

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
	if c.Cookie() != "" {
		req.Header.Set("Cookie", c.Cookie())
	}
	return req, nil
}

// Do performs APIC client http requests
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v %v: %d %v", req.Method, req.URL.String(), resp.StatusCode, err)
	}
	defer resp.Body.Close()
	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

// SetCookie sets the value of the APIC authentication cookie
func (c *Client) SetCookie(r *http.Response) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "APIC-cookie" {
			c.cookie = cookie.String()
		}
	}
}

// Cookie returns the APIC authentication cookie.
// Returns an empty string if it has not been set.
func (c *Client) Cookie() string {
	return c.cookie
}

// loginRequest is the JSON request for authenticating with the APIC
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

// Login authenticates a new APIC session, setting the authentication cookie
func (c *Client) Login() error {
	var lr loginRequest
	lr.Name = c.Username
	lr.Pwd = c.Password

	req, err := c.NewRequest(http.MethodPost, loginPath, lr)
	if err != nil {
		return fmt.Errorf("login for %s: %v", lr.Name, err)
	}

	var la loginAttributes
	resp, err := c.Do(req, &la)
	if err != nil {
		return err
	}

	// set auth cookie
	c.SetCookie(resp)

	return nil
}

// ErrorResponse reports any errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Imdata   []struct {     // Imdata is an array of errors
		ImdataError struct {
			ErrorAttributes struct {
				Code string `json:"code"` // APIC specific error code
				Text string `json:"text"` // Error Text
			} `json:"attributes"`
		} `json:"error"`
	} `json:"imdata"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v (APIC Code: %+v)",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode,
		// while Imdata is an array of errors, we only want the first one (I think)
		r.Imdata[0].ImdataError.ErrorAttributes.Text,
		r.Imdata[0].ImdataError.ErrorAttributes.Code,
	)
}

// CheckResponse checks the API response for errors, and returns them if
// present.
//
// A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have a JSON
// response body that maps to ErrorResponse. Any other response body
// will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	json.NewDecoder(r.Body).Decode(errorResponse)
	return errorResponse
}

// Mapper wraps basic key/value methods.
// This is useful for building hashmaps of ACI objects.
type Mapper interface {
	Key() string
	Value() Mapper
}
