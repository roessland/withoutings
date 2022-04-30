package withings

import (
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

const APIBase = "https://wbsapi.withings.net"

type Client struct {
	HTTPClient   *http.Client
	OAuth2Config oauth2.Config
}

func NewClient(clientID, clientSecret, redirectURL string) *Client {
	c := Client{}
	c.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout:   5 * time.Second,
			MaxIdleConns:          5,
			MaxConnsPerHost:       10,
			IdleConnTimeout:       5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
		Timeout: time.Second * 10,
	}

	c.OAuth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       OAuth2Scopes,
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  OAuth2AuthURL,
			TokenURL: OAuth2TokenURL,
		},
	}

	return &c
}
