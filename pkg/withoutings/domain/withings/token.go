package withings

import "time"

type GetAccessTokenResponse struct {
	Status int
	Body   Token `json:"body"`
}

type Token struct {
	UserID       string    `json:"userid,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresIn    int       `json:"expires_in,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	CSRFToken    string    `json:"csrf_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}
