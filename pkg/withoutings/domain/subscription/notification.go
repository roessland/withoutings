package subscription

type RawNotification struct {
	RawNotificationID int64
	Source            string
	Data              string
	Status            RawNotificationStatus
}

type RawNotificationStatus string

const RawNotificationStatusPending RawNotificationStatus = "pending"
const RawNotificationStatusProcessed RawNotificationStatus = "processed"

func NewRawNotification(source string, data string) RawNotification {
	return RawNotification{
		Source: source,
		Data:   data,
		Status: RawNotificationStatusPending,
	}
}
