package handlebars_test

import "github.com/BlindGarret/echorend"

type MockTemplateGatherer struct {
	templates []echorend.RawTemplateData
	err       error
}

func NewMockTemplateGatherer() *MockTemplateGatherer {
	return &MockTemplateGatherer{
		templates: make([]echorend.RawTemplateData, 0),
	}
}

func (m *MockTemplateGatherer) MustGather() []echorend.RawTemplateData {
	if m.err != nil {
		panic(m.err)
	}
	return m.templates
}

func (m *MockTemplateGatherer) Gather() ([]echorend.RawTemplateData, error) {
	return m.templates, m.err
}

func (m *MockTemplateGatherer) AddTemplate(template echorend.RawTemplateData) {
	m.templates = append(m.templates, template)
}

func (m *MockTemplateGatherer) SetError(err error) {
	m.err = err
}
