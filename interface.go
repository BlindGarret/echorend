package echorend

import "github.com/labstack/echo/v4"

type RawTemplateData struct {
	TemplateName string
	TemplateData string
}

type RawTemplateGatherer interface {
	MustGather() []RawTemplateData
	Gather() ([]RawTemplateData, error)
}

type Renderer interface {
	echo.Renderer
	Setup() error
	MustSetup()
}
