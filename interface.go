package echorend

type RawTemplateData struct {
	TemplateName string
	TemplateData string
}

type RawTemplateGatherer interface {
	MustGather() []RawTemplateData
	Gather() ([]RawTemplateData, error)
}
