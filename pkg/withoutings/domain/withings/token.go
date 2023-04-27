package withings

import "time"

type Token struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Scope        string
	CSRFToken    string
	TokenType    string
	Expiry       time.Time
}
