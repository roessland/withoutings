package subscription

import "errors"

// Services that can be called to fetch data from the Withings API.
// https://developer.withings.com/developer-guide/v3/data-api/keep-user-data-up-to-date/

type NotificationDataService string

const NotificationDataServiceMeasureGetMeas NotificationDataService = "Measure - Getmeas"
const NotificationDataServiceMeasurev2Getactivity NotificationDataService = "Measure v2 - Getactivity"
const NotificationDataServiceMeasurev2Getintradayactivity NotificationDataService = "Measure v2 - Getintradayactivity"
const NotificationDataServiceSleepv2Get NotificationDataService = "Sleep v2 - Get"
const NotificationDataServiceSleepv2Getsummary NotificationDataService = "Sleep v2 - Getsummary"
const NotificationDataServiceHeartv2List NotificationDataService = "Heart v2 - List"

func NewNotificationDataService(s string) (NotificationDataService, error) {
	switch s {
	case "Measure - Getmeas":
		return NotificationDataServiceMeasureGetMeas, nil
	case "Measure v2 - Getactivity":
		return NotificationDataServiceMeasurev2Getactivity, nil
	case "Measure v2 - Getintradayactivity":
		return NotificationDataServiceMeasurev2Getintradayactivity, nil
	case "Sleep v2 - Get":
		return NotificationDataServiceSleepv2Get, nil
	case "Sleep v2 - Getsummary":
		return NotificationDataServiceSleepv2Getsummary, nil
	case "Heart v2 - List":
		return NotificationDataServiceHeartv2List, nil
	default:
		return "", errors.New("unknown notification data service")
	}
}

func MustNewNotificationDataService(s string) NotificationDataService {
	d, err := NewNotificationDataService(s)
	if err != nil {
		panic(err)
	}
	return d
}
