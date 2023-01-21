package templates

import (
	"embed"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/clients/withingsapi"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"html/template"
	"io"
)

//go:embed templates
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

type SleepSummariesVars struct {
	Error     string
	Token     *withingsapi.Token
	SleepData interface{}
}

func (t Templates) RenderSleepSummaries(w io.Writer, sleepData *sleep.GetSleepSummaryOutput, err string) error {
	return t.template.ExecuteTemplate(w, "sleepsummaries.gohtml", SleepSummariesVars{
		SleepData: sleepData,
		Error:     err,
	})
}

type RefreshAccessTokenVars struct {
	Error            string
	Token, PrevToken *withingsapi.Token
}

func (t Templates) RenderRefreshAccessToken(w io.Writer, token, prevToken *withingsapi.Token) error {
	return t.template.ExecuteTemplate(w, "refreshaccesstoken.gohtml", RefreshAccessTokenVars{
		Token:     token,
		PrevToken: prevToken,
	})
}
