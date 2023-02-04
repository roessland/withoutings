package withings

// https://developer.withings.com/api-reference#operation/notify-subscribe

// NewNotifySubscribeParams creates new NewNotifySubscribeParams with some defaults.
func NewNotifySubscribeParams() NotifySubscribeParams {
	return NotifySubscribeParams{
		Action: "subscribe",
	}
}

// NotifySubscribeParams are the parameters for Notify - Subscribe
type NotifySubscribeParams struct {
	Action      string `json:"action" url:"action"`
	Callbackurl string `json:"callbackurl" url:"callbackurl"`
	Appli       int    `json:"appli" url:"appli"`
	Comment     string `json:"comment" url:"comment"`
}

type NotifySubscribeResponse struct {
	Status int    `json:"status"`
	Raw    []byte `json:"-"`
}
