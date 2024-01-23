package presenters

import (
	"github.com/MikeMwita/go-strict/models"
	"html/template"
	"net/http"
)

// WebTemplate  renders the linting results as HTML

type WebTemplate struct {
	template *template.Template // the HTML template
}

// Render renders the linting results as HTML to the given response writer
func (wt *WebTemplate) Render(w http.ResponseWriter, results []*models.LintResult) error {
	// execute the HTML template with the results
	err := wt.template.Execute(w, results)
	if err != nil {
		return err
	}
	return nil
}

func NewWebTemplate(templatePath string) (*WebTemplate, error) {
	// parse the HTML template from the given path
	template, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &WebTemplate{
		template: template,
	}, nil
}
