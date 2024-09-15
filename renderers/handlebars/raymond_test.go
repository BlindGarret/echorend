package handlebars_test

import (
	"testing"

	"github.com/BlindGarret/echorend"
	"github.com/BlindGarret/echorend/renderers/handlebars"
)

func TestHandlebarsRenderer_Interface_CompliesWithRenderer(t *testing.T) {
	renderer := handlebars.NewHandlebarsRenderer(nil, nil)
	_, ok := interface{}(renderer).(echorend.Renderer)
	if !ok {
		t.Errorf("HandlebarsRenderer does not comply with the Renderer interface")
	}
}
