package subscription

type NotificationCategory struct {
	Appli       int
	Scope       string
	Description string
}

// TODO scopes are wrong in DB, add a migration to fix

var NotificationCategoryByAppli = map[int]NotificationCategory{
	1:  {Appli: 1, Scope: "user.metrics", Description: "New weight-related data"},
	2:  {Appli: 2, Scope: "user.metrics", Description: "New temperature related data"},
	4:  {Appli: 4, Scope: "user.metrics", Description: "New pressure related data"},
	16: {Appli: 16, Scope: "users.activity", Description: "New activity-related data"},
	44: {Appli: 44, Scope: "users.activity", Description: "New sleep-related data"},
	46: {Appli: 46, Scope: "user.info", Description: "New action on user profile"},
	50: {Appli: 50, Scope: "user.sleepevents", Description: "New bed in event"},
	51: {Appli: 51, Scope: "user.sleepevents", Description: "New bed out event"},
	52: {Appli: 52, Scope: "user.sleepevents", Description: "New inflate done event"},
	53: {Appli: 53, Scope: "n/a", Description: "No account associated"},
	54: {Appli: 54, Scope: "user.metrics", Description: "New ECG data"},
	55: {Appli: 55, Scope: "user.metrics", Description: "ECG measure failed event"},
	58: {Appli: 58, Scope: "user.metrics", Description: "New glucose data"},
}
