package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	dialTimeout         = 60 * time.Second
	httpClientTimeout   = 120 * time.Second
	tlsHandshakeTimeout = 60 * time.Second

	// ContentTypeJSON json content type header
	ContentTypeJSON = "application/json"
	// ContentTypeForm form data content type header
	ContentTypeForm = "application/x-www-form-urlencoded"
)

type Client struct {
	HTTPClient *http.Client
	Addr       string
	Username   string
	Password   string
}

//TODo use existing api.ErrorResponse here?
// ClientError is used for non-200 API responses
type ClientError struct {
	Status     string
	StatusCode int
	Message    string
}

//TODo replace this
// HTTPError is included in an HTTPResponse
type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewClientError creates a ClientError
func NewClientError(status string, statusCode int, message string) ClientError {
	return ClientError{
		Status:     status,
		StatusCode: statusCode,
		Message:    strings.TrimRight(message, "\n"),
	}
}

func (e ClientError) Error() string {
	return e.Message
}

// ReceivedHTTPResponse parsed a HTTPResponse received by the Client, for the V2 API
type ReceivedHTTPResponse struct {
	Error *HTTPError      `json:"error,omitempty"`
	Data  json.RawMessage `json:"data"`
}

// NewClient creates a Client
func NewClient(addr string) *Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: dialTimeout,
		}).Dial,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   httpClientTimeout,
	}
	addr = strings.TrimRight(addr, "/")
	addr += "/"

	return &Client{
		Addr:       addr,
		HTTPClient: httpClient,
	}
}

// Get makes a GET request to an endpoint and unmarshals the response to obj.
// If the response is not 200 OK, returns an error
func (c *Client) Get(endpoint string, obj interface{}) error {
	resp, err := c.get(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return NewClientError(resp.Status, resp.StatusCode, string(body))
	}

	if obj == nil {
		return nil
	}

	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()
	return d.Decode(obj)
}

// Put makes a PUT request with provided body to an endpoint and unmarshals the response to obj.
// If the response is not 200 OK, returns an error
func (c *Client) Put(endpoint string, body, obj interface{}) error {
	resp, err := c.put(endpoint, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return NewClientError(resp.Status, resp.StatusCode, string(body))
	}

	if obj == nil {
		return nil
	}

	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()
	return d.Decode(obj)
}

// get makes a GET request to an endpoint. Caller must close response body.
func (c *Client) get(endpoint string) (*http.Response, error) {
	return c.makeRequestWithoutBody(endpoint, http.MethodGet)
}

// post makes a PUT request to an endpoint with provided body. Caller must close response body.
func (c *Client) put(endpoint string, body interface{}) (*http.Response, error) {
	return c.makeRequestWithBody(body, endpoint, http.MethodPut)
}

// makeRequestWithoutBody makes a `method` request to an endpoint. Caller must close response body.
func (c *Client) makeRequestWithoutBody(endpoint, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

// makeRequestWithBody makes a `method` request to an endpoint. Caller must close response body.
func (c *Client) makeRequestWithBody(body interface{}, endpoint, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, body.(io.Reader))
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}
