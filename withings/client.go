package withings

import (
	"github.com/roessland/withoutings/withings/openapi"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

const APIBase = "https://wbsapi.withings.net"

type Client struct {
	HTTPClient   *http.Client
	OAuth2Config oauth2.Config
	API2         *openapi.Client
}

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

	var err error
	c.API2, err = openapi.NewClient(APIBase, openapi.WithHTTPClient(c.HTTPClient))
	if err != nil {
		log.Fatal(err)
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
