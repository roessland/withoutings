package templates

import (
	"embed"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"runtime"
)

//go:embed templates/*.gohtml
var embeddedFS embed.FS

// FS is either disk or embedded files.
var FS fs.FS

// init sets FS if the templates are available on disk.
func init() {
	var err error

	// Set FS to use embedded files.
	FS, err = fs.Sub(embeddedFS, "templates")
	if err != nil {
		panic(err.Error())
	}

	// If the templates are available on disk, set FS to use disk files instead.
	_, templatesGoPath, _, _ := runtime.Caller(0)
	templatesDir := path.Dir(templatesGoPath)
	stat, err := os.Stat(path.Join(templatesDir, "templates/base.gohtml"))
	if err == nil && stat.Size() > 0 {
		FS = os.DirFS(path.Join(templatesDir, "templates"))
		logrus.Info("Using disk files for templates")
	}
}

type Templates struct {
}

func NewTemplates() *Templates {
	t := &Templates{}
	return t
}

type HomePageVars struct {
	Error   string
	Account *account.Account
}

func (t *Templates) RenderHomePage(w io.Writer, account_ *account.Account) error {
	tmpl, err := template.New("base.gohtml").ParseFS(FS, "base.gohtml", "homepage.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", HomePageVars{
		Account: account_,
	})
}

type RefreshAccessTokenVars struct {
	Token *withings.Token
	Error string
}

func (t *Templates) RenderRefreshAccessToken(w io.Writer, token *withings.Token, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(FS, "base.gohtml", "refreshaccesstoken.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", RefreshAccessTokenVars{
		Token: token,
		Error: errMsg,
	})
}

type SleepSummariesVars struct {
	Error     string
	Token     *withings.Token
	SleepData interface{}
}

func (t *Templates) RenderSleepSummaries(w io.Writer, sleepData *sleep.GetSleepSummaryOutput, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(FS, "base.gohtml", "sleepsummaries.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", SleepSummariesVars{
		SleepData: sleepData,
		Error:     errMsg,
	})
}

type SubscriptionsPageVars struct {
	Error         string
	Subscriptions []*subscription.Subscription
	Categories    []subscription.NotificationCategory
}

func (t *Templates) RenderSubscriptionsPage(w io.Writer, subscriptions []*subscription.Subscription, categories []subscription.NotificationCategory, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(FS, "base.gohtml", "subscriptionspage.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", SubscriptionsPageVars{
		Subscriptions: subscriptions,
		Categories:    categories,
		Error:         errMsg,
	})
}

type SubscriptionsWithingsPageVars struct {
	Error                 string
	WithingsSubscriptions []SubscriptionsWithingsPageItem
}

type SubscriptionsWithingsPageItem struct {
	Appli            int
	AppliDescription string
	Exists           bool
	Comment          string
}

func (t *Templates) RenderSubscriptionsWithingsPage(w io.Writer, withingsSubscriptions []SubscriptionsWithingsPageItem, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(FS, "base.gohtml", "subscriptionswithingspage.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", SubscriptionsWithingsPageVars{
		WithingsSubscriptions: withingsSubscriptions,
		Error:                 errMsg,
	})
}
