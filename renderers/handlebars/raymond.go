package handlebars

import (
	"bytes"
	"fmt"
	"io"

	"github.com/BlindGarret/echorend"
	"github.com/aymerick/raymond"
	"github.com/labstack/echo/v4"
)

// HandlebarsRenderer is a renderer that uses the raymond library to render Handlebars templates.
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

// Setup initializes the renderer by gathering templates from the view and partial gatherers and parsing them for render calls.
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

// MustSetup initializes the renderer by gathering templates from the view and partial gatherers
// and parsing them for render calls. If an error occurs, it panics.
func (r *HandlebarsRenderer) MustSetup() {
	if err := r.Setup(); err != nil {
		panic(err)
	}
}

// Render renders a template with the given name and daata to the IO writer.
// this function is designed to slot directly into echo as a renderer
func (r *HandlebarsRenderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
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

func (r *HandlebarsRenderer) CheckRenders() []error {
	errs := make([]error, 0)
	for name, _ := range r.templates {
		buf := new(bytes.Buffer)
		err := r.Render(buf, name, nil, nil)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
