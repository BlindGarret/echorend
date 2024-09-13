package glob_test

import (
	"testing"

	"github.com/BlindGarret/echorend/glob"
)

func TestGlobGatherer_HappyPathFlatDirectory_ReturnsAsExpected(t *testing.T) {
	templateDir := "templates"
	mockFileAccess := NewMemoryFileAccess()
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		TemplateDir: &templateDir,
		FileAccess:  mockFileAccess,
		Extensions:  []string{".html"},
	})
	mockFileAccess.RegisterGlob("templates/*.html", []string{"templates/file1.html", "templates/file2.html"}, nil)
	mockFileAccess.RegisterFile("templates/file1.html", []byte("file1"), nil)
	mockFileAccess.RegisterFile("templates/file2.html", []byte("file2"), nil)

	files, err := gatherer.Gather()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(files))
	}
	if files[0].TemplateName != "file1" {
		t.Errorf("expected template name file1, got %s", files[0].TemplateName)
	}
	if files[0].TemplateData != "file1" {
		t.Errorf("expected template data file1, got %s", files[0].TemplateData)
	}
	if files[1].TemplateName != "file2" {
		t.Errorf("expected template name file2, got %s", files[1].TemplateName)
	}
	if files[1].TemplateData != "file2" {
		t.Errorf("expected template data file2, got %s", files[1].TemplateData)
	}
}
