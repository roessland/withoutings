# Domain-driven design

## Dependency graph

```mermaid
flowchart LR;


subgraph withingsapidomain
    WithingsAPIClient
end

subgraph infrastructure
    PostgreSQL
    WithingsAPI
end

subgraph accountdomain
    account.Account
    account.Repo
end

subgraph subscriptionsdomain
    subscription.Subscription
    subscription.RawNotification
    subscription.Repo
end

subgraph adapters
    AccountPgRepo-->account.Repo
    AccountPgRepo-->PostgreSQL
    SubscriptionPgRepo-->subscription.Repo
    SubscriptionPgRepo-->PostgreSQL
    WithingsAPIDefaultClient-->WithingsAPI
    WithingsAPIDefaultClient-->WithingsAPIClient
end

subgraph commands
    CreateOrUpdateAccount-->account.Repo
    CreateOrUpdateAccount-->account.Account
    SubscribeAccount-->account.Account
    SubscribeAccount-->WithingsAPIClient
end

subgraph queries
    AccountByWithingsUserID-->account.Account
    AccountByWithingsUserID-->account.Repo
end

subgraph services
    App-->AccountByWithingsUserID
    App-->SubscribeAccount
    App-->CreateOrUpdateAccount
    App-->account.Repo
    App-->AccountPgRepo
    App-->WithingsAPIClient
    App-->WithingsAPIDefaultClient
end

subgraph handlers
    direction LR
    Homepage
    Health
    Callback
    Logout
    Login
    SleepSummaries
   
end
handlers-->App


```