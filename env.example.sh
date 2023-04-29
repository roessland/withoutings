# See config.go for explanations.
export WOT_LISTEN_ADDR='127.0.0.1:3628';
export WOT_SESSION_SECRET=asdfasdf
export WOT_WEBSITE_URL='https://wot.mywebsite.com/'
export WOT_WITHINGS_CLIENT_ID=asdfasdf
export WOT_WITHINGS_CLIENT_SECRET=asdfasdf
export WOT_WITHINGS_REDIRECT_URL=https://withings.example.com/auth/callback
# Register {WOT_WEBSITE_URL}withings/webhooks/{secret}
# in Withings Developer Dashboard.
export WOT_WITHINGS_WEBHOOK_SECRET=supersecret
export WOT_DATABASE_URL='postgres://wotrw:<pass>@127.0.0.1:5432/wot?sslmode=disable'
export WOT_DATABASE_URL_SA='postgres://wotsa:<pass>@127.0.0.1:5432/wot?sslmode=disable'