package withings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var OAuth2Scopes = []string{"user.info,user.activity,user.metrics,user.sleepevents"}
var OAuth2AuthURL = "https://account.withings.com/oauth2_user/authorize2"
var OAuth2TokenURL = "https://wbsapi.withings.net/v2/oauth2"

func (c *HTTPClient) GetAccessToken(ctx context.Context, authCode string) (*withings.Token, error) {
	v := url.Values{}
	v.Set("action", "requesttoken")
	v.Set("client_id", c.OAuth2Config.ClientID)
	v.Set("client_secret", c.OAuth2Config.ClientSecret)
	v.Set("grant_type", "authorization_code")
	v.Set("code", authCode)
	v.Set("redirect_uri", c.OAuth2Config.RedirectURL)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.OAuth2Config.Endpoint.TokenURL,
		strings.NewReader(v.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
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

	var response *withings.GetAccessTokenResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	response.Body.Expiry = time.Now().UTC().Add(time.Second * time.Duration(response.Body.ExpiresIn))

	if response.Body.AccessToken == "" {
		fmt.Println("body", string(body))
		// {"status":503,"body":{},"error":"Invalid Params: invalid code"}
		return nil, errors.New("oauth2: server response missing access_token")
	}
	return &response.Body, nil
}

func (c *HTTPClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*withings.Token, error) {
	v := url.Values{}
	v.Set("action", "requesttoken")
	v.Set("grant_type", "refresh_token")
	v.Set("client_id", c.OAuth2Config.ClientID)
	v.Set("client_secret", c.OAuth2Config.ClientSecret)
	v.Set("refresh_token", refreshToken)

	log := logging.MustGetLoggerFromContext(ctx)
	log.WithField("Post body", v.Encode()).Info()

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.OAuth2Config.Endpoint.TokenURL,
		strings.NewReader(v.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("creating token refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot refresh token: %w", err)
	}
	if code := resp.StatusCode; code < 200 || code > 299 {
		return nil, &oauth2.RetrieveError{
			Response: resp,
			Body:     body,
		}
	}

	var response *withings.GetAccessTokenResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	response.Body.Expiry = time.Now().UTC().Add(time.Second * time.Duration(response.Body.ExpiresIn))

	if response.Body.AccessToken == "" {
		fmt.Println("body", string(body))
		return nil, errors.New("oauth2: server response missing access_token")
	}
	return &response.Body, nil
}

func (c *HTTPClient) AuthCodeURL(nonce string) string {
	return c.OAuth2Config.AuthCodeURL(nonce)
}