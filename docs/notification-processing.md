# Notification processing

## Timeline
- Withings calls our webhook
- RawNotification is created in DB
- RawNotificationReceived event is emitted
- Notification is created in DB
- NotificationReceived event is emitted
- Notification data is fetched
- Notification data is added to the notification in DB
- NotificationDataFetched is emited

## Raw notification status finite state machine

```mermaid
flowchart TD;
null --> |POST /withings/webhooks/secret| pending
pending -->|WithingsRawNotificationReceived| processed
processed -->|WithingsRawNotificationReceived| processed
null --> |WithingsRawNotificationReceived| null
```
TODO: Implement all these state transitions in the code

## Notification data finite state machine


"awaiting_fetch"
"fetched"
"fetch_failed"


`notification.data_status` flowchart:

```mermaid
flowchart TD;
null --> |WithingsRawNotificationReceived| awaiting_fetch
awaiting_fetch -->|NotificationDataReceived| fetched
fetched -->|NotificationDataReceived| fetched
awaiting_fetch -->|NotificationDataReceived| fetch_failed
fetch_failed -->|WithingsRawNotificationReceived| fetch_failed
fetch_failed -->|WithingsRawNotificationReceived| fetched
fetch_failed -->|NotificationDataReceived| fetched

```

TODO: Implement all these state transitions in the code