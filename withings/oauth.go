package withings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var OAuth2Scopes = []string{"user.activity,user.metrics"}
var OAuth2AuthURL string = "https://account.withings.com/oauth2_user/authorize2"
var OAuth2TokenURL = "https://wbsapi.withings.net/v2/oauth2"

type GetAccessTokenResponse struct {
	Status int
	Body   Token `json:"body"`
}

type Token struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	CSRFToken    string `json:"csrf_token"`
	TokenType    string `json:"token_type,omitempty"`
	Expiry       time.Time
}

func (c *Client) GetAccessToken(ctx context.Context, authCode string) (*Token, error) {

	v := url.Values{}
	v.Set("action", "requesttoken")
	v.Set("client_id", c.OAuth2Config.ClientID)
	v.Set("client_secret", c.OAuth2Config.ClientSecret)
	v.Set("grant_type", "authorization_code")
	v.Set("code", authCode)
	v.Set("redirect_uri", c.OAuth2Config.RedirectURL)

	req, err := http.NewRequest(http.MethodPost, OAuth2TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %w", err)
	}
	if code := resp.StatusCode; code < 200 || code > 299 {
		return nil, &oauth2.RetrieveError{
			Response: resp,
			Body:     body,
		}
	}

	var response *GetAccessTokenResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Body.AccessToken == "" {
		return nil, errors.New("oauth2: server response missing access_token")
	}
	return &response.Body, nil
}
