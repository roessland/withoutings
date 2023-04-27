package withings

// https://developer.withings.com/api-reference#operation/notify-list

// NewNotifyListParams creates new NotifyListParams with some defaults.
func NewNotifyListParams(appli int) NotifyListParams {
	return NotifyListParams{
		Action: "list",
		Appli:  appli,
	}
}

// NotifyListParams are the parameters for Notify - List
type NotifyListParams struct {
	Action string `json:"action" url:"action"`
	Appli  int    `json:"appli" url:"appli"`
}

type NotifyListResponse struct {
	Status int            `json:"status"`
	Body   NotifyListBody `json:"body"`
	Raw    []byte         `json:"-"`
}

type NotifyListBody struct {
	Profiles []NotifyListProfile `json:"profiles"`
}

type NotifyListProfile struct {
	Appli       int    `json:"appli"`
	CallbackURL string `json:"callbackurl"`
	Expires     int64  `json:"expires"`
	Comment     string `json:"comment"`
}