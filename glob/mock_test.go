package glob_test

type globResponse struct {
	content []string
	err     error
}

type fileResponse struct {
	content []byte
	err     error
}

// MemoryFileAccess is a mock implementation of FileAccess that stores its responses in memory.
type MemoryFileAccess struct {
	globs map[string]globResponse
	files map[string]fileResponse
}

// NewMemoryFileAccess creates a new MemoryFileAccess.
func NewMemoryFileAccess() *MemoryFileAccess {
	return &MemoryFileAccess{
		globs: make(map[string]globResponse),
		files: make(map[string]fileResponse),
	}
}

func (m *MemoryFileAccess) Glob(pattern string) ([]string, error) {
	return m.globs[pattern].content, m.globs[pattern].err
}

func (m *MemoryFileAccess) ReadFile(filename string) ([]byte, error) {
	return m.files[filename].content, m.files[filename].err
}

func (m *MemoryFileAccess) RegisterGlob(pattern string, content []string, err error) {
	m.globs[pattern] = globResponse{content: content, err: err}
}

func (m *MemoryFileAccess) RegisterFile(filename string, content []byte, err error) {
	m.files[filename] = fileResponse{content: content, err: err}
}
