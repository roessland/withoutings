package templates

import (
	"context"
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"runtime"

	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

//go:embed templates/*.gohtml
var embeddedFS embed.FS

type Templates struct {
	FS     fs.FS
	source string
}

func (t *Templates) Source() string {
	return t.source
}

type Config struct {
	EmbeddedOnly bool
}

// NewTemplates creates a new templates instance.
// Uses templates from disk if available, otherwise uses embedded templates.
// embeddedOnly: Set to true to force the use of embedded templates (never load from disk).
func NewTemplates(templatesConfig Config) *Templates {
	t := &Templates{}

	// Set FS to use embedded files.
	var err error
	t.FS, err = fs.Sub(embeddedFS, "templates")
	t.source = "embedded"
	if err != nil {
		panic(err.Error())
	}

	// If embeddedOnly is set, don't allow loading templates from disk.
	if templatesConfig.EmbeddedOnly {
		return t
	}

	// If the templates are available on disk, set FS to use disk files instead.
	_, templatesGoPath, _, _ := runtime.Caller(0)
	templatesDir := path.Dir(templatesGoPath)
	stat, err := os.Stat(path.Join(templatesDir, "templates/base.gohtml"))
	if err == nil && stat.Size() > 0 {
		t.FS = os.DirFS(path.Join(templatesDir, "templates"))
		t.source = "disk"
	}

	return t
}

type HomePageVars struct {
	Account *account.Account
	Error   string
	Context TemplateContext
}

func (t *Templates) RenderHomePage(ctx context.Context, w io.Writer, account_ *account.Account) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "homepage.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", HomePageVars{
		Context: extractTemplateContext(ctx),
		Account: account_,
	})
}

type RefreshAccessTokenVars struct {
	Token   *withings.Token
	Error   string
	Context TemplateContext
}

func (t *Templates) RenderRefreshAccessToken(ctx context.Context, w io.Writer, token *withings.Token, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "refreshaccesstoken.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", RefreshAccessTokenVars{
		Context: extractTemplateContext(ctx),
		Token:   token,
		Error:   errMsg,
	})
}

type SleepSummariesVars struct {
	SleepData interface{}
	Token     *withings.Token
	Error     string
	Context   TemplateContext
}

func (t *Templates) RenderSleepSummaries(ctx context.Context, w io.Writer, sleepData *sleep.GetSleepSummaryOutput, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "sleepsummaries.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", SleepSummariesVars{
		Context:   extractTemplateContext(ctx),
		SleepData: sleepData,
		Error:     errMsg,
	})
}

type SubscriptionsPageVars struct {
	Context       TemplateContext
	Error         string
	Subscriptions []*subscription.Subscription
	Categories    []subscription.NotificationCategory
}

func (t *Templates) RenderSubscriptionsPage(ctx context.Context, w io.Writer, subscriptions []*subscription.Subscription, categories []subscription.NotificationCategory, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "subscriptionspage.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", SubscriptionsPageVars{
		Context:       extractTemplateContext(ctx),
		Subscriptions: subscriptions,
		Categories:    categories,
		Error:         errMsg,
	})
}

type SubscriptionsWithingsPageVars struct {
	Context               TemplateContext
	Error                 string
	WithingsSubscriptions []SubscriptionsWithingsPageItem
}

type SubscriptionsWithingsPageItem struct {
	AppliDescription string
	Comment          string
	Appli            int
	Exists           bool
}

func (t *Templates) RenderSubscriptionsWithingsPage(ctx context.Context, w io.Writer, withingsSubscriptions []SubscriptionsWithingsPageItem, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "subscriptionswithingspage.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", SubscriptionsWithingsPageVars{
		Context:               extractTemplateContext(ctx),
		WithingsSubscriptions: withingsSubscriptions,
		Error:                 errMsg,
	})
}

type TemplateTestVars struct {
	Error   string
	Title   string
	Content string
	Context TemplateContext
}

func (t *Templates) RenderTemplateTest(ctx context.Context, w io.Writer) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "templatetest.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", TemplateTestVars{
		Context: TemplateContext{},
		Error:   "ThisIsTheError",
		Title:   "ThisIsTheTitle",
		Content: "ThisIsTheContent",
	})
}
