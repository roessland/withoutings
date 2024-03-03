# Account registration sequence

## Initial registration

Flow when user has never logged in before and an account needs to be created.

- User has never logged in before.
- Users clicks login link
- User logs in with Withings auth provider
- User is redirected to Withoutings callback
- Webapp retrieves access token from Withings
- Webapp creates an account in the database
- Webapp stores user ID in session

```mermaid
sequenceDiagram
participant User
participant Webapp
participant Withings API

User->>Webapp: GET /auth/login
User->>Webapp: POST /auth/redirect-to-withings-login
activate Webapp
note over Webapp: Generates nonce.<br>Saves in session.
Webapp-->>User: 302 Found<br>Redirect: Withings Auth URL
deactivate Webapp
Note over User: Saves session cookie.<br>Redirects to Withings Auth URL.

Note over User: Logs in with Withings OAuth. <br>Is redirected to callback.


User->>Webapp: GET /auth/callback<br>{code: <code>}
note right of User: Request payload:<code>
activate Webapp
    Note over Webapp: Checks nonce.
    
    Webapp->>Withings API: POST https://wbsapi.withings.net/v2/oauth2
    activate Withings API
    Note right of Webapp: Request payload:<br>action=requesttoken<br>client_id=<client_id><br>client_secret=<client_secret><br>grant_type=authorization_code<br>code=<code><br>redirect_uri=<redirect_url>
    Withings API->>Webapp: 200 OK
    deactivate Withings API
    
    Note left of Withings API: Response payload<br>userid<br>access_token<br>refresh_token<br>expires_in<br>scope<br>csrf_token
    
    Note over Webapp: Deletes stored nonce<br>from session.
    
    Note over Webapp: Creates account in database.
    
    Note over Webapp:Stores user ID in session.
    
    Webapp->>User: 302 Found<br>Redirect to homepage

deactivate Webapp
Note over User: Saves session cookie
```
fdsfdsf