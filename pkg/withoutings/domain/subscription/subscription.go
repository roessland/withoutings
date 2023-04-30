package subscription

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"regexp"
)

type Subscription struct {
	subscriptionUUID uuid.UUID
	accountUUID      uuid.UUID
	appli            int
	callbackURL      string
	comment          string
	webhookSecret    string
	status           Status
}

type Status string

const StatusPending Status = "pending"
const StatusActive Status = "active"
const StatusUnlinked Status = "unlinked"
const StatusUserDeleted Status = "user-deleted"

// NewSubscription returns a new subscription.
// TODO should return (sub, error)
func NewSubscription(
	subscriptionUUID uuid.UUID,
	accountUUID uuid.UUID,
	appli int,
	callbackURL string,
	comment string,
	webhookSecret string,
	status Status,
) Subscription {
	return Subscription{
		subscriptionUUID: subscriptionUUID,
		accountUUID:      accountUUID,
		appli:            appli,
		callbackURL:      callbackURL,
		comment:          comment,
		webhookSecret:    webhookSecret,
		status:           status,
	}
}

func (s Subscription) UUID() uuid.UUID {
	return s.subscriptionUUID
}

func (s Subscription) AccountUUID() uuid.UUID {
	return s.accountUUID
}

func (s Subscription) Appli() int {
	return s.appli
}

func (s Subscription) CallbackURL() string {
	return s.callbackURL
}

func (s Subscription) Comment() string {
	return s.comment
}

func (s Subscription) WebhookSecret() string {
	return s.webhookSecret
}

func (s Subscription) Status() Status {
	return s.status
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
