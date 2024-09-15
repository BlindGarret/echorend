package glob_test

import (
	"errors"
	"testing"

	"github.com/BlindGarret/echorend"
	"github.com/BlindGarret/echorend/gatherers/glob"
)

func TestGlobGatherer_Interface_CompliesWithRawTemplateGatherer(t *testing.T) {
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{})
	_, ok := interface{}(gatherer).(echorend.RawTemplateGatherer)
	if !ok {
		t.Fatalf("GlobGatherer does not comply with RawTemplateGatherer interface")
	}
}

func TestGlobGatherer_HappyPathNoTLD_ReturnsAsExpected(t *testing.T) {
	templateDir := "templates"
	mockFileAccess := NewMemoryFileAccess()
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		TemplateDir: &templateDir,
		FileAccess:  mockFileAccess,
		Extensions:  []string{".html"},
	})
	files := []struct {
		filePath     string
		expectedName string
		expectedData string
		flat         bool
	}{
		{"templates/file1.html", "file1", "file1", true},
		{"templates/file2.html", "file2", "file2", true},
		{"templates/nested/file3.html", "nested/file3", "file3", false},
	}
	flatFilePaths := make([]string, 0)
	nestedFilePaths := make([]string, 0)
	for _, f := range files {
		if f.flat {
			flatFilePaths = append(flatFilePaths, f.filePath)
		} else {
			nestedFilePaths = append(nestedFilePaths, f.filePath)
		}
		mockFileAccess.RegisterFile(f.filePath, []byte(f.expectedData), nil)
	}
	mockFileAccess.RegisterGlob("templates/*.html", flatFilePaths, nil)
	mockFileAccess.RegisterGlob("templates/**/*.html", nestedFilePaths, nil)

	fs, err := gatherer.Gather()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fs) != len(files) {
		t.Fatalf("expected %d templates, got %d", len(files), len(fs))
	}
	for i, f := range fs {
		if f.TemplateName != files[i].expectedName {
			t.Errorf("expected template name %s, got %s", files[i].expectedName, f.TemplateName)
		}
		if f.TemplateData != files[i].expectedData {
			t.Errorf("expected template data %s, got %s", files[i].expectedData, f.TemplateData)
		}
	}
}

func TestGlobGatherer_HappyPathWithTLD_ReturnsAsExpected(t *testing.T) {
	templateDir := "templates"
	mockFileAccess := NewMemoryFileAccess()
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		TemplateDir:     &templateDir,
		FileAccess:      mockFileAccess,
		IncludeTLDInKey: true,
		Extensions:      []string{".html"},
	})
	files := []struct {
		filePath     string
		expectedName string
		expectedData string
		flat         bool
	}{
		{"templates/file1.html", "templates/file1", "file1", true},
		{"templates/file2.html", "templates/file2", "file2", true},
		{"templates/nested/file3.html", "templates/nested/file3", "file3", false},
	}
	flatFilePaths := make([]string, 0)
	nestedFilePaths := make([]string, 0)
	for _, f := range files {
		if f.flat {
			flatFilePaths = append(flatFilePaths, f.filePath)
		} else {
			nestedFilePaths = append(nestedFilePaths, f.filePath)
		}
		mockFileAccess.RegisterFile(f.filePath, []byte(f.expectedData), nil)
	}
	mockFileAccess.RegisterGlob("templates/*.html", flatFilePaths, nil)
	mockFileAccess.RegisterGlob("templates/**/*.html", nestedFilePaths, nil)

	fs, err := gatherer.Gather()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fs) != len(files) {
		t.Fatalf("expected %d templates, got %d", len(files), len(fs))
	}
	for i, f := range fs {
		if f.TemplateName != files[i].expectedName {
			t.Errorf("expected template name %s, got %s", files[i].expectedName, f.TemplateName)
		}
		if f.TemplateData != files[i].expectedData {
			t.Errorf("expected template data %s, got %s", files[i].expectedData, f.TemplateData)
		}
	}
}

