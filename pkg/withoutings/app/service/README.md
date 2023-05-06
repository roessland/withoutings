# Domain services

Domain services are used to implement domain logic that does not fit into a single entity. They are stateless and are
used by application services (commands and queries).

Input argument and return types should come from the domain.

## Withings service

The Withings service is used to interact with the Withings API. It encapsulates the Withings API repository, and
provides automated access token refreshing that is seamless to the application.

It depends on the Withings API client, and the account repository (domain included).