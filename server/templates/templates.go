package templates

import (
	"embed"
	"github.com/roessland/withoutings/withings"
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
	Token *withings.Token
}

func (t Templates) RenderHomePage(w io.Writer, token *withings.Token) error {
	return t.template.ExecuteTemplate(w, "homepage.tmpl", HomePageVars{
		Token: token,
	})
}