func TestGlobGatherer_ErrorGlobbingRecursive_ReturnsError(t *testing.T) {
	templateDir := "templates"
	mockFileAccess := NewMemoryFileAccess()
	expectedErr := errors.New("test error")
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		TemplateDir: &templateDir,
		FileAccess:  mockFileAccess,
		Extensions:  []string{".html"},
	})
	mockFileAccess.RegisterGlob("templates/*.html", []string{}, nil)
	mockFileAccess.RegisterGlob("templates/**/*.html", nil, expectedErr)

	_, err := gatherer.Gather()

	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestGlobGatherer_ErrorGlobbing_ReturnsError(t *testing.T) {
	templateDir := "templates"
	mockFileAccess := NewMemoryFileAccess()
	expectedErr := errors.New("test error")
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		TemplateDir: &templateDir,
		FileAccess:  mockFileAccess,
		Extensions:  []string{".html"},
	})
	mockFileAccess.RegisterGlob("templates/*.html", nil, expectedErr)

	_, err := gatherer.Gather()

	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestGlobGatherer_ErrorReadingFile_ReturnsError(t *testing.T) {
	templateDir := "templates"
	mockFileAccess := NewMemoryFileAccess()
	expectedErr := errors.New("test error")
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		TemplateDir: &templateDir,
		FileAccess:  mockFileAccess,
		Extensions:  []string{".html"},
	})
	mockFileAccess.RegisterGlob("templates/*.html", []string{"templates/file1.html"}, nil)
	mockFileAccess.RegisterFile("templates/file1.html", nil, expectedErr)

	_, err := gatherer.Gather()

	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestGlobGatherer_HappyPathWithRelativePaths_ReturnsAsExpected(t *testing.T) {
	templateDir := "./templates"
	mockFileAccess := NewMemoryFileAccess()
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		TemplateDir: &templateDir,
		FileAccess:  mockFileAccess,
		Extensions:  []string{".html"},
	})
	files := []struct {
		filePath     string
		expectedName string
		expectedData string
		flat         bool
	}{
		{"./templates/file1.html", "file1", "file1", true},
		{"./templates/file2.html", "file2", "file2", true},
		{"./templates/nested/file3.html", "nested/file3", "file3", false},
	}
	flatFilePaths := make([]string, 0)
	nestedFilePaths := make([]string, 0)
	for _, f := range files {
		if f.flat {
			flatFilePaths = append(flatFilePaths, f.filePath)
		} else {
			nestedFilePaths = append(nestedFilePaths, f.filePath)
		}
		mockFileAccess.RegisterFile(f.filePath, []byte(f.expectedData), nil)
	}
	mockFileAccess.RegisterGlob("./templates/*.html", flatFilePaths, nil)
	mockFileAccess.RegisterGlob("./templates/**/*.html", nestedFilePaths, nil)

	fs, err := gatherer.Gather()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fs) != len(files) {
		t.Fatalf("expected %d templates, got %d", len(files), len(fs))
	}
	for i, f := range fs {
		if f.TemplateName != files[i].expectedName {
			t.Errorf("expected template name %s, got %s", files[i].expectedName, f.TemplateName)
		}
		if f.TemplateData != files[i].expectedData {
			t.Errorf("expected template data %s, got %s", files[i].expectedData, f.TemplateData)
		}
	}
}

func TestGlobGatherer_HappyPathMustGatherWithDefaultTLD_ReturnsAsExpected(t *testing.T) {
	mockFileAccess := NewMemoryFileAccess()
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		FileAccess: mockFileAccess,
		Extensions: []string{".html"},
	})
	files := []struct {
		filePath     string
		expectedName string
		expectedData string
		flat         bool
	}{
		{"templates/views/file1.html", "file1", "file1", true},
		{"templates/views/file2.html", "file2", "file2", true},
		{"templates/views/nested/file3.html", "nested/file3", "file3", false},
	}
	flatFilePaths := make([]string, 0)
	nestedFilePaths := make([]string, 0)
	for _, f := range files {
		if f.flat {
			flatFilePaths = append(flatFilePaths, f.filePath)
		} else {
			nestedFilePaths = append(nestedFilePaths, f.filePath)
		}
		mockFileAccess.RegisterFile(f.filePath, []byte(f.expectedData), nil)
	}
	mockFileAccess.RegisterGlob("templates/views/*.html", flatFilePaths, nil)
	mockFileAccess.RegisterGlob("templates/views/**/*.html", nestedFilePaths, nil)

	fs := gatherer.MustGather()

	if len(fs) != len(files) {
		t.Fatalf("expected %d templates, got %d", len(files), len(fs))
	}
	for i, f := range fs {
		if f.TemplateName != files[i].expectedName {
			t.Errorf("expected template name %s, got %s", files[i].expectedName, f.TemplateName)
		}
		if f.TemplateData != files[i].expectedData {
			t.Errorf("expected template data %s, got %s", files[i].expectedData, f.TemplateData)
		}
	}
}

func TestGlobGatherer_ErrorMustGather_Panics(t *testing.T) {
	mockFileAccess := NewMemoryFileAccess()
	gatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
		FileAccess: mockFileAccess,
		Extensions: []string{".html"},
	})
	mockFileAccess.RegisterGlob("templates/views/*.html", nil, errors.New("test error"))

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got nil")
		}
	}()

	gatherer.MustGather()
}
