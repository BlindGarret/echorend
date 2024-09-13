package glob

import (
	"path/filepath"
	"strings"

	"github.com/BlindGarret/echorend"
	"github.com/BlindGarret/echorend/externals"
)

type GlobGathererConfig struct {
	TemplateDir     *string
	FileAccess      externals.FileAccess
	IncludeTLDInKey bool
	Extensions      []string
}

type globGatherer struct {
	config GlobGathererConfig
}

func NewGlobGatherer(config GlobGathererConfig) echorend.RawTemplateGatherer {
	return &globGatherer{
		config: defaultGlobGathererConfig(config),
	}
}

func (g *globGatherer) MustGather() []echorend.RawTemplateData {
	templates, err := g.Gather()
	if err != nil {
		panic(err)
	}
	return templates
}

func (g *globGatherer) Gather() ([]echorend.RawTemplateData, error) {
	templates := make([]echorend.RawTemplateData, 0)

	for _, extension := range g.config.Extensions {
		files, err := getTemplateFiles(extension, *g.config.TemplateDir, g.config.FileAccess)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			templateName := getTemplateName(file, *g.config.TemplateDir)
			if g.config.IncludeTLDInKey {
				templateName = *g.config.TemplateDir + "/" + templateName
			}
			bs, err := g.config.FileAccess.ReadFile(file)
			if err != nil {
				return nil, err
			}
			data := echorend.RawTemplateData{
				TemplateName: templateName,
				TemplateData: string(bs),
			}
			templates = append(templates, data)
		}
	}

	return templates, nil
}

func defaultGlobGathererConfig(config GlobGathererConfig) GlobGathererConfig {
	if config.TemplateDir == nil {
		tld := "templates/views"
		config.TemplateDir = &tld
	}

	if config.FileAccess == nil {
		config.FileAccess = &externals.StdFileAccess{}
	}

	return config
}

func getTemplateFiles(extension string, topLevelDirectory string, fileAccess externals.FileAccess) ([]string, error) {
	files, err := fileAccess.Glob(topLevelDirectory + "/*" + extension)
	if err != nil {
		return nil, err
	}

	recursiveMatches, err := fileAccess.Glob(topLevelDirectory + "/**/*" + extension)
	if err != nil {
		return nil, err
	}
	files = append(files, recursiveMatches...)
	return files, nil
}

func getTemplateName(path string, tld string) string {
	path = strings.Replace(path, tld, "", 1)
	if strings.HasPrefix(path, "/") {
		// Remove leading slash
		path = path[1:]
	}
	return path[:len(path)-len(filepath.Ext(path))]
}
