package withings

import (
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

// DefaultAPIBase is the base URL to the Withings API.
const DefaultAPIBase = "https://wbsapi.withings.net"

// Client sends requests to the Withings API.
type Client struct {
	HTTPClient   *http.Client
	OAuth2Config *oauth2.Config
	APIBase      string
}

// AuthenticatedClient is a client with an access token.
type AuthenticatedClient struct {
	Client
	AccessToken string
}

// NewClient returns a client.
func NewClient(clientID, clientSecret, redirectURL string) *Client {
	c := Client{}

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
func (c *Client) WithAccessToken(accessToken string) *AuthenticatedClient {
	var ac AuthenticatedClient
	ac.Client = *c
	ac.AccessToken = accessToken
	return &ac
}

// Do sends a request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

// Do sends a request with an authorization header.
func (c *AuthenticatedClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	return c.HTTPClient.Do(req)
}
