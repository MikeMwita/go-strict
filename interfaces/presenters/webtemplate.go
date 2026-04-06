package presenters

import (
	"html/template"
	"net/http"

	"github.com/MikeMwita/go-strict/models"
)

// WebTemplate renders a LintReport as HTML.
type WebTemplate struct {
	template *template.Template
}

func (wt *WebTemplate) Render(w http.ResponseWriter, report *models.LintReport) error {
	return wt.template.Execute(w, report)
}

func NewWebTemplate(templatePath string) (*WebTemplate, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	return &WebTemplate{template: t}, nil
}
