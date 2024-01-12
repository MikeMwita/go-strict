package presenters

import (
	"github.com/MikeMwita/go-strict/datamodels"
	"html/template"
	"net/http"
)

// WebTemplate is a presenter that renders the linting results as HTML

type WebTemplate struct {
	template *template.Template // the HTML template
}

// NewWebTemplate creates a new WebTemplate

func NewWebTemplate(templatePath string) (*WebTemplate, error) {
	// parse the HTML template from the given path
	template, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	// return a new WebTemplate
	return &WebTemplate{
		template: template,
	}, nil
}

// Render renders the linting results as HTML to the given response writer

func (wt *WebTemplate) Render(w http.ResponseWriter, results []*datamodels.LintResult) error {
	// execute the HTML template with the results
	err := wt.template.Execute(w, results)
	if err != nil {
		return err
	}
	// return nil if no error
	return nil
}
