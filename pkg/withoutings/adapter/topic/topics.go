package topic

const WithingsRawNotificationReceived = "withings_raw_notification_received"

const WithingsNotificationReceived = "withings_notification_received"

const WithingsNotificationDataFetched = "withings_notification_data_fetched"

const WithingsNotificationDataFetchFailed = "withings_notification_data_fetch_failed"

var AllTopics = []string{
	WithingsRawNotificationReceived,
	WithingsNotificationReceived,
	WithingsNotificationDataFetched,
	WithingsNotificationDataFetchFailed,
}
