package handlebars

import (
	"fmt"
	"io"

	"github.com/BlindGarret/echorend"
	"github.com/aymerick/raymond"
	"github.com/labstack/echo/v4"
)

type HandlebarsRenderer struct {
	templates       map[string]*raymond.Template
	viewGatherer    echorend.RawTemplateGatherer
	partialGatherer echorend.RawTemplateGatherer
}

func NewHandlebarsRenderer(
	viewGatherer echorend.RawTemplateGatherer,
	partialsGatherer echorend.RawTemplateGatherer,
) *HandlebarsRenderer {
	return &HandlebarsRenderer{
		templates:       make(map[string]*raymond.Template),
		viewGatherer:    viewGatherer,
		partialGatherer: partialsGatherer,
	}
}

func (r *HandlebarsRenderer) Setup() error {
	if r.viewGatherer != nil {
		views, err := r.viewGatherer.Gather()
		if err != nil {
			return err
		}
		for _, view := range views {
			tmpl, err := raymond.Parse(view.TemplateData)
			if err != nil {
				return err
			}
			r.templates[view.TemplateName] = tmpl
		}
	}

	if r.partialGatherer != nil {
		partials, err := r.partialGatherer.Gather()
		if err != nil {
			return err
		}
		for _, partial := range partials {
			tmpl, err := raymond.Parse(partial.TemplateData)
			if err != nil {
				return err
			}
			raymond.RegisterPartialTemplate(partial.TemplateName, tmpl)
		}
	}

	return nil
}

func (r *HandlebarsRenderer) MustSetup() {
	if err := r.Setup(); err != nil {
		panic(err)
	}
}

func (r *HandlebarsRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}

	str, err := tmpl.Exec(data)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(str))
	return err
}
