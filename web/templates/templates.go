package templates

import (
	"embed"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"html/template"
	"io"
)

//go:embed templates/*.gohtml
var fs embed.FS

type Templates struct {
	template *template.Template
}

func LoadTemplates() Templates {
	templates := Templates{}
	t, err := template.ParseFS(fs, "*/**")
	if err != nil {
		panic(err)
	}
	t.Option("missingkey=error")
	templates.template = t
	return templates
}

type HomePageVars struct {
	Error   string
	Account *account.Account
}

func (t Templates) RenderHomePage(w io.Writer, account_ *account.Account) error {
	return t.template.ExecuteTemplate(w, "homepage.gohtml", HomePageVars{
		Account: account_,
	})
}

type RefreshAccessTokenVars struct {
	Token *withings.Token
	Error string
}

func (t Templates) RenderRefreshAccessToken(w io.Writer, token *withings.Token, err string) error {
	return t.template.ExecuteTemplate(w, "refreshaccesstoken.gohtml", RefreshAccessTokenVars{
		Token: token,
		Error: err,
	})
}

type SleepSummariesVars struct {
	Error     string
	Token     *withings.Token
	SleepData interface{}
}

func (t Templates) RenderSleepSummaries(w io.Writer, sleepData *sleep.GetSleepSummaryOutput, err string) error {
	return t.template.ExecuteTemplate(w, "sleepsummaries.gohtml", SleepSummariesVars{
		SleepData: sleepData,
		Error:     err,
	})
}

type SubscriptionsPageVars struct {
	Error         string
	Subscriptions []subscription.Subscription
}

func (t Templates) RenderSubscriptionsPage(w io.Writer, subscriptions []subscription.Subscription) error {
	return t.template.ExecuteTemplate(w, "subscriptionspage.gohtml", SubscriptionsPageVars{
		Subscriptions: subscriptions,
	})
}
