package account

import "time"

type Account struct {
	AccountID                 int64
	WithingsUserID            string
	WithingsAccessToken       string
	WithingsRefreshToken      string
	WithingsAccessTokenExpiry time.Time
	WithingsScopes            string
}
