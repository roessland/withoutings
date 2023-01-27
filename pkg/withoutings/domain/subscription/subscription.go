package subscription

type Subscription struct {
	SubscriptionID int64
	AccountID      int64
	Appli          int
	CallbackURL    string
	Comment        string
}

func NewSubscription(accountID int64) {

}
