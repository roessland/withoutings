# Domain-driven design

## Aggregates


```mermaid
flowchart LR;

subgraph WithoutingsAggregate
    Account
    Account-->|has many|Subscriptions
    Account-->|has a|Session
    Account-->|has many|RawNotifications
end
```