package echorend

import "github.com/labstack/echo/v4"

// RawTemplateData is a data struct for passing around template data
type RawTemplateData struct {
	TemplateName string
	TemplateData string
}

// RawTemplateGatherer is the interface for implementing Gatherers for the renderer to use during setup.
type RawTemplateGatherer interface {
	MustGather() []RawTemplateData
	Gather() ([]RawTemplateData, error)
}

// Renderer is the interface for implementing renderers for the echo framework.
type Renderer interface {
	echo.Renderer
	Setup() error
	MustSetup()
}
