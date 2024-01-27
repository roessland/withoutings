package withings

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// DefaultAPIBase is the base URL to the Withings API.
const DefaultAPIBase = "https://wbsapi.withings.net"

// HTTPClient sends requests to the Withings API.
type HTTPClient struct {
	HTTPClient   *http.Client
	OAuth2Config *oauth2.Config
	APIBase      string
}

// AuthenticatedClient is a client with an access token.
type AuthenticatedClient struct {
	HTTPClient
	AccessToken string
}

// NewClient returns a client.
func NewClient(clientID, clientSecret, redirectURL string) *HTTPClient {
	c := HTTPClient{}

	c.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout:   25 * time.Second,
			MaxIdleConns:          5,
			MaxConnsPerHost:       10,
			IdleConnTimeout:       25 * time.Second,
			ResponseHeaderTimeout: 25 * time.Second,
		},
		Timeout: time.Second * 35,
	}

	c.OAuth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       OAuth2Scopes,
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  OAuth2AuthURL,
			TokenURL: OAuth2TokenURL,
		},
	}

	c.APIBase = DefaultAPIBase

	return &c
}

// WithAccessToken returns an authenticated version of a client
func (c *HTTPClient) WithAccessToken(accessToken string) *AuthenticatedClient {
	var ac AuthenticatedClient
	ac.HTTPClient = *c
	ac.AccessToken = accessToken
	return &ac
}

// NewRequest creates a standard request to the Withings API.
// Params can be either a struct, as accepted by github.com/google/go-querystring/query,
// or a string, which will be used directly, without further encoding.
func (c *HTTPClient) NewRequest(endpoint string, params any, body []byte) (*http.Request, error) {
	var encodedParams string

	q, err := query.Values(params)
	if err != nil {
		fmt.Println("error encoding params", err)
		encodedParams = fmt.Sprintf("%v", params)
	} else {
		encodedParams = q.Encode()
	}

	url := fmt.Sprintf("%s%s?%s", c.APIBase, endpoint, encodedParams)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Do sends a request
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	var err error
	log := logging.MustGetLoggerFromContext(req.Context())

	var reqBody []byte
	if req.Body != nil {
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	log.WithFields(logrus.Fields{
		"event":           "withings-api-request.started",
		"url":             req.URL.String(),
		"request.body":    string(reqBody),
		"request.headers": req.Header,
		"request.method":  req.Method,
	}).Info()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(respBody))

	log.WithFields(logrus.Fields{
		"event":            "withings-api-request.finished",
		"response.body":    string(respBody),
		"response.headers": resp.Header,
		"response.status":  resp.StatusCode,
	}).Info()

	return resp, nil
}

// Do sends a request with an authorization header.
func (c *AuthenticatedClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	return c.HTTPClient.Do(req)
}
