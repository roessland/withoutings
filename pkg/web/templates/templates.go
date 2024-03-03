package templates

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"runtime"
	"time"

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

type LoginPageVars struct {
	Context TemplateContext
	Error   string
}

func (t *Templates) RenderLoginPage(ctx context.Context, w io.Writer, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "login.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", LoginPageVars{
		Context: extractTemplateContext(ctx),
		Error:   errMsg,
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

type MeasureGetmeasPageVars struct {
	Context             TemplateContext
	Error               string
	GetmeasResponseBody string
}

func (t *Templates) RenderMeasureGetmeas(ctx context.Context, w io.Writer, rawResp string, errMsg string) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "measuregetmeaspage.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	return tmpl.ExecuteTemplate(w, "base.gohtml", MeasureGetmeasPageVars{
		Context:             extractTemplateContext(ctx),
		GetmeasResponseBody: rawResp,
		Error:               errMsg,
	})
}

type NotificationsPageVars struct {
	Context       TemplateContext
	Error         string
	Notifications []NotificationPageNotification
}

type NotificationPageNotification struct {
	NotificationUUID string
	ReceivedAt       string
	Params           string
	DataStatus       string
	FetchedAt        string
	Data             []NotificationPageNotificationData
	Appli            string
	AppliDescription string
}

type NotificationPageNotificationData struct {
	Service    string
	Data       string
	DataPretty string
}

func (t *Templates) RenderNotifications(
	ctx context.Context,
	w io.Writer,
	notifications []*subscription.Notification,
	notificationData [][]*subscription.NotificationData,
	errMsg string,
) error {
	tmpl, err := template.New("base.gohtml").ParseFS(t.FS, "base.gohtml", "notifications.gohtml")
	if err != nil {
		panic(err.Error())
	}
	tmpl.Option("missingkey=error")
	log := logging.GetOrCreateLoggerFromContext(ctx)
	log.WithField("notifications", fmt.Sprintf(`%v`, notifications)).Debug("Rendering notifications")
	tmplNotifications := make([]NotificationPageNotification, 0)

	for i, n := range notifications {
		tn := NotificationPageNotification{
			NotificationUUID: n.UUID().String(),
			ReceivedAt:       n.ReceivedAt().Format(time.RFC3339),
			Params:           n.Params(),
			DataStatus:       string(n.DataStatus()),
			FetchedAt:        "see below",
			Data:             mapNotificationDataToTemplateData(notificationData[i]),
			Appli:            "<appli missing>",
			AppliDescription: "<AppliDescription missing>",
		}
		if n.FetchedAt() != nil {
			tn.FetchedAt = n.FetchedAt().Format(time.RFC3339)
		}

		if params, err := subscription.ParseNotificationParams(n.Params()); err == nil {
			desc := subscription.NotificationCategoryByAppli[params.Appli].Description
			tn.Appli = params.AppliStr
			tn.AppliDescription = desc
		}
		tmplNotifications = append(tmplNotifications, tn)
	}
	return tmpl.ExecuteTemplate(w, "base.gohtml", NotificationsPageVars{
		Context:       extractTemplateContext(ctx),
		Notifications: tmplNotifications,
		Error:         errMsg,
	})
}

func mapNotificationDataToTemplateData(data []*subscription.NotificationData) []NotificationPageNotificationData {
	tmplData := make([]NotificationPageNotificationData, 0)
	for _, d := range data {
		td := NotificationPageNotificationData{
			Service: string(d.Service()),
			Data:    string(d.Data()),
		}
		if d.Data() != nil {
			var out bytes.Buffer
			_ = json.Indent(&out, d.Data(), "", "  ")
			td.DataPretty = out.String()
		}
		tmplData = append(tmplData, td)
	}
	return tmplData
}
