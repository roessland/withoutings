package subscription

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
)

type Subscription struct {
	SubscriptionID int64
	AccountID      int64
	Appli          int
	CallbackURL    string
	Comment        string
	WebhookSecret  string
	Status         Status
}

type Status string

const StatusPending Status = "pending"
const StatusActive Status = "active"
const StatusUnlinked Status = "unlinked"
const StatusUserDeleted Status = "user-deleted"

func NewSubscription(accountID int64, appli int, callbackURL string, webhookSecret string) Subscription {
	return Subscription{
		AccountID:     accountID,
		Appli:         appli,
		CallbackURL:   callbackURL,
		Status:        StatusPending,
		WebhookSecret: webhookSecret,
	}
}

var b64nonAlphaNumeric = regexp.MustCompile(`[=+/]`)

func RandomWebhookSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return b64nonAlphaNumeric.ReplaceAllString(base64.RawURLEncoding.EncodeToString(b), "")
}
