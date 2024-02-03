package subscription

// FSM for notification data status:
// * awaiting_fetch -> fetched
// * awaiting_fetch -> fetch_failed

type NotificationDataStatus string

const NotificationDataStatusAwaitingFetch NotificationDataStatus = "awaiting_fetch"
const NotificationDataStatusFetched NotificationDataStatus = "fetched"
const NotificationDataStatusFetchFailed NotificationDataStatus = "fetch_failed"
