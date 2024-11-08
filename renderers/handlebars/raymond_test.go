package handlebars_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/BlindGarret/echorend"
	"github.com/BlindGarret/echorend/renderers/handlebars"
)

func renderToString(name string, data interface{}, renderer *handlebars.HandlebarsRenderer) (string, error) {
	buf := new(bytes.Buffer)

	err := renderer.Render(buf, name, data, nil)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestHandlebarsRenderer_Interface_CompliesWithRenderer(t *testing.T) {
	renderer := handlebars.NewHandlebarsRenderer(nil, nil)
	_, ok := interface{}(renderer).(echorend.Renderer)
	if !ok {
		t.Errorf("HandlebarsRenderer does not comply with the Renderer interface")
	}
}

func TestHandlebarsRendererSetup_CalledWithValidTemplates_NoErr(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	viewGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-view1",
		TemplateData: "<HTML></HTML>",
	})
	partialGatherer := NewMockTemplateGatherer()
	partialGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-partial1",
		TemplateData: "<h1>test</h1>",
	})

	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	err := renderer.Setup()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHandlebarsRendererSetup_ViewGathererErrors_ReturnsError(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	expectedErr := errors.New("test error")
	viewGatherer.SetError(expectedErr)
	partialGatherer := NewMockTemplateGatherer()

	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	err := renderer.Setup()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestHandlebarsRendererSetup_PartialGathererErrors_ReturnsError(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	partialGatherer := NewMockTemplateGatherer()
	expectedErr := errors.New("test error")
	partialGatherer.SetError(expectedErr)

	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	err := renderer.Setup()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestHandlebarsRendererSetup_ViewGathererHasBadTemplates_ReturnsError(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	viewGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-view2",
		TemplateData: "<HTML>{{herp}</HTML>",
	})
	partialGatherer := NewMockTemplateGatherer()
	partialGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-partial2",
		TemplateData: "<h1>test</h1>",
	})

	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	err := renderer.Setup()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestHandlebarsRendererSetup_PartialGathererHasBadTemplates_ReturnsError(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	viewGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-view3",
		TemplateData: "<HTML></HTML>",
	})
	partialGatherer := NewMockTemplateGatherer()
	partialGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-partial3",
		TemplateData: "<h1>{{herp}</h1>",
	})

	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	err := renderer.Setup()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestHandlebarsRendererMustSetup_CalledWithValidTemplates_NoPanic(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	viewGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-view4",
		TemplateData: "<HTML></HTML>",
	})
	partialGatherer := NewMockTemplateGatherer()
	partialGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-partial4",
		TemplateData: "<h1>test</h1>",
	})

	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	renderer.MustSetup()
}

func TestHandlebarsRendererMustSetup_ViewGathererErrors_Panics(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	expectedErr := errors.New("test error")
	viewGatherer.SetError(expectedErr)
	partialGatherer := NewMockTemplateGatherer()

	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got nil")
		}
	}()
	renderer.MustSetup()
}

func TestHandlebarsRendererRender_TemplateNotFound_Errors(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	partialGatherer := NewMockTemplateGatherer()
	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	renderer.MustSetup()

	_, err := renderToString("test-view5", nil, renderer)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestHandlebarsRendererRender_TemplateFound_NoErr(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	viewGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-view6",
		TemplateData: "<HTML></HTML>",
	})
	partialGatherer := NewMockTemplateGatherer()
	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	renderer.MustSetup()

	_, err := renderToString("test-view6", nil, renderer)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHandlebarsRendererRender_TemplateMissingData_Errors(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	viewGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-view7",
		TemplateData: "<HTML>{{>non-existant-partial}}</HTML>",
	})
	partialGatherer := NewMockTemplateGatherer()
	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	renderer.MustSetup()

	_, err := renderToString("test-view7", nil, renderer)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestHandlebarsCheckRenders_RunWithValidTemplates_ReturnsNoErrors(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	viewGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-view8",
		TemplateData: "<HTML>{{> test-partial8}}</HTML>",
	})
	partialGatherer := NewMockTemplateGatherer()
	partialGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-partial8",
		TemplateData: "<h1>test</h1>",
	})
	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	renderer.MustSetup()

	errs := renderer.CheckRenders()

	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d", len(errs))
		for _, err := range errs {
			t.Errorf("Error: %v", err)
		}
	}
}

func TestHandlebarsCheckRenders_RunWithInvalidTemplates_ReturnsErrors(t *testing.T) {
	viewGatherer := NewMockTemplateGatherer()
	viewGatherer.AddTemplate(echorend.RawTemplateData{
		TemplateName: "test-view9",
		TemplateData: "<HTML>{{> test-partial9}}</HTML>",
	})
	partialGatherer := NewMockTemplateGatherer()
	renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
	renderer.MustSetup()

	errs := renderer.CheckRenders()

	if len(errs) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errs))
	}
}
