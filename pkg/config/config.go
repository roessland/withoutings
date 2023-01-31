package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// SessionSecret is used to encrypt cookies.
	SessionSecret []byte `envconfig:"SESSION_SECRET"`

	// WebsiteURL is the public URL where the website is accessible.
	WebsiteURL string `envconfig:"WEBSITE_URL"`

	// WithingsClientID generated in Withings Developer Dashboard.
	WithingsClientID string `envconfig:"WITHINGS_CLIENT_ID"`

	// WithingsClientSecret generated in Withings Developer Dashboard.
	WithingsClientSecret string `envconfig:"WITHINGS_CLIENT_SECRET"`

	// WithingsRedirectURL is where OAuth2 logins from Withings are redirected to.
	// E.g. https://withings.example.com/auth/callback .
	// This must be added in Withings Developer Dashboard.
	WithingsRedirectURL string `envconfig:"WITHINGS_REDIRECT_URL"`

	// WithingsWebhookSecret adds extra security by obscurity to the
	// incoming webhook handler. Add the URL
	// "https://withings.example.com/withings/webhooks/{secret}"
	// as an additional callback URL in Withings Developer Dashboard.
	WithingsWebhookSecret string `envconfig:"WITHINGS_WEBHOOK_SECRET"`

	// DatabaseURL is a PostgreSQL connection string, e.g.
	// "postgres://wotsa:<pass>@127.0.0.1:5432/wot?sslmode=disable".
	// If you have separate admin/superuser and read/write users,
	// this should use the read/write user with fewer permissions.
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func LoadFromEnv() (*Config, error) {
	c := Config{}
	err := envconfig.Process("WOT", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
