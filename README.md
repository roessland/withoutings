# Withoutings

Demo application for talking with Withings API.

## Setup

Go to the [Withings Developer Dashboard](https://developer.withings.com/dashboard/).
Create a new application.

## Required environment variables

```bash
SESSION_SECRET=<random uuid, generated once> \
WITHINGS_CLIENT_ID=<your app id> \
WITHINGS_CLIENT_SECRET=<your app secret> \
WITHINGS_REDIRECT_URL=<your app callback URL> \
go run cmd/main.go
```

## Forward remote port to local port

To make callback to a remote server call your development app you can
use SSH port forwarding.

```bash
ssh -R 3628:127.0.0.1:3628 -N -f myuser@withings.mywebsite.com
```

So Withings calls `https://withings.mywebsite.com/auth/callback` which is
forwarded to port 3628 on the server (e.g. using Caddy or nginx), which
is again forwarded to port 3628 in your development environment.