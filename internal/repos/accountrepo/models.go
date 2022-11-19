// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package accountrepo

import (
	"database/sql"
	"time"
)

type Account struct {
	AccountID                 int64
	WithingsUserID            string
	WithingsAccessToken       string
	WithingsRefreshToken      string
	WithingsAccessTokenExpiry time.Time
	WithingsScopes            string
}

type Subscription struct {
	SubscriptionID int32
	AccountID      sql.NullInt64
	Appli          int32
	Callbackurl    string
	Comment        string
}